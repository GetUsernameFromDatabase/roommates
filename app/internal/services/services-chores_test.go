package services

import (
	"errors"
	"slices"
	"testing"
	"time"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestChoreNextDue(t *testing.T) {
	day := func(d int) time.Time { return time.Date(2026, 11, d, 0, 0, 0, 0, time.UTC) } // the 16th is a Monday
	monThu := weekdayBit(time.Monday) | weekdayBit(time.Thursday)
	every := func(n int16) *int16 { return &n }
	unit := func(u string) *string { return &u }
	for _, tc := range []struct {
		name string
		got  time.Time
		want time.Time
	}{
		{"done late Tue", choreNextDue(nil, nil, &monThu, day(16), day(17)), day(19)},
		{"done late Sat", choreNextDue(nil, nil, &monThu, day(19), day(21)), day(23)},
		{"done on time Thu", choreNextDue(nil, nil, &monThu, day(19), day(19)), day(23)},
		{"done on time Mon", choreNextDue(nil, nil, &monThu, day(16), day(16)), day(19)},
		{"done early Fri for Mon", choreNextDue(nil, nil, &monThu, day(23), day(20)), day(26)},
		{"first due on a ticked today", nextWeekday(day(19), monThu), day(19)},
		{"first due on Fri", nextWeekday(day(20), monThu), day(23)},
		{"every 3 days", choreNextDue(every(3), unit(g.ChoreUnitDay), nil, day(20), day(17)), day(20)},
		{"every 2 weeks", choreNextDue(every(2), unit(g.ChoreUnitWeek), nil, day(17), day(17)), time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)},
		{"every 2 months", choreNextDue(every(2), unit(g.ChoreUnitMonth), nil, day(17), day(19)), time.Date(2027, 1, 19, 0, 0, 0, 0, time.UTC)},
		{"monthly from the 31st", choreNextDue(every(1), unit(g.ChoreUnitMonth), nil, day(17), time.Date(2027, 1, 31, 0, 0, 0, 0, time.UTC)), time.Date(2027, 2, 28, 0, 0, 0, 0, time.UTC)},
	} {
		if !tc.got.Equal(tc.want) {
			t.Errorf("%s: got %s, want %s", tc.name, tc.got.Format(time.DateOnly), tc.want.Format(time.DateOnly))
		}
	}
	if got := weekdaysOf(monThu | weekdayBit(time.Sunday)); !slices.Equal(got, []time.Weekday{time.Monday, time.Thursday, time.Sunday}) {
		t.Errorf("weekdaysOf: got %v", got)
	}
}

func TestChoreClean(t *testing.T) {
	today := dateOf(time.Date(2026, 11, 19, 0, 0, 0, 0, time.UTC))
	a, b, stranger := uid(1), uid(2), uid(3)
	active := []dbqueries.ListActiveMembersRow{{UserID: a, Username: "a", AvatarKey: nil, Role: g.RoleOwner, JoinedAt: pgtype.Timestamptz{}}, {UserID: b, Username: "b", AvatarKey: nil, Role: g.RoleOwner, JoinedAt: pgtype.Timestamptz{}}}
	ok := func() ChoreInput {
		return ChoreInput{Title: " Bins ", Rotation: g.RotationFixed, Unit: g.ChoreUnitWeek, Every: 2, Weekdays: nil, FirstDue: "2026-11-19", Members: []pgtype.UUID{a, b}, Assignee: b, Shown: ChoreTurn{}}
	}
	f, err := ok().clean(today, pgtype.Date{}, active)
	if err != nil || f.title != "Bins" || *f.every != 2 || *f.unit != g.ChoreUnitWeek || f.weekdays != nil || !f.due.Equal(today.Time) {
		t.Fatalf("ok: got %+v %v", f, err)
	}
	in := ok()
	in.Unit, in.Weekdays = g.ChoreUnitWeekdays, []time.Weekday{time.Monday}
	if f, err := in.clean(today, pgtype.Date{}, active); err != nil || f.every != nil || *f.weekdays != 1 || f.due.Format(time.DateOnly) != "2026-11-23" {
		t.Fatalf("weekdays: got %+v %v", f, err)
	}
	// the current due date may stay even when it is past
	in = ok()
	in.FirstDue = "2026-11-10"
	if _, err := in.clean(today, dateOf(time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC)), active); err != nil {
		t.Fatalf("current past due: %v", err)
	}
	for _, tc := range []struct {
		field string
		edit  func(*ChoreInput)
	}{
		{FieldChoreTitle, func(in *ChoreInput) { in.Title = "  " }},
		{FieldChoreTitle, func(in *ChoreInput) { in.Title = string(make([]rune, MaxChoreTitleLength+1)) + "x" }},
		{FieldChoreRotation, func(in *ChoreInput) { in.Rotation = "random" }},
		{FieldChoreSchedule, func(in *ChoreInput) { in.Unit = "year" }},
		{FieldChoreSchedule, func(in *ChoreInput) { in.Unit = g.ChoreUnitMonth; in.Every = 13 }},
		{FieldChoreSchedule, func(in *ChoreInput) { in.Every = 0 }},
		{FieldChoreSchedule, func(in *ChoreInput) { in.Every = 53 }},
		{FieldChoreSchedule, func(in *ChoreInput) { in.Unit = g.ChoreUnitDay; in.Every = 366 }},
		{FieldChoreSchedule, func(in *ChoreInput) { in.Unit = g.ChoreUnitWeekdays }},
		{FieldChoreFirstDue, func(in *ChoreInput) { in.FirstDue = "2026-11-18" }},
		{FieldChoreFirstDue, func(in *ChoreInput) { in.FirstDue = "tomorrow" }},
		{FieldChoreMembers, func(in *ChoreInput) { in.Members = nil }},
		{FieldChoreMembers, func(in *ChoreInput) { in.Members = []pgtype.UUID{a, a} }},
		{FieldChoreMembers, func(in *ChoreInput) { in.Members = []pgtype.UUID{a, stranger} }},
		{FieldChoreAssignee, func(in *ChoreInput) { in.Members = []pgtype.UUID{a} }},
	} {
		in := ok()
		tc.edit(&in)
		if _, err := in.clean(today, pgtype.Date{}, active); !errors.Is(err, FieldError{Field: tc.field}) {
			t.Errorf("%s: got %v", tc.field, err)
		}
	}
}

// choreTestHousehold is newTestHousehold with a third member
func choreTestHousehold(t *testing.T, s *DefaultService) (ownerM, memberM, thirdM, strangerM Membership) {
	t.Helper()
	ownerM, memberM, strangerM = newTestHousehold(t, s)
	overview, err := s.GetHouseholdOverview(t.Context(), ownerM)
	if err != nil {
		t.Fatal(err)
	}
	third := newTestUser(t, s)
	if _, err := s.JoinWithInvite(t.Context(), overview.Invites[0].Code, third); err != nil {
		t.Fatal(err)
	}
	return ownerM, memberM, mustMembership(t, s, ownerM.HouseholdID, third), strangerM
}

func everyDayChore(title, rotation string, assignee pgtype.UUID, members ...pgtype.UUID) ChoreInput {
	return ChoreInput{
		Title:    title,
		Rotation: rotation,
		Unit:     g.ChoreUnitWeekdays,
		Every:    0,
		Weekdays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday},
		FirstDue: "",
		Members:  members,
		Assignee: assignee,
		Shown:    ChoreTurn{},
	}
}

func mustCreateChore(t *testing.T, s *DefaultService, actor Membership, in ChoreInput) pgtype.UUID {
	t.Helper()
	id, err := s.CreateChore(t.Context(), actor, in)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func mustChore(t *testing.T, s *DefaultService, actor Membership, id pgtype.UUID) ChoreDetail {
	t.Helper()
	c, err := s.GetChore(t.Context(), actor, id)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// wantTurn checks whose turn chore id is and how many days from today it is due
func wantTurn(t *testing.T, s *DefaultService, actor Membership, id, assignee pgtype.UUID, inDays int, step string) Chore {
	t.Helper()
	c := mustChore(t, s, actor, id).Chore
	if c.AssigneeID != assignee || !c.NextDue.Equal(c.Today.AddDate(0, 0, inDays)) {
		t.Fatalf("%s: got turn of %v due %s (today %s), want %v in %d days", step, c.AssigneeID, c.NextDue.Format(time.DateOnly), c.Today.Format(time.DateOnly), assignee, inDays)
	}
	return c
}

func TestChoreFixedRotation(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, thirdM, strangerM := choreTestHousehold(t, s)
	id := mustCreateChore(t, s, ownerM, everyDayChore("Bins", g.RotationFixed, pgtype.UUID{}, ownerM.UserID, memberM.UserID, thirdM.UserID))
	c := wantTurn(t, s, ownerM, id, ownerM.UserID, 0, "created")

	// anyone can do it, the turn still moves on from whose it was
	turn := c.Turn()
	if err := s.DoneChore(ctx, memberM, id, turn); err != nil {
		t.Fatal(err)
	}
	wantTurn(t, s, ownerM, id, memberM.UserID, 1, "done")
	// a second click on the same turn does nothing
	wantErr(t, "done twice", s.DoneChore(ctx, memberM, id, turn), ErrChoreTurnMoved)
	c = wantTurn(t, s, ownerM, id, memberM.UserID, 1, "done twice")

	if err := s.SkipChore(ctx, memberM, id, c.Turn()); err != nil {
		t.Fatal(err)
	}
	// skipped before its day, the next ticked day after it
	c = wantTurn(t, s, ownerM, id, thirdM.UserID, 2, "skipped")
	if err := s.PassChore(ctx, thirdM, id, c.Turn()); err != nil {
		t.Fatal(err)
	}
	// a pass keeps the date, a second one on the same page is still stale
	wantErr(t, "passed twice", s.PassChore(ctx, thirdM, id, c.Turn()), ErrChoreTurnMoved)
	wantTurn(t, s, ownerM, id, ownerM.UserID, 2, "passed")

	d := mustChore(t, s, ownerM, id)
	kinds := historyKinds(d.History)
	if !slices.Equal(kinds, []string{g.ChorePassed, g.ChoreSkipped, g.ChoreDone}) || len(d.Chore.Members) != 3 || d.Chore.Members[1].LastDone == nil || d.Chore.Members[0].LastDone != nil {
		t.Fatalf("history: got %v, members %+v", kinds, d.Chore.Members)
	}
	wantKinds(t, s, ownerM, g.ActChoreAdded, g.ActChoreDone, g.ActChoreSkipped, g.ActChorePassed)

	// another household's chore is not found
	_, err := s.GetChore(ctx, strangerM, id)
	wantErr(t, "stranger get", err, g.ErrorNotFound)
	wantErr(t, "stranger done", s.DoneChore(ctx, strangerM, id, c.Turn()), g.ErrorNotFound)
	wantErr(t, "stranger delete", s.DeleteChore(ctx, strangerM, id), g.ErrorNotFound)

	// the rota moves
	wantErr(t, "move stranger", s.MoveChoreMember(ctx, ownerM, id, strangerM.UserID, true), g.ErrorNotFound)
	if err := s.MoveChoreMember(ctx, ownerM, id, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}
	if err := s.MoveChoreMember(ctx, ownerM, id, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}
	if got := rotaOf(mustChore(t, s, ownerM, id).Chore); !slices.Equal(got, []pgtype.UUID{memberM.UserID, ownerM.UserID, thirdM.UserID}) {
		t.Fatalf("moved: got %v", got)
	}

	// an edit keeps the turn of someone still in the rota
	in := everyDayChore("Bins out", g.RotationFixed, pgtype.UUID{}, thirdM.UserID, ownerM.UserID)
	if err := s.UpdateChore(ctx, memberM, id, in); err != nil {
		t.Fatal(err)
	}
	if c := wantTurn(t, s, ownerM, id, ownerM.UserID, 2, "edited"); c.Title != "Bins out" || len(c.Members) != 2 {
		t.Fatalf("edited: got %+v", c)
	}
	in.Members = []pgtype.UUID{thirdM.UserID}
	if err := s.UpdateChore(ctx, memberM, id, in); err != nil {
		t.Fatal(err)
	}
	c = wantTurn(t, s, ownerM, id, thirdM.UserID, 2, "edited without the assignee")
	// a rota of one cannot pass
	wantErr(t, "pass alone", s.PassChore(ctx, thirdM, id, c.Turn()), g.ErrorInvalidInput)

	if err := s.DeleteChore(ctx, memberM, id); err != nil {
		t.Fatal(err)
	}
	wantKinds(t, s, ownerM, g.ActChoreDeleted)
}

// an edit form opened before a Done does not undo it when saved
func TestChoreStaleEdit(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, _, _ := choreTestHousehold(t, s)
	today, err := s.HouseholdToday(ctx, ownerM)
	if err != nil {
		t.Fatal(err)
	}
	in := everyDayChore("Bins", g.RotationFixed, ownerM.UserID, ownerM.UserID, memberM.UserID)
	in.Unit, in.Every, in.Weekdays, in.FirstDue = g.ChoreUnitWeek, 1, nil, today.Format(time.DateOnly)
	id := mustCreateChore(t, s, ownerM, in)
	c := wantTurn(t, s, ownerM, id, ownerM.UserID, 0, "created")

	// the form as the edit page shows it
	in.FirstDue, in.Assignee, in.Shown = c.NextDue.Format(time.DateOnly), c.AssigneeID, c.Turn()
	if err := s.DoneChore(ctx, ownerM, id, c.Turn()); err != nil {
		t.Fatal(err)
	}
	wantTurn(t, s, ownerM, id, memberM.UserID, 7, "done")
	in.Title = "Bins out"
	if err := s.UpdateChore(ctx, memberM, id, in); err != nil {
		t.Fatal(err)
	}
	wantTurn(t, s, ownerM, id, memberM.UserID, 7, "stale edit")

	// the stale form takes off whose turn it is now: its shown person gets the turn, no error about
	// someone it does not show
	in.Members = []pgtype.UUID{ownerM.UserID}
	if err := s.UpdateChore(ctx, ownerM, id, in); err != nil {
		t.Fatal(err)
	}
	wantTurn(t, s, ownerM, id, ownerM.UserID, 7, "stale edit without the turn's person")

	// a field changed on the form still changes
	in.Members = []pgtype.UUID{ownerM.UserID, memberM.UserID}
	c = mustChore(t, s, ownerM, id).Chore
	in.Shown = c.Turn()
	in.Assignee, in.FirstDue = memberM.UserID, c.Today.AddDate(0, 0, 3).Format(time.DateOnly)
	if err := s.UpdateChore(ctx, memberM, id, in); err != nil {
		t.Fatal(err)
	}
	wantTurn(t, s, ownerM, id, memberM.UserID, 3, "changed")
}

func historyKinds(h []ChoreLogEntry) []string {
	var kinds []string
	for _, e := range h {
		kinds = append(kinds, e.Kind)
	}
	return kinds
}

func rotaOf(c Chore) []pgtype.UUID {
	var ids []pgtype.UUID
	for _, m := range c.Members {
		ids = append(ids, m.UserID)
	}
	return ids
}

func TestChoreLeastRecent(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, thirdM, _ := choreTestHousehold(t, s)
	id := mustCreateChore(t, s, ownerM, everyDayChore("Dishes", g.RotationLeastRecent, pgtype.UUID{}, ownerM.UserID, memberM.UserID, thirdM.UserID))
	step := func(f func(Membership, pgtype.UUID, ChoreTurn) error, actor Membership, want pgtype.UUID, name string) {
		t.Helper()
		c := mustChore(t, s, ownerM, id).Chore
		if err := f(actor, id, c.Turn()); err != nil {
			t.Fatal(err)
		}
		if got := mustChore(t, s, ownerM, id).Chore.AssigneeID; got != want {
			t.Fatalf("%s: got turn of %v, want %v", name, got, want)
		}
	}
	done := func(a Membership, id pgtype.UUID, turn ChoreTurn) error { return s.DoneChore(ctx, a, id, turn) }
	skip := func(a Membership, id pgtype.UUID, turn ChoreTurn) error { return s.SkipChore(ctx, a, id, turn) }
	wantTurn(t, s, ownerM, id, ownerM.UserID, 0, "created")
	step(done, ownerM, memberM.UserID, "owner done")
	// a skip is not done, the skipper stays first in line
	step(skip, memberM, memberM.UserID, "member skipped")
	step(done, memberM, thirdM.UserID, "member done")
	step(done, thirdM, ownerM.UserID, "third done")
	// the third does the owner's turn, the owner has done it longest ago
	step(done, thirdM, ownerM.UserID, "third did the owner's")
}

func TestChoreManual(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, _, strangerM := choreTestHousehold(t, s)
	id := mustCreateChore(t, s, ownerM, everyDayChore("Plants", g.RotationManual, memberM.UserID, ownerM.UserID, memberM.UserID))
	c := wantTurn(t, s, ownerM, id, memberM.UserID, 0, "created")
	wantErr(t, "pass manual", s.PassChore(ctx, memberM, id, c.Turn()), g.ErrorInvalidInput)
	if err := s.DoneChore(ctx, memberM, id, c.Turn()); err != nil {
		t.Fatal(err)
	}
	wantTurn(t, s, ownerM, id, memberM.UserID, 1, "done stays")

	unassigned := mustCreateChore(t, s, ownerM, everyDayChore("Fridge", g.RotationManual, pgtype.UUID{}, ownerM.UserID))
	wantTurn(t, s, ownerM, unassigned, pgtype.UUID{}, 0, "unassigned")
	wantErr(t, "assign stranger", s.AssignChore(ctx, ownerM, unassigned, strangerM.UserID), g.ErrorInvalidInput)
	if err := s.AssignChore(ctx, memberM, unassigned, ownerM.UserID); err != nil {
		t.Fatal(err)
	}
	wantTurn(t, s, ownerM, unassigned, ownerM.UserID, 0, "assigned")
}

func TestChoreLeaver(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, thirdM, _ := choreTestHousehold(t, s)
	id := mustCreateChore(t, s, ownerM, everyDayChore("Bins", g.RotationFixed, memberM.UserID, ownerM.UserID, memberM.UserID, thirdM.UserID))
	if err := s.Leave(ctx, memberM); err != nil {
		t.Fatal(err)
	}
	wantTurn(t, s, ownerM, id, thirdM.UserID, 0, "member left")
	if err := s.RemoveMember(ctx, ownerM, thirdM.UserID); err != nil {
		t.Fatal(err)
	}
	c := wantTurn(t, s, ownerM, id, ownerM.UserID, 0, "third removed")
	if got := rotaOf(c); !slices.Equal(got, []pgtype.UUID{ownerM.UserID}) {
		t.Fatalf("rota after leaving: got %v", got)
	}
	// a leaver cannot be put in a rota
	_, err := s.CreateChore(ctx, ownerM, everyDayChore("Floor", g.RotationFixed, pgtype.UUID{}, ownerM.UserID, memberM.UserID))
	wantErr(t, "leaver in rota", err, FieldError{Field: FieldChoreMembers})
}

func TestNotifyChoresDue(t *testing.T) {
	s := newTestService(t)
	locales.LoadLocales()
	ctx := t.Context()
	ownerM, memberM, _ := newTestHousehold(t, s)
	device := subscribeTestDevice(t, s, ownerM.UserID, "en")
	// the reminder hour 0 has always passed
	if err := s.UpdateHousehold(ctx, ownerM, "Flat", ownerM.Currency, "Europe/Tallinn", 0); err != nil {
		t.Fatal(err)
	}
	// other households' chores can be due too, count only this device
	notify := func(step string, want int) {
		t.Helper()
		if err := s.NotifyChoresDue(ctx); err != nil {
			t.Fatal(err)
		}
		n := 0
		for _, p := range takePushes(s) {
			if p.endpoint == device {
				n++
			}
		}
		if n != want {
			t.Fatalf("%s: got %d pushes, want %d", step, n, want)
		}
	}
	id := mustCreateChore(t, s, memberM, everyDayChore("Bins", g.RotationFixed, ownerM.UserID, ownerM.UserID, memberM.UserID))
	notify("due today", 1)
	notify("again", 0)

	// overdue and never pushed: one push, not a due push and a nudge
	c := mustChore(t, s, ownerM, id).Chore
	yesterday := dateOf(c.Today.AddDate(0, 0, -1))
	if err := s.db.SetChoreTurn(ctx, dbqueries.SetChoreTurnParams{NextDue: yesterday, AssigneeID: ownerM.UserID, Handed: false, ID: id, HouseholdID: ownerM.HouseholdID}); err != nil {
		t.Fatal(err)
	}
	notify("overdue", 1)
	notify("overdue again", 0)

	// pushed on the day, nudged the day after
	if err := s.db.SetChoreTurn(ctx, dbqueries.SetChoreTurnParams{NextDue: yesterday, AssigneeID: ownerM.UserID, Handed: false, ID: id, HouseholdID: ownerM.HouseholdID}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.ClaimChoreNotified(ctx, dbqueries.ClaimChoreNotifiedParams{Overdue: false, ID: id, NextDue: yesterday, AssigneeID: ownerM.UserID}); err != nil {
		t.Fatal(err)
	}
	notify("nudge", 1)
	notify("nudge again", 0)

	// a pass and an assignment push the new assignee, which stands for the due push and the nudge
	if err := s.AssignChore(ctx, ownerM, id, memberM.UserID); err != nil {
		t.Fatal(err)
	}
	takePushes(s)
	if err := s.PassChore(ctx, memberM, id, ChoreTurn{Due: yesterday, Assignee: memberM.UserID}); err != nil {
		t.Fatal(err)
	}
	if got := takePushes(s); len(got) != 1 || got[0].endpoint != device || got[0].msg.Tag != "chore-"+id.String() {
		t.Fatalf("pass: got %v", got)
	}
	notify("after the pass", 0)
	// assigning whose turn it already is pushes nobody
	if err := s.AssignChore(ctx, memberM, id, ownerM.UserID); err != nil {
		t.Fatal(err)
	}
	notify("assigned again", 0)
}
