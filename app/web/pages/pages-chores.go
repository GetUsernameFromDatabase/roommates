package pages

import (
	"context"
	"slices"
	"strings"
	"time"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgtype"
)

// ChoreForm is what a chore form shows besides the chore: what was typed, the error of the last
// submit and where it belongs, and whether the new chore section is open
type ChoreForm struct {
	Input   services.ChoreInput
	Message string
	// a services.FieldChore* name, "" for an error that belongs to no field
	Field string
	Saved bool
	Open  bool
}

// choreWeek is the weekdays in the order the form shows them
var choreWeek = []time.Weekday{
	time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday,
}

var choreWeekdayLong = map[time.Weekday]locales.LK{
	time.Monday:    locales.LK_ChoresWeekdayMon,
	time.Tuesday:   locales.LK_ChoresWeekdayTue,
	time.Wednesday: locales.LK_ChoresWeekdayWed,
	time.Thursday:  locales.LK_ChoresWeekdayThu,
	time.Friday:    locales.LK_ChoresWeekdayFri,
	time.Saturday:  locales.LK_ChoresWeekdaySat,
	time.Sunday:    locales.LK_ChoresWeekdaySun,
}

var choreWeekdayShort = map[time.Weekday]locales.LK{
	time.Monday:    locales.LK_ChoresWeekdayShortMon,
	time.Tuesday:   locales.LK_ChoresWeekdayShortTue,
	time.Wednesday: locales.LK_ChoresWeekdayShortWed,
	time.Thursday:  locales.LK_ChoresWeekdayShortThu,
	time.Friday:    locales.LK_ChoresWeekdayShortFri,
	time.Saturday:  locales.LK_ChoresWeekdayShortSat,
	time.Sunday:    locales.LK_ChoresWeekdayShortSun,
}

// choreEveryValue is the schedule control's value for an "every N" chore, its unit is the select's
const choreEveryValue = "every"

var choreScheduleOptions = []segmentOption{
	{Value: g.ChoreUnitWeekdays, Label: locales.LK_ChoresScheduleWeekdaysOption},
	{Value: choreEveryValue, Label: locales.LK_ChoresScheduleEveryOption},
}

var choreUnitKeys = map[string]locales.LK{
	g.ChoreUnitDay:   locales.LK_ChoresUnitDay,
	g.ChoreUnitWeek:  locales.LK_ChoresUnitWeek,
	g.ChoreUnitMonth: locales.LK_ChoresUnitMonth,
}

var choreRotationKeys = map[string][2]locales.LK{
	g.RotationFixed:       {locales.LK_ChoresRotationFixed, locales.LK_ChoresRotationFixedHint},
	g.RotationLeastRecent: {locales.LK_ChoresRotationLeastRecent, locales.LK_ChoresRotationLeastRecentHint},
	g.RotationManual:      {locales.LK_ChoresRotationManual, locales.LK_ChoresRotationManualHint},
}

func choreURL(householdID pgtype.UUID, parts ...string) templ.SafeURL {
	return householdURL(householdID, append([]string{"chores"}, parts...)...)
}

// choreDay is the calendar day of t, so a due date and the household's today compare by day
func choreDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// choreDue is when the chore is next due, "Today", "Tomorrow" or the date, overdue for a past day
func choreDue(ctx context.Context, c services.Chore) (label string, overdue bool) {
	days := int(choreDay(c.NextDue).Sub(choreDay(c.Today)).Hours() / 24)
	date := utils.FormatDate(ctx, c.NextDue)
	switch {
	case days < 0:
		return utils.T(ctx, locales.LK_ChoresDueOverdue, date), true
	case days == 0:
		return utils.T(ctx, locales.LK_ChoresDueToday), false
	case days == 1:
		return utils.T(ctx, locales.LK_ChoresDueTomorrow), false
	}
	return utils.T(ctx, locales.LK_ChoresDueDate, date), false
}

// choreSchedule is a short line of how often the chore comes round
func choreSchedule(ctx context.Context, c services.Chore) string {
	switch {
	case c.Unit == g.ChoreUnitDay && c.Every == 1:
		return utils.T(ctx, locales.LK_ChoresScheduleDay)
	case c.Unit == g.ChoreUnitDay:
		return utils.T(ctx, locales.LK_ChoresScheduleDays, c.Every)
	case c.Unit == g.ChoreUnitWeek && c.Every == 1:
		return utils.T(ctx, locales.LK_ChoresScheduleWeek)
	case c.Unit == g.ChoreUnitWeek:
		return utils.T(ctx, locales.LK_ChoresScheduleWeeks, c.Every)
	case c.Unit == g.ChoreUnitMonth && c.Every == 1:
		return utils.T(ctx, locales.LK_ChoresScheduleMonth)
	case c.Unit == g.ChoreUnitMonth:
		return utils.T(ctx, locales.LK_ChoresScheduleMonths, c.Every)
	}
	var names []string
	for _, d := range choreWeek {
		for _, w := range c.Weekdays {
			if w == d {
				names = append(names, utils.T(ctx, choreWeekdayShort[d]))
			}
		}
	}
	return utils.T(ctx, locales.LK_ChoresScheduleWeekdays, strings.Join(names, ", "))
}

// choreTurn is whose turn it is, the viewer's own as "Your turn"
func choreTurn(ctx context.Context, actor services.Membership, c services.Chore) string {
	switch {
	case !c.AssigneeID.Valid:
		return utils.T(ctx, locales.LK_ChoresTurnNone)
	case c.AssigneeID == actor.UserID:
		return utils.T(ctx, locales.LK_ChoresTurnYou)
	}
	return utils.T(ctx, locales.LK_ChoresTurn, c.AssigneeName)
}

// choreCanPass says whether the turn can be handed on: not in a manual chore or a rota of one
func choreCanPass(c services.Chore) bool {
	return c.Rotation != g.RotationManual && len(c.Members) > 1
}

// choreLogLine is "Mari did it · 12.11.", for another's turn "Mari did it for Jaan · 12.11."
func choreLogLine(ctx context.Context, e services.ChoreLogEntry) string {
	who := e.UserName
	if who == "" {
		who = utils.T(ctx, locales.LK_ActivitySomeone)
	}
	date := utils.FormatDate(ctx, e.At)
	forOther := e.AssigneeName != "" && e.AssigneeName != who
	switch {
	case e.Kind == g.ChorePassed:
		return utils.T(ctx, locales.LK_ChoresLogPassed, who, date)
	case e.Kind == g.ChoreSkipped && forOther:
		return utils.T(ctx, locales.LK_ChoresLogSkippedFor, who, e.AssigneeName, date)
	case e.Kind == g.ChoreSkipped:
		return utils.T(ctx, locales.LK_ChoresLogSkipped, who, date)
	case forOther:
		return utils.T(ctx, locales.LK_ChoresLogDoneFor, who, e.AssigneeName, date)
	}
	return utils.T(ctx, locales.LK_ChoresLogDone, who, date)
}

// choreMemberPicked says whether userID is in the form's rota
func choreMemberPicked(in services.ChoreInput, userID pgtype.UUID) bool {
	return slices.Contains(in.Members, userID)
}

// choreMemberOrder is the rota in its order and then the other members, the form posts the
// ticked ones in this order, so saving keeps the order set with the arrows
func choreMemberOrder(in services.ChoreInput, members []dbqueries.ListActiveMembersRow) []dbqueries.ListActiveMembersRow {
	pos := func(m dbqueries.ListActiveMembersRow) int {
		if i := slices.Index(in.Members, m.UserID); i >= 0 {
			return i
		}
		return len(in.Members)
	}
	out := slices.Clone(members)
	slices.SortStableFunc(out, func(a, b dbqueries.ListActiveMembersRow) int { return pos(a) - pos(b) })
	return out
}

// choreSwapList is what an action on a card swaps: the list, and the toast outside it
func choreSwapList(action templ.SafeURL) templ.Attributes {
	return templ.Attributes{
		"hx-post":         string(action),
		"hx-target":       "#chores",
		"hx-select":       "#chores",
		"hx-select-oob":   "#chores-notice",
		"hx-swap":         "outerHTML",
		"hx-disabled-elt": "find button",
	}
}

// choreSwapRota is [choreSwapList] for the rota of the chore page, the edit form comes along
// so its members' order and assignee match and a later save does not undo the change
func choreSwapRota(action templ.SafeURL) templ.Attributes {
	return templ.Attributes{
		"hx-post":         string(action),
		"hx-target":       "#chore-rota",
		"hx-select":       "#chore-rota",
		"hx-select-oob":   "#chore-notice, #chore-edit",
		"hx-swap":         "outerHTML",
		"hx-disabled-elt": "find button",
	}
}
