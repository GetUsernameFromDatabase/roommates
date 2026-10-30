package services

import (
	"context"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	MaxShoppingListNameLength = 64
	MaxShoppingItemTextLength = 200
	MaxShoppingItemQtyLength  = 32

	MaxShoppingItemDescriptionLength = 500
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
// and category is one of [g.ExpenseCategories]
func (s *DefaultService) CreateShoppingList(ctx context.Context, actor Membership, name, category string) (pgtype.UUID, error) {
	name, ok := trimToLength(name, 1, MaxShoppingListNameLength)
	if !ok || !slices.Contains(g.ExpenseCategories, category) {
		return pgtype.UUID{}, g.ErrorInvalidInput
	}
	return s.db.InsertShoppingList(ctx, dbqueries.InsertShoppingListParams{
		HouseholdID: actor.HouseholdID,
		Name:        name,
		Category:    category,
		CreatedBy:   actor.UserID,
	})
}

type ShoppingList struct {
	List     dbqueries.GetShoppingListRow
	Items    []dbqueries.ListShoppingItemsRow
	BuyAgain []dbqueries.ListBuyAgainRow
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
	if err != nil {
		return l, err
	}
	l.BuyAgain, err = s.db.ListBuyAgain(ctx, dbqueries.ListBuyAgainParams{
		HouseholdID: actor.HouseholdID,
		ListID:      listID,
	})
	return l, err
}

// RenameShoppingList returns [g.ErrorInvalidInput] unless the name has 1..[MaxShoppingListNameLength] runes
func (s *DefaultService) RenameShoppingList(ctx context.Context, actor Membership, listID pgtype.UUID, name string) error {
	name, ok := trimToLength(name, 1, MaxShoppingListNameLength)
	if !ok {
		return g.ErrorInvalidInput
	}
	return rowsAffected(s.db.SetShoppingListName(ctx, dbqueries.SetShoppingListNameParams{
		ID:          listID,
		HouseholdID: actor.HouseholdID,
		Name:        name,
	}))
}

// SetShoppingListCategory returns [g.ErrorInvalidInput] unless category is one of [g.ExpenseCategories]
func (s *DefaultService) SetShoppingListCategory(ctx context.Context, actor Membership, listID pgtype.UUID, category string) error {
	if !slices.Contains(g.ExpenseCategories, category) {
		return g.ErrorInvalidInput
	}
	return rowsAffected(s.db.SetShoppingListCategory(ctx, dbqueries.SetShoppingListCategoryParams{
		ID:          listID,
		HouseholdID: actor.HouseholdID,
		Category:    category,
	}))
}

// ShoppingExpenseDraft is the title and category of a new expense for a list's bought items,
// the item's text for one item or the list's name for the whole bag (no itemID)
func (s *DefaultService) ShoppingExpenseDraft(ctx context.Context, actor Membership, listID, itemID pgtype.UUID) (title, category string, err error) {
	l, err := s.db.GetShoppingList(ctx, dbqueries.GetShoppingListParams{
		ID:          listID,
		HouseholdID: actor.HouseholdID,
	})
	if err != nil {
		return "", "", noRowsAs(err, g.ErrorNotFound)
	}
	if !itemID.Valid {
		return l.Name, l.Category, nil
	}
	it, err := s.GetShoppingItem(ctx, actor, listID, itemID)
	// an item's name can be longer than an expense title
	text := []rune(it.Text)
	return string(text[:min(len(text), MaxExpenseTitleLength)]), l.Category, err
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

// AddShoppingItem puts an item on the list, a hidden item of the same name comes back with its
// description and files, returns [g.ErrorInvalidInput] for an empty or too long text or qty
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
	revived, err := s.db.ReviveShoppingItem(ctx, dbqueries.ReviveShoppingItemParams{
		ListID:      listID,
		Qty:         qty,
		CreatedBy:   actor.UserID,
		HouseholdID: actor.HouseholdID,
		Text:        text,
	})
	if err != nil || revived > 0 {
		return err
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

// SetShoppingItemDescription returns [g.ErrorInvalidInput] over [MaxShoppingItemDescriptionLength] runes
func (s *DefaultService) SetShoppingItemDescription(ctx context.Context, actor Membership, listID, itemID pgtype.UUID, description string) error {
	// a form sends line breaks as two characters, its maxlength counts them as one
	description = strings.ReplaceAll(description, "\r\n", "\n")
	description, ok := trimToLength(description, 0, MaxShoppingItemDescriptionLength)
	if !ok {
		return g.ErrorInvalidInput
	}
	return rowsAffected(s.db.SetShoppingItemDescription(ctx, dbqueries.SetShoppingItemDescriptionParams{
		Description: description,
		ID:          itemID,
		ListID:      listID,
		HouseholdID: actor.HouseholdID,
	}))
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

// ClearCheckedShoppingItems lets any member remove the bought items, they stay to be bought again
func (s *DefaultService) ClearCheckedShoppingItems(ctx context.Context, actor Membership, listID pgtype.UUID) error {
	return s.db.HideCheckedShoppingItems(ctx, dbqueries.HideCheckedShoppingItemsParams{
		ListID:      listID,
		HouseholdID: actor.HouseholdID,
	})
}
