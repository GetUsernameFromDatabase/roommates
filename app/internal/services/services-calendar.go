package services

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/push"
	"roommates/internal/routes"
	"roommates/internal/utils"
	"slices"
	"strings"
	"time"
	// the release image has no zoneinfo, household time zones need it
	_ "time/tzdata"

	"github.com/jackc/pgx/v5/pgtype"
)

// firstOccurrenceHorizon is how far past the start to look for an event's first occurrence,
// it covers the longest step, 12 months
const firstOccurrenceHorizon = 400

// EventInput is the event form
type EventInput struct {
	Title    string
	Location string
	Notes    string
	// [time.DateOnly], a weekdays series starts on the first ticked day from it
	StartDate string
	AllDay    bool
	// "15:04", not for all-day
	StartTime string
	// [time.DateOnly], empty for the start date
	EndDate string
	// "15:04", optional for a timed event, without it the event ends on its start date
	EndTime string
	// "" for no repeat, else a g.ChoreUnit*
	Unit string
	// with [g.ChoreEveryUnits]
	Every int
	// with [g.ChoreUnitWeekdays]
	Weekdays []time.Weekday
	// [time.DateOnly], empty to repeat forever
	RepeatUntil string
	// minutes before the start, one of [g.RemindTimed] or for all-day [g.RemindAllDay], nil for none
	Remind *int
	// empty for the whole household
	Members []pgtype.UUID
}

type EventMember struct {
	UserID pgtype.UUID
	Name   string
}

// Occurrence is one occurrence of an event, dates are household local dates as [time.DateOnly] parses them
type Occurrence struct {
	EventID pgtype.UUID
	// the day it starts and the day it ends
	On     time.Time
	EndOn  time.Time
	AllDay bool
	// "15:04" wall clock, empty for all-day; EndClock is empty without an end time
	StartClock string
	EndClock   string
	// in the household's time zone: all-day from midnight On to midnight after EndOn,
	// a timed one without an end time ends at Start
	Start    time.Time
	End      time.Time
	Title    string
	Location string
	Repeats  bool
	// empty for the whole household
	Members []EventMember
	// the actor created it or owns the household, see [DefaultService.UpdateEvent]
	CanEdit bool
}

type EventDetail struct {
	Occurrence
	Notes string
	// edited alone out of a series, it cannot repeat
	Detached bool
	// the edit form filled in for this occurrence
	Input EventInput
	Files []AttachedFile
}

// what a [CalendarItem] is
const (
	CalendarKindEvent     = "event"
	CalendarKindChore     = "chore"
	CalendarKindRecurring = "recurring"
)

// CalendarItem is one line of a day in the calendar, an event spanning days has one on each
type CalendarItem struct {
	Kind string
	// the day it shows on
	Date     time.Time
	Title    string
	Location string
	AllDay   bool
	// "15:04", Clock when an event starts on Date, EndClock when it ends on Date
	Clock    string
	EndClock string
	// the event, chore or recurring expense
	ID pgtype.UUID
	// the first day of an event's occurrence, ?on= of its page
	On           time.Time
	AssigneeName string
	AmountCents  int64
	// a chore or recurring expense due before today, shown on today
	Overdue bool
}

// what [DefaultService.UpdateEvent] and [DefaultService.DeleteEvent] change of a series,
// a one-off always changes alone
const (
	EventScopeThis      = "this"
	EventScopeFollowing = "following"
	EventScopeAll       = "all"
)

var eventScopes = []string{EventScopeThis, EventScopeFollowing, EventScopeAll}

const (
	MaxEventTitleLength    = 100
	MaxEventLocationLength = 100
	MaxEventNotesLength    = 500
	// MaxCalendarDays is the longest window [DefaultService.ListCalendar] lists
	MaxCalendarDays = 92
	// a reminder later than this, after downtime, is dropped
	remindLate    = 15 * time.Minute
	clockLayout   = "15:04"
	minutesPerDay = 24 * 60
)

// what a [FieldError] of [DefaultService.CreateEvent] and [DefaultService.UpdateEvent] names
const (
	FieldEventTitle    = "title"
	FieldEventLocation = "location"
	FieldEventNotes    = "notes"
	FieldEventStart    = "start"
	FieldEventEnd      = "end"
	// a repeating event that lasts past its next occurrence, shown at the end
	FieldEventLength      = "length"
	FieldEventSchedule    = "schedule"
	FieldEventRepeatUntil = "repeat_until"
	FieldEventRemind      = "remind"
	FieldEventMembers     = "members"
)

// ErrorTooManyEvents is a new event in a household that has [g.MaxUpcomingEvents] upcoming ones
var ErrorTooManyEvents = errors.New("too many upcoming events")

// eventFields is a checked [EventInput]
type eventFields struct {
	title, location, notes string
	start, end             time.Time
	startTime, endTime     pgtype.Time
	every                  *int16
	unit                   *string
	weekdays               *int16
	until                  pgtype.Date
	remind                 *int16
}

func (f eventFields) repeats() bool { return f.every != nil || f.weekdays != nil }

// clean checks in against the household's active members
func (in EventInput) clean(active []dbqueries.ListActiveMembersRow) (eventFields, error) {
	var f eventFields
	var ok bool
	if f.title, ok = trimToLength(in.Title, 1, MaxEventTitleLength); !ok {
		return f, FieldError{Field: FieldEventTitle}
	}
	if f.location, ok = trimToLength(in.Location, 0, MaxEventLocationLength); !ok {
		return f, FieldError{Field: FieldEventLocation}
	}
	if f.notes, ok = trimToLength(in.Notes, 0, MaxEventNotesLength); !ok {
		return f, FieldError{Field: FieldEventNotes}
	}
	var err error
	if f.start, err = parseDate(in.StartDate); err != nil {
		return f, FieldError{Field: FieldEventStart}
	}
	f.end = f.start
	if strings.TrimSpace(in.EndDate) != "" {
		if f.end, err = parseDate(in.EndDate); err != nil {
			return f, FieldError{Field: FieldEventEnd}
		}
	}
	if !in.AllDay {
		if f.startTime, err = parseClock(in.StartTime); err != nil {
			return f, FieldError{Field: FieldEventStart}
		}
		if strings.TrimSpace(in.EndTime) != "" {
			if f.endTime, err = parseClock(in.EndTime); err != nil {
				return f, FieldError{Field: FieldEventEnd}
			}
		}
	}
	if f.end.Before(f.start) ||
		(f.endTime.Valid && f.end.Equal(f.start) && f.endTime.Microseconds <= f.startTime.Microseconds) ||
		(!in.AllDay && !f.endTime.Valid && !f.end.Equal(f.start)) {
		return f, FieldError{Field: FieldEventEnd}
	}
	switch in.Unit {
	case "":
	case g.ChoreUnitWeekdays:
		var mask int16
		for _, d := range in.Weekdays {
			if d < time.Sunday || d > time.Saturday {
				return f, FieldError{Field: FieldEventSchedule}
			}
			mask |= weekdayBit(d)
		}
		if mask == 0 {
			return f, FieldError{Field: FieldEventSchedule}
		}
		f.weekdays = &mask
	case g.ChoreUnitDay, g.ChoreUnitWeek, g.ChoreUnitMonth:
		if in.Every < 1 || in.Every > choreMaxEvery[in.Unit] {
			return f, FieldError{Field: FieldEventSchedule}
		}
		every, unit := int16(in.Every), in.Unit
		f.every, f.unit = &every, &unit
	default:
		return f, FieldError{Field: FieldEventSchedule}
	}
	if f.repeats() && strings.TrimSpace(in.RepeatUntil) != "" {
		until, err := parseDate(in.RepeatUntil)
		if err != nil {
			return f, FieldError{Field: FieldEventRepeatUntil}
		}
		f.until = dateOf(until)
	}
	if in.Remind != nil {
		choices := g.RemindTimed
		if in.AllDay {
			choices = g.RemindAllDay
		}
		if !slices.Contains(choices, *in.Remind) {
			return f, FieldError{Field: FieldEventRemind}
		}
		r := int16(*in.Remind)
		f.remind = &r
	}
	for i, m := range in.Members {
		if slices.Contains(in.Members[:i], m) || !slices.ContainsFunc(active, func(a dbqueries.ListActiveMembersRow) bool { return a.UserID == m }) {
			return f, FieldError{Field: FieldEventMembers}
		}
	}
	return f, nil
}

// settle moves a weekdays series to its first ticked day, the end with it, and checks that
// it repeats until that day or later and that an occurrence ends before the next one starts,
// overlapping ones would all be listed on every day they cover; call it once the start is final
func (f *eventFields) settle() error {
	if f.weekdays != nil {
		first := nextWeekday(f.start, *f.weekdays)
		f.end = f.end.AddDate(0, 0, daysBetween(f.start, first))
		f.start = first
	}
	if f.until.Valid && f.until.Time.Before(f.start) {
		return FieldError{Field: FieldEventRepeatUntil}
	}
	if f.repeats() && f.minutes() > f.stepDays()*minutesPerDay {
		return FieldError{Field: FieldEventLength}
	}
	return nil
}

// minutes is how long an occurrence of f lasts, an all-day one to the end of its last day
func (f eventFields) minutes() int {
	m := daysBetween(f.start, f.end) * minutesPerDay
	if !f.startTime.Valid {
		return m + minutesPerDay
	}
	if f.endTime.Valid {
		m += int((f.endTime.Microseconds - f.startTime.Microseconds) / int64(time.Minute/time.Microsecond))
	}
	return m
}

// stepDays is the fewest days from one occurrence of a repeating f to the next,
// a month counts as its shortest, 28 days
func (f eventFields) stepDays() int {
	if f.weekdays != nil {
		// the shortest gap from a ticked weekday to the next, round the week
		gap := 7
		for d := range 7 {
			if *f.weekdays&(1<<d) == 0 {
				continue
			}
			for n := 1; n < gap; n++ {
				if *f.weekdays&(1<<((d+n)%7)) != 0 {
					gap = n
					break
				}
			}
		}
		return gap
	}
	n := int(*f.every)
	switch *f.unit {
	case g.ChoreUnitWeek:
		return 7 * n
	case g.ChoreUnitMonth:
		return 28 * n
	}
	return n
}

func parseDate(s string) (time.Time, error) {
	return time.Parse(time.DateOnly, strings.TrimSpace(s))
}

func parseClock(s string) (pgtype.Time, error) {
	t, err := time.Parse(clockLayout, strings.TrimSpace(s))
	return pgtype.Time{Microseconds: int64(t.Hour()*60+t.Minute()) * int64(time.Minute/time.Microsecond), Valid: err == nil}, err
}

// clockOf is t as "15:04", empty when it is not valid
func clockOf(t pgtype.Time) string {
	if !t.Valid {
		return ""
	}
	m := t.Microseconds / int64(time.Minute/time.Microsecond)
	return fmt.Sprintf("%02d:%02d", m/60, m%60)
}

// wallClock is the instant of day at t in loc, a time skipped by DST takes [time.Date]'s normalisation
func wallClock(day time.Time, t pgtype.Time, loc *time.Location) time.Time {
	m := int(t.Microseconds / int64(time.Minute/time.Microsecond))
	return time.Date(day.Year(), day.Month(), day.Day(), m/60, m%60, 0, 0, loc)
}

// DateOnly is t's calendar date as [time.DateOnly] parses it, keys of the skip sets compare with ==
func DateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// daysBetween counts the days from a to b, both from [DateOnly] or the database;
// in seconds, as [time.Time.Sub] stops at about 292 years
func daysBetween(a, b time.Time) int {
	return int((b.Unix() - a.Unix()) / int64(24*time.Hour/time.Second))
}

func eventRepeats(e dbqueries.Event) bool { return e.Every != nil || e.Weekdays != nil }

// eventSpan is how many days after its first day an occurrence of e ends
func eventSpan(e dbqueries.Event) int { return daysBetween(e.StartDate.Time, e.EndDate.Time) }

// occurrences are the first days of e's occurrences that overlap the days [from, to), without
// those in skip: anchored at its start date, every N months clamped to the month's end, day and
// week steps jump straight to the window so an old series does not walk its whole past
func occurrences(e dbqueries.Event, skip map[time.Time]bool, from, to time.Time) []time.Time {
	start := e.StartDate.Time
	lo := from.AddDate(0, 0, -eventSpan(e))
	if lo.Before(start) {
		lo = start
	}
	hi := to
	if e.RepeatUntil.Valid && e.RepeatUntil.Time.Before(hi) {
		hi = e.RepeatUntil.Time.AddDate(0, 0, 1)
	}
	var days []time.Time
	add := func(d time.Time) {
		if !d.Before(lo) && !skip[d] {
			days = append(days, d)
		}
	}
	switch {
	case e.Weekdays != nil:
		for d := lo; d.Before(hi); d = d.AddDate(0, 0, 1) {
			if *e.Weekdays&weekdayBit(d.Weekday()) != 0 {
				add(d)
			}
		}
	case e.Every == nil:
		if start.Before(to) {
			add(start)
		}
	case *e.EveryUnit == g.ChoreUnitMonth:
		n := int(*e.Every)
		months := (lo.Year()-start.Year())*12 + int(lo.Month()-start.Month())
		for k := months / n; ; k++ {
			d := addIntervals(start, g.IntervalMonth, k*n)
			if !d.Before(hi) {
				break
			}
			add(d)
		}
	default:
		step := int(*e.Every)
		if *e.EveryUnit == g.ChoreUnitWeek {
			step *= 7
		}
		k := (daysBetween(start, lo) + step - 1) / step
		for d := start.AddDate(0, 0, k*step); d.Before(hi); d = d.AddDate(0, 0, step) {
			add(d)
		}
	}
	return days
}

// remindAt is when the reminder of e's occurrence on on is due: minutes before a timed start,
// the day before at the same wall clock, or for all-day at the household's reminder hour
func remindAt(e dbqueries.Event, on time.Time, loc *time.Location, reminderHour int) time.Time {
	minutes := int(*e.RemindMinutes)
	if minutes == g.RemindDayBefore {
		on = on.AddDate(0, 0, -1)
	}
	if !e.StartTime.Valid {
		return time.Date(on.Year(), on.Month(), on.Day(), reminderHour, 0, 0, 0, loc)
	}
	if minutes == g.RemindDayBefore {
		minutes = 0
	}
	return wallClock(on, e.StartTime, loc).Add(-time.Duration(minutes) * time.Minute)
}

func occurrenceOf(e dbqueries.Event, on time.Time, loc *time.Location, actor Membership, members []EventMember) Occurrence {
	o := Occurrence{
		EventID:    e.ID,
		On:         on,
		EndOn:      on.AddDate(0, 0, eventSpan(e)),
		AllDay:     !e.StartTime.Valid,
		StartClock: clockOf(e.StartTime),
		EndClock:   clockOf(e.EndTime),
		Start:      time.Date(on.Year(), on.Month(), on.Day(), 0, 0, 0, 0, loc),
		End:        time.Time{},
		Title:      e.Title,
		Location:   e.Location,
		Repeats:    eventRepeats(e),
		Members:    members,
		CanEdit:    canChangeEvent(actor, e) == nil,
	}
	switch {
	case o.AllDay:
		o.End = time.Date(o.EndOn.Year(), o.EndOn.Month(), o.EndOn.Day()+1, 0, 0, 0, 0, loc)
	case e.EndTime.Valid:
		o.Start, o.End = wallClock(on, e.StartTime, loc), wallClock(o.EndOn, e.EndTime, loc)
	default:
		o.Start = wallClock(on, e.StartTime, loc)
		o.End = o.Start
	}
	return o
}

// inputOf is the form of e filled in for its occurrence on on
func inputOf(e dbqueries.Event, on time.Time, members []EventMember) EventInput {
	in := EventInput{
		Title:       e.Title,
		Location:    e.Location,
		Notes:       e.Notes,
		StartDate:   on.Format(time.DateOnly),
		AllDay:      !e.StartTime.Valid,
		StartTime:   clockOf(e.StartTime),
		EndDate:     on.AddDate(0, 0, eventSpan(e)).Format(time.DateOnly),
		EndTime:     clockOf(e.EndTime),
		Unit:        "",
		Every:       0,
		Weekdays:    nil,
		RepeatUntil: "",
		Remind:      nil,
		Members:     utils.Map(members, func(m EventMember) pgtype.UUID { return m.UserID }),
	}
	switch {
	case e.Weekdays != nil:
		in.Unit, in.Weekdays = g.ChoreUnitWeekdays, weekdaysOf(*e.Weekdays)
	case e.Every != nil:
		in.Unit, in.Every = *e.EveryUnit, int(*e.Every)
	}
	if e.RepeatUntil.Valid {
		in.RepeatUntil = e.RepeatUntil.Time.Format(time.DateOnly)
	}
	if e.RemindMinutes != nil {
		r := int(*e.RemindMinutes)
		in.Remind = &r
	}
	return in
}

// canChangeEvent is [g.ErrorForbidden] unless the actor created e or owns the household
func canChangeEvent(actor Membership, e dbqueries.Event) error {
	if e.CreatedBy != actor.UserID && actor.Role != g.RoleOwner {
		return g.ErrorForbidden
	}
	return nil
}

// lockEventToChange locks event id for the actor to edit or delete, [g.ErrorNotFound] for
// another household's, see [canChangeEvent]
func lockEventToChange(ctx context.Context, q *dbqueries.Queries, actor Membership, id pgtype.UUID) (dbqueries.Event, error) {
	e, err := q.LockEvent(ctx, dbqueries.LockEventParams{ID: id, HouseholdID: actor.HouseholdID})
	if err != nil {
		return e, noRowsAs(err, g.ErrorNotFound)
	}
	return e, canChangeEvent(actor, e)
}

// eventSkips are the days each series of ids skips, by event
func eventSkips(ctx context.Context, q *dbqueries.Queries, householdIDs, ids []pgtype.UUID) (map[pgtype.UUID]map[time.Time]bool, error) {
	rows, err := q.ListEventSkips(ctx, dbqueries.ListEventSkipsParams{HouseholdIds: householdIDs, EventIds: ids})
	skips := map[pgtype.UUID]map[time.Time]bool{}
	for _, r := range rows {
		if skips[r.EventID] == nil {
			skips[r.EventID] = map[time.Time]bool{}
		}
		skips[r.EventID][r.OccursOn.Time] = true
	}
	return skips, err
}

// checkOccurrence is [g.ErrorNotFound] unless e has an occurrence starting on on
func checkOccurrence(ctx context.Context, q *dbqueries.Queries, e dbqueries.Event, on time.Time) error {
	skips, err := eventSkips(ctx, q, []pgtype.UUID{e.HouseholdID}, []pgtype.UUID{e.ID})
	if err != nil {
		return err
	}
	if !slices.Contains(occurrences(e, skips[e.ID], on, on.AddDate(0, 0, 1)), on) {
		return g.ErrorNotFound
	}
	return nil
}

// insertEvent adds f with its members, detached from a series when detachedFrom is valid
func insertEvent(ctx context.Context, q *dbqueries.Queries, householdID, createdBy pgtype.UUID, f eventFields,
	detachedFrom pgtype.UUID, detachedOn pgtype.Date, members []pgtype.UUID) (pgtype.UUID, error) {
	id, err := q.InsertEvent(ctx, dbqueries.InsertEventParams{
		HouseholdID:   householdID,
		Title:         f.title,
		Location:      f.location,
		Notes:         f.notes,
		StartDate:     dateOf(f.start),
		StartTime:     f.startTime,
		EndDate:       dateOf(f.end),
		EndTime:       f.endTime,
		Every:         f.every,
		EveryUnit:     f.unit,
		Weekdays:      f.weekdays,
		RepeatUntil:   f.until,
		RemindMinutes: f.remind,
		DetachedFrom:  detachedFrom,
		DetachedOn:    detachedOn,
		CreatedBy:     createdBy,
	})
	if err != nil {
		return id, err
	}
	return id, setEventMembers(ctx, q, householdID, id, members)
}

// copyEvent adds f as a copy of e, made by e's creator and with e's files
func copyEvent(ctx context.Context, q *dbqueries.Queries, e dbqueries.Event, f eventFields, detachedFrom pgtype.UUID, detachedOn pgtype.Date, members []pgtype.UUID) (pgtype.UUID, error) {
	id, err := insertEvent(ctx, q, e.HouseholdID, e.CreatedBy, f, detachedFrom, detachedOn, members)
	if err != nil {
		return id, err
	}
	return id, q.CopyEventFileRefs(ctx, dbqueries.CopyEventFileRefsParams{ToID: id, EventID: e.ID, HouseholdID: e.HouseholdID})
}

func setEventMembers(ctx context.Context, q *dbqueries.Queries, householdID, id pgtype.UUID, members []pgtype.UUID) error {
	if err := q.DeleteEventMembers(ctx, dbqueries.DeleteEventMembersParams{EventID: id, HouseholdID: householdID}); err != nil {
		return err
	}
	for _, m := range members {
		if err := q.InsertEventMember(ctx, dbqueries.InsertEventMemberParams{EventID: id, HouseholdID: householdID, UserID: m}); err != nil {
			return err
		}
	}
	return nil
}

// dropEdits deletes the cancelled and the alone-edited occurrences of series e, from from on
// when it is valid: their dates are gone with a new rule
func dropEdits(ctx context.Context, q *dbqueries.Queries, e dbqueries.Event, from pgtype.Date) error {
	if err := q.DeleteEventExceptions(ctx, dbqueries.DeleteEventExceptionsParams{EventID: e.ID, HouseholdID: e.HouseholdID, FromDate: from}); err != nil {
		return err
	}
	return q.DeleteDetachedEvents(ctx, dbqueries.DeleteDetachedEventsParams{EventID: e.ID, HouseholdID: e.HouseholdID, FromDate: from})
}

func samePtr[T comparable](a, b *T) bool {
	return (a == nil) == (b == nil) && (a == nil || *a == *b)
}

// sameRule says whether f repeats like e, the end date aside
func sameRule(e dbqueries.Event, f eventFields) bool {
	return samePtr(e.Every, f.every) && samePtr(e.EveryUnit, f.unit) && samePtr(e.Weekdays, f.weekdays)
}

func (s *DefaultService) householdLocation(ctx context.Context, householdID pgtype.UUID) (dbqueries.GetHouseholdRow, *time.Location, error) {
	h, err := s.db.GetHousehold(ctx, householdID)
	if err != nil {
		return h, nil, noRowsAs(err, g.ErrorNotFound)
	}
	loc, err := time.LoadLocation(h.TimeZone)
	return h, loc, err
}

// ListCalendar lists the household's events, chores and recurring expenses on the household local
// days [from, to), at most [MaxCalendarDays], by day, all-day first and then by time: an event on
// each day it covers, a chore or recurring expense on its (next) due day or today when it is overdue
func (s *DefaultService) ListCalendar(ctx context.Context, actor Membership, from, to time.Time) ([]CalendarItem, error) {
	from, to = DateOnly(from), DateOnly(to)
	if !from.Before(to) || daysBetween(from, to) > MaxCalendarDays {
		return nil, g.ErrorInvalidInput
	}
	events, err := s.db.ListEventsInWindow(ctx, dbqueries.ListEventsInWindowParams{HouseholdID: actor.HouseholdID, ToDate: dateOf(to), FromDate: dateOf(from)})
	if err != nil {
		return nil, err
	}
	skips, err := eventSkips(ctx, s.db, []pgtype.UUID{actor.HouseholdID}, utils.Map(events, func(e dbqueries.Event) pgtype.UUID { return e.ID }))
	if err != nil {
		return nil, err
	}
	var items []CalendarItem
	item := func(kind string, d time.Time, title string, id pgtype.UUID) CalendarItem {
		return CalendarItem{Kind: kind, Date: d, Title: title, Location: "", AllDay: false, Clock: "", EndClock: "", ID: id, On: time.Time{}, AssigneeName: "", AmountCents: 0, Overdue: false}
	}
	for _, e := range events {
		for _, on := range occurrences(e, skips[e.ID], from, to) {
			end := on.AddDate(0, 0, eventSpan(e))
			d := on
			if d.Before(from) {
				d = from
			}
			for ; !d.After(end) && d.Before(to); d = d.AddDate(0, 0, 1) {
				it := item(CalendarKindEvent, d, e.Title, e.ID)
				it.Location, it.AllDay, it.On = e.Location, !e.StartTime.Valid, on
				if d.Equal(on) {
					it.Clock = clockOf(e.StartTime)
				}
				if d.Equal(end) {
					it.EndClock = clockOf(e.EndTime)
				}
				items = append(items, it)
			}
		}
	}
	inWindow := func(d time.Time) bool { return !d.Before(from) && d.Before(to) }
	chores, err := s.ListChores(ctx, actor)
	if err != nil {
		return nil, err
	}
	for _, c := range chores {
		d, overdue := c.NextDue, c.NextDue.Before(c.Today)
		if overdue {
			d = c.Today
		}
		if inWindow(d) {
			it := item(CalendarKindChore, d, c.Title, c.ID)
			it.AssigneeName, it.Overdue = c.AssigneeName, overdue
			items = append(items, it)
		}
	}
	recurring, err := s.ListRecurring(ctx, actor)
	if err != nil {
		return nil, err
	}
	var today time.Time
	if len(recurring) > 0 {
		if today, err = s.HouseholdToday(ctx, actor); err != nil {
			return nil, err
		}
	}
	for _, r := range recurring {
		// a draft not added yet waits on today, like an overdue chore
		d, overdue := r.NextDue.Time, r.NextDue.Time.Before(today)
		if overdue {
			d = today
		}
		if inWindow(d) {
			it := item(CalendarKindRecurring, d, r.Title, r.ID)
			it.AmountCents, it.Overdue = r.AmountCents, overdue
			items = append(items, it)
		}
	}
	// no clock sorts first, "15:04" sorts by time
	slices.SortStableFunc(items, func(a, b CalendarItem) int {
		return cmp.Or(a.Date.Compare(b.Date), cmp.Compare(a.Clock, b.Clock), cmp.Compare(a.Title, b.Title))
	})
	return items, nil
}

// GetOccurrence returns event id's occurrence starting on on, the first one left when on is zero,
// [g.ErrorNotFound] when there is none
func (s *DefaultService) GetOccurrence(ctx context.Context, actor Membership, id pgtype.UUID, on time.Time) (EventDetail, error) {
	if !id.Valid {
		return EventDetail{}, g.ErrorNotFound
	}
	e, err := s.db.GetEvent(ctx, dbqueries.GetEventParams{ID: id, HouseholdID: actor.HouseholdID})
	if err != nil {
		return EventDetail{}, noRowsAs(err, g.ErrorNotFound)
	}
	if on.IsZero() {
		// the first day may be cancelled or edited alone
		skips, err := eventSkips(ctx, s.db, []pgtype.UUID{e.HouseholdID}, []pgtype.UUID{e.ID})
		if err != nil {
			return EventDetail{}, err
		}
		days := occurrences(e, skips[e.ID], e.StartDate.Time, e.StartDate.Time.AddDate(0, 0, firstOccurrenceHorizon))
		if len(days) == 0 {
			return EventDetail{}, g.ErrorNotFound
		}
		on = days[0]
	} else {
		on = DateOnly(on)
		if err := checkOccurrence(ctx, s.db, e, on); err != nil {
			return EventDetail{}, err
		}
	}
	_, loc, err := s.householdLocation(ctx, actor.HouseholdID)
	if err != nil {
		return EventDetail{}, err
	}
	rows, err := s.db.ListEventMembers(ctx, dbqueries.ListEventMembersParams{HouseholdID: actor.HouseholdID, EventIds: []pgtype.UUID{id}})
	if err != nil {
		return EventDetail{}, err
	}
	members := utils.Map(rows, func(r dbqueries.ListEventMembersRow) EventMember {
		return EventMember{UserID: r.UserID, Name: r.Username}
	})
	files, err := s.AttachedFiles(ctx, actor, FileTarget{NoteID: pgtype.UUID{}, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}, EventID: id})
	return EventDetail{
		Occurrence: occurrenceOf(e, on, loc, actor, members),
		Notes:      e.Notes,
		Detached:   e.DetachedFrom.Valid,
		Input:      inputOf(e, on, members),
		Files:      files,
	}, err
}

// CreateEvent returns a [FieldError] for a bad field, [ErrorTooManyEvents] at the cap
func (s *DefaultService) CreateEvent(ctx context.Context, actor Membership, in EventInput) (pgtype.UUID, error) {
	var id pgtype.UUID
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		// members are read and the cap counted under the lock a leaver and other creates take
		if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		active, err := q.ListActiveMembers(ctx, actor.HouseholdID)
		if err != nil {
			return err
		}
		f, err := in.clean(active)
		if err != nil {
			return err
		}
		if err := f.settle(); err != nil {
			return err
		}
		today, err := q.HouseholdToday(ctx, actor.HouseholdID)
		if err != nil {
			return err
		}
		n, err := q.CountUpcomingEvents(ctx, dbqueries.CountUpcomingEventsParams{HouseholdID: actor.HouseholdID, Today: today})
		if err != nil {
			return err
		}
		if n >= g.MaxUpcomingEvents {
			return ErrorTooManyEvents
		}
		if id, err = insertEvent(ctx, q, actor.HouseholdID, actor.UserID, f, pgtype.UUID{}, pgtype.Date{}, in.Members); err != nil {
			return err
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActEventAdded, SubjectID: id, TargetUserID: pgtype.UUID{}, Title: f.title, AmountCents: nil, Currency: ""})
	})
	return id, err
}

// UpdateEvent saves in for event id's occurrence on on, the form shows that occurrence: scope this
// cancels it and adds in as a one-off detached from the series, following ends the series the day
// before and starts in as a new one, all changes the series and moves it by as many days as in
// moved the occurrence. A new start or rule drops the cancelled and alone-edited occurrences it
// covers, an unchanged one hands those of following to the new series. Copies keep the creator
// and the files. A one-off changes alone whatever scope says.
//
// returns [g.ErrorNotFound] when on is not an occurrence, [g.ErrorForbidden] unless the actor
// created the event or owns the household, a [FieldError] for a bad field
func (s *DefaultService) UpdateEvent(ctx context.Context, actor Membership, id pgtype.UUID, on time.Time, scope string, in EventInput) error {
	if !slices.Contains(eventScopes, scope) {
		return g.ErrorInvalidInput
	}
	on = DateOnly(on)
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		// the household is locked before the event like a leaver does, see [cleanChore]
		if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		e, err := lockEventToChange(ctx, q, actor, id)
		if err != nil {
			return err
		}
		if err := checkOccurrence(ctx, q, e, on); err != nil {
			return err
		}
		active, err := q.ListActiveMembers(ctx, actor.HouseholdID)
		if err != nil {
			return err
		}
		f, err := in.clean(active)
		if err != nil {
			return err
		}
		if e.DetachedFrom.Valid && f.repeats() {
			return FieldError{Field: FieldEventSchedule}
		}
		if !eventRepeats(e) || (scope == EventScopeFollowing && on.Equal(e.StartDate.Time)) {
			scope = EventScopeAll
		}
		switch {
		case scope == EventScopeThis:
			f.every, f.unit, f.weekdays, f.until = nil, nil, nil, pgtype.Date{}
		case scope == EventScopeAll && eventRepeats(e):
			span := daysBetween(f.start, f.end)
			f.start = e.StartDate.Time.AddDate(0, 0, daysBetween(on, f.start))
			f.end = f.start.AddDate(0, 0, span)
		}
		if err := f.settle(); err != nil {
			return err
		}
		switch scope {
		case EventScopeThis:
			if err := q.InsertEventException(ctx, dbqueries.InsertEventExceptionParams{EventID: e.ID, HouseholdID: e.HouseholdID, OccursOn: dateOf(on)}); err != nil {
				return err
			}
			_, err := copyEvent(ctx, q, e, f, e.ID, dateOf(on), in.Members)
			return err
		case EventScopeFollowing:
			if err := q.SetEventRepeatUntil(ctx, dbqueries.SetEventRepeatUntilParams{RepeatUntil: dateOf(on.AddDate(0, 0, -1)), ID: e.ID, HouseholdID: e.HouseholdID}); err != nil {
				return err
			}
			next, err := copyEvent(ctx, q, e, f, pgtype.UUID{}, pgtype.Date{}, in.Members)
			if err != nil {
				return err
			}
			if !f.start.Equal(on) || !sameRule(e, f) {
				return dropEdits(ctx, q, e, dateOf(on))
			}
			if err := q.MoveEventExceptions(ctx, dbqueries.MoveEventExceptionsParams{ToID: next, EventID: e.ID, HouseholdID: e.HouseholdID, FromDate: dateOf(on)}); err != nil {
				return err
			}
			return q.MoveDetachedEvents(ctx, dbqueries.MoveDetachedEventsParams{ToID: next, EventID: e.ID, HouseholdID: e.HouseholdID, FromDate: dateOf(on)})
		}
		// the old dates are gone with a new start or rule
		changed := eventRepeats(e) && (!f.start.Equal(e.StartDate.Time) || !sameRule(e, f))
		if err := q.UpdateEvent(ctx, dbqueries.UpdateEventParams{
			Title:         f.title,
			Location:      f.location,
			Notes:         f.notes,
			StartDate:     dateOf(f.start),
			StartTime:     f.startTime,
			EndDate:       dateOf(f.end),
			EndTime:       f.endTime,
			Every:         f.every,
			EveryUnit:     f.unit,
			Weekdays:      f.weekdays,
			RepeatUntil:   f.until,
			RemindMinutes: f.remind,
			ID:            e.ID,
			HouseholdID:   e.HouseholdID,
		}); err != nil {
			return err
		}
		if err := setEventMembers(ctx, q, e.HouseholdID, e.ID, in.Members); err != nil {
			return err
		}
		// a claim is per day, so a moved start or a new reminder would not remind again
		if f.startTime != e.StartTime || !f.start.Equal(e.StartDate.Time) || !samePtr(f.remind, e.RemindMinutes) {
			if err := q.DeleteEventReminders(ctx, dbqueries.DeleteEventRemindersParams{EventID: e.ID, HouseholdID: e.HouseholdID}); err != nil {
				return err
			}
		}
		if changed {
			return dropEdits(ctx, q, e, pgtype.Date{})
		}
		return nil
	})
}

// DeleteEvent deletes event id's occurrence on on: scope this cancels it, following ends the
// series the day before with the occurrences edited alone after it, all deletes the series;
// a one-off goes whatever scope says, errors as [DefaultService.UpdateEvent]
func (s *DefaultService) DeleteEvent(ctx context.Context, actor Membership, id pgtype.UUID, on time.Time, scope string) error {
	if !slices.Contains(eventScopes, scope) {
		return g.ErrorInvalidInput
	}
	on = DateOnly(on)
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		// the household is locked before the event like a leaver does, see [cleanChore]
		if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		e, err := lockEventToChange(ctx, q, actor, id)
		if err != nil {
			return err
		}
		if err := checkOccurrence(ctx, q, e, on); err != nil {
			return err
		}
		switch {
		case scope == EventScopeThis && eventRepeats(e):
			err = q.InsertEventException(ctx, dbqueries.InsertEventExceptionParams{EventID: e.ID, HouseholdID: e.HouseholdID, OccursOn: dateOf(on)})
		case scope == EventScopeFollowing && eventRepeats(e) && on.After(e.StartDate.Time):
			if err = q.SetEventRepeatUntil(ctx, dbqueries.SetEventRepeatUntilParams{RepeatUntil: dateOf(on.AddDate(0, 0, -1)), ID: e.ID, HouseholdID: e.HouseholdID}); err == nil {
				err = q.DeleteDetachedEvents(ctx, dbqueries.DeleteDetachedEventsParams{EventID: e.ID, HouseholdID: e.HouseholdID, FromDate: dateOf(on)})
			}
		default:
			// only a deleted event is logged, a series that loses some days is still there
			if err = rowsAffected(q.DeleteEvent(ctx, dbqueries.DeleteEventParams{ID: e.ID, HouseholdID: e.HouseholdID})); err == nil {
				err = logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActEventDeleted, SubjectID: e.ID, TargetUserID: pgtype.UUID{}, Title: e.Title, AmountCents: nil, Currency: ""})
			}
		}
		return err
	})
}

// NotifyCalendarDue pushes the reminders of event occurrences that are due, run it every minute
func (s *DefaultService) NotifyCalendarDue(ctx context.Context) error {
	return s.notifyCalendarDue(ctx, time.Now())
}

// shortcut: every run expands all events with a reminder around today, store each event's next
// reminder time when households grow
func (s *DefaultService) notifyCalendarDue(ctx context.Context, now time.Time) error {
	// any household's yesterday to tomorrow, no time zone is more than a day off UTC
	today := DateOnly(now.UTC())
	events, err := s.db.ListEventsToRemind(ctx, dbqueries.ListEventsToRemindParams{ToDate: dateOf(today.AddDate(0, 0, 2)), FromDate: dateOf(today.AddDate(0, 0, -2))})
	if err != nil || len(events) == 0 {
		return err
	}
	var households []pgtype.UUID
	for _, e := range events {
		if !slices.Contains(households, e.HouseholdID) {
			households = append(households, e.HouseholdID)
		}
	}
	skips, err := eventSkips(ctx, s.db, households, utils.Map(events, func(e dbqueries.Event) pgtype.UUID { return e.ID }))
	if err != nil {
		return err
	}
	type zone struct {
		loc  *time.Location
		hour int
	}
	zones := map[pgtype.UUID]zone{}
	for _, e := range events {
		z, ok := zones[e.HouseholdID]
		if !ok {
			h, loc, err := s.householdLocation(ctx, e.HouseholdID)
			if err != nil {
				// one household's zone must not stop the others' reminders
				log.Error().Err(err).Str("household", e.HouseholdID.String()).Msg("could not load the time zone for calendar reminders")
			}
			z = zone{loc: loc, hour: int(h.ReminderHour)}
			zones[e.HouseholdID] = z
		}
		if z.loc == nil {
			continue
		}
		local := DateOnly(now.In(z.loc))
		for _, on := range occurrences(e, skips[e.ID], local.AddDate(0, 0, -1), local.AddDate(0, 0, 2)) {
			if at := remindAt(e, on, z.loc, z.hour); at.After(now) || !at.After(now.Add(-remindLate)) {
				continue
			}
			// claimed first, so two runs at once cannot both push
			n, err := s.db.ClaimEventReminder(ctx, dbqueries.ClaimEventReminderParams{EventID: e.ID, HouseholdID: e.HouseholdID, OccursOn: dateOf(on)})
			if err != nil {
				return err
			}
			if n == 0 {
				continue
			}
			if err := s.pushEventReminder(ctx, e, on); err != nil {
				return err
			}
		}
	}
	return nil
}

// pushEventReminder pushes the reminder of e's occurrence on on to its members, or everyone
func (s *DefaultService) pushEventReminder(ctx context.Context, e dbqueries.Event, on time.Time) error {
	members, err := s.db.ListEventMembers(ctx, dbqueries.ListEventMembersParams{HouseholdID: e.HouseholdID, EventIds: []pgtype.UUID{e.ID}})
	if err != nil {
		return err
	}
	to := utils.Map(members, func(m dbqueries.ListEventMembersRow) pgtype.UUID { return m.UserID })
	if len(to) == 0 {
		active, err := s.db.ListActiveMembers(ctx, e.HouseholdID)
		if err != nil {
			return err
		}
		to = utils.Map(active, func(m dbqueries.ListActiveMembersRow) pgtype.UUID { return m.UserID })
	}
	// the household's reminder, the creator only stands in for notifyAsync's actor
	from := Membership{HouseholdID: e.HouseholdID, UserID: e.CreatedBy, Role: "", Currency: "", DeletedAt: pgtype.Timestamptz{}, UnreadMessages: 0, PendingSettlements: 0}
	s.notifyAsync(ctx, from, g.NotifyCalendar, to, func(ctx context.Context, _ pgtype.UUID, _, household string) push.Message {
		when := utils.FormatDate(ctx, on)
		if e.StartTime.Valid {
			when += " " + clockOf(e.StartTime)
		}
		return push.Message{
			Title: utils.T(ctx, locales.LK_CalendarPushTitle, household),
			Body:  utils.T(ctx, locales.LK_CalendarPushBody, e.Title, when),
			URL:   routes.URL(routes.CalendarEvent, e.HouseholdID, e.ID) + "?on=" + on.Format(time.DateOnly),
			Tag:   "event-" + e.ID.String() + "-" + on.Format(time.DateOnly),
		}
	})
	return nil
}
