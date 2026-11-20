package pages

import (
	"context"
	"slices"
	"strconv"
	"strings"
	"time"

	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgtype"
)

// EventForm is what an event form shows besides the event: what was typed, the error of the last
// submit and where it belongs, whether the new event section is open and the scope that was picked
type EventForm struct {
	Input   services.EventInput
	Message string
	// a services.FieldEvent* name, "" for an error that belongs to no field
	Field string
	Saved bool
	Open  bool
	// a services.EventScope*, "" for the default
	Scope string
}

// CalendarView is what the calendar page shows: the phone's agenda from today and the desktop's
// month with its picked day
type CalendarView struct {
	Today  time.Time
	Agenda []services.CalendarItem
	// the first day "load more" lists
	AgendaNext time.Time
	// the first day of the month shown and the day picked in it
	Month      time.Time
	Day        time.Time
	MonthItems []services.CalendarItem
}

func eventURL(householdID, eventID pgtype.UUID, on time.Time) templ.SafeURL {
	u := string(householdURL(householdID, "calendar", "events", eventID.String()))
	if !on.IsZero() {
		u += "?on=" + on.Format(services.DateLayout)
	}
	return templ.SafeURL(u)
}

// calendarItemURL is where the item is: its event on the day it starts, its chore or its schedule
func calendarItemURL(householdID pgtype.UUID, it services.CalendarItem) templ.SafeURL {
	switch it.Kind {
	case services.CalendarKindChore:
		return householdURL(householdID, "chores", it.ID.String())
	case services.CalendarKindRecurring:
		return householdURL(householdID, "recurring", it.ID.String())
	}
	return eventURL(householdID, it.ID, it.On)
}

// calendarMoreURL lists the agenda from day, last is the last day already shown, if any
func calendarMoreURL(householdID pgtype.UUID, from, last time.Time) templ.SafeURL {
	u := string(householdURL(householdID, "calendar", "agenda")) + "?from=" + from.Format(services.DateLayout)
	if !last.IsZero() {
		u += "&last=" + last.Format(services.DateLayout)
	}
	return templ.SafeURL(u)
}

// calendarDay is the items of one day
type calendarDay struct {
	Day   time.Time
	Items []services.CalendarItem
}

// calendarDays groups items, sorted by day, by their day, days without items are not there
func calendarDays(items []services.CalendarItem) []calendarDay {
	var days []calendarDay
	for _, it := range items {
		if n := len(days); n > 0 && days[n-1].Day.Equal(it.Date) {
			days[n-1].Items = append(days[n-1].Items, it)
			continue
		}
		days = append(days, calendarDay{Day: it.Date, Items: []services.CalendarItem{it}})
	}
	return days
}

// calendarByDay is the items of each day by its services.DateLayout
func calendarByDay(items []services.CalendarItem) map[string][]services.CalendarItem {
	days := map[string][]services.CalendarItem{}
	for _, it := range items {
		k := it.Date.Format(services.DateLayout)
		days[k] = append(days[k], it)
	}
	return days
}

// calendarAgendaMonth is the days of the agenda in one month, Heading says whether it is named:
// not when an earlier block of the agenda already showed that month
type calendarAgendaMonth struct {
	Month   time.Time
	Heading bool
	Days    []calendarDay
}

func calendarAgendaMonths(items []services.CalendarItem, last time.Time) []calendarAgendaMonth {
	var months []calendarAgendaMonth
	for _, d := range calendarDays(items) {
		m := monthOf(d.Day)
		if n := len(months); n > 0 && months[n-1].Month.Equal(m) {
			months[n-1].Days = append(months[n-1].Days, d)
			continue
		}
		months = append(months, calendarAgendaMonth{Month: m, Heading: len(months) > 0 || last.IsZero() || !monthOf(last).Equal(m), Days: []calendarDay{d}})
	}
	return months
}

// calendarLast is the day of the last item, or last when there are none
func calendarLast(items []services.CalendarItem, last time.Time) time.Time {
	if len(items) == 0 {
		return last
	}
	return items[len(items)-1].Date
}

func monthOf(day time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), 1, 0, 0, 0, 0, time.UTC)
}

var calendarMonthKeys = [...]locales.LK{
	locales.LK_CalendarMonth1, locales.LK_CalendarMonth2, locales.LK_CalendarMonth3, locales.LK_CalendarMonth4,
	locales.LK_CalendarMonth5, locales.LK_CalendarMonth6, locales.LK_CalendarMonth7, locales.LK_CalendarMonth8,
	locales.LK_CalendarMonth9, locales.LK_CalendarMonth10, locales.LK_CalendarMonth11, locales.LK_CalendarMonth12,
}

// calendarMonthName is "November 2026"
func calendarMonthName(ctx context.Context, month time.Time) string {
	return utils.T(ctx, calendarMonthKeys[month.Month()-1]) + " " + strconv.Itoa(month.Year())
}

// calendarWeekdays are the columns of the month grid
var calendarWeekdays = []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}

// CalendarGrid is the first day of the month grid and the day after its last: whole weeks from
// Monday around month
func CalendarGrid(month time.Time) (from, to time.Time) {
	next := month.AddDate(0, 1, 0)
	return month.AddDate(0, 0, -(int(month.Weekday())+6)%7), next.AddDate(0, 0, (7-(int(next.Weekday())+6)%7)%7)
}

// calendarWeeks are the days of the month grid by week
func calendarWeeks(month time.Time) [][]time.Time {
	from, to := CalendarGrid(month)
	var weeks [][]time.Time
	for d := from; d.Before(to); {
		week := make([]time.Time, 7)
		for i := range week {
			week[i], d = d, d.AddDate(0, 0, 1)
		}
		weeks = append(weeks, week)
	}
	return weeks
}

// calendarMonthURL is the calendar on month, with day picked when it is set
func calendarMonthURL(householdID pgtype.UUID, month, day time.Time) templ.SafeURL {
	u := string(householdURL(householdID, "calendar")) + "?month=" + month.Format("2006-01")
	if !day.IsZero() {
		u += "&day=" + day.Format(services.DateLayout)
	}
	return templ.SafeURL(u)
}

// calendarRail is what the agenda shows under a day's number: "Today" or the short weekday
func calendarRail(ctx context.Context, day, today time.Time) string {
	if day.Equal(today) {
		return utils.T(ctx, locales.LK_ChoresDueToday)
	}
	return utils.T(ctx, choreWeekdayShort[day.Weekday()])
}

// calendarChipClock is the time a timed event shows in the month grid, "" for an all-day one
func calendarChipClock(it services.CalendarItem) string {
	switch {
	case it.AllDay:
		return ""
	case it.Clock != "":
		return it.Clock
	case it.EndClock != "":
		return "–" + it.EndClock
	}
	return ""
}

// calendarDayHeading is "Today", "Tomorrow" or "Friday, 21.11.2026"
func calendarDayHeading(ctx context.Context, day, today time.Time) string {
	switch daysFrom(today, day) {
	case 0:
		return utils.T(ctx, locales.LK_ChoresDueToday)
	case 1:
		return utils.T(ctx, locales.LK_ChoresDueTomorrow)
	}
	return calendarFullDay(ctx, day)
}

func calendarFullDay(ctx context.Context, day time.Time) string {
	return utils.T(ctx, choreWeekdayLong[day.Weekday()]) + ", " + utils.FormatDate(ctx, day)
}

func daysFrom(from, to time.Time) int {
	return int(choreDay(to).Sub(choreDay(from)).Hours() / 24)
}

// calendarItemTime is when an event is on the item's day: "18:30–20:00", "All day" or for a day an
// event goes on "until 20:00"
func calendarItemTime(ctx context.Context, it services.CalendarItem) string {
	switch {
	case it.Kind != services.CalendarKindEvent:
		return ""
	case it.AllDay:
		return utils.T(ctx, locales.LK_CalendarAllDay)
	case it.Clock != "" && it.EndClock != "":
		return it.Clock + "–" + it.EndClock
	case it.Clock != "":
		return it.Clock
	case it.EndClock != "":
		return utils.T(ctx, locales.LK_CalendarUntilClock, it.EndClock)
	}
	return utils.T(ctx, locales.LK_CalendarAllDay)
}

// calendarWhen is the date and time of an occurrence: "Fri 21.11.2026 18:30–20:00", over days
// "21.11.2026 18:30 – 22.11.2026 10:00", all-day "21.11.2026, all day"
func calendarWhen(ctx context.Context, o services.Occurrence) string {
	start, end := utils.FormatDate(ctx, o.On), utils.FormatDate(ctx, o.EndOn)
	oneDay := o.On.Equal(o.EndOn)
	switch {
	case o.AllDay && oneDay:
		return calendarFullDay(ctx, o.On) + ", " + utils.T(ctx, locales.LK_CalendarAllDay)
	case o.AllDay:
		return start + " – " + end + ", " + utils.T(ctx, locales.LK_CalendarAllDay)
	case oneDay && o.EndClock != "":
		return calendarFullDay(ctx, o.On) + " " + o.StartClock + "–" + o.EndClock
	case oneDay:
		return calendarFullDay(ctx, o.On) + " " + o.StartClock
	}
	return start + " " + o.StartClock + " – " + end + " " + o.EndClock
}

// calendarRepeat is how often the event comes round, "" for a one-off
func calendarRepeat(ctx context.Context, in services.EventInput) string {
	if in.Unit == "" {
		return ""
	}
	// the chore's wording of a schedule
	var c services.Chore
	c.Unit, c.Every, c.Weekdays = in.Unit, in.Every, in.Weekdays
	s := choreSchedule(ctx, c)
	if in.RepeatUntil != "" {
		if until, err := time.Parse(services.DateLayout, in.RepeatUntil); err == nil {
			s += ", " + utils.T(ctx, locales.LK_CalendarUntil, utils.FormatDate(ctx, until))
		}
	}
	return s
}

// calendarWho is who the event is for
func calendarWho(ctx context.Context, o services.Occurrence) string {
	if len(o.Members) == 0 {
		return utils.T(ctx, locales.LK_CalendarWholeHousehold)
	}
	names := utils.Map(o.Members, func(m services.EventMember) string { return utils.MemberName(ctx, m.Name) })
	return strings.Join(names, ", ")
}

var calendarRemindKeys = map[int]locales.LK{
	g.RemindAtStart:   locales.LK_CalendarRemindStart,
	g.Remind15Min:     locales.LK_CalendarRemind15,
	g.RemindHour:      locales.LK_CalendarRemind60,
	g.RemindDayBefore: locales.LK_CalendarRemindDay,
}

// calendarRemindLabel is a reminder choice, "On the day" for an all-day event's 0
func calendarRemindLabel(ctx context.Context, allDay bool, minutes int) string {
	if allDay && minutes == g.RemindAtStart {
		return utils.T(ctx, locales.LK_CalendarRemindOnDay)
	}
	return utils.T(ctx, calendarRemindKeys[minutes])
}

// calendarRemindValue is the reminder control's value of the timed or the all-day set
func calendarRemindValue(in services.EventInput, allDay bool) string {
	choices := g.RemindTimed
	if allDay {
		choices = g.RemindAllDay
	}
	if in.AllDay == allDay && in.Remind != nil && slices.Contains(choices, *in.Remind) {
		return strconv.Itoa(*in.Remind)
	}
	return "none"
}

// calendarScopes are the scopes of an edit or a delete of a series in the order they are offered
var calendarScopes = []string{services.EventScopeThis, services.EventScopeFollowing, services.EventScopeAll}

var calendarScopeKeys = map[string]locales.LK{
	services.EventScopeThis:      locales.LK_CalendarScopeThis,
	services.EventScopeFollowing: locales.LK_CalendarScopeFollowing,
	services.EventScopeAll:       locales.LK_CalendarScopeAll,
}

// calendarEndOn says whether the form has an end on another day than the start
func calendarEndOn(in services.EventInput) bool {
	return in.EndDate != "" && in.EndDate != in.StartDate
}

// calendarDateOr is date, or fallback when there is none to show yet
func calendarDateOr(date, fallback string) string {
	if date == "" {
		return fallback
	}
	return date
}
