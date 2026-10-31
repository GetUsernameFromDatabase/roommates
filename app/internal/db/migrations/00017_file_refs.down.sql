DROP TRIGGER delete_unreferenced_file ON file_refs;
DROP FUNCTION delete_unreferenced_file;
DROP FUNCTION file_visible;
-- -----------------------------------------------------------------------------
ALTER TABLE files
ADD COLUMN note_id UUID,
  ADD COLUMN shopping_item_id UUID,
  ADD COLUMN expense_id UUID,
  ADD COLUMN message_id UUID,
  ADD CONSTRAINT files_check CHECK (
    num_nonnulls(note_id, shopping_item_id, expense_id, message_id) <= 1
  ),
  ADD FOREIGN KEY (note_id, household_id) REFERENCES notes (id, household_id) ON DELETE CASCADE,
  ADD FOREIGN KEY (shopping_item_id, household_id) REFERENCES shopping_items (id, household_id) ON DELETE CASCADE,
  ADD FOREIGN KEY (expense_id, household_id) REFERENCES expenses (id, household_id) ON DELETE CASCADE,
  ADD FOREIGN KEY (message_id, household_id) REFERENCES messages (id, household_id) ON DELETE CASCADE;
CREATE INDEX idx_files_note_id ON files (note_id);
CREATE INDEX idx_files_shopping_item_id ON files (shopping_item_id);
CREATE INDEX idx_files_expense_id ON files (expense_id);
CREATE INDEX idx_files_message_id ON files (message_id);
CREATE TABLE message_file_refs (
  message_id UUID NOT NULL,
  file_id UUID NOT NULL,
  household_id UUID NOT NULL,
  PRIMARY KEY (message_id, file_id),
  FOREIGN KEY (message_id, household_id) REFERENCES messages (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (file_id, household_id) REFERENCES files (id, household_id) ON DELETE CASCADE
);
CREATE INDEX idx_message_file_refs_file_id ON message_file_refs (file_id);
-- -----------------------------------------------------------------------------
-- a file's first place owns it again unless it is kept, its other chat messages refer to it,
-- other places lose it
UPDATE files f
SET note_id = r.note_id,
  shopping_item_id = r.shopping_item_id,
  expense_id = r.expense_id,
  message_id = r.message_id
FROM (
    SELECT DISTINCT ON (file_id) *
    FROM file_refs
    ORDER BY file_id,
      created_at
  ) r
WHERE r.file_id = f.id
  AND NOT f.kept;
INSERT INTO message_file_refs (message_id, file_id, household_id)
SELECT r.message_id,
  r.file_id,
  r.household_id
FROM file_refs r
  JOIN files f ON f.id = r.file_id
WHERE r.message_id IS NOT NULL
  AND r.message_id IS DISTINCT FROM f.message_id;
DROP TABLE file_refs;
ALTER TABLE files DROP COLUMN kept;
