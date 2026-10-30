-- -----------------------------------------------------------------------------
-- what a list's expenses are filed under
ALTER TABLE shopping_lists
ADD COLUMN category TEXT NOT NULL DEFAULT 'groceries' CHECK (
    category IN ('groceries', 'utilities', 'rent', 'household', 'other')
  );
-- -----------------------------------------------------------------------------
-- cleared_at hides a bought item from its list, it can be bought again with its
-- description and files, expense_id is the expense it was bought under, several items
-- can share one receipt and deleting the expense only drops the link
ALTER TABLE shopping_items
ADD COLUMN description TEXT NOT NULL DEFAULT '',
  ADD COLUMN cleared_at TIMESTAMP WITH TIME ZONE,
  ADD COLUMN expense_id UUID,
  ADD FOREIGN KEY (expense_id, household_id) REFERENCES expenses (id, household_id) ON DELETE SET NULL (expense_id);
CREATE INDEX idx_shopping_items_expense_id ON shopping_items (expense_id);
CREATE INDEX idx_shopping_items_cleared ON shopping_items (household_id, cleared_at DESC)
WHERE cleared_at IS NOT NULL;
