package services

import (
	"context"
	g "roommates/internal/config"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestShoppingList(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	var err error

	_, err = s.CreateShoppingList(ctx, ownerM, "   ", "groceries")
	wantErr(t, "empty list name", err, g.ErrorInvalidInput)
	listID, err := s.CreateShoppingList(ctx, ownerM, " Groceries ", "groceries")
	if err != nil {
		t.Fatal(err)
	}

	wantErr(t, "empty item", s.AddShoppingItem(ctx, memberM, listID, " ", ""), g.ErrorInvalidInput)
	wantErr(t, "long qty", s.AddShoppingItem(ctx, memberM, listID, "milk", strings.Repeat("x", MaxShoppingItemQtyLength+1)), g.ErrorInvalidInput)
	if err := s.AddShoppingItem(ctx, memberM, listID, "milk", "2 l"); err != nil {
		t.Fatal(err)
	}
	if err := s.AddShoppingItem(ctx, ownerM, listID, "bread", ""); err != nil {
		t.Fatal(err)
	}

	// another household sees nothing and changes nothing
	_, err = s.GetShoppingList(ctx, strangerM, listID)
	wantErr(t, "stranger reads list", err, g.ErrorNotFound)
	wantErr(t, "stranger adds item", s.AddShoppingItem(ctx, strangerM, listID, "x", ""), g.ErrorNotFound)
	wantErr(t, "stranger deletes list", s.DeleteShoppingList(ctx, strangerM, listID), g.ErrorNotFound)

	l, err := s.GetShoppingList(ctx, ownerM, listID)
	if err != nil {
		t.Fatal(err)
	}
	if l.List.Name != "Groceries" || len(l.Items) != 2 || l.Items[0].Text != "milk" || l.Items[0].Qty != "2 l" {
		t.Fatalf("unexpected list %+v", l)
	}
	milk, bread := l.Items[0].ID, l.Items[1].ID

	wantErr(t, "stranger checks item", s.SetShoppingItemChecked(ctx, strangerM, listID, milk, true), g.ErrorNotFound)
	wantErr(t, "stranger claims item", s.SetShoppingItemClaim(ctx, strangerM, listID, milk, true), g.ErrorNotFound)
	wantErr(t, "owner deletes member's item", s.DeleteShoppingItem(ctx, ownerM, listID, milk), g.ErrorNotFound)
	wantErr(t, "member deletes owner's list", s.DeleteShoppingList(ctx, memberM, listID), g.ErrorNotFound)

	// checking moves the item after the open ones, a second check from
	// a stale page keeps it checked by the first one
	if err := s.SetShoppingItemClaim(ctx, ownerM, listID, bread, true); err != nil {
		t.Fatal(err)
	}
	if err := s.SetShoppingItemChecked(ctx, memberM, listID, milk, true); err != nil {
		t.Fatal(err)
	}
	if err := s.SetShoppingItemChecked(ctx, ownerM, listID, milk, true); err != nil {
		t.Fatal(err)
	}
	l, _ = s.GetShoppingList(ctx, ownerM, listID)
	if l.Items[0].ID != bread || l.Items[0].ClaimedBy != ownerM.UserID ||
		!l.Items[1].CheckedAt.Valid || l.Items[1].CheckedByName == nil {
		t.Fatalf("unexpected items after claim and check %+v", l.Items)
	}

	// the member takes over the claim, the owner's stale "Not buying" keeps it theirs
	if err := s.SetShoppingItemClaim(ctx, memberM, listID, bread, true); err != nil {
		t.Fatal(err)
	}
	if err := s.SetShoppingItemClaim(ctx, ownerM, listID, bread, false); err != nil {
		t.Fatal(err)
	}
	l, _ = s.GetShoppingList(ctx, ownerM, listID)
	if l.Items[0].ClaimedBy != memberM.UserID {
		t.Fatalf("stale unclaim dropped someone else's claim %+v", l.Items[0])
	}
	if err := s.SetShoppingItemClaim(ctx, memberM, listID, bread, false); err != nil {
		t.Fatal(err)
	}
	l, _ = s.GetShoppingList(ctx, ownerM, listID)
	if l.Items[0].ClaimedBy.Valid {
		t.Fatal("unclaim should drop the own claim")
	}

	if err := s.ClearCheckedShoppingItems(ctx, ownerM, listID); err != nil {
		t.Fatal(err)
	}
	l, _ = s.GetShoppingList(ctx, ownerM, listID)
	if len(l.Items) != 1 || l.Items[0].ID != bread {
		t.Fatalf("clear should remove only checked items %+v", l.Items)
	}

	if err := s.DeleteShoppingItem(ctx, ownerM, listID, bread); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteShoppingList(ctx, ownerM, listID); err != nil {
		t.Fatal(err)
	}
	_, err = s.GetShoppingList(ctx, ownerM, listID)
	wantErr(t, "deleted list", err, g.ErrorNotFound)
}

// newBoughtItems makes a list with the items, all bought, and returns them in the order given
func newBoughtItems(t *testing.T, s *DefaultService, m Membership, category string, texts ...string) (pgtype.UUID, []pgtype.UUID) {
	t.Helper()
	ctx := t.Context()
	listID, err := s.CreateShoppingList(ctx, m, "Shop", category)
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range texts {
		if err := s.AddShoppingItem(ctx, m, listID, text, ""); err != nil {
			t.Fatal(err)
		}
	}
	l, err := s.GetShoppingList(ctx, m, listID)
	if err != nil {
		t.Fatal(err)
	}
	ids := make([]pgtype.UUID, 0, len(l.Items))
	for _, it := range l.Items {
		ids = append(ids, it.ID)
		if err := s.SetShoppingItemChecked(ctx, m, listID, it.ID, true); err != nil {
			t.Fatal(err)
		}
	}
	return listID, ids
}

// expenseIDs maps the list's visible items to their expense
func expenseIDs(t *testing.T, s *DefaultService, m Membership, listID pgtype.UUID) map[pgtype.UUID]pgtype.UUID {
	t.Helper()
	l, err := s.GetShoppingList(t.Context(), m, listID)
	if err != nil {
		t.Fatal(err)
	}
	ids := map[pgtype.UUID]pgtype.UUID{}
	for _, it := range l.Items {
		ids[it.ID] = it.ExpenseID
	}
	return ids
}

func TestShoppingRename(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	listID, _ := newBoughtItems(t, s, ownerM, "groceries")

	wantErr(t, "empty name", s.RenameShoppingList(ctx, memberM, listID, "  "), g.ErrorInvalidInput)
	long := strings.Repeat("x", MaxShoppingListNameLength+1)
	wantErr(t, "long name", s.RenameShoppingList(ctx, memberM, listID, long), g.ErrorInvalidInput)
	wantErr(t, "stranger renames", s.RenameShoppingList(ctx, strangerM, listID, "x"), g.ErrorNotFound)
	if err := s.RenameShoppingList(ctx, memberM, listID, " Rimi "); err != nil {
		t.Fatal(err)
	}
	l, err := s.GetShoppingList(ctx, ownerM, listID)
	if err != nil || l.List.Name != "Rimi" {
		t.Fatalf("after rename %q %v", l.List.Name, err)
	}
}

func TestShoppingCategory(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, strangerM := newTestHousehold(t, s)

	_, err := s.CreateShoppingList(ctx, ownerM, "Tools", "fun")
	wantErr(t, "unknown category", err, g.ErrorInvalidInput)
	listID, ids := newBoughtItems(t, s, ownerM, "household", "Nails")

	title, category, err := s.ShoppingExpenseDraft(ctx, ownerM, listID, ids[0])
	if err != nil || title != "Nails" || category != "household" {
		t.Fatalf("item draft %q %q %v", title, category, err)
	}
	title, _, err = s.ShoppingExpenseDraft(ctx, ownerM, listID, pgtype.UUID{})
	if err != nil || title != "Shop" {
		t.Fatalf("bag draft %q %v", title, err)
	}
	_, _, err = s.ShoppingExpenseDraft(ctx, strangerM, listID, ids[0])
	wantErr(t, "stranger drafts", err, g.ErrorNotFound)

	// an item can have a longer name than an expense may have a title
	longList, longIDs := newBoughtItems(t, s, ownerM, "groceries", strings.Repeat("õ", MaxExpenseTitleLength+1))
	title, _, err = s.ShoppingExpenseDraft(ctx, ownerM, longList, longIDs[0])
	if err != nil || title != strings.Repeat("õ", MaxExpenseTitleLength) {
		t.Fatalf("long item draft %q %v", title, err)
	}

	wantErr(t, "unknown category", s.SetShoppingListCategory(ctx, memberM, listID, "fun"), g.ErrorInvalidInput)
	wantErr(t, "stranger recategorises", s.SetShoppingListCategory(ctx, strangerM, listID, "rent"), g.ErrorNotFound)
	if err := s.SetShoppingListCategory(ctx, memberM, listID, "other"); err != nil {
		t.Fatal(err)
	}
	if _, category, _ = s.ShoppingExpenseDraft(ctx, ownerM, listID, ids[0]); category != "other" {
		t.Fatalf("category after change %q", category)
	}
}

func TestShoppingExpenseLink(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, _, strangerM := newTestHousehold(t, s)
	input := func(m Membership, listID, itemID pgtype.UUID) ExpenseInput {
		return ExpenseInput{
			Title:      "Shop",
			Amount:     "10",
			Category:   "groceries",
			IncurredOn: "2026-10-01",
			Custom:     false,
			Members:    []pgtype.UUID{m.UserID},
			Shares:     nil,
			ListID:     listID,
			ItemID:     itemID,
		}
	}

	// same name, different items: the link goes by id, not by name
	listID, ids := newBoughtItems(t, s, ownerM, "groceries", "Milk", "Milk", "Bread")
	first, second, bread := ids[0], ids[1], ids[2]
	one, err := s.CreateExpense(ctx, ownerM, input(ownerM, listID, first))
	if err != nil {
		t.Fatal(err)
	}
	linked := expenseIDs(t, s, ownerM, listID)
	if linked[first] != one || linked[second].Valid || linked[bread].Valid {
		t.Fatalf("one item link %+v", linked)
	}

	// another household's expense links nothing here, its expense is still made
	if _, err := s.CreateExpense(ctx, strangerM, input(strangerM, listID, pgtype.UUID{})); err != nil {
		t.Fatal(err)
	}
	if linked = expenseIDs(t, s, ownerM, listID); linked[second].Valid {
		t.Fatalf("stranger linked %+v", linked)
	}

	// the whole bag takes only what has no expense yet, an open item stays out
	if err := s.AddShoppingItem(ctx, ownerM, listID, "Eggs", ""); err != nil {
		t.Fatal(err)
	}
	bag, err := s.CreateExpense(ctx, ownerM, input(ownerM, listID, pgtype.UUID{}))
	if err != nil {
		t.Fatal(err)
	}
	linked = expenseIDs(t, s, ownerM, listID)
	if linked[first] != one || linked[second] != bag || linked[bread] != bag {
		t.Fatalf("bag link %+v", linked)
	}
	for id, expenseID := range linked {
		if id != first && id != second && id != bread && expenseID.Valid {
			t.Fatalf("open item linked %+v", linked)
		}
	}

	// deleting the expense only drops the link
	if err := s.DeleteExpense(ctx, ownerM, bag); err != nil {
		t.Fatal(err)
	}
	if linked = expenseIDs(t, s, ownerM, listID); linked[first] != one || linked[bread].Valid || len(linked) != 4 {
		t.Fatalf("after deleting the expense %+v", linked)
	}

	// bought again after an untick is a new purchase, without the old receipt
	if err := s.SetShoppingItemChecked(ctx, ownerM, listID, first, false); err != nil {
		t.Fatal(err)
	}
	if err := s.SetShoppingItemChecked(ctx, ownerM, listID, first, true); err != nil {
		t.Fatal(err)
	}
	if linked = expenseIDs(t, s, ownerM, listID); linked[first].Valid {
		t.Fatalf("after unticking %+v", linked)
	}
}

// a stale page can untick an item someone removed meanwhile
func TestShoppingHiddenItemNotOpen(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, _, _ := newTestHousehold(t, s)
	listID, ids := newBoughtItems(t, s, ownerM, "groceries", "Milk")
	if err := s.ClearCheckedShoppingItems(ctx, ownerM, listID); err != nil {
		t.Fatal(err)
	}
	if err := s.SetShoppingItemChecked(ctx, ownerM, listID, ids[0], false); err != nil {
		t.Fatal(err)
	}
	o, err := s.GetHouseholdOverview(ctx, ownerM)
	if err != nil {
		t.Fatal(err)
	}
	if len(o.ShoppingLists) != 1 || o.ShoppingLists[0].OpenItems != 0 {
		t.Fatalf("hidden item counted %+v", o.ShoppingLists)
	}
}

func TestShoppingBuyAgain(t *testing.T) {
	s := newTestService(t)
	ctx := t.Context()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	listID, ids := newBoughtItems(t, s, ownerM, "groceries", "Milk", "Typo")
	milk, typo := ids[0], ids[1]
	otherList, _ := newBoughtItems(t, s, ownerM, "groceries")

	long := strings.Repeat("x", MaxShoppingItemDescriptionLength+1)
	wantErr(t, "long description", s.SetShoppingItemDescription(ctx, memberM, listID, milk, long), g.ErrorInvalidInput)
	wantErr(t, "stranger describes", s.SetShoppingItemDescription(ctx, strangerM, listID, milk, "x"), g.ErrorNotFound)
	// a form's line breaks are two characters, maxlength counted them as one
	full := strings.Repeat("x", MaxShoppingItemDescriptionLength-2) + "\r\n" + "x"
	if err := s.SetShoppingItemDescription(ctx, memberM, listID, milk, full); err != nil {
		t.Fatal(err)
	}
	if err := s.SetShoppingItemDescription(ctx, memberM, listID, milk, " the blue one "); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CreateExpense(ctx, ownerM, ExpenseInput{
		Title: "Milk", Amount: "2", Category: "groceries", IncurredOn: "2026-10-01", Custom: false,
		Members: []pgtype.UUID{ownerM.UserID}, Shares: nil, ListID: listID, ItemID: milk,
	}); err != nil {
		t.Fatal(err)
	}

	// a deleted item is gone for good, a removed one waits to be bought again
	if err := s.DeleteShoppingItem(ctx, ownerM, listID, typo); err != nil {
		t.Fatal(err)
	}
	if err := s.ClearCheckedShoppingItems(ctx, memberM, listID); err != nil {
		t.Fatal(err)
	}
	l, err := s.GetShoppingList(ctx, ownerM, otherList)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.BuyAgain) != 1 || l.BuyAgain[0].Text != "Milk" {
		t.Fatalf("buy again %+v", l.BuyAgain)
	}
	strangerList, _ := newBoughtItems(t, s, strangerM, "groceries")
	if l, _ = s.GetShoppingList(ctx, strangerM, strangerList); len(l.BuyAgain) != 0 {
		t.Fatalf("another household's items %+v", l.BuyAgain)
	}

	// any spelling brings the same item back on this list as a new purchase
	if err := s.AddShoppingItem(ctx, memberM, otherList, "MILK", "3 l"); err != nil {
		t.Fatal(err)
	}
	l, err = s.GetShoppingList(ctx, ownerM, otherList)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Items) != 1 || len(l.BuyAgain) != 0 {
		t.Fatalf("after buying again %+v", l)
	}
	it := l.Items[0]
	if it.ID != milk || it.Text != "Milk" || it.Qty != "3 l" || it.Description != "the blue one" ||
		it.CheckedAt.Valid || it.ClaimedBy.Valid || it.ExpenseID.Valid || it.CreatedBy != memberM.UserID {
		t.Fatalf("revived item %+v", it)
	}
	if l, _ = s.GetShoppingList(ctx, ownerM, listID); len(l.Items) != 0 {
		t.Fatalf("the old list still has it %+v", l.Items)
	}
}
