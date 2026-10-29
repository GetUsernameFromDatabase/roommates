package pages

import (
	"roommates/internal/services"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestItemClass(t *testing.T) {
	ben, you, other := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}, pgtype.UUID{Bytes: [16]byte{2}, Valid: true}, pgtype.UUID{Bytes: [16]byte{3}, Valid: true}
	actor := services.Membership{HouseholdID: pgtype.UUID{}, UserID: you, Role: "", Currency: ""}
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
