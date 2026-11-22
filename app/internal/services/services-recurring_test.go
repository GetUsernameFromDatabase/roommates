package services

import (
	"errors"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestNextDue(t *testing.T) {
	day := func(s string) time.Time {
		d, err := time.Parse(time.DateOnly, s)
		if err != nil {
			t.Fatal(err)
		}
		return d
	}
	for _, tc := range []struct {
		anchor, interval, after, want string
	}{
		{"2026-01-31", g.IntervalMonth, "2026-01-31", "2026-02-28"},
		{"2026-01-31", g.IntervalMonth, "2026-02-28", "2026-03-31"},
		{"2028-01-31", g.IntervalMonth, "2028-01-31", "2028-02-29"},
		{"2026-01-31", g.IntervalMonth, "2026-03-31", "2026-04-30"},
		{"2026-11-30", g.IntervalMonth, "2026-12-30", "2027-01-30"},
		{"2028-02-29", g.IntervalYear, "2028-02-29", "2029-02-28"},
		{"2028-02-29", g.IntervalYear, "2031-02-28", "2032-02-29"},
		{"2026-12-28", g.IntervalWeek, "2026-12-28", "2027-01-04"},
		// a missed stretch: the first one after, not one interval on
		{"2026-01-01", g.IntervalWeek, "2026-01-20", "2026-01-22"},
	} {
		got := nextDue(day(tc.anchor), tc.interval, day(tc.after))
		if got.Format(time.DateOnly) != tc.want {
			t.Errorf("%s %s after %s: got %s, want %s", tc.anchor, tc.interval, tc.after, got.Format(time.DateOnly), tc.want)
		}
	}
}

func daysAgo(n int) string {
	return time.Now().AddDate(0, 0, -n).Format(time.DateOnly)
}

func recurringInput(title, date string, members ...pgtype.UUID) RecurringInput {
	return RecurringInput{
		Title:      title,
		Amount:     "10",
		Category:   "rent",
		IncurredOn: date,
		Custom:     false,
		Members:    members,
		Shares:     nil,
		ListID:     pgtype.UUID{},
		ItemID:     pgtype.UUID{},
		Interval:   g.IntervalWeek,
	}
}

func wantField(t *testing.T, step string, err error, field string) {
	t.Helper()
	var fe FieldError
	if !errors.As(err, &fe) || fe.Field != field {
		t.Fatalf("%s: got %v, want a FieldError for %s", step, err, field)
	}
}

// mustRecurring returns the only schedule of the household titled title
func mustRecurring(t *testing.T, s *DefaultService, actor Membership, title string) RecurringRow {
	t.Helper()
	rows, err := s.ListRecurring(t.Context(), actor)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range rows {
		if r.Title == title {
			r, err := s.GetRecurring(t.Context(), actor, r.ID)
			if err != nil {
				t.Fatal(err)
			}
			return r
		}
	}
	t.Fatalf("no schedule %s in %v", title, rows)
	return RecurringRow{}
}

func TestRecurringDrafts(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, strangerM := newTestHousehold(t, s)

	if _, err := s.CreateRecurring(ctx, ownerM, recurringInput("Rent", daysAgo(14), ownerM.UserID, memberM.UserID)); err != nil {
		t.Fatal(err)
	}
	r := mustRecurring(t, s, ownerM, "Rent")
	if r.Anchor.Time.Format(time.DateOnly) != daysAgo(14) || r.NextDue.Time.Format(time.DateOnly) != daysAgo(7) || len(r.Shares) != 2 {
		t.Fatalf("schedule: got %+v", r)
	}
	if _, err := s.GetRecurring(ctx, strangerM, r.ID); !errors.Is(err, g.ErrorNotFound) {
		t.Fatalf("stranger read the schedule: %v", err)
	}
	if due, err := s.ListDueRecurring(ctx, ownerM); err != nil || len(due) != 1 {
		t.Fatalf("owner's due: got %v %v", due, err)
	}
	if due, err := s.ListDueRecurring(ctx, memberM); err != nil || len(due) != 0 {
		t.Fatalf("member's due: got %v %v", due, err)
	}

	draft := r.Input().ExpenseInput
	draft.RecurringID, draft.RecurringDue = r.ID, r.NextDue
	if _, err := s.CreateExpense(ctx, memberM, draft); !errors.Is(err, g.ErrorNotFound) {
		t.Fatalf("member added the owner's draft: %v", err)
	}
	if _, err := s.CreateExpense(ctx, ownerM, draft); err != nil {
		t.Fatal(err)
	}
	_, err := s.CreateExpense(ctx, ownerM, draft)
	wantField(t, "second add", err, FieldRecurring)
	r = mustRecurring(t, s, ownerM, "Rent")
	if r.NextDue.Time.Format(time.DateOnly) != daysAgo(0) {
		t.Fatalf("after add: next due %v", r.NextDue)
	}

	wantField(t, "stale skip", s.SkipRecurring(ctx, ownerM, r.ID, draft.RecurringDue), FieldRecurring)
	if err := s.SkipRecurring(ctx, memberM, r.ID, r.NextDue); !errors.Is(err, g.ErrorNotFound) {
		t.Fatalf("member skipped the owner's draft: %v", err)
	}
	if err := s.SkipRecurring(ctx, ownerM, r.ID, r.NextDue); err != nil {
		t.Fatal(err)
	}
	if due, err := s.ListDueRecurring(ctx, ownerM); err != nil || len(due) != 0 {
		t.Fatalf("after skip: due %v %v", due, err)
	}
}

func TestRecurringChanges(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, _ := newTestHousehold(t, s)

	// each pays only their own share, so nobody owes anything and the member can leave
	if _, err := s.CreateRecurring(ctx, ownerM, recurringInput("Internet", daysAgo(1), ownerM.UserID)); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CreateRecurring(ctx, memberM, recurringInput("Gym", daysAgo(1), memberM.UserID)); err != nil {
		t.Fatal(err)
	}
	internet, gym := mustRecurring(t, s, ownerM, "Internet"), mustRecurring(t, s, ownerM, "Gym")

	if err := s.UpdateRecurring(ctx, memberM, internet.ID, recurringInput("Internet", daysAgo(1), ownerM.UserID)); !errors.Is(err, g.ErrorNotFound) {
		t.Fatalf("member edited the owner's schedule: %v", err)
	}
	if err := s.DeleteRecurring(ctx, memberM, internet.ID); !errors.Is(err, g.ErrorNotFound) {
		t.Fatalf("member deleted the owner's schedule: %v", err)
	}
	// an owner edits the member's, it stays theirs
	in := recurringInput("Gym", daysAgo(1), memberM.UserID)
	in.Interval = g.IntervalMonth
	if err := s.UpdateRecurring(ctx, ownerM, gym.ID, in); err != nil {
		t.Fatal(err)
	}
	if gym = mustRecurring(t, s, ownerM, "Gym"); gym.Interval != g.IntervalMonth || gym.PayerID != memberM.UserID {
		t.Fatalf("owner's edit: got %+v", gym)
	}
	// a new interval alone counts from the next due day, the old anchor is on the old step
	if err := s.SkipRecurring(ctx, memberM, gym.ID, gym.NextDue); err != nil {
		t.Fatal(err)
	}
	gym = mustRecurring(t, s, ownerM, "Gym")
	next := gym.NextDue.Time
	in = gym.Input()
	in.Interval = g.IntervalYear
	if err := s.UpdateRecurring(ctx, memberM, gym.ID, in); err != nil {
		t.Fatal(err)
	}
	if gym = mustRecurring(t, s, ownerM, "Gym"); !gym.Anchor.Time.Equal(next) || !gym.NextDue.Time.Equal(next) {
		t.Fatalf("new interval: got anchor %v next %v, want both %v", gym.Anchor, gym.NextDue, next)
	}
	// a new date starts it again from that day
	if err := s.UpdateRecurring(ctx, ownerM, internet.ID, recurringInput("Internet", daysAgo(3), ownerM.UserID, memberM.UserID)); err != nil {
		t.Fatal(err)
	}
	internet = mustRecurring(t, s, ownerM, "Internet")
	if internet.Anchor.Time.Format(time.DateOnly) != daysAgo(3) || internet.NextDue.Time.Format(time.DateOnly) != daysAgo(3) || len(internet.Shares) != 2 {
		t.Fatalf("new date: got %+v", internet)
	}

	if err := s.Leave(ctx, memberM); err != nil {
		t.Fatal(err)
	}
	internet = mustRecurring(t, s, ownerM, "Internet")
	if len(internet.Shares) != 1 || internet.Shares[0].UserID != ownerM.UserID {
		t.Fatalf("the leaver's share was kept: %+v", internet.Shares)
	}
	if rows, err := s.ListRecurring(ctx, ownerM); err != nil || len(rows) != 1 {
		t.Fatalf("the leaver's schedule was kept: %v %v", rows, err)
	}
	if err := s.DeleteRecurring(ctx, ownerM, internet.ID); err != nil {
		t.Fatal(err)
	}
}

func TestNotifyRecurringDue(t *testing.T) {
	s := newTestService(t)
	locales.LoadLocales()
	ctx := t.Context()
	ownerM, _, _ := newTestHousehold(t, s)
	device := subscribeTestDevice(t, s, ownerM.UserID, "en")
	if _, err := s.CreateRecurring(ctx, ownerM, recurringInput("Rent", daysAgo(21), ownerM.UserID)); err != nil {
		t.Fatal(err)
	}
	// other households' schedules can be due too, count only this device
	notify := func(step string, want int) {
		t.Helper()
		if err := s.NotifyRecurringDue(ctx); err != nil {
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
	notify("due", 1)
	notify("again", 0)
	r := mustRecurring(t, s, ownerM, "Rent")
	if err := s.SkipRecurring(ctx, ownerM, r.ID, r.NextDue); err != nil {
		t.Fatal(err)
	}
	notify("next due", 1)
	// skipped, then edited back to the date it was pushed about: that is a new due, pushed again
	r = mustRecurring(t, s, ownerM, "Rent")
	if err := s.SkipRecurring(ctx, ownerM, r.ID, r.NextDue); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateRecurring(ctx, ownerM, r.ID, recurringInput("Rent", r.NextDue.Time.Format(time.DateOnly), ownerM.UserID)); err != nil {
		t.Fatal(err)
	}
	notify("moved back", 1)
}
