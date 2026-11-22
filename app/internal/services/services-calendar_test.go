package services

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/routes"

	"github.com/jackc/pgx/v5/pgtype"
)

func calDay(s string) time.Time {
	d, err := time.Parse(DateLayout, s)
	if err != nil {
		panic(err)
	}
	return d
}

func testEvent(start, end string) dbqueries.Event {
	return dbqueries.Event{
		ID: uid(9), HouseholdID: uid(8), Title: "e", Location: "", Notes: "",
		StartDate: dateOf(calDay(start)), StartTime: pgtype.Time{}, EndDate: dateOf(calDay(end)), EndTime: pgtype.Time{},
		Every: nil, EveryUnit: nil, Weekdays: nil, RepeatUntil: pgtype.Date{}, RemindMinutes: nil,
		DetachedFrom: pgtype.UUID{}, DetachedOn: pgtype.Date{}, CreatedBy: uid(1), CreatedAt: pgtype.Timestamptz{},
	}
}

func repeating(e dbqueries.Event, every int16, unit string) dbqueries.Event {
	e.Every, e.EveryUnit = &every, &unit
	return e
}

func mustClock(s string) pgtype.Time {
	c, err := parseClock(s)
	if err != nil {
		panic(err)
	}
	return c
}

func TestEventOccurrences(t *testing.T) {
	monThu := weekdayBit(time.Monday) | weekdayBit(time.Thursday) // the 16th of November 2026 is a Monday
	weekdays := func(e dbqueries.Event, mask int16) dbqueries.Event { e.Weekdays = &mask; return e }
	until := func(e dbqueries.Event, d string) dbqueries.Event { e.RepeatUntil = dateOf(calDay(d)); return e }
	day, week, month := g.ChoreUnitDay, g.ChoreUnitWeek, g.ChoreUnitMonth
	for _, tc := range []struct {
		name     string
		e        dbqueries.Event
		skip     []string
		from, to string
		want     []string
	}{
		{"one-off", testEvent("2026-11-20", "2026-11-20"), nil, "2026-11-16", "2026-11-23", []string{"2026-11-20"}},
		{"one-off on the end", testEvent("2026-11-23", "2026-11-23"), nil, "2026-11-16", "2026-11-23", nil},
		{"one-off before", testEvent("2026-11-15", "2026-11-15"), nil, "2026-11-16", "2026-11-23", nil},
		{"multi-day over the start", testEvent("2026-11-14", "2026-11-16"), nil, "2026-11-16", "2026-11-23", []string{"2026-11-14"}},
		{"multi-day ended before", testEvent("2026-11-14", "2026-11-15"), nil, "2026-11-16", "2026-11-23", nil},
		{"weekdays", weekdays(testEvent("2026-11-16", "2026-11-16"), monThu), nil, "2026-11-16", "2026-11-30",
			[]string{"2026-11-16", "2026-11-19", "2026-11-23", "2026-11-26"}},
		{"weekdays until", until(weekdays(testEvent("2026-11-16", "2026-11-16"), monThu), "2026-11-23"), nil, "2026-11-16", "2026-11-30",
			[]string{"2026-11-16", "2026-11-19", "2026-11-23"}},
		{"weekdays started long ago", weekdays(testEvent("2026-01-01", "2026-01-01"), monThu), nil, "2026-11-17", "2026-11-24",
			[]string{"2026-11-19", "2026-11-23"}},
		{"monthly from the 31st", repeating(testEvent("2027-01-31", "2027-01-31"), 1, month), nil, "2027-01-01", "2027-06-01",
			[]string{"2027-01-31", "2027-02-28", "2027-03-31", "2027-04-30", "2027-05-31"}},
		{"every 2 months from the 31st", repeating(testEvent("2027-01-31", "2027-01-31"), 2, month), nil, "2027-03-01", "2027-08-01",
			[]string{"2027-03-31", "2027-05-31", "2027-07-31"}},
		{"every 3 days", repeating(testEvent("2026-11-01", "2026-11-01"), 3, day), nil, "2026-11-16", "2026-11-23",
			[]string{"2026-11-16", "2026-11-19", "2026-11-22"}},
		{"every 2 weeks", repeating(testEvent("2026-11-02", "2026-11-02"), 2, week), nil, "2026-11-16", "2026-12-01",
			[]string{"2026-11-16", "2026-11-30"}},
		{"old daily series", repeating(testEvent("2000-01-01", "2000-01-01"), 1, day), nil, "2026-11-16", "2026-11-19",
			[]string{"2026-11-16", "2026-11-17", "2026-11-18"}},
		{"cancelled and detached days", repeating(testEvent("2026-11-16", "2026-11-16"), 1, day), []string{"2026-11-17", "2026-11-18"}, "2026-11-16", "2026-11-20",
			[]string{"2026-11-16", "2026-11-19"}},
		{"until before the window", until(repeating(testEvent("2026-11-01", "2026-11-01"), 1, day), "2026-11-10"), nil, "2026-11-16", "2026-11-20", nil},
		{"weekends over the start", repeating(testEvent("2026-11-14", "2026-11-15"), 1, week), nil, "2026-11-15", "2026-11-22",
			[]string{"2026-11-14", "2026-11-21"}},
		{"weekends until", until(repeating(testEvent("2026-11-14", "2026-11-15"), 1, week), "2026-11-14"), nil, "2026-11-15", "2026-11-22",
			[]string{"2026-11-14"}},
	} {
		skip := map[time.Time]bool{}
		for _, d := range tc.skip {
			skip[calDay(d)] = true
		}
		var got []string
		for _, d := range occurrences(tc.e, skip, calDay(tc.from), calDay(tc.to)) {
			got = append(got, d.Format(DateLayout))
		}
		if !slices.Equal(got, tc.want) {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestEventTimes(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Fatal(err)
	}
	// summer time starts on 2027-03-28 at 03:00, the wall clock stays, the instant moves
	e := repeating(testEvent("2027-03-21", "2027-03-21"), 1, g.ChoreUnitWeek)
	e.StartTime, e.EndTime = mustClock("10:00"), mustClock("11:30")
	winter, summer := occurrenceOf(e, calDay("2027-03-21"), loc, Membership{}, nil), occurrenceOf(e, calDay("2027-03-28"), loc, Membership{}, nil)
	if winter.Start.UTC().Hour() != 8 || summer.Start.UTC().Hour() != 7 || summer.End.Sub(summer.Start) != 90*time.Minute || summer.StartClock != "10:00" || summer.EndClock != "11:30" {
		t.Fatalf("DST: got %v - %v and %v - %v", winter.Start, winter.End, summer.Start, summer.End)
	}
	allDay := testEvent("2027-03-28", "2027-03-29")
	if o := occurrenceOf(allDay, calDay("2027-03-28"), loc, Membership{}, nil); !o.AllDay || o.End.Sub(o.Start) != 47*time.Hour || o.EndOn != calDay("2027-03-29") {
		t.Fatalf("all-day: got %+v", o)
	}
	remind := func(e dbqueries.Event, minutes int16, on string) string {
		e.RemindMinutes = &minutes
		return remindAt(e, calDay(on), loc, 9).UTC().Format(time.RFC3339)
	}
	for _, tc := range []struct{ name, got, want string }{
		{"an hour before", remind(e, g.RemindHour, "2027-03-28"), "2027-03-28T06:00:00Z"},
		{"15 minutes before", remind(e, g.Remind15Min, "2027-03-21"), "2027-03-21T07:45:00Z"},
		{"the day before at the same wall clock", remind(e, g.RemindDayBefore, "2027-03-28"), "2027-03-27T08:00:00Z"},
		{"at the start", remind(e, g.RemindAtStart, "2027-03-28"), "2027-03-28T07:00:00Z"},
		{"all-day on the day", remind(allDay, g.RemindAtStart, "2027-03-28"), "2027-03-28T06:00:00Z"},
		{"all-day the day before", remind(allDay, g.RemindDayBefore, "2027-03-28"), "2027-03-27T07:00:00Z"},
	} {
		if tc.got != tc.want {
			t.Errorf("%s: got %s, want %s", tc.name, tc.got, tc.want)
		}
	}
	// 03:30 does not exist that night, it is what time.Date makes of it
	e.StartTime, e.EndTime = mustClock("03:30"), pgtype.Time{}
	if o := occurrenceOf(e, calDay("2027-03-28"), loc, Membership{}, nil); !o.Start.Equal(time.Date(2027, 3, 28, 3, 30, 0, 0, loc)) || !o.End.Equal(o.Start) {
		t.Fatalf("skipped wall clock: got %v", o.Start)
	}
	if got := clockOf(mustClock(" 9:05 ")); got != "09:05" {
		t.Fatalf("clock: got %q", got)
	}
}

func TestEventClean(t *testing.T) {
	a, b, stranger := uid(1), uid(2), uid(3)
	active := []dbqueries.ListActiveMembersRow{{UserID: a, Username: "a", AvatarKey: nil, Role: g.RoleOwner, JoinedAt: pgtype.Timestamptz{}}, {UserID: b, Username: "b", AvatarKey: nil, Role: g.RoleMember, JoinedAt: pgtype.Timestamptz{}}}
	hour := g.RemindHour
	ok := func() EventInput {
		return EventInput{Title: " Dinner ", Location: " Home ", Notes: "", StartDate: "2026-11-18", AllDay: false, StartTime: "19:00", EndDate: "", EndTime: "21:00",
			Unit: "", Every: 0, Weekdays: nil, RepeatUntil: "2026-12-01", Remind: &hour, Members: []pgtype.UUID{a}}
	}
	f, err := ok().clean(active)
	if err != nil || f.title != "Dinner" || f.location != "Home" || !f.start.Equal(calDay("2026-11-18")) || !f.end.Equal(f.start) ||
		clockOf(f.startTime) != "19:00" || clockOf(f.endTime) != "21:00" || f.repeats() || f.until.Valid || *f.remind != 60 {
		t.Fatalf("ok: got %+v %v", f, err)
	}
	// a weekdays series starts on its first ticked day, the 18th is a Wednesday
	in := ok()
	in.Unit, in.Weekdays, in.EndDate = g.ChoreUnitWeekdays, []time.Weekday{time.Monday, time.Friday}, "2026-11-19"
	if f, err = in.clean(active); err == nil {
		err = f.settle()
	}
	if err != nil || !f.start.Equal(calDay("2026-11-20")) || !f.end.Equal(calDay("2026-11-21")) || *f.weekdays != weekdayBit(time.Monday)|weekdayBit(time.Friday) || !f.until.Valid {
		t.Fatalf("weekdays: got %+v %v", f, err)
	}
	in = ok()
	in.AllDay, in.StartTime, in.EndTime, in.EndDate = true, "x", "y", "2026-11-20"
	dayBefore := g.RemindDayBefore
	in.Remind = &dayBefore
	if f, err = in.clean(active); err != nil || f.startTime.Valid || f.endTime.Valid || !f.end.Equal(calDay("2026-11-20")) {
		t.Fatalf("all-day: got %+v %v", f, err)
	}
	in = ok()
	in.Unit, in.Every, in.RepeatUntil = g.ChoreUnitDay, 1, "2026-11-17"
	if f, err = in.clean(active); err == nil {
		err = f.settle()
	}
	if !errors.Is(err, FieldError{Field: FieldEventRepeatUntil}) {
		t.Fatalf("until before the start: got %v", err)
	}
	seven, fifteen := 7, g.Remind15Min
	for _, tc := range []struct {
		name  string
		field string
		edit  func(*EventInput)
	}{
		{"blank title", FieldEventTitle, func(in *EventInput) { in.Title = "  " }},
		{"long title", FieldEventTitle, func(in *EventInput) { in.Title = strings.Repeat("ä", MaxEventTitleLength+1) }},
		{"long location", FieldEventLocation, func(in *EventInput) { in.Location = strings.Repeat("x", MaxEventLocationLength+1) }},
		{"long notes", FieldEventNotes, func(in *EventInput) { in.Notes = strings.Repeat("x", MaxEventNotesLength+1) }},
		{"bad date", FieldEventStart, func(in *EventInput) { in.StartDate = "18.11.2026" }},
		{"bad time", FieldEventStart, func(in *EventInput) { in.StartTime = "25:00" }},
		{"end before start", FieldEventEnd, func(in *EventInput) { in.EndTime = "18:00" }},
		{"end at start", FieldEventEnd, func(in *EventInput) { in.EndTime = "19:00" }},
		{"end day before", FieldEventEnd, func(in *EventInput) { in.EndDate = "2026-11-17" }},
		{"days without an end time", FieldEventEnd, func(in *EventInput) { in.EndTime, in.EndDate = "", "2026-11-19" }},
		{"unit", FieldEventSchedule, func(in *EventInput) { in.Unit = "year" }},
		{"every 0", FieldEventSchedule, func(in *EventInput) { in.Unit, in.Every = g.ChoreUnitDay, 0 }},
		{"every 13 months", FieldEventSchedule, func(in *EventInput) { in.Unit, in.Every = g.ChoreUnitMonth, 13 }},
		{"no weekdays", FieldEventSchedule, func(in *EventInput) { in.Unit = g.ChoreUnitWeekdays }},
		{"bad until", FieldEventRepeatUntil, func(in *EventInput) { in.Unit, in.Every, in.RepeatUntil = g.ChoreUnitWeek, 1, "soon" }},
		{"remind 7", FieldEventRemind, func(in *EventInput) { in.Remind = &seven }},
		{"all-day 15 minutes", FieldEventRemind, func(in *EventInput) { in.AllDay, in.Remind = true, &fifteen }},
		{"member twice", FieldEventMembers, func(in *EventInput) { in.Members = []pgtype.UUID{a, a} }},
		{"stranger", FieldEventMembers, func(in *EventInput) { in.Members = []pgtype.UUID{b, stranger} }},
	} {
		in := ok()
		tc.edit(&in)
		if _, err := in.clean(active); !errors.Is(err, FieldError{Field: tc.field}) {
			t.Errorf("%s: got %v, want %s", tc.name, err, tc.field)
		}
	}
}

func eventInput(title string, start time.Time) EventInput {
	return EventInput{Title: title, Location: "", Notes: "", StartDate: start.Format(DateLayout), AllDay: true, StartTime: "", EndDate: "", EndTime: "",
		Unit: "", Every: 0, Weekdays: nil, RepeatUntil: "", Remind: nil, Members: nil}
}

func dailyEvent(title string, start time.Time) EventInput {
	in := eventInput(title, start)
	in.Unit, in.Every = g.ChoreUnitDay, 1
	return in
}

func eventTarget(id pgtype.UUID) FileTarget {
	return FileTarget{NoteID: pgtype.UUID{}, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}, EventID: id}
}

func mustCreateEvent(t *testing.T, s *DefaultService, actor Membership, in EventInput) pgtype.UUID {
	t.Helper()
	id, err := s.CreateEvent(t.Context(), actor, in)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

// calendarEvents is the events of the n days from from, a day's titles joined by +
func calendarEvents(t *testing.T, s *DefaultService, actor Membership, from time.Time, n int) ([]string, []CalendarItem) {
	t.Helper()
	items, err := s.ListCalendar(t.Context(), actor, from, from.AddDate(0, 0, n))
	if err != nil {
		t.Fatal(err)
	}
	days := make([]string, n)
	var events []CalendarItem
	for _, it := range items {
		if it.Kind == CalendarKindEvent {
			i := daysBetween(from, it.Date)
			days[i] = strings.Trim(days[i]+"+"+it.Title, "+")
			events = append(events, it)
		}
	}
	return days, events
}

func wantCalendar(t *testing.T, s *DefaultService, actor Membership, from time.Time, step string, want ...string) []CalendarItem {
	t.Helper()
	got, events := calendarEvents(t, s, actor, from, len(want))
	if !slices.Equal(got, want) {
		t.Fatalf("%s: got %q, want %q", step, got, want)
	}
	return events
}

// eventOn is the ID of the event shown on day d
func eventOn(items []CalendarItem, d time.Time, title string) pgtype.UUID {
	for _, it := range items {
		if it.Date.Equal(d) && it.Title == title {
			return it.ID
		}
	}
	return pgtype.UUID{}
}

func mustOccurrence(t *testing.T, s *DefaultService, actor Membership, id pgtype.UUID, on time.Time) EventDetail {
	t.Helper()
	d, err := s.GetOccurrence(t.Context(), actor, id, on)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestEventScopes(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, _ := newTestHousehold(t, s)
	today, err := s.HouseholdToday(ctx, ownerM)
	if err != nil {
		t.Fatal(err)
	}
	day := func(i int) time.Time { return today.AddDate(0, 0, i) }
	in := dailyEvent("Yoga", day(0))
	in.Members = []pgtype.UUID{memberM.UserID}
	id := mustCreateEvent(t, s, ownerM, in)
	if _, err := s.UploadFile(ctx, ownerM, eventTarget(id), "plan.txt", strings.NewReader("plan"), 4); err != nil {
		t.Fatal(err)
	}
	wantCalendar(t, s, ownerM, day(0), "created", "Yoga", "Yoga", "Yoga", "Yoga", "Yoga")

	// this: a one-off on the day with the series' creator, members and files
	edit := in
	edit.Title, edit.StartDate = "Yoga late", day(1).Format(DateLayout)
	if err := s.UpdateEvent(ctx, ownerM, id, day(1), EventScopeThis, edit); err != nil {
		t.Fatal(err)
	}
	items := wantCalendar(t, s, ownerM, day(0), "edited this", "Yoga", "Yoga late", "Yoga", "Yoga", "Yoga")
	late := eventOn(items, day(1), "Yoga late")
	d := mustOccurrence(t, s, memberM, late, day(1))
	if !d.Detached || d.Repeats || d.CanEdit || len(d.Files) != 1 || len(d.Members) != 1 || d.Members[0].UserID != memberM.UserID || d.Input.Unit != "" {
		t.Fatalf("detached: got %+v", d)
	}
	wantErr(t, "detached repeats", s.UpdateEvent(ctx, ownerM, late, day(1), EventScopeAll, dailyEvent("Yoga late", day(1))), FieldError{Field: FieldEventSchedule})
	_, err = s.GetOccurrence(ctx, ownerM, id, day(1))
	wantErr(t, "edited day of the series", err, g.ErrorNotFound)

	// this on delete cancels the day
	for _, i := range []int{2, 4} {
		if err := s.DeleteEvent(ctx, ownerM, id, day(i), EventScopeThis); err != nil {
			t.Fatal(err)
		}
	}
	wantCalendar(t, s, ownerM, day(0), "deleted this", "Yoga", "Yoga late", "", "Yoga", "")
	wantErr(t, "cancelled day", s.DeleteEvent(ctx, ownerM, id, day(2), EventScopeThis), g.ErrorNotFound)

	// following with the same rule and day: a new series, the cancelled day stays cancelled
	edit = in
	edit.Title, edit.StartDate = "Pilates", day(3).Format(DateLayout)
	if err := s.UpdateEvent(ctx, ownerM, id, day(3), EventScopeFollowing, edit); err != nil {
		t.Fatal(err)
	}
	items = wantCalendar(t, s, ownerM, day(0), "edited following", "Yoga", "Yoga late", "", "Pilates", "", "Pilates")
	pilates := eventOn(items, day(3), "Pilates")
	if d := mustOccurrence(t, s, ownerM, pilates, day(3)); pilates == id || len(d.Files) != 1 || !d.Repeats || d.Input.RepeatUntil != "" {
		t.Fatalf("following: got %+v", d)
	}
	if d := mustOccurrence(t, s, ownerM, id, day(0)); d.Input.RepeatUntil != day(2).Format(DateLayout) {
		t.Fatalf("ended series: got %+v", d.Input)
	}

	// all with a new rule drops the series' cancelled and alone-edited days
	edit = in
	edit.Every, edit.RepeatUntil = 2, day(2).Format(DateLayout)
	if err := s.UpdateEvent(ctx, ownerM, id, day(0), EventScopeAll, edit); err != nil {
		t.Fatal(err)
	}
	wantCalendar(t, s, ownerM, day(0), "new rule", "Yoga", "", "Yoga", "Pilates", "")
	_, err = s.GetOccurrence(ctx, ownerM, late, day(1))
	wantErr(t, "dropped one-off", err, g.ErrorNotFound)

	// all from a later occurrence moves the whole series by as many days
	edit.StartDate, edit.RepeatUntil = day(3).Format(DateLayout), day(3).Format(DateLayout)
	if err := s.UpdateEvent(ctx, ownerM, id, day(2), EventScopeAll, edit); err != nil {
		t.Fatal(err)
	}
	wantCalendar(t, s, ownerM, day(0), "moved", "", "Yoga", "", "Pilates+Yoga", "")

	// following on the first day is all
	if err := s.DeleteEvent(ctx, ownerM, pilates, day(3), EventScopeFollowing); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteEvent(ctx, ownerM, id, day(3), EventScopeFollowing); err != nil {
		t.Fatal(err)
	}
	wantCalendar(t, s, ownerM, day(0), "deleted following", "", "Yoga", "", "", "")
	if err := s.DeleteEvent(ctx, ownerM, id, day(1), EventScopeAll); err != nil {
		t.Fatal(err)
	}
	wantCalendar(t, s, ownerM, day(0), "deleted all", "", "", "", "", "")
	wantKinds(t, s, ownerM, g.ActEventAdded, g.ActEventDeleted)
}

func TestEventRights(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, thirdM, strangerM := choreTestHousehold(t, s)
	today, err := s.HouseholdToday(ctx, ownerM)
	if err != nil {
		t.Fatal(err)
	}
	in := dailyEvent("Party", today)
	id := mustCreateEvent(t, s, memberM, in)

	wantErr(t, "third edits", s.UpdateEvent(ctx, thirdM, id, today, EventScopeAll, in), g.ErrorForbidden)
	wantErr(t, "third deletes", s.DeleteEvent(ctx, thirdM, id, today, EventScopeThis), g.ErrorForbidden)
	_, err = s.UploadFile(ctx, thirdM, eventTarget(id), "x.txt", strings.NewReader("x"), 1)
	wantErr(t, "third uploads", err, g.ErrorForbidden)
	_, err = s.GetOccurrence(ctx, strangerM, id, today)
	wantErr(t, "stranger reads", err, g.ErrorNotFound)
	wantErr(t, "stranger edits", s.UpdateEvent(ctx, strangerM, id, today, EventScopeAll, in), g.ErrorNotFound)
	wantErr(t, "bad scope", s.DeleteEvent(ctx, memberM, id, today, "some"), g.ErrorInvalidInput)
	if d := mustOccurrence(t, s, thirdM, id, time.Time{}); d.CanEdit || !d.On.Equal(today) {
		t.Fatalf("third reads the first occurrence: got %+v", d.Occurrence)
	}

	// an owner edits anyone's, the copy is still the creator's
	edit := in
	edit.Title, edit.StartDate = "Party later", today.AddDate(0, 0, 1).Format(DateLayout)
	if err := s.UpdateEvent(ctx, ownerM, id, today.AddDate(0, 0, 1), EventScopeThis, edit); err != nil {
		t.Fatal(err)
	}
	_, items := calendarEvents(t, s, ownerM, today, 2)
	copyID := eventOn(items, today.AddDate(0, 0, 1), "Party later")
	if !mustOccurrence(t, s, memberM, copyID, time.Time{}).CanEdit || mustOccurrence(t, s, thirdM, copyID, time.Time{}).CanEdit || !mustOccurrence(t, s, ownerM, id, today).CanEdit {
		t.Fatal("rights of the copy")
	}
	if err := s.DeleteEvent(ctx, memberM, id, today, EventScopeAll); err != nil {
		t.Fatal(err)
	}
}

// without a day, an event opens on its first day left, links of the activity and files have none
func TestEventFirstDayLeft(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, _, _ := newTestHousehold(t, s)
	today, err := s.HouseholdToday(ctx, ownerM)
	if err != nil {
		t.Fatal(err)
	}
	id := mustCreateEvent(t, s, ownerM, dailyEvent("Walk", today))
	if err := s.DeleteEvent(ctx, ownerM, id, today, EventScopeThis); err != nil {
		t.Fatal(err)
	}
	if d := mustOccurrence(t, s, ownerM, id, time.Time{}); !d.On.Equal(today.AddDate(0, 0, 1)) {
		t.Fatalf("got %v, want the day after %v", d.On, today)
	}
}

func TestEventCap(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, _, _ := newTestHousehold(t, s)
	today, err := s.HouseholdToday(ctx, ownerM)
	if err != nil {
		t.Fatal(err)
	}
	// one below the cap, and an ended series that does not count
	if _, err := s.pool.Exec(ctx, `INSERT INTO events (household_id, title, start_date, end_date, created_by)
		SELECT $1, 'e', $2::date, $2::date, $3 FROM generate_series(1, $4)`, ownerM.HouseholdID, dateOf(today), ownerM.UserID, g.MaxUpcomingEvents-1); err != nil {
		t.Fatal(err)
	}
	if _, err := s.pool.Exec(ctx, `INSERT INTO events (household_id, title, start_date, end_date, every, every_unit, repeat_until, created_by)
		VALUES ($1, 'ended', $2::date - 30, $2::date - 30, 1, 'day', $2::date - 1, $3)`, ownerM.HouseholdID, dateOf(today), ownerM.UserID); err != nil {
		t.Fatal(err)
	}
	id := mustCreateEvent(t, s, ownerM, dailyEvent("Last", today))
	_, err = s.CreateEvent(ctx, ownerM, eventInput("Over", today.AddDate(0, 0, -1)))
	wantErr(t, "over the cap", err, ErrorTooManyEvents)
	// edits never fail on the cap
	if err := s.UpdateEvent(ctx, ownerM, id, today.AddDate(0, 0, 1), EventScopeThis, eventInput("Moved", today.AddDate(0, 0, 1))); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateEvent(ctx, ownerM, id, today.AddDate(0, 0, 2), EventScopeFollowing, dailyEvent("Split", today.AddDate(0, 0, 2))); err != nil {
		t.Fatal(err)
	}
}

func TestNotifyCalendarDue(t *testing.T) {
	s := newTestService(t)
	locales.LoadLocales()
	ctx := t.Context()
	ownerM, memberM, _ := newTestHousehold(t, s)
	device := subscribeTestDevice(t, s, ownerM.UserID, "en")
	if err := s.UpdateHousehold(ctx, ownerM, "Flat", ownerM.Currency, "Europe/Tallinn", 7); err != nil {
		t.Fatal(err)
	}
	loc, err := time.LoadLocation("Europe/Tallinn")
	if err != nil {
		t.Fatal(err)
	}
	today, err := s.HouseholdToday(ctx, ownerM)
	if err != nil {
		t.Fatal(err)
	}
	tomorrow := today.AddDate(0, 0, 1)
	at := func(h, m int) time.Time {
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), h, m, 0, 0, loc)
	}
	// other households' reminders can be due too, count only this device
	notify := func(step string, now time.Time, want int) []sentPush {
		t.Helper()
		if err := s.notifyCalendarDue(ctx, now); err != nil {
			t.Fatal(err)
		}
		var mine []sentPush
		for _, p := range takePushes(s) {
			if p.endpoint == device {
				mine = append(mine, p)
			}
		}
		if len(mine) != want {
			t.Fatalf("%s: got %d pushes, want %d", step, len(mine), want)
		}
		return mine
	}
	timed := func(title, start string, remind int, members ...pgtype.UUID) pgtype.UUID {
		in := eventInput(title, tomorrow)
		in.AllDay, in.StartTime, in.Remind, in.Members = false, start, &remind, members
		return mustCreateEvent(t, s, memberM, in)
	}

	dentist := timed("Dentist", "10:00", g.RemindHour)
	notify("too early", at(8, 59), 0)
	got := notify("due", at(9, 1), 1)
	on := tomorrow.Format(DateLayout)
	if got[0].msg.URL != routes.URL(routes.CalendarEvent, ownerM.HouseholdID, dentist)+"?on="+on || got[0].msg.Tag != "event-"+dentist.String()+"-"+on {
		t.Fatalf("push: got %+v", got[0].msg)
	}
	notify("claimed", at(9, 2), 0)

	// a moved start reminds again
	moved := eventInput("Dentist", tomorrow)
	hour := g.RemindHour
	moved.AllDay, moved.StartTime, moved.Remind = false, "11:00", &hour
	if err := s.UpdateEvent(ctx, memberM, dentist, tomorrow, EventScopeAll, moved); err != nil {
		t.Fatal(err)
	}
	notify("moved", at(10, 1), 1)

	timed("Gym", "13:00", g.RemindAtStart)
	notify("over 15 minutes late", at(13, 16), 0)

	timed("Members only", "15:00", g.RemindAtStart, memberM.UserID)
	notify("not a member", at(15, 0), 0)

	// all-day: the day before at the reminder hour
	trip := eventInput("Trip", tomorrow.AddDate(0, 0, 1))
	dayBefore := g.RemindDayBefore
	trip.Remind = &dayBefore
	mustCreateEvent(t, s, memberM, trip)
	notify("all-day early", at(6, 59), 0)
	notify("all-day", at(7, 0), 1)
	notify("all-day again", at(7, 1), 0)
}
