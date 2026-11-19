package services

import (
	"context"
	"errors"
	"testing"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"

	"github.com/jackc/pgx/v5/pgtype"
)

// editHelpers are the steps of the expense edit tests, the owner pays every expense
func editHelpers(t *testing.T, s *DefaultService, ownerM Membership) (
	create func(shares map[pgtype.UUID]string, amount string) pgtype.UUID,
	edit func(actor Membership, id pgtype.UUID, shares map[pgtype.UUID]string, amount string) error,
	net func(m Membership) int64,
	fromEdit func(id pgtype.UUID) []dbqueries.ListDebtItemsRow,
) {
	t.Helper()
	ctx := context.Background()
	input := func(shares map[pgtype.UUID]string, amount string) ExpenseInput {
		return ExpenseInput{
			Title:      "Bill",
			Amount:     amount,
			Category:   "other",
			IncurredOn: "2026-10-01",
			Custom:     true,
			Members:    nil,
			Shares:     shares,
			ListID:     pgtype.UUID{},
			ItemID:     pgtype.UUID{},
		}
	}
	create = func(shares map[pgtype.UUID]string, amount string) pgtype.UUID {
		t.Helper()
		id, err := s.CreateExpense(ctx, ownerM, input(shares, amount))
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
	edit = func(actor Membership, id pgtype.UUID, shares map[pgtype.UUID]string, amount string) error {
		return s.UpdateExpense(ctx, actor, id, input(shares, amount))
	}
	// what m is owed in all, the overview's netting
	net = func(m Membership) int64 {
		t.Helper()
		debts, err := s.Debts(ctx, m)
		if err != nil {
			t.Fatal(err)
		}
		var n int64
		for _, d := range debts {
			switch m.UserID {
			case d.To:
				n += d.Cents
			case d.From:
				n -= d.Cents
			}
		}
		return n
	}
	fromEdit = func(id pgtype.UUID) []dbqueries.ListDebtItemsRow {
		t.Helper()
		open, err := s.db.ListDebtItems(ctx, dbqueries.ListDebtItemsParams{HouseholdID: ownerM.HouseholdID, UserID: pgtype.UUID{}})
		if err != nil {
			t.Fatal(err)
		}
		var made []dbqueries.ListDebtItemsRow
		for _, st := range open {
			if st.FromEdit && st.EditExpenseID == id {
				made = append(made, st)
			}
		}
		return made
	}
	return create, edit, net, fromEdit
}

func TestEditExpenseOverpaidInApp(t *testing.T) {
	s := newTestService(t)
	ownerM, memberM, _ := newTestHousehold(t, s)
	_, pay, wantDebt, wantPaid := settlementHelpers(t, s, ownerM, memberM)
	create, edit, _, fromEdit := editHelpers(t, s, ownerM)

	id := create(map[pgtype.UUID]string{memberM.UserID: "10"}, "10")
	pay(memberM, ownerM, "10")
	wantPaid("paid in the app", id, true)
	if err := edit(ownerM, id, map[pgtype.UUID]string{memberM.UserID: "6"}, "6"); err != nil {
		t.Fatal(err)
	}
	wantDebt("lowered, the owner owes the rest back", -400)
	wantPaid("the payment still covers the share", id, true)
	if made := fromEdit(id); len(made) != 0 {
		t.Fatalf("nothing was paid outside the app, yet %+v", made)
	}
}

func TestEditExpenseOverpaidOutside(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)
	_, _, wantDebt, wantPaid := settlementHelpers(t, s, ownerM, memberM)
	create, edit, _, fromEdit := editHelpers(t, s, ownerM)

	id := create(map[pgtype.UUID]string{memberM.UserID: "10"}, "10")
	if err := s.SetExpenseSharePaid(ctx, ownerM, id, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}
	if err := edit(ownerM, id, map[pgtype.UUID]string{memberM.UserID: "6"}, "6"); err != nil {
		t.Fatal(err)
	}
	wantDebt("lowered, the owner owes the rest back", -400)
	wantPaid("the cash covers the share", id, true)
	made := fromEdit(id)
	// a settlement is listed as a debt the other way, so the member who paid is ToUser
	if len(made) != 1 || made[0].ToUser != memberM.UserID || made[0].AmountCents != 400 || made[0].Method != nil {
		t.Fatalf("payment from the edit = %+v", made)
	}
}

func TestEditExpenseUnderpaid(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)
	_, pay, wantDebt, wantPaid := settlementHelpers(t, s, ownerM, memberM)
	create, edit, _, _ := editHelpers(t, s, ownerM)

	id := create(map[pgtype.UUID]string{memberM.UserID: "10"}, "10")
	if err := s.SetExpenseSharePaid(ctx, ownerM, id, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}
	if err := edit(ownerM, id, map[pgtype.UUID]string{memberM.UserID: "15"}, "15"); err != nil {
		t.Fatal(err)
	}
	wantDebt("raised, the member owes the rest", 500)
	wantPaid("the cash does not cover the share", id, false)
	pay(memberM, ownerM, "5")
	wantDebt("the rest paid", 0)
	wantPaid("the rest paid", id, true)
}

func TestEditExpenseParticipants(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)
	_, _, wantDebt, _ := settlementHelpers(t, s, ownerM, memberM)
	create, edit, _, _ := editHelpers(t, s, ownerM)

	// a paid share removed: what the member paid is owed back
	id := create(map[pgtype.UUID]string{ownerM.UserID: "5", memberM.UserID: "5"}, "10")
	if err := s.SetExpenseSharePaid(ctx, ownerM, id, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}
	if err := edit(ownerM, id, map[pgtype.UUID]string{ownerM.UserID: "10"}, "10"); err != nil {
		t.Fatal(err)
	}
	wantDebt("paid share removed", -500)
	if e, err := s.GetExpense(ctx, ownerM, id); err != nil || len(e.Shares) != 1 {
		t.Fatalf("after removing = %+v, %v", e, err)
	}

	// someone added gets an unpaid share
	alone := create(map[pgtype.UUID]string{ownerM.UserID: "8"}, "8")
	if err := edit(ownerM, alone, map[pgtype.UUID]string{ownerM.UserID: "4", memberM.UserID: "4"}, "8"); err != nil {
		t.Fatal(err)
	}
	wantDebt("member added", -100)
}

func TestEditExpenseRefused(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	create, edit, _, _ := editHelpers(t, s, ownerM)
	shares := map[pgtype.UUID]string{memberM.UserID: "10"}
	id := create(shares, "10")

	wantErr(t, "member edits payer's expense", edit(memberM, id, shares, "10"), g.ErrorNotFound)
	wantErr(t, "stranger edits", edit(strangerM, id, shares, "10"), g.ErrorNotFound)
	var fe FieldError
	if err := edit(ownerM, id, map[pgtype.UUID]string{memberM.UserID: "9"}, "10"); !errors.As(err, &fe) || fe.Field != FieldShares {
		t.Fatalf("shares not adding up: %v", err)
	}

	// the household moved to another currency, the old expense keeps its amount
	if _, err := s.pool.Exec(ctx, "UPDATE households SET currency = 'USD' WHERE id = $1", ownerM.HouseholdID); err != nil {
		t.Fatal(err)
	}
	if err := edit(ownerM, id, map[pgtype.UUID]string{memberM.UserID: "12"}, "12"); !errors.As(err, &fe) || fe.Field != FieldAmount {
		t.Fatalf("amount of an old currency expense: %v", err)
	}
	// the disabled amount is not sent, the rest still changes
	if err := s.UpdateExpense(ctx, ownerM, id, ExpenseInput{Title: "Renamed", Amount: "", Category: "rent", IncurredOn: "2026-10-02",
		Custom: false, Members: nil, Shares: nil, ListID: pgtype.UUID{}, ItemID: pgtype.UUID{}}); err != nil {
		t.Fatal(err)
	}
	e, err := s.GetExpense(ctx, ownerM, id)
	if err != nil || e.Expense.Title != "Renamed" || e.Expense.AmountCents != 1000 || e.Expense.Currency == "USD" ||
		e.AmountLock("USD") != AmountLockCurrency {
		t.Fatalf("after editing an old currency expense = %+v, %v", e, err)
	}
}

// whatever the share was paid with, each balance moves by exactly the change of the share
func TestEditExpenseMovesBalanceByTheChange(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	for _, paid := range []string{"unpaid", "in app", "outside", "partly in app"} {
		for _, newShare := range []int64{0, 300, 600, 900} {
			ownerM, memberM, _ := newTestHousehold(t, s)
			_, pay, _, _ := settlementHelpers(t, s, ownerM, memberM)
			create, edit, net, _ := editHelpers(t, s, ownerM)
			id := create(map[pgtype.UUID]string{ownerM.UserID: "4", memberM.UserID: "6"}, "10")
			switch paid {
			case "in app":
				pay(memberM, ownerM, "6")
			case "outside":
				if err := s.SetExpenseSharePaid(ctx, ownerM, id, memberM.UserID, true); err != nil {
					t.Fatal(err)
				}
			case "partly in app":
				pay(memberM, ownerM, "2")
				if err := s.SetExpenseSharePaid(ctx, ownerM, id, memberM.UserID, true); err != nil {
					t.Fatal(err)
				}
			}
			ownerBefore, memberBefore := net(ownerM), net(memberM)

			shares := map[pgtype.UUID]string{ownerM.UserID: "4"}
			if newShare > 0 {
				shares[memberM.UserID] = FormatCentsInput(newShare)
			}
			if err := edit(ownerM, id, shares, FormatCentsInput(400+newShare)); err != nil {
				t.Fatalf("%s to %d: %v", paid, newShare, err)
			}
			change := newShare - 600
			if got := net(memberM) - memberBefore; got != -change {
				t.Errorf("%s to %d: member's balance moved %d, want %d", paid, newShare, got, -change)
			}
			if got := net(ownerM) - ownerBefore; got != change {
				t.Errorf("%s to %d: owner's balance moved %d, want %d", paid, newShare, got, change)
			}
		}
	}
}

func TestHouseholdSettled(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)
	expense, pay, _, _ := settlementHelpers(t, s, ownerM, memberM)
	want := func(step string, settled bool) {
		t.Helper()
		got, err := householdSettled(ctx, s.db, ownerM.HouseholdID)
		if err != nil || got != settled {
			t.Fatalf("%s: settled %v, %v, want %v", step, got, err, settled)
		}
	}

	want("nothing yet", true)
	expense(ownerM, memberM, "10")
	want("member owes", false)
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "10", "cash", ""); err != nil {
		t.Fatal(err)
	}
	want("payment pending", false)
	pending, err := s.ListPendingSettlements(ctx, ownerM)
	if err != nil || len(pending) != 1 {
		t.Fatalf("pending = %+v, %v", pending, err)
	}
	if err := s.ConfirmSettlement(ctx, ownerM, pending[0].ID); err != nil {
		t.Fatal(err)
	}
	want("paid", true)
	// shares both ways that net to zero are settled too
	expense(ownerM, memberM, "3")
	expense(memberM, ownerM, "3")
	want("netted", true)
	pay(memberM, ownerM, "1")
	want("overpaid", false)
}
