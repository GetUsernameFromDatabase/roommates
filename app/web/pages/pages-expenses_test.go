package pages

import (
	"context"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"
	"strings"
	"testing"

	"github.com/invopop/ctxi18n"
	"github.com/jackc/pgx/v5/pgtype"
)

// the hidden fields are what ties a new expense to the shopping items it paid for
func TestExpenseEditorCarriesTheItemLink(t *testing.T) {
	locales.LoadLocales()
	ctx, err := ctxi18n.WithLocale(context.Background(), string(locales.Default))
	if err != nil {
		t.Fatal(err)
	}
	list, item := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}, pgtype.UUID{Bytes: [16]byte{2}, Valid: true}
	actor := services.Membership{HouseholdID: pgtype.UUID{}, UserID: pgtype.UUID{}, Role: "", Currency: "", UnreadMessages: 0}
	in := services.ExpenseInput{Title: "", Amount: "", Category: "groceries", IncurredOn: "", Custom: false,
		Members: nil, Shares: nil, ListID: list, ItemID: item}
	var b strings.Builder
	if err := ExpenseEditor(actor, nil, in, "").Render(ctx, &b); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`name="list" value="` + list.String() + `"`,
		`name="item" value="` + item.String() + `"`,
	} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("the form does not have %s", want)
		}
	}
}

// the stamp tells everyone nothing is owed, so it must not show over a debt or a payment
func TestHouseholdExpensesSettledStamp(t *testing.T) {
	locales.LoadLocales()
	ctx, err := ctxi18n.WithLocale(context.Background(), string(locales.Default))
	if err != nil {
		t.Fatal(err)
	}
	you, ben := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}, pgtype.UUID{Bytes: [16]byte{2}, Valid: true}
	actor := services.Membership{HouseholdID: pgtype.UUID{}, UserID: you, Role: "", Currency: "EUR", UnreadMessages: 0}
	expenses := []dbqueries.ListExpensesRow{{ID: pgtype.UUID{}, Category: "other", Title: "Bread", AmountCents: 1,
		IncurredOn: pgtype.Date{}, PayerID: ben, PayerName: "Ben"}}
	debts := []services.Debt{{From: you, FromName: "", To: ben, ToName: "Ben", Cents: 1, Items: nil}}
	pending := []dbqueries.ListPendingSettlementsRow{{ID: pgtype.UUID{}, FromUser: you, FromName: "", ToUser: ben, ToName: "Ben",
		AmountCents: 1, CreatedAt: pgtype.Timestamptz{}}}
	for _, tc := range []struct {
		name string
		o    services.ExpensesOverview
		want bool
	}{
		{"all paid", services.ExpensesOverview{Debts: nil, Pending: nil, Expenses: expenses}, true},
		{"a debt", services.ExpensesOverview{Debts: debts, Pending: nil, Expenses: expenses}, false},
		{"a pending payment", services.ExpensesOverview{Debts: nil, Pending: pending, Expenses: expenses}, false},
		{"nothing yet", services.ExpensesOverview{Debts: nil, Pending: nil, Expenses: nil}, false},
	} {
		var b strings.Builder
		if err := HouseholdExpenses(actor, tc.o, "").Render(ctx, &b); err != nil {
			t.Fatal(err)
		}
		if got := strings.Contains(b.String(), utils.T(ctx, locales.LK_ExpensesSettled)); got != tc.want {
			t.Errorf("%s: stamp shown %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestYoursFirst(t *testing.T) {
	you, ben, cleo, dan := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}, pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
		pgtype.UUID{Bytes: [16]byte{3}, Valid: true}, pgtype.UUID{Bytes: [16]byte{4}, Valid: true}
	actor := services.Membership{HouseholdID: pgtype.UUID{}, UserID: you, Role: "", Currency: "", UnreadMessages: 0}
	debt := func(from, to pgtype.UUID, cents int64) services.Debt {
		return services.Debt{From: from, FromName: "", To: to, ToName: "", Cents: cents, Items: nil}
	}
	debts := []services.Debt{debt(ben, cleo, 1), debt(you, ben, 2), debt(cleo, dan, 3), debt(dan, you, 4)}
	got := yoursFirst(debts, actor)
	for i, want := range []int64{2, 4, 1, 3} {
		if got[i].Cents != want {
			t.Errorf("place %d: got the debt of %d, want %d", i, got[i].Cents, want)
		}
	}
	if debts[0].Cents != 1 {
		t.Error("the overview's slice was reordered")
	}
}

func TestItemClass(t *testing.T) {
	ben, you, other := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}, pgtype.UUID{Bytes: [16]byte{2}, Valid: true}, pgtype.UUID{Bytes: [16]byte{3}, Valid: true}
	actor := services.Membership{HouseholdID: pgtype.UUID{}, UserID: you, Role: "", Currency: "", UnreadMessages: 0}
	youOwe := services.Debt{From: you, FromName: "", To: ben, ToName: "", Cents: 1, Items: nil}
	owesYou := services.Debt{From: ben, FromName: "", To: you, ToName: "", Cents: 1, Items: nil}
	others := services.Debt{From: ben, FromName: "", To: other, ToName: "", Cents: 1, Items: nil}
	for _, tc := range []struct {
		name  string
		d     services.Debt
		cents int64
		want  string
	}{
		{"you owe more", youOwe, 5, "text-against"},
		{"you paid some off", youOwe, -5, "text-favour"},
		{"owed more to you", owesYou, 5, "text-favour"},
		{"owed less to you", owesYou, -5, "text-against"},
		{"not yours", others, 5, ""},
	} {
		if got := itemClass(actor, tc.d, tc.cents); got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}
