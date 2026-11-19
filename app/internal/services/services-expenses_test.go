package services

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/utils"

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

// the form shows the error at the field clean names
func TestExpenseCleanField(t *testing.T) {
	a, b := testUUID(1), testUUID(2)
	member := func(id pgtype.UUID) dbqueries.ListActiveMembersRow {
		return dbqueries.ListActiveMembersRow{UserID: id, Username: "", AvatarKey: nil, Role: "", JoinedAt: pgtype.Timestamptz{}}
	}
	members := []dbqueries.ListActiveMembersRow{member(a), member(b)}
	valid := func() ExpenseInput {
		return ExpenseInput{Title: "Milk", Amount: "10", Category: g.ExpenseCategories[0], IncurredOn: "2026-11-19",
			Custom: false, Members: []pgtype.UUID{a, b}, Shares: map[pgtype.UUID]string{}, ListID: pgtype.UUID{}, ItemID: pgtype.UUID{}}
	}
	for _, tc := range []struct {
		name  string
		edit  func(*ExpenseInput)
		field string
	}{
		{"ok", func(*ExpenseInput) {}, ""},
		{"title", func(in *ExpenseInput) { in.Title = " " }, FieldTitle},
		{"amount", func(in *ExpenseInput) { in.Amount = "abc" }, FieldAmount},
		{"date", func(in *ExpenseInput) { in.IncurredOn = "" }, FieldDate},
		{"nobody ticked", func(in *ExpenseInput) { in.Members = nil }, FieldShares},
		{"custom sum", func(in *ExpenseInput) { in.Custom, in.Shares[a] = true, "3" }, FieldShares},
		{"custom amount", func(in *ExpenseInput) { in.Custom, in.Shares[a] = true, "x" }, FieldShares},
		{"category", func(in *ExpenseInput) { in.Category = "nope" }, "none"},
	} {
		in := valid()
		tc.edit(&in)
		_, _, err := in.clean(a, members)
		var fe FieldError
		switch {
		case tc.field == "" && err != nil:
			t.Errorf("%s: error %v", tc.name, err)
		case tc.field == "none" && (!errors.Is(err, g.ErrorInvalidInput) || errors.As(err, &fe)):
			t.Errorf("%s: error %v, want plain invalid input", tc.name, err)
		case tc.field != "" && tc.field != "none" && (!errors.As(err, &fe) || fe.Field != tc.field || !errors.Is(err, g.ErrorInvalidInput)):
			t.Errorf("%s: error %v, want field %s", tc.name, err, tc.field)
		}
	}
}

func TestNetDebts(t *testing.T) {
	a, b, c := testUUID(1), testUUID(2), testUUID(3)
	row := func(from, to pgtype.UUID, cents int64) dbqueries.ListDebtItemsRow {
		return dbqueries.ListDebtItemsRow{FromUser: from, FromName: from.String(), ToUser: to, ToName: to.String(),
			ExpenseID: pgtype.UUID{}, Title: "", OnDate: pgtype.Date{}, AmountCents: cents,
			Method: nil, Note: "", FromEdit: false, EditExpenseID: pgtype.UUID{}}
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
	if in := e.CopyInput([]pgtype.UUID{ownerM.UserID, memberM.UserID}); in.Custom || len(in.Members) != 2 {
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
	wantErr(t, "settle with self", s.CreateSettlement(ctx, memberM, memberM.UserID, "1", "cash", ""), g.ErrorInvalidInput)
	wantErr(t, "settle with stranger", s.CreateSettlement(ctx, memberM, strangerM.UserID, "1", "cash", ""), g.ErrorInvalidInput)
	var fe FieldError
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "1", "cheque", ""); !errors.As(err, &fe) || fe.Field != FieldMethod {
		t.Fatalf("unknown method: %v", err)
	}
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "1", "cash", strings.Repeat("x", MaxSettlementNoteLength+1)); !errors.As(err, &fe) || fe.Field != FieldNote {
		t.Fatalf("long note: %v", err)
	}
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "3", "transfer", " rent "); err != nil {
		t.Fatal(err)
	}
	wantDebt("pending settlement", 500)
	pending, err := s.ListPendingSettlements(ctx, ownerM)
	if err != nil || len(pending) != 1 {
		t.Fatalf("pending = %+v, %v", pending, err)
	}
	if p := pending[0]; p.Method == nil || *p.Method != "transfer" || p.Note != "rent" || p.Currency != ownerM.Currency {
		t.Fatalf("pending settlement = %+v", p)
	}
	wantErr(t, "sender confirms", s.ConfirmSettlement(ctx, memberM, pending[0].ID), g.ErrorNotFound)
	wantErr(t, "stranger rejects", s.DeleteSettlement(ctx, strangerM, pending[0].ID), g.ErrorNotFound)
	if err := s.ConfirmSettlement(ctx, ownerM, pending[0].ID); err != nil {
		t.Fatal(err)
	}
	wantDebt("confirmed settlement", 200)
	wantErr(t, "confirmed cannot be deleted", s.DeleteSettlement(ctx, memberM, pending[0].ID), g.ErrorNotFound)

	// overpaying flips the debt
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "2,50", "cash", ""); err != nil {
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
	if in := pizza.CopyInput([]pgtype.UUID{ownerM.UserID, memberM.UserID}); err != nil || len(pizza.Files) != 1 || in.Custom || !slices.Equal(in.Members, []pgtype.UUID{memberM.UserID}) {
		t.Fatalf("pizza = %+v, %v", pizza, err)
	}
	wantErr(t, "member deletes payer's expense", s.DeleteExpense(ctx, memberM, pizzaID), g.ErrorNotFound)
	if err := s.DeleteExpense(ctx, ownerM, pizzaID); err != nil {
		t.Fatal(err)
	}
	files, _ := s.ListFiles(ctx, ownerM, FileFilter{})
	if slices.ContainsFunc(files, func(f dbqueries.ListFilesRow) bool { return f.ExpenseID == pizzaID }) {
		t.Error("file of a deleted expense is still listed")
	}

	// files can stay when their expense goes: only what it alone holds is kept, and a refused delete keeps nothing
	keepingID, err := s.CreateExpense(ctx, ownerM, custom)
	if err != nil {
		t.Fatal(err)
	}
	onKeeping := FileTarget{NoteID: pgtype.UUID{}, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: keepingID, MessageID: pgtype.UUID{}}
	alone, err := s.UploadFile(ctx, memberM, onKeeping, "alone.txt", strings.NewReader("bill"), 4)
	if err != nil {
		t.Fatal(err)
	}
	noteID, err := s.CreateNote(ctx, ownerM, NoteInput{Title: "Bills", BodyMD: "", Private: false})
	if err != nil {
		t.Fatal(err)
	}
	onNote := FileTarget{NoteID: noteID, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}}
	both, err := s.UploadFile(ctx, ownerM, onNote, "both.txt", strings.NewReader("bill"), 4)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.AttachFile(ctx, ownerM, onKeeping, both); err != nil {
		t.Fatal(err)
	}
	wantErr(t, "member deletes keeping files", s.DeleteExpenseKeepingFiles(ctx, memberM, keepingID), g.ErrorNotFound)
	if f, err := s.GetFile(ctx, ownerM, alone); err != nil || f.Kept {
		t.Fatalf("after the refused delete %+v, %v", f, err)
	}
	if err := s.DeleteExpenseKeepingFiles(ctx, ownerM, keepingID); err != nil {
		t.Fatal(err)
	}
	if f, err := s.GetFile(ctx, ownerM, alone); err != nil || !f.Kept || f.RefCount != 0 {
		t.Fatalf("file only the expense held %+v, %v", f, err)
	}
	if f, err := s.GetFile(ctx, ownerM, both); err != nil || f.Kept || f.RefCount != 1 {
		t.Fatalf("file also on a note %+v, %v", f, err)
	}

	// back to the member being owed 0.50, paid back by the owner
	if err := s.CreateSettlement(ctx, ownerM, memberM.UserID, "0.5", "cash", ""); err != nil {
		t.Fatal(err)
	}
	pending, _ = s.ListPendingSettlements(ctx, memberM)
	if err := s.ConfirmSettlement(ctx, memberM, pending[0].ID); err != nil {
		t.Fatal(err)
	}
	wantDebt("all settled", 0)

	// a pending settlement would turn into a debt with a former member once confirmed
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "1", "cash", ""); err != nil {
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

// settlementHelpers are the steps of the settlement tests between the owner and the member
func settlementHelpers(t *testing.T, s *DefaultService, ownerM, memberM Membership) (
	expense func(payer, debtor Membership, amount string) pgtype.UUID,
	pay func(from, to Membership, amount string),
	wantDebt func(step string, cents int64),
	wantPaid func(step string, id pgtype.UUID, want bool),
) {
	t.Helper()
	ctx := context.Background()

	// payer paid an expense that is all on debtor
	expense = func(payer, debtor Membership, amount string) pgtype.UUID {
		t.Helper()
		id, err := s.CreateExpense(ctx, payer, ExpenseInput{
			Title:      "Bill",
			Amount:     amount,
			Category:   "other",
			IncurredOn: "2026-10-01",
			Custom:     true,
			Members:    nil,
			Shares:     map[pgtype.UUID]string{debtor.UserID: amount},
			ListID:     pgtype.UUID{},
			ItemID:     pgtype.UUID{},
		})
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
	pay = func(from, to Membership, amount string) {
		t.Helper()
		if err := s.CreateSettlement(ctx, from, to.UserID, amount, "transfer", ""); err != nil {
			t.Fatal(err)
		}
		pending, err := s.ListPendingSettlements(ctx, to)
		if err != nil || len(pending) != 1 {
			t.Fatalf("pending = %+v, %v", pending, err)
		}
		if err := s.ConfirmSettlement(ctx, to, pending[0].ID); err != nil {
			t.Fatal(err)
		}
	}
	// positive when the member owes the owner
	wantDebt = func(step string, cents int64) {
		t.Helper()
		debts, err := s.Debts(ctx, memberM)
		if err != nil {
			t.Fatal(err)
		}
		var got int64
		for _, d := range debts {
			got = d.Cents
			if d.From == ownerM.UserID {
				got = -d.Cents
			}
		}
		if len(debts) > 1 || got != cents {
			t.Fatalf("%s: debts = %+v, want %d", step, debts, cents)
		}
	}
	wantPaid = func(step string, id pgtype.UUID, want bool) {
		t.Helper()
		e, err := s.GetExpense(ctx, ownerM, id)
		if err != nil {
			t.Fatal(err)
		}
		for _, sh := range e.Shares {
			if sh.AmountCents > 0 && sh.PaidAt.Valid != want {
				t.Fatalf("%s: share = %+v, want paid %v", step, sh, want)
			}
		}
	}
	return expense, pay, wantDebt, wantPaid
}

// a payment can close a debt while one of its shares is being marked: the confirm waits
// for the mark and leaves that share alone, and the two never wait for each other
func TestConfirmWaitsForMarkInFlight(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)
	expense, pay, wantDebt, wantPaid := settlementHelpers(t, s, ownerM, memberM)

	older, newer := expense(ownerM, memberM, "15"), expense(ownerM, memberM, "15")
	pay(memberM, ownerM, "20")
	if err := s.CreateSettlement(ctx, memberM, ownerM.UserID, "25", "transfer", ""); err != nil {
		t.Fatal(err)
	}
	pending, err := s.ListPendingSettlements(ctx, ownerM)
	if err != nil || len(pending) != 1 {
		t.Fatalf("pending = %+v, %v", pending, err)
	}

	// the mark in flight, as far as SetExpenseSharePaid gets before it spends anything
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op after commit
	q := s.db.WithTx(tx)
	sh, err := q.GetShareToMark(ctx, dbqueries.GetShareToMarkParams{
		ExpenseID:   newer,
		UserID:      memberM.UserID,
		HouseholdID: ownerM.HouseholdID,
	})
	if err != nil {
		t.Fatal(err)
	}
	var markPID int32
	if err := tx.QueryRow(ctx, "SELECT pg_backend_pid()").Scan(&markPID); err != nil {
		t.Fatal(err)
	}

	confirmed := make(chan error, 1)
	go func() { confirmed <- s.ConfirmSettlement(ctx, ownerM, pending[0].ID) }()

	// the mark goes on only once the confirm is stuck behind it
	deadline := time.Now().Add(5 * time.Second)
	for blocked := false; !blocked; time.Sleep(10 * time.Millisecond) {
		if time.Now().After(deadline) {
			t.Fatal("the confirm never waited for the mark")
		}
		if err := s.pool.QueryRow(ctx,
			"SELECT EXISTS (SELECT 1 FROM pg_stat_activity WHERE $1::int = ANY(pg_blocking_pids(pid)))", markPID,
		).Scan(&blocked); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := payShare(ctx, q, ownerM.HouseholdID, newer, memberM.UserID, ownerM.UserID, sh.AmountCents); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	if err := <-confirmed; err != nil {
		t.Fatal(err)
	}

	wantPaid("the share marked in flight", newer, true)
	wantPaid("the share the confirm closed", older, true)
	// 45 paid against 30 of shares, a second spend on the marked share would leave 0
	wantDebt("after both", -1500)
}

// an expense waits for a member leaving in flight and is then refused, the leaver owes nothing
func TestExpenseWaitsForLeave(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)

	// the leave in flight, as far as Leave gets before it ends the membership
	tx, pid := lockHouseholdInTx(t, s, ownerM.HouseholdID)
	created := make(chan error, 1)
	go func() {
		_, err := s.CreateExpense(ctx, ownerM, ExpenseInput{
			Title:      "Rent",
			Amount:     "10",
			Category:   "other",
			IncurredOn: "2026-10-01",
			Custom:     false,
			Members:    []pgtype.UUID{ownerM.UserID, memberM.UserID},
			Shares:     nil,
			ListID:     pgtype.UUID{},
			ItemID:     pgtype.UUID{},
		})
		created <- err
	}()
	waitBlockedBy(t, s, pid)
	if err := s.db.WithTx(tx).EndMembership(ctx, dbqueries.EndMembershipParams{
		HouseholdID: memberM.HouseholdID,
		UserID:      memberM.UserID,
	}); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	wantErr(t, "expense on a member who left meanwhile", <-created, g.ErrorInvalidInput)
	if debts, err := s.Debts(ctx, ownerM); err != nil || len(debts) != 0 {
		t.Fatalf("debts = %+v, %v", debts, err)
	}
}

// the shares are listed by name, the form split the leftover cents in its member order
func TestCopyInputFormOrder(t *testing.T) {
	a, b, c := testUUID(1), testUUID(2), testUUID(3)
	share := func(id pgtype.UUID, cents int64) dbqueries.ListExpenseSharesRow {
		return dbqueries.ListExpenseSharesRow{UserID: id, Username: id.String(), AmountCents: cents, PaidAt: pgtype.Timestamptz{}, Active: true}
	}
	for _, tc := range []struct {
		order  []pgtype.UUID
		shares []dbqueries.ListExpenseSharesRow
		custom bool
	}{
		{[]pgtype.UUID{c, a}, []dbqueries.ListExpenseSharesRow{share(a, 500), share(c, 501)}, false},
		{[]pgtype.UUID{a, c}, []dbqueries.ListExpenseSharesRow{share(a, 501), share(c, 500)}, false},
		// the extra cent went to the second person, so it was typed by hand
		{[]pgtype.UUID{a, c}, []dbqueries.ListExpenseSharesRow{share(a, 500), share(c, 501)}, true},
		{[]pgtype.UUID{a, c}, []dbqueries.ListExpenseSharesRow{share(a, 300), share(c, 701)}, true},
	} {
		e := Expense{
			Expense: dbqueries.GetExpenseRow{ID: testUUID(9), Category: "other", Title: "Bill", AmountCents: 1001, Currency: "EUR",
				IncurredOn: pgtype.Date{}, PayerID: b, PayerName: "b"},
			Shares: tc.shares,
			Files:  nil,
		}
		if in := e.CopyInput(tc.order); in.Custom != tc.custom || len(in.Members) != 2 {
			t.Errorf("CopyInput of %+v = %+v, want custom %v", tc.shares, in, tc.custom)
		}
	}
}

func TestExpenseFilter(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)
	expense := func(title string, debtor Membership) pgtype.UUID {
		t.Helper()
		id, err := s.CreateExpense(ctx, ownerM, ExpenseInput{
			Title:      title,
			Amount:     "10",
			Category:   "other",
			IncurredOn: "2026-10-01",
			Custom:     true,
			Members:    nil,
			Shares:     map[pgtype.UUID]string{debtor.UserID: "10"},
			ListID:     pgtype.UUID{},
			ItemID:     pgtype.UUID{},
		})
		if err != nil {
			t.Fatal(err)
		}
		return id
	}
	want := func(what string, filter ExpenseFilter, titles ...string) {
		t.Helper()
		o, err := s.GetExpensesOverview(ctx, ownerM, filter)
		if err != nil {
			t.Fatal(err)
		}
		got := utils.Map(o.Expenses, func(e dbqueries.ListExpensesRow) string { return e.Title })
		slices.Sort(got)
		if !slices.Equal(got, titles) {
			t.Errorf("%s: got %v, want %v", what, got, titles)
		}
	}

	paid := expense("Rent March", memberM)
	expense("Rent April", memberM)
	// only the payer shares it, nobody owes anything
	expense("Own lunch", ownerM)
	if err := s.SetExpenseSharePaid(ctx, ownerM, paid, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}

	want("default", ExpenseFilter{Query: "", Show: ""}, "Rent April")
	want("unknown show is the default", ExpenseFilter{Query: "", Show: "x"}, "Rent April")
	want("search stays in unpaid", ExpenseFilter{Query: " rent ", Show: ""}, "Rent April")
	want("all", ExpenseFilter{Query: "", Show: ExpenseShowAll}, "Own lunch", "Rent April", "Rent March")
	want("search all in any case", ExpenseFilter{Query: "rENT m", Show: ExpenseShowAll}, "Rent March")
	want("no match", ExpenseFilter{Query: "water", Show: ExpenseShowAll})
}

// shares that net to zero let someone leave, then nobody can mark them paid
func TestExpenseFilterLeftMember(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, _ := newTestHousehold(t, s)
	expense, _, _, _ := settlementHelpers(t, s, ownerM, memberM)
	expense(ownerM, memberM, "10")
	expense(memberM, ownerM, "10")
	unpaid := ExpenseFilter{Query: "", Show: ""}
	if o, err := s.GetExpensesOverview(ctx, ownerM, unpaid); err != nil || len(o.Expenses) != 2 {
		t.Fatalf("before leaving = %+v, %v", o.Expenses, err)
	}

	if err := s.Leave(ctx, memberM); err != nil {
		t.Fatal(err)
	}
	if o, err := s.GetExpensesOverview(ctx, ownerM, unpaid); err != nil || len(o.Expenses) != 0 {
		t.Fatalf("after leaving = %+v, %v", o.Expenses, err)
	}
	if o, err := s.GetExpensesOverview(ctx, ownerM, ExpenseFilter{Query: "", Show: ExpenseShowAll}); err != nil || len(o.Expenses) != 2 {
		t.Fatalf("all after leaving = %+v, %v", o.Expenses, err)
	}
}

func TestSettlementCoversShares(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	expense, pay, wantDebt, wantPaid := settlementHelpers(t, s, ownerM, memberM)

	first, second := expense(ownerM, memberM, "15"), expense(ownerM, memberM, "15")
	pay(memberM, ownerM, "20")
	wantDebt("partial payment", 1000)
	wantPaid("partial payment leaves the shares", first, false)
	debts, _ := s.Debts(ctx, memberM)
	if d := debts[0]; d.Credit() != 2000 || !d.Covered(d.Items[0]) || !d.Covered(d.Items[1]) {
		t.Fatalf("debt = %+v, credit %d", d, d.Credit())
	}

	wantErr(t, "stranger marks", s.SetExpenseSharePaid(ctx, strangerM, first, memberM.UserID, true), g.ErrorNotFound)
	if err := s.SetExpenseSharePaid(ctx, memberM, first, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}
	wantDebt("debtor marks a covered share", 1000)
	wantPaid("debtor marks a covered share", first, true)
	if err := s.SetExpenseSharePaid(ctx, memberM, first, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}
	wantDebt("marking twice spends once", 1000)
	wantErr(t, "only 5 left of the payment", s.SetExpenseSharePaid(ctx, memberM, second, memberM.UserID, true), g.ErrorNotFound)
	wantPaid("the refused mark changes nothing", second, false)
	wantDebt("the refused mark changes nothing", 1000)
	wantErr(t, "debtor unmarks", s.SetExpenseSharePaid(ctx, memberM, first, memberM.UserID, false), g.ErrorNotFound)

	if err := s.SetExpenseSharePaid(ctx, ownerM, first, memberM.UserID, false); err != nil {
		t.Fatal(err)
	}
	wantDebt("unmarking gives the money back", 1000)
	if err := s.SetExpenseSharePaid(ctx, ownerM, first, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}
	wantDebt("the payer's mark spends the payment, the debt does not flip", 1000)

	pay(memberM, ownerM, "10")
	wantDebt("paid in full", 0)
	wantPaid("paid in full", first, true)
	wantPaid("paid in full", second, true)
	if o, err := s.GetExpensesOverview(ctx, ownerM, ExpenseFilter{Query: "", Show: ""}); err != nil || len(o.Expenses) != 0 || !o.HasExpenses {
		t.Fatalf("all paid, default overview = %+v, %v", o, err)
	}
	if o, err := s.GetExpensesOverview(ctx, ownerM, ExpenseFilter{Query: "", Show: ExpenseShowAll}); err != nil || len(o.Expenses) != 2 || !o.Expenses[0].FullyPaid {
		t.Fatalf("overview of all = %+v, %v", o, err)
	}

	// shares both ways are netted, paying the difference closes them all
	theirs, ours := expense(ownerM, memberM, "30"), expense(memberM, ownerM, "10")
	pay(memberM, ownerM, "20")
	wantDebt("netted", 0)
	wantPaid("netted, the debtor's share", theirs, true)
	e, _ := s.GetExpense(ctx, ownerM, ours)
	if e.Shares[0].AmountCents == 0 || !e.Shares[0].PaidAt.Valid {
		t.Fatalf("netted, the other share = %+v", e.Shares)
	}

	// overpaying marks the shares and leaves the rest as money owed back
	over := expense(ownerM, memberM, "5")
	pay(memberM, ownerM, "7")
	wantPaid("overpaid", over, true)
	wantDebt("overpaid", -200)

	// a share paid only in part by payments counts the rest as paid outside the app
	partly := expense(ownerM, memberM, "9")
	wantDebt("the overpayment lowers a new share", 700)
	if err := s.SetExpenseSharePaid(ctx, ownerM, partly, memberM.UserID, true); err != nil {
		t.Fatal(err)
	}
	wantDebt("the payer marks more than the payments", 0)

	// overpaying with shares both ways keeps the balance, the payment does not cover the 12
	bigger, smaller := expense(ownerM, memberM, "12"), expense(memberM, ownerM, "10")
	wantDebt("shares both ways", 200)
	pay(memberM, ownerM, "5")
	wantDebt("overpaid against shares both ways", -300)
	wantPaid("overpaid, the share it does not cover", bigger, false)
	wantPaid("overpaid, the other way", smaller, false)
}
