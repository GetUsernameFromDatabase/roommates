package services

import (
	"context"
	"slices"
	"strconv"
	"strings"
	"testing"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// the form and validation read the Go list, the foreign keys the table
func TestExpenseCategoriesMatchTable(t *testing.T) {
	s := newTestService(t)
	rows, err := s.pool.Query(t.Context(), "SELECT id FROM expense_categories")
	if err != nil {
		t.Fatal(err)
	}
	inTable, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		t.Fatal(err)
	}
	inConfig := slices.Clone(g.ExpenseCategories)
	slices.Sort(inTable)
	slices.Sort(inConfig)
	if !slices.Equal(inTable, inConfig) {
		t.Fatalf("expense_categories %v, config %v", inTable, inConfig)
	}
}

func TestParseCents(t *testing.T) {
	for in, want := range map[string]int64{
		"12": 1200, "12.5": 1250, "12,50": 1250, " 0.01 ": 1, "10000000": MaxAmountCents,
	} {
		if got, err := ParseCents(in); err != nil || got != want {
			t.Errorf("ParseCents(%q) = %d, %v, want %d", in, got, err, want)
		}
	}
	for _, in := range []string{"", "0", "0,00", "-1", "1.234", "1e3", "1 000", "10000000.01", "abc"} {
		if _, err := ParseCents(in); err == nil {
			t.Errorf("ParseCents(%q) has no error", in)
		}
	}
	for _, cents := range []int64{1, 105, 1250, MaxAmountCents} {
		if got, _ := ParseCents(FormatCentsInput(cents)); got != cents {
			t.Errorf("ParseCents(FormatCentsInput(%d)) = %d", cents, got)
		}
	}
}

func TestSplitEqual(t *testing.T) {
	a, b, c := testUUID(1), testUUID(2), testUUID(3)
	for _, tc := range []struct {
		total  int64
		people []pgtype.UUID
		payer  pgtype.UUID
		want   []int64
	}{
		{1000, []pgtype.UUID{a, b, c}, b, []int64{333, 334, 333}},
		{1000, []pgtype.UUID{a, c}, b, []int64{500, 500}},
		{1001, []pgtype.UUID{a, c}, b, []int64{501, 500}},
		{2, []pgtype.UUID{a, b, c}, a, []int64{2, 0, 0}},
	} {
		got := splitEqual(tc.total, tc.people, tc.payer)
		if !slices.Equal(got, tc.want) {
			t.Errorf("splitEqual(%d) = %v, want %v", tc.total, got, tc.want)
		}
	}
}

func TestNetDebts(t *testing.T) {
	a, b, c := testUUID(1), testUUID(2), testUUID(3)
	row := func(from, to pgtype.UUID, cents int64) dbqueries.ListDebtItemsRow {
		return dbqueries.ListDebtItemsRow{FromUser: from, FromName: from.String(), ToUser: to, ToName: to.String(),
			ExpenseID: pgtype.UUID{}, Title: "", OnDate: pgtype.Date{}, AmountCents: cents}
	}
	got := netDebts([]dbqueries.ListDebtItemsRow{
		row(a, b, 500), row(b, a, 200), // a owes b 300
		row(c, a, 100), row(a, c, 100), // settled
		row(b, c, 20), row(c, b, 70), // c owes b 50, the pair is keyed b, c so it gets flipped
	})
	type summary struct {
		from, to pgtype.UUID
		cents    int64
		items    string
	}
	var sums []summary
	for _, d := range got {
		var items []string
		for _, it := range d.Items {
			items = append(items, strconv.FormatInt(it.Cents, 10))
		}
		sums = append(sums, summary{d.From, d.To, d.Cents, strings.Join(items, " ")})
	}
	want := []summary{
		{a, b, 300, "500 -200"},
		{c, b, 50, "-20 70"},
	}
	if !slices.Equal(sums, want) {
		t.Errorf("netDebts = %+v, want %+v", sums, want)
	}
}

func testUUID(b byte) pgtype.UUID {
	return pgtype.UUID{Bytes: [16]byte{15: b}, Valid: true}
}

func TestExpenses(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	equal := ExpenseInput{
		Title:      " Electricity ",
		Amount:     "10,01",
		Category:   "utilities",
		IncurredOn: "2026-10-01",
		Custom:     false,
		Members:    []pgtype.UUID{ownerM.UserID, memberM.UserID},
		Shares:     nil,
		ListID:     pgtype.UUID{},
		ItemID:     pgtype.UUID{},
	}

	bad := equal
	bad.Category = "fun"
	_, err := s.CreateExpense(ctx, ownerM, bad)
	wantErr(t, "unknown category", err, g.ErrorInvalidInput)
	bad = equal
	bad.Members = []pgtype.UUID{ownerM.UserID, strangerM.UserID}
	_, err = s.CreateExpense(ctx, ownerM, bad)
	wantErr(t, "stranger shares", err, g.ErrorInvalidInput)
	bad = equal
	bad.Custom, bad.Shares = true, map[pgtype.UUID]string{ownerM.UserID: "5", memberM.UserID: "5"}
	_, err = s.CreateExpense(ctx, ownerM, bad)
	wantErr(t, "custom shares do not add up", err, g.ErrorInvalidInput)

	expenseID, err := s.CreateExpense(ctx, ownerM, equal)
	if err != nil {
		t.Fatal(err)
	}
	e, err := s.GetExpense(ctx, memberM, expenseID)
	if err != nil {
		t.Fatal(err)
	}
	if e.Expense.Title != "Electricity" || e.Expense.AmountCents != 1001 || len(e.Shares) != 2 {
		t.Fatalf("expense = %+v", e)
	}
	if in := e.CopyInput(); in.Custom || len(in.Members) != 2 {
		t.Errorf("copy of an equal split = %+v", in)
	}
	_, err = s.GetExpense(ctx, strangerM, expenseID)
	wantErr(t, "stranger reads expense", err, g.ErrorNotFound)

	// the payer keeps the leftover cent, the member owes 5.00
	wantDebt := func(step string, cents int64) {
		t.Helper()
		debts, err := s.Debts(ctx, memberM)
		if err != nil {
			t.Fatal(err)
		}
		switch {
		case cents == 0 && len(debts) == 0:
		case len(debts) == 1 && debts[0].From == memberM.UserID && debts[0].To == ownerM.UserID && debts[0].Cents == cents:
		default:
			t.Fatalf("%s: debts = %+v, want member owes %d", step, debts, cents)
		}
	}
	wantDebt("after expense", 500)
	if debts, _ := s.Debts(ctx, memberM); len(debts[0].Items) != 1 || debts[0].Items[0].Row.Title != "Electricity" || debts[0].Items[0].Row.ExpenseID != expenseID {
		t.Fatalf("breakdown = %+v", debts[0].Items)
	}

	wantErr(t, "member leaves owing", s.Leave(ctx, memberM), g.ErrorBalanceNotZero)
	wantErr(t, "owner removes debtor", s.RemoveMember(ctx, ownerM, memberM.UserID), g.ErrorBalanceNotZero)

	// only the payer marks a share paid
	wantErr(t, "debtor marks own share", s.SetExpenseSharePaid(ctx, memberM, expenseID, memberM.UserID, true), g.ErrorNotFound)
	if err := s.SetExpenseSharePaid(ctx, ownerM, expenseID, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}
	wantDebt("share paid", 0)
	if err := s.SetExpenseSharePaid(ctx, ownerM, expenseID, memberM.UserID, false); err != nil {
		t.Fatal(err)
	}
	wantDebt("share unpaid again", 500)

	// a settlement counts once the receiver confirms it
	wantErr(t, "settle with self", s.CreateSettlement(ctx, memberM, memberM.UserID, "1"), g.ErrorInvalidInput)
	wantErr(t, "settle with stranger", s.CreateSettlement(ctx, memberM, strangerM.UserID, "1"), g.ErrorInvalidInput)
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "3"); err != nil {
		t.Fatal(err)
	}
	wantDebt("pending settlement", 500)
	pending, err := s.ListPendingSettlements(ctx, ownerM)
	if err != nil || len(pending) != 1 {
		t.Fatalf("pending = %+v, %v", pending, err)
	}
	wantErr(t, "sender confirms", s.ConfirmSettlement(ctx, memberM, pending[0].ID), g.ErrorNotFound)
	wantErr(t, "stranger rejects", s.DeleteSettlement(ctx, strangerM, pending[0].ID), g.ErrorNotFound)
	if err := s.ConfirmSettlement(ctx, ownerM, pending[0].ID); err != nil {
		t.Fatal(err)
	}
	wantDebt("confirmed settlement", 200)
	wantErr(t, "confirmed cannot be deleted", s.DeleteSettlement(ctx, memberM, pending[0].ID), g.ErrorNotFound)

	// overpaying flips the debt
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "2,50"); err != nil {
		t.Fatal(err)
	}
	pending, _ = s.ListPendingSettlements(ctx, ownerM)
	if err := s.ConfirmSettlement(ctx, ownerM, pending[0].ID); err != nil {
		t.Fatal(err)
	}
	debts, _ := s.Debts(ctx, ownerM)
	if len(debts) != 1 || debts[0].From != ownerM.UserID || debts[0].Cents != 50 {
		t.Fatalf("after overpaying: debts = %+v", debts)
	}
	wantErr(t, "member leaves while owed", s.Leave(ctx, memberM), g.ErrorBalanceNotZero)

	// a custom split where the payer pays nothing settles the overpayment
	custom := equal
	custom.Title, custom.Amount, custom.Custom = "Pizza", "0,50", true
	custom.Shares = map[pgtype.UUID]string{memberM.UserID: "0.50", ownerM.UserID: " "}
	pizzaID, err := s.CreateExpense(ctx, ownerM, custom)
	if err != nil {
		t.Fatal(err)
	}
	wantDebt("settled up", 0)

	// files attach to an expense and go with it
	onExpense := FileTarget{NoteID: pgtype.UUID{}, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: pizzaID, MessageID: pgtype.UUID{}}
	_, err = s.UploadFile(ctx, strangerM, onExpense, "bill.txt", strings.NewReader("bill"), 4)
	wantErr(t, "stranger attaches", err, g.ErrorNotFound)
	if _, err := s.UploadFile(ctx, memberM, onExpense, "bill.txt", strings.NewReader("bill"), 4); err != nil {
		t.Fatal(err)
	}
	pizza, err := s.GetExpense(ctx, ownerM, pizzaID)
	// one share is an equal split between one person
	if in := pizza.CopyInput(); err != nil || len(pizza.Files) != 1 || in.Custom || !slices.Equal(in.Members, []pgtype.UUID{memberM.UserID}) {
		t.Fatalf("pizza = %+v, %v", pizza, err)
	}
	wantErr(t, "member deletes payer's expense", s.DeleteExpense(ctx, memberM, pizzaID), g.ErrorNotFound)
	if err := s.DeleteExpense(ctx, ownerM, pizzaID); err != nil {
		t.Fatal(err)
	}
	files, _ := s.ListFiles(ctx, ownerM)
	if slices.ContainsFunc(files, func(f dbqueries.ListFilesRow) bool { return f.ExpenseID == pizzaID }) {
		t.Error("file of a deleted expense is still listed")
	}

	// back to the member being owed 0.50, paid back by the owner
	if err := s.CreateSettlement(ctx, ownerM, memberM.UserID, "0.5"); err != nil {
		t.Fatal(err)
	}
	pending, _ = s.ListPendingSettlements(ctx, memberM)
	if err := s.ConfirmSettlement(ctx, memberM, pending[0].ID); err != nil {
		t.Fatal(err)
	}
	wantDebt("all settled", 0)

	// a pending settlement would turn into a debt with a former member once confirmed
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "1"); err != nil {
		t.Fatal(err)
	}
	wantErr(t, "member leaves with a pending settlement", s.Leave(ctx, memberM), g.ErrorBalanceNotZero)
	pending, _ = s.ListPendingSettlements(ctx, memberM)
	if err := s.DeleteSettlement(ctx, memberM, pending[0].ID); err != nil {
		t.Fatal(err)
	}
	if err := s.Leave(ctx, memberM); err != nil {
		t.Fatal(err)
	}

	// the settled balance of a former member does not change anymore
	wantErr(t, "unpay former member's share", s.SetExpenseSharePaid(ctx, ownerM, expenseID, memberM.UserID, false), g.ErrorNotFound)
	wantErr(t, "delete expense former member shared", s.DeleteExpense(ctx, ownerM, expenseID), g.ErrorNotFound)
	wantDebt("after leaving", 0)
}
