package pages

import (
	"testing"

	"roommates/internal/db/dbqueries"
)

func TestSplitBought(t *testing.T) {
	var open, bought dbqueries.ListShoppingItemsRow
	bought.CheckedAt.Valid = true
	cases := []struct {
		name             string
		items            []dbqueries.ListShoppingItemsRow
		wantOpen, wantIn int
	}{
		{"empty", nil, 0, 0},
		{"all open", []dbqueries.ListShoppingItemsRow{open, open}, 2, 0},
		{"all bought", []dbqueries.ListShoppingItemsRow{bought}, 0, 1},
		{"both", []dbqueries.ListShoppingItemsRow{open, open, bought}, 2, 1},
	}
	for _, c := range cases {
		o, b := splitBought(c.items)
		if len(o) != c.wantOpen || len(b) != c.wantIn {
			t.Errorf("%s: got %d open and %d bought, want %d and %d", c.name, len(o), len(b), c.wantOpen, c.wantIn)
		}
	}
}
