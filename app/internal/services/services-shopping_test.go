package services

import (
	"context"
	g "roommates/internal/config"
	"strings"
	"testing"
)

func TestShoppingList(t *testing.T) {
	s := newTestService(t)
	ctx := context.Background()
	ownerM, memberM, strangerM := newTestHousehold(t, s)
	var err error

	_, err = s.CreateShoppingList(ctx, ownerM, "   ")
	wantErr(t, "empty list name", err, g.ErrorInvalidInput)
	listID, err := s.CreateShoppingList(ctx, ownerM, " Groceries ")
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
