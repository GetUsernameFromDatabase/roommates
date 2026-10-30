DROP INDEX idx_shopping_items_cleared;
DROP INDEX idx_shopping_items_expense_id;
ALTER TABLE shopping_items DROP COLUMN expense_id,
  DROP COLUMN cleared_at,
  DROP COLUMN description;
ALTER TABLE shopping_lists DROP COLUMN category;
