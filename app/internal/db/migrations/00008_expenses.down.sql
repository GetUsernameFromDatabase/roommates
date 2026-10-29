DELETE FROM files
WHERE expense_id IS NOT NULL;
ALTER TABLE files DROP CONSTRAINT files_check,
  DROP COLUMN expense_id,
  ADD CONSTRAINT files_check CHECK (num_nonnulls(note_id, shopping_item_id) <= 1);
DROP TABLE settlements;
DROP TABLE expense_shares;
DROP TABLE expenses;
