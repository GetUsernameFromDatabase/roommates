package services

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"

	"github.com/invopop/ctxi18n"
	"github.com/jackc/pgx/v5/pgtype"
)

func localeCtx(t *testing.T, lang string) context.Context {
	t.Helper()
	locales.LoadLocales()
	ctx, err := ctxi18n.WithLocale(t.Context(), lang)
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}

func uid(b byte) pgtype.UUID { return pgtype.UUID{Bytes: [16]byte{b}, Valid: true} }

func TestCSVAmountAndText(t *testing.T) {
	for lang, want := range map[string][2]string{"et": {"12,50", "-0,05"}, "en": {"12.50", "-0.05"}} {
		ctx := localeCtx(t, lang)
		if got := csvAmount(ctx, 1250); got != want[0] {
			t.Errorf("%s: 1250 = %q", lang, got)
		}
		if got := csvAmount(ctx, -5); got != want[1] {
			t.Errorf("%s: -5 = %q", lang, got)
		}
	}
	for in, want := range map[string]string{"=SUM(A1)": "'=SUM(A1)", "+1": "'+1", "-1": "'-1", "@x": "'@x", "\tx": "'\tx", "milk": "milk", "": ""} {
		if got := csvText(in); got != want {
			t.Errorf("csvText(%q) = %q, want %q", in, got, want)
		}
	}
}

func exportFixture() ([]dbqueries.ListExportExpensesRow, []dbqueries.ListExportSharesRow) {
	day := pgtype.Date{Time: time.Date(2026, 11, 19, 0, 0, 0, 0, time.UTC), InfinityModifier: pgtype.Finite, Valid: true}
	e1, e2 := uid(101), uid(102)
	expenses := []dbqueries.ListExportExpensesRow{
		{ID: e1, Category: "groceries", Title: "=HYPERLINK(1)", AmountCents: 1250, Currency: "EUR", IncurredOn: day, PayerID: uid(1), PayerName: "Ann"},
		{ID: e2, Category: "rent", Title: "Rent", AmountCents: 10000, Currency: "EUR", IncurredOn: day, PayerID: uid(2), PayerName: "Bob"},
	}
	shares := []dbqueries.ListExportSharesRow{
		{ExpenseID: e1, UserID: uid(1), Username: "Ann", AmountCents: 650, Paid: true},
		{ExpenseID: e1, UserID: uid(2), Username: "Bob", AmountCents: 600, Paid: false},
		{ExpenseID: e2, UserID: uid(2), Username: "Bob", AmountCents: 10000, Paid: true},
	}
	return expenses, shares
}

func TestExpenseRows(t *testing.T) {
	expenses, shares := exportFixture()
	ctx := localeCtx(t, "et")

	wide := expenseRows(ctx, ExportWide, expenses, shares)
	if len(wide) != 3 || len(wide[0]) != 8 {
		t.Fatalf("wide has %d rows, header %v", len(wide), wide[0])
	}
	if wide[0][6] != "Ann" || wide[0][7] != "Bob" {
		t.Errorf("member columns %v", wide[0][6:])
	}
	if got := wide[1][0] + "|" + wide[1][1] + "|" + strings.Join(wide[1][4:], "|"); got != "2026-11-19|'=HYPERLINK(1)|12,50|EUR|6,50|6,00" {
		t.Errorf("wide row 1 = %q", got)
	}
	if wide[2][6] != "" || wide[2][7] != "100,00" {
		t.Errorf("wide row 2 = %v", wide[2])
	}

	long := expenseRows(localeCtx(t, "en"), ExportLong, expenses, shares)
	if len(long) != 4 {
		t.Fatalf("long has %d rows", len(long))
	}
	if long[2][6] != "Bob" || long[2][7] != "6.00" || long[1][4] != "12.50" {
		t.Errorf("long rows %v %v", long[1], long[2])
	}
	if long[1][8] == long[2][8] {
		t.Errorf("paid cells equal: %v %v", long[1], long[2])
	}

	out, err := writeCSV(ctx, wide)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(out, utf8BOM) || !strings.Contains(string(out), ";12,50;EUR;6,50;6,00\n") {
		t.Errorf("et csv: %q", out)
	}
	en, err := writeCSV(localeCtx(t, "en"), long)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(en), ",12.50,EUR,Ann,6.50,") {
		t.Errorf("en csv: %q", en)
	}
}

func TestBalanceRows(t *testing.T) {
	ctx := localeCtx(t, "en")
	ann, bob, cy := uid(1), uid(2), uid(3)
	// bob owes ann 5.00 and 3.00, ann owes cy 2.00
	debt := func(from pgtype.UUID, fromName string, to pgtype.UUID, toName string, cents int64) dbqueries.ListDebtItemsRow {
		return dbqueries.ListDebtItemsRow{FromUser: from, FromName: fromName, ToUser: to, ToName: toName, ExpenseID: pgtype.UUID{}, Title: "", OnDate: pgtype.Date{}, AmountCents: cents,
			Method: nil, Note: "", FromEdit: false, EditExpenseID: pgtype.UUID{}}
	}
	items := []dbqueries.ListDebtItemsRow{
		debt(bob, "Bob", ann, "Ann", 500),
		debt(bob, "Bob", ann, "Ann", 300),
		debt(ann, "Ann", cy, "Cy", 200),
	}
	person := func(id pgtype.UUID, name string) dbqueries.ListActiveMembersRow {
		return dbqueries.ListActiveMembersRow{UserID: id, Username: name, AvatarKey: nil, Role: g.RoleMember, JoinedAt: pgtype.Timestamptz{}}
	}
	members := []dbqueries.ListActiveMembersRow{person(ann, "Ann"), person(bob, "Bob"), person(cy, "Cy")}

	owner := Membership{HouseholdID: uid(9), UserID: ann, Role: g.RoleOwner, Currency: "EUR", UnreadMessages: 0, PendingSettlements: 0, DeletedAt: pgtype.Timestamptz{}}
	rows := balanceRows(ctx, owner, items, members)
	want := [][2]string{{"Ann", "6.00"}, {"Bob", "-8.00"}, {"Cy", "2.00"}}
	if len(rows) != 4 {
		t.Fatalf("owner rows %v", rows)
	}
	for i, w := range want {
		if rows[i+1][0] != w[0] || rows[i+1][1] != w[1] || rows[i+1][2] != "EUR" {
			t.Errorf("owner row %d = %v, want %v", i, rows[i+1], w)
		}
	}

	// the member only sees their own pairs: bob owes ann, ann owes cy
	member := Membership{HouseholdID: uid(9), UserID: ann, Role: g.RoleMember, Currency: "EUR", UnreadMessages: 0, PendingSettlements: 0, DeletedAt: pgtype.Timestamptz{}}
	rows = balanceRows(ctx, member, items, nil)
	if len(rows) != 3 || rows[1][0] != "Bob" || rows[1][1] != "8.00" || rows[2][0] != "Cy" || rows[2][1] != "-2.00" {
		t.Errorf("member rows %v", rows)
	}
}

func TestExportScope(t *testing.T) {
	s := newTestService(t)
	ctx := localeCtx(t, "en")
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	in := func(title string, members ...pgtype.UUID) ExpenseInput {
		return ExpenseInput{Title: title, Amount: "10", Category: "other", IncurredOn: "2026-11-19",
			Custom: false, Members: members, Shares: nil, ListID: pgtype.UUID{}, ItemID: pgtype.UUID{}}
	}
	// the owner pays alone, and for both
	if _, err := s.CreateExpense(ctx, ownerM, in("Solo", ownerM.UserID)); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CreateExpense(ctx, ownerM, in("Shared", ownerM.UserID, memberM.UserID)); err != nil {
		t.Fatal(err)
	}
	count := func(m Membership, format string) int {
		out, err := s.ExportExpenses(ctx, m, format)
		if err != nil {
			t.Fatal(err)
		}
		return strings.Count(string(out), "\n") - 1
	}
	if n := count(ownerM, ExportWide); n != 2 {
		t.Errorf("owner wide has %d rows", n)
	}
	if n := count(memberM, ExportWide); n != 1 {
		t.Errorf("member wide has %d rows", n)
	}
	// long: one row per share
	if n := count(ownerM, ExportLong); n != 3 {
		t.Errorf("owner long has %d rows", n)
	}
	if n := count(strangerM, ExportWide); n != 0 {
		t.Errorf("stranger wide has %d rows", n)
	}
	if _, err := s.ExportExpenses(ctx, ownerM, "x"); err == nil {
		t.Error("unknown format")
	}
	if _, err := s.ExportBalances(ctx, memberM); err != nil {
		t.Fatal(err)
	}
}
