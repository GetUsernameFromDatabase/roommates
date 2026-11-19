package services

import (
	"errors"
	"testing"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"

	"github.com/jackc/pgx/v5/pgtype"
)

func notificationOn(t *testing.T, s *DefaultService, m Membership, kind string) bool {
	t.Helper()
	list, err := s.ListNotificationSettings(t.Context(), m)
	if err != nil || len(list) != len(g.NotifyKinds) {
		t.Fatalf("list: %v %v", list, err)
	}
	for i, n := range list {
		if n.Kind != g.NotifyKinds[i] {
			t.Fatalf("order: got %s at %d", n.Kind, i)
		}
		if n.Kind == kind {
			return n.On
		}
	}
	return false
}

func TestNotificationSettings(t *testing.T) {
	s := newTestService(t)
	locales.LoadLocales()
	ctx := t.Context()
	ownerM, memberM, _ := newTestHousehold(t, s)
	ownerDevice := subscribeTestDevice(t, s, ownerM.UserID, "en")
	subscribeTestDevice(t, s, memberM.UserID, "en")

	if !notificationOn(t, s, ownerM, g.NotifyChat) {
		t.Fatal("kinds are on by default")
	}
	if err := s.SetNotification(ctx, ownerM, "nope", false); !errors.Is(err, g.ErrorNotFound) {
		t.Fatalf("unknown kind: %v", err)
	}

	// muting chat stops the chat push only
	if err := s.SetNotification(ctx, ownerM, g.NotifyChat, false); err != nil {
		t.Fatal(err)
	}
	if notificationOn(t, s, ownerM, g.NotifyChat) || !notificationOn(t, s, ownerM, g.NotifyMoney) || !notificationOn(t, s, memberM, g.NotifyChat) {
		t.Fatal("only the owner's chat is off")
	}
	if err := s.PostMessage(ctx, memberM, "hi", pgtype.UUID{}, nil, nil); err != nil {
		t.Fatal(err)
	}
	wantPushes(t, "muted chat", takePushes(s))
	if _, err := s.CreateExpense(ctx, memberM, ExpenseInput{
		Title: "Rent", Amount: "10", Category: "rent", IncurredOn: "2026-10-01", Custom: false,
		Members: []pgtype.UUID{ownerM.UserID, memberM.UserID}, Shares: nil, ListID: pgtype.UUID{}, ItemID: pgtype.UUID{},
	}); err != nil {
		t.Fatal(err)
	}
	wantPushes(t, "expense still arrives", takePushes(s), ownerDevice)

	// a mute is for one household
	otherID, err := s.CreateHousehold(ctx, ownerM.UserID, "Second")
	if err != nil {
		t.Fatal(err)
	}
	otherM := mustMembership(t, s, otherID, ownerM.UserID)
	if !notificationOn(t, s, otherM, g.NotifyChat) {
		t.Fatal("the mute leaked into another household")
	}

	// on again deletes the row, and the push is back
	if err := s.SetNotification(ctx, ownerM, g.NotifyChat, true); err != nil {
		t.Fatal(err)
	}
	rows, err := s.db.ListNotificationMutes(ctx, dbqueries.ListNotificationMutesParams{
		HouseholdID: ownerM.HouseholdID, UserID: ownerM.UserID,
	})
	if err != nil || len(rows) != 0 {
		t.Fatalf("mute rows: %v %v", rows, err)
	}
	if err := s.PostMessage(ctx, memberM, "hi again", pgtype.UUID{}, nil, nil); err != nil {
		t.Fatal(err)
	}
	wantPushes(t, "chat back", takePushes(s), ownerDevice)
}
