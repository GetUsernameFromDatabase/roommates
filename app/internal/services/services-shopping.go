package services

import (
	"context"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	MaxShoppingListNameLength = 64
	MaxShoppingItemTextLength = 200
	MaxShoppingItemQtyLength  = 32
)

// trimToLength trims s and reports whether it has min..max runes
func trimToLength(s string, minLength, maxLength int) (string, bool) {
	s = strings.TrimSpace(s)
	n := utf8.RuneCountInString(s)
	return s, n >= minLength && n <= maxLength
}

// rowsAffected turns 0 affected rows into [g.ErrorNotFound]
func rowsAffected(n int64, err error) error {
	if err == nil && n == 0 {
		return g.ErrorNotFound
	}
	return err
}

// CreateShoppingList returns [g.ErrorInvalidInput] unless the name has 1..[MaxShoppingListNameLength] runes
func (s *DefaultService) CreateShoppingList(ctx context.Context, actor Membership, name string) (pgtype.UUID, error) {
	name, ok := trimToLength(name, 1, MaxShoppingListNameLength)
	if !ok {
		return pgtype.UUID{}, g.ErrorInvalidInput
	}
	return s.db.InsertShoppingList(ctx, dbqueries.InsertShoppingListParams{
		HouseholdID: actor.HouseholdID,
		Name:        name,
		CreatedBy:   actor.UserID,
	})
}

type ShoppingList struct {
	List  dbqueries.GetShoppingListRow
	Items []dbqueries.ListShoppingItemsRow
}

// GetShoppingList returns [g.ErrorNotFound] if the list is not in the actor's household
func (s *DefaultService) GetShoppingList(ctx context.Context, actor Membership, listID pgtype.UUID) (ShoppingList, error) {
	var l ShoppingList
	var err error
	l.List, err = s.db.GetShoppingList(ctx, dbqueries.GetShoppingListParams{
		ID:          listID,
		HouseholdID: actor.HouseholdID,
	})
	if err != nil {
		return l, noRowsAs(err, g.ErrorNotFound)
	}
	l.Items, err = s.db.ListShoppingItems(ctx, dbqueries.ListShoppingItemsParams{
		ListID:      listID,
		HouseholdID: actor.HouseholdID,
	})
	return l, err
}

// DeleteShoppingList returns [g.ErrorNotFound] unless the actor created the list
func (s *DefaultService) DeleteShoppingList(ctx context.Context, actor Membership, listID pgtype.UUID) error {
	return rowsAffected(s.db.DeleteShoppingList(ctx, dbqueries.DeleteShoppingListParams{
		ID:          listID,
		HouseholdID: actor.HouseholdID,
		CreatedBy:   actor.UserID,
	}))
}

// -----------------------------------------------------------------------------

// AddShoppingItem returns [g.ErrorInvalidInput] for an empty or too long text or qty
func (s *DefaultService) AddShoppingItem(ctx context.Context, actor Membership, listID pgtype.UUID, text, qty string) error {
	text, okText := trimToLength(text, 1, MaxShoppingItemTextLength)
	qty, okQty := trimToLength(qty, 0, MaxShoppingItemQtyLength)
	if !okText || !okQty {
		return g.ErrorInvalidInput
	}
	if _, err := s.db.GetShoppingList(ctx, dbqueries.GetShoppingListParams{
		ID:          listID,
		HouseholdID: actor.HouseholdID,
	}); err != nil {
		return noRowsAs(err, g.ErrorNotFound)
	}
	return s.db.InsertShoppingItem(ctx, dbqueries.InsertShoppingItemParams{
		ListID:      listID,
		HouseholdID: actor.HouseholdID,
		Text:        text,
		Qty:         qty,
		CreatedBy:   actor.UserID,
	})
}

type ShoppingItem = dbqueries.GetShoppingItemRow

// GetShoppingItem returns [g.ErrorNotFound] unless the item is on the actor's household list listID
func (s *DefaultService) GetShoppingItem(ctx context.Context, actor Membership, listID, itemID pgtype.UUID) (ShoppingItem, error) {
	it, err := s.db.GetShoppingItem(ctx, dbqueries.GetShoppingItemParams{
		ID:          itemID,
		ListID:      listID,
		HouseholdID: actor.HouseholdID,
	})
	return it, noRowsAs(err, g.ErrorNotFound)
}

// SetShoppingItemClaim is "I'm buying", claim false drops only the actor's own claim
func (s *DefaultService) SetShoppingItemClaim(ctx context.Context, actor Membership, listID, itemID pgtype.UUID, claim bool) error {
	return rowsAffected(s.db.SetShoppingItemClaim(ctx, dbqueries.SetShoppingItemClaimParams{
		Claim:       claim,
		UserID:      actor.UserID,
		ID:          itemID,
		ListID:      listID,
		HouseholdID: actor.HouseholdID,
	}))
}

func (s *DefaultService) SetShoppingItemChecked(ctx context.Context, actor Membership, listID, itemID pgtype.UUID, checked bool) error {
	return rowsAffected(s.db.SetShoppingItemChecked(ctx, dbqueries.SetShoppingItemCheckedParams{
		Checked:     checked,
		UserID:      actor.UserID,
		ID:          itemID,
		ListID:      listID,
		HouseholdID: actor.HouseholdID,
	}))
}

// DeleteShoppingItem returns [g.ErrorNotFound] unless the actor created the item
func (s *DefaultService) DeleteShoppingItem(ctx context.Context, actor Membership, listID, itemID pgtype.UUID) error {
	return rowsAffected(s.db.DeleteShoppingItem(ctx, dbqueries.DeleteShoppingItemParams{
		ID:          itemID,
		ListID:      listID,
		HouseholdID: actor.HouseholdID,
		CreatedBy:   actor.UserID,
	}))
}

// ClearCheckedShoppingItems lets any member remove the bought items
func (s *DefaultService) ClearCheckedShoppingItems(ctx context.Context, actor Membership, listID pgtype.UUID) error {
	return s.db.DeleteCheckedShoppingItems(ctx, dbqueries.DeleteCheckedShoppingItemsParams{
		ListID:      listID,
		HouseholdID: actor.HouseholdID,
	})
}
