package services

import (
	"context"
	"errors"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/push"
	"roommates/internal/routes"
	"roommates/internal/utils"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ChoreInput is the chore form
type ChoreInput struct {
	Title    string
	Rotation string // g.Rotation*
	Unit     string // g.ChoreUnit*
	// with [g.ChoreEveryUnits]
	Every int
	// with [g.ChoreUnitWeekdays]
	Weekdays []time.Weekday
	// [DateLayout], the first due date of an Every chore, a weekdays chore starts on the next ticked day
	FirstDue string
	// the rota, in order
	Members []pgtype.UUID
	// whose turn it is first, or the manual assignee, not valid for unassigned
	Assignee pgtype.UUID
}

type ChoreMember struct {
	UserID pgtype.UUID
	Name   string
	// when they last did the chore, nil for never
	LastDone *time.Time
}

type Chore struct {
	ID       pgtype.UUID
	Title    string
	Rotation string
	Unit     string
	Every    int
	Weekdays []time.Weekday
	NextDue  time.Time
	// the household's local date, NextDue before it is overdue
	Today time.Time
	// not valid when it is nobody's turn
	AssigneeID   pgtype.UUID
	AssigneeName string
	// the rota, in order
	Members []ChoreMember
}

// ChoreTurn is the turn a page showed, Done, Skip and Pass act only while it is still the chore's:
// a second click or a stale page gets [ErrChoreTurnMoved]
type ChoreTurn struct {
	Due      pgtype.Date
	Assignee pgtype.UUID
}

// Turn is the chore's current turn
func (c Chore) Turn() ChoreTurn {
	return ChoreTurn{Due: dateOf(c.NextDue), Assignee: c.AssigneeID}
}

// ErrChoreTurnMoved is a Done, Skip or Pass on a turn that is no longer the chore's
var ErrChoreTurnMoved = errors.New("the chore's turn moved on")

type ChoreLogEntry struct {
	Kind         string // g.Chore{Done,Skipped,Passed}
	UserName     string
	AssigneeName string
	// someone did, skipped or passed another member's turn
	ForOther bool
	Due      time.Time
	At       time.Time
}

type ChoreDetail struct {
	Chore Chore
	// the last [ChoreHistoryLength] entries, newest first
	History []ChoreLogEntry
}

// ChoreHistoryLength is how many log entries a chore's page shows
const ChoreHistoryLength = 5

// MaxChoreTitleLength is the longest chore title in runes
const MaxChoreTitleLength = 64

// what a [FieldError] of [DefaultService.CreateChore] and [DefaultService.UpdateChore] names
const (
	FieldChoreTitle    = "title"
	FieldChoreRotation = "rotation"
	FieldChoreSchedule = "schedule"
	FieldChoreFirstDue = "first_due"
	FieldChoreMembers  = "members"
	FieldChoreAssignee = "assignee"
)

// the longest Every of each unit, a year
var choreMaxEvery = map[string]int{g.ChoreUnitDay: 365, g.ChoreUnitWeek: 52, g.ChoreUnitMonth: 12}

var choreActivity = map[string]string{g.ChoreDone: g.ActChoreDone, g.ChoreSkipped: g.ActChoreSkipped, g.ChorePassed: g.ActChorePassed}

// weekdayBit is d's bit in a weekdays mask, Monday is bit 0
func weekdayBit(d time.Weekday) int16 {
	return 1 << ((int(d) + 6) % 7)
}

// weekdaysOf decodes a weekdays mask, Monday first
func weekdaysOf(mask int16) []time.Weekday {
	var days []time.Weekday
	for b := range 7 {
		if mask&(1<<b) != 0 {
			days = append(days, time.Weekday((b+1)%7))
		}
	}
	return days
}

// nextWeekday is the first day from from on, from itself included, that mask ticks
func nextWeekday(from time.Time, mask int16) time.Time {
	for i := range 7 {
		if d := from.AddDate(0, 0, i); mask&weekdayBit(d.Weekday()) != 0 {
			return d
		}
	}
	return from
}

// choreNextDue is when a chore due on due and done or skipped today is due again: its next ticked
// weekday after both (done early it still moves past due), or every days (weeks, months) from today,
// a month without today's day gives its last day
func choreNextDue(every *int16, unit *string, weekdays *int16, due, today time.Time) time.Time {
	if weekdays != nil {
		if due.After(today) {
			today = due
		}
		return nextWeekday(today.AddDate(0, 0, 1), *weekdays)
	}
	n := int(*every)
	switch *unit {
	case g.ChoreUnitMonth:
		return addIntervals(today, g.IntervalMonth, n)
	case g.ChoreUnitWeek:
		n *= 7
	}
	return today.AddDate(0, 0, n)
}

// choreFields is a checked [ChoreInput]
type choreFields struct {
	title    string
	every    *int16
	unit     *string
	weekdays *int16
	// the first due date, the next ticked weekday from today for a weekdays chore
	due time.Time
}

// clean checks in against the household's local date today and its active members;
// FirstDue may be before today when it is current, the chore's due date as it is
func (in ChoreInput) clean(today, current pgtype.Date, active []dbqueries.ListActiveMembersRow) (choreFields, error) {
	f := choreFields{title: "", every: nil, unit: nil, weekdays: nil, due: time.Time{}}
	var ok bool
	if f.title, ok = trimToLength(in.Title, 1, MaxChoreTitleLength); !ok {
		return f, FieldError{Field: FieldChoreTitle}
	}
	if !slices.Contains(g.Rotations, in.Rotation) {
		return f, FieldError{Field: FieldChoreRotation}
	}
	switch in.Unit {
	case g.ChoreUnitWeekdays:
		var mask int16
		for _, d := range in.Weekdays {
			if d < time.Sunday || d > time.Saturday {
				return f, FieldError{Field: FieldChoreSchedule}
			}
			mask |= weekdayBit(d)
		}
		if mask == 0 {
			return f, FieldError{Field: FieldChoreSchedule}
		}
		f.weekdays, f.due = &mask, nextWeekday(today.Time, mask)
	case g.ChoreUnitDay, g.ChoreUnitWeek, g.ChoreUnitMonth:
		if in.Every < 1 || in.Every > choreMaxEvery[in.Unit] {
			return f, FieldError{Field: FieldChoreSchedule}
		}
		every, unit := int16(in.Every), in.Unit
		f.every, f.unit = &every, &unit
		d, err := time.Parse(DateLayout, strings.TrimSpace(in.FirstDue))
		isCurrent := current.Valid && d.Equal(current.Time)
		if err != nil || (d.Before(today.Time) && !isCurrent) {
			return f, FieldError{Field: FieldChoreFirstDue}
		}
		f.due = d
	default:
		return f, FieldError{Field: FieldChoreSchedule}
	}
	if len(in.Members) == 0 {
		return f, FieldError{Field: FieldChoreMembers}
	}
	for i, m := range in.Members {
		if slices.Contains(in.Members[:i], m) || !slices.ContainsFunc(active, func(a dbqueries.ListActiveMembersRow) bool { return a.UserID == m }) {
			return f, FieldError{Field: FieldChoreMembers}
		}
	}
	if in.Assignee.Valid && !slices.Contains(in.Members, in.Assignee) {
		return f, FieldError{Field: FieldChoreAssignee}
	}
	return f, nil
}

// cleanChore checks in under the household lock a leaver takes, see [DefaultService.Leave]
func cleanChore(ctx context.Context, q *dbqueries.Queries, actor Membership, in ChoreInput, current pgtype.Date) (choreFields, error) {
	if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
		return choreFields{}, err
	}
	today, err := q.HouseholdToday(ctx, actor.HouseholdID)
	if err != nil {
		return choreFields{}, err
	}
	active, err := q.ListActiveMembers(ctx, actor.HouseholdID)
	if err != nil {
		return choreFields{}, err
	}
	return in.clean(today, current, active)
}

// setChoreRota replaces the rota of chore id with members in order
func setChoreRota(ctx context.Context, q *dbqueries.Queries, actor Membership, id pgtype.UUID, members []pgtype.UUID) error {
	if err := q.DeleteChoreMembers(ctx, dbqueries.DeleteChoreMembersParams{ChoreID: id, HouseholdID: actor.HouseholdID}); err != nil {
		return err
	}
	for i, m := range members {
		if err := q.InsertChoreMember(ctx, dbqueries.InsertChoreMemberParams{ChoreID: id, HouseholdID: actor.HouseholdID, UserID: m, Position: int32(i)}); err != nil {
			return err
		}
	}
	return nil
}

// firstAssignee is assignee, or when it is not valid whose turn the rota rule gives first
// (nobody for a manual chore)
func firstAssignee(ctx context.Context, q *dbqueries.Queries, id, assignee pgtype.UUID) (pgtype.UUID, error) {
	if assignee.Valid {
		return assignee, nil
	}
	return q.ChoreNextAssignee(ctx, dbqueries.ChoreNextAssigneeParams{ChoreID: id, After: pgtype.UUID{}, ExcludeAfter: false})
}

// lockChore locks the actor's chore id, [g.ErrorNotFound] for another household's
func lockChore(ctx context.Context, q *dbqueries.Queries, actor Membership, id pgtype.UUID) (dbqueries.Chore, error) {
	c, err := q.LockChore(ctx, dbqueries.LockChoreParams{ID: id, HouseholdID: actor.HouseholdID})
	return c, noRowsAs(err, g.ErrorNotFound)
}

// HouseholdToday is the household's local date, the server's clock may be in another zone
func (s *DefaultService) HouseholdToday(ctx context.Context, actor Membership) (time.Time, error) {
	d, err := s.db.HouseholdToday(ctx, actor.HouseholdID)
	return d.Time, err
}

// ListChores returns the household's chores, soonest due first
func (s *DefaultService) ListChores(ctx context.Context, actor Membership) ([]Chore, error) {
	return s.listChores(ctx, actor, pgtype.UUID{})
}

func (s *DefaultService) listChores(ctx context.Context, actor Membership, id pgtype.UUID) ([]Chore, error) {
	rows, err := s.db.ListChores(ctx, dbqueries.ListChoresParams{HouseholdID: actor.HouseholdID, ID: id})
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	members, err := s.db.ListChoreMembers(ctx, dbqueries.ListChoreMembersParams{HouseholdID: actor.HouseholdID, ChoreID: id})
	if err != nil {
		return nil, err
	}
	rotas := map[pgtype.UUID][]ChoreMember{}
	for _, m := range members {
		var last *time.Time
		if m.LastDone.Valid {
			last = &m.LastDone.Time
		}
		rotas[m.ChoreID] = append(rotas[m.ChoreID], ChoreMember{UserID: m.UserID, Name: m.Username, LastDone: last})
	}
	return utils.Map(rows, func(r dbqueries.ListChoresRow) Chore {
		c := Chore{
			ID:           r.ID,
			Title:        r.Title,
			Rotation:     r.Rotation,
			Unit:         g.ChoreUnitWeekdays,
			Every:        0,
			Weekdays:     nil,
			NextDue:      r.NextDue.Time,
			Today:        r.Today.Time,
			AssigneeID:   r.AssigneeID,
			AssigneeName: r.AssigneeName,
			Members:      rotas[r.ID],
		}
		if r.Weekdays != nil {
			c.Weekdays = weekdaysOf(*r.Weekdays)
		} else {
			c.Unit, c.Every = *r.EveryUnit, int(*r.Every)
		}
		return c
	}), nil
}

// GetChore returns [g.ErrorNotFound] if the chore is not in the actor's household
func (s *DefaultService) GetChore(ctx context.Context, actor Membership, id pgtype.UUID) (ChoreDetail, error) {
	if !id.Valid {
		return ChoreDetail{}, g.ErrorNotFound
	}
	chores, err := s.listChores(ctx, actor, id)
	if err != nil {
		return ChoreDetail{}, err
	}
	if len(chores) == 0 {
		return ChoreDetail{}, g.ErrorNotFound
	}
	log, err := s.db.ListChoreLog(ctx, dbqueries.ListChoreLogParams{ChoreID: id, HouseholdID: actor.HouseholdID, Lim: ChoreHistoryLength})
	return ChoreDetail{
		Chore: chores[0],
		History: utils.Map(log, func(l dbqueries.ListChoreLogRow) ChoreLogEntry {
			return ChoreLogEntry{Kind: l.Kind, UserName: l.UserName, AssigneeName: l.AssigneeName, ForOther: l.ForOther, Due: l.Due.Time, At: l.CreatedAt.Time}
		}),
	}, err
}

// CreateChore returns a [FieldError] for a bad field
func (s *DefaultService) CreateChore(ctx context.Context, actor Membership, in ChoreInput) (pgtype.UUID, error) {
	var id pgtype.UUID
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		f, err := cleanChore(ctx, q, actor, in, pgtype.Date{})
		if err != nil {
			return err
		}
		if id, err = q.InsertChore(ctx, dbqueries.InsertChoreParams{
			HouseholdID: actor.HouseholdID,
			Title:       f.title,
			Rotation:    in.Rotation,
			Every:       f.every,
			EveryUnit:   f.unit,
			Weekdays:    f.weekdays,
			NextDue:     dateOf(f.due),
			CreatedBy:   actor.UserID,
		}); err != nil {
			return err
		}
		if err := setChoreRota(ctx, q, actor, id, in.Members); err != nil {
			return err
		}
		assignee, err := firstAssignee(ctx, q, id, in.Assignee)
		if err != nil {
			return err
		}
		if err := q.SetChoreTurn(ctx, dbqueries.SetChoreTurnParams{NextDue: dateOf(f.due), AssigneeID: assignee, Handed: false, ID: id, HouseholdID: actor.HouseholdID}); err != nil {
			return err
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActChoreAdded, SubjectID: id, TargetUserID: pgtype.UUID{}, Title: f.title, AmountCents: nil, Currency: ""})
	})
	return id, err
}

// UpdateChore returns a [FieldError] for a bad field, [g.ErrorNotFound] for another household's chore
func (s *DefaultService) UpdateChore(ctx context.Context, actor Membership, id pgtype.UUID, in ChoreInput) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		// the household is locked before the chore like a leaver does, see [cleanChore]
		if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		c, err := lockChore(ctx, q, actor, id)
		if err != nil {
			return err
		}
		f, err := cleanChore(ctx, q, actor, in, c.NextDue)
		if err != nil {
			return err
		}
		// a weekdays chore keeps its turn's date unless its days changed
		if f.weekdays != nil && c.Weekdays != nil && *f.weekdays == *c.Weekdays {
			f.due = c.NextDue.Time
		}
		if err := q.UpdateChore(ctx, dbqueries.UpdateChoreParams{
			Title:       f.title,
			Rotation:    in.Rotation,
			Every:       f.every,
			EveryUnit:   f.unit,
			Weekdays:    f.weekdays,
			NextDue:     dateOf(f.due),
			ID:          id,
			HouseholdID: actor.HouseholdID,
		}); err != nil {
			return err
		}
		if err := setChoreRota(ctx, q, actor, id, in.Members); err != nil {
			return err
		}
		assignee := in.Assignee
		if !assignee.Valid && slices.Contains(in.Members, c.AssigneeID) {
			assignee = c.AssigneeID
		}
		if assignee, err = firstAssignee(ctx, q, id, assignee); err != nil || assignee == c.AssigneeID {
			return err
		}
		return q.SetChoreTurn(ctx, dbqueries.SetChoreTurnParams{NextDue: dateOf(f.due), AssigneeID: assignee, Handed: false, ID: id, HouseholdID: actor.HouseholdID})
	})
}

func (s *DefaultService) DeleteChore(ctx context.Context, actor Membership, id pgtype.UUID) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		c, err := lockChore(ctx, q, actor, id)
		if err != nil {
			return err
		}
		if err := rowsAffected(q.DeleteChore(ctx, dbqueries.DeleteChoreParams{ID: id, HouseholdID: actor.HouseholdID})); err != nil {
			return err
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActChoreDeleted, SubjectID: id, TargetUserID: pgtype.UUID{}, Title: c.Title, AmountCents: nil, Currency: ""})
	})
}

// DoneChore marks turn done by the actor
func (s *DefaultService) DoneChore(ctx context.Context, actor Membership, id pgtype.UUID, turn ChoreTurn) error {
	return s.moveChoreTurn(ctx, actor, id, turn, g.ChoreDone)
}

// SkipChore drops turn's due date
func (s *DefaultService) SkipChore(ctx context.Context, actor Membership, id pgtype.UUID, turn ChoreTurn) error {
	return s.moveChoreTurn(ctx, actor, id, turn, g.ChoreSkipped)
}

// PassChore hands turn on with the same due date, [g.ErrorInvalidInput] for a manual chore or a rota of one
func (s *DefaultService) PassChore(ctx context.Context, actor Membership, id pgtype.UUID, turn ChoreTurn) error {
	return s.moveChoreTurn(ctx, actor, id, turn, g.ChorePassed)
}

// moveChoreTurn logs kind (a g.Chore* value) for turn and hands the turn on, a pass with the
// same date, a done or skip with the next one; [ErrChoreTurnMoved] when turn is not current
func (s *DefaultService) moveChoreTurn(ctx context.Context, actor Membership, id pgtype.UUID, turn ChoreTurn, kind string) error {
	pass := kind == g.ChorePassed
	var c dbqueries.Chore
	var next pgtype.UUID
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		// the rota is read under the lock a leaver takes, see [DefaultService.Leave]
		if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		var err error
		if c, err = lockChore(ctx, q, actor, id); err != nil {
			return err
		}
		// a second click or a stale page; a pass keeps the date, so the assignee tells it too
		if !turn.Due.Valid || !c.NextDue.Time.Equal(turn.Due.Time) || c.AssigneeID != turn.Assignee {
			return ErrChoreTurnMoved
		}
		if pass && c.Rotation == g.RotationManual {
			return g.ErrorInvalidInput
		}
		nextDue := c.NextDue
		if !pass {
			today, err := q.HouseholdToday(ctx, actor.HouseholdID)
			if err != nil {
				return err
			}
			nextDue = dateOf(choreNextDue(c.Every, c.EveryUnit, c.Weekdays, c.NextDue.Time, today.Time))
		}
		// logged first, a least_recent rota goes by the done rows
		if err := q.InsertChoreLog(ctx, dbqueries.InsertChoreLogParams{ChoreID: id, HouseholdID: actor.HouseholdID, Kind: kind, UserID: actor.UserID, AssigneeID: c.AssigneeID, Due: c.NextDue}); err != nil {
			return err
		}
		if next, err = q.ChoreNextAssignee(ctx, dbqueries.ChoreNextAssigneeParams{ChoreID: id, After: c.AssigneeID, ExcludeAfter: pass}); err != nil {
			return err
		}
		if pass && !next.Valid {
			return g.ErrorInvalidInput
		}
		if err := q.SetChoreTurn(ctx, dbqueries.SetChoreTurnParams{NextDue: nextDue, AssigneeID: next, Handed: pass, ID: id, HouseholdID: actor.HouseholdID}); err != nil {
			return err
		}
		target := c.AssigneeID
		if pass {
			target = next
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: choreActivity[kind], SubjectID: id, TargetUserID: target, Title: c.Title, AmountCents: nil, Currency: ""})
	})
	if err == nil && pass {
		s.pushChoreTurn(ctx, actor, c, next, locales.LK_ChorePushPassedBody)
	}
	return err
}

// AssignChore gives the turn to userID, who must be in the rota
func (s *DefaultService) AssignChore(ctx context.Context, actor Membership, id, userID pgtype.UUID) error {
	var c dbqueries.Chore
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		var err error
		if c, err = lockChore(ctx, q, actor, id); err != nil {
			return err
		}
		rota, err := q.ListChoreMembers(ctx, dbqueries.ListChoreMembersParams{HouseholdID: actor.HouseholdID, ChoreID: id})
		if err != nil {
			return err
		}
		if !slices.ContainsFunc(rota, func(m dbqueries.ListChoreMembersRow) bool { return m.UserID == userID }) {
			return g.ErrorInvalidInput
		}
		// the same person keeps the pushes they already had
		if userID == c.AssigneeID {
			return nil
		}
		return q.SetChoreTurn(ctx, dbqueries.SetChoreTurnParams{NextDue: c.NextDue, AssigneeID: userID, Handed: true, ID: id, HouseholdID: actor.HouseholdID})
	})
	if err == nil && userID != c.AssigneeID {
		s.pushChoreTurn(ctx, actor, c, userID, locales.LK_ChorePushAssignedBody)
	}
	return err
}

// pushChoreTurn tells to that the actor gave them chore c's turn, body gets the actor and the chore
func (s *DefaultService) pushChoreTurn(ctx context.Context, actor Membership, c dbqueries.Chore, to pgtype.UUID, body locales.LK) {
	if !to.Valid || to == actor.UserID {
		return
	}
	s.notifyAsync(ctx, actor, g.NotifyChores, []pgtype.UUID{to}, func(ctx context.Context, _ pgtype.UUID, actorName, household string) push.Message {
		return push.Message{
			Title: utils.T(ctx, locales.LK_ChorePushDueTitle, household),
			Body:  utils.T(ctx, body, actorName, c.Title),
			URL:   routes.URL(routes.Chore, c.HouseholdID, c.ID),
			Tag:   "chore-" + c.ID.String(),
		}
	})
}

// MoveChoreMember moves userID one place up or down in the rota
func (s *DefaultService) MoveChoreMember(ctx context.Context, actor Membership, id, userID pgtype.UUID, up bool) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		if _, err := lockChore(ctx, q, actor, id); err != nil {
			return err
		}
		rota, err := q.ListChoreMembers(ctx, dbqueries.ListChoreMembersParams{HouseholdID: actor.HouseholdID, ChoreID: id})
		if err != nil {
			return err
		}
		i := slices.IndexFunc(rota, func(m dbqueries.ListChoreMembersRow) bool { return m.UserID == userID })
		if i < 0 {
			return g.ErrorNotFound
		}
		j := i + 1
		if up {
			j = i - 1
		}
		if j < 0 || j >= len(rota) {
			return nil
		}
		for _, p := range [][2]int{{i, j}, {j, i}} {
			if err := q.SetChoreMemberPosition(ctx, dbqueries.SetChoreMemberPositionParams{Position: rota[p[1]].Position, ChoreID: id, HouseholdID: actor.HouseholdID, UserID: rota[p[0]].UserID}); err != nil {
				return err
			}
		}
		return nil
	})
}

// NotifyChoresDue pushes the assignees of due chores at their household's reminder hour, and once more the day after
func (s *DefaultService) NotifyChoresDue(ctx context.Context) error {
	rows, err := s.db.ListChoresToNotify(ctx)
	if err != nil {
		return err
	}
	for _, r := range rows {
		title, body := locales.LK_ChorePushDueTitle, locales.LK_ChorePushDueBody
		// claimed first, so two runs at once cannot both push
		var n int64
		if r.Due {
			n, err = s.db.ClaimChoreNotified(ctx, dbqueries.ClaimChoreNotifiedParams{Overdue: r.Overdue, ID: r.ID, NextDue: r.NextDue, AssigneeID: r.AssigneeID})
		} else {
			title, body = locales.LK_ChorePushNudgeTitle, locales.LK_ChorePushNudgeBody
			n, err = s.db.ClaimChoreNudged(ctx, dbqueries.ClaimChoreNudgedParams{ID: r.ID, NextDue: r.NextDue, AssigneeID: r.AssigneeID})
		}
		if err != nil {
			return err
		}
		if n == 0 {
			continue
		}
		assignee, err := s.db.GetActiveMembership(ctx, dbqueries.GetActiveMembershipParams{HouseholdID: r.HouseholdID, UserID: r.AssigneeID})
		if err != nil {
			// left since the list was read, the trigger moved the turn on
			if err = noRowsAs(err, nil); err != nil {
				return err
			}
			continue
		}
		s.notifyAsync(ctx, assignee, g.NotifyChores, []pgtype.UUID{r.AssigneeID}, func(ctx context.Context, _ pgtype.UUID, _, household string) push.Message {
			return push.Message{
				Title: utils.T(ctx, title, household),
				Body:  utils.T(ctx, body, r.Title, utils.FormatDate(ctx, r.NextDue.Time)),
				URL:   routes.URL(routes.Chore, r.HouseholdID, r.ID),
				Tag:   "chore-" + r.ID.String(),
			}
		})
	}
	return nil
}
