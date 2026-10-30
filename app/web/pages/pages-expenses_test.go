package pages

import (
	"context"
	"roommates/internal/locales"
	"roommates/internal/services"
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
