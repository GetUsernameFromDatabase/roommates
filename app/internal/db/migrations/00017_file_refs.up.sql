-- -----------------------------------------------------------------------------
-- every attachment is a reference, a file lasts as long as something refers to it,
-- kept files, those of the Files page, last until they are deleted
ALTER TABLE files
ADD COLUMN kept BOOLEAN NOT NULL DEFAULT FALSE;
CREATE TABLE file_refs (
  file_id UUID NOT NULL,
  household_id UUID NOT NULL,
  note_id UUID,
  shopping_item_id UUID,
  expense_id UUID,
  message_id UUID,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (file_id, household_id) REFERENCES files (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (note_id, household_id) REFERENCES notes (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (shopping_item_id, household_id) REFERENCES shopping_items (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (expense_id, household_id) REFERENCES expenses (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (message_id, household_id) REFERENCES messages (id, household_id) ON DELETE CASCADE,
  CHECK (
    num_nonnulls(note_id, shopping_item_id, expense_id, message_id) = 1
  ),
  UNIQUE (file_id, note_id),
  UNIQUE (file_id, shopping_item_id),
  UNIQUE (file_id, expense_id),
  UNIQUE (file_id, message_id)
);
CREATE INDEX idx_file_refs_note_id ON file_refs (note_id);
CREATE INDEX idx_file_refs_shopping_item_id ON file_refs (shopping_item_id);
CREATE INDEX idx_file_refs_expense_id ON file_refs (expense_id);
CREATE INDEX idx_file_refs_message_id ON file_refs (message_id);
CREATE TRIGGER notify_file_refs
AFTER INSERT OR UPDATE OR DELETE ON file_refs FOR EACH ROW EXECUTE FUNCTION notify_household('files');
-- -----------------------------------------------------------------------------
-- the created_at of the file, so where it was uploaded stays its first place
INSERT INTO file_refs (
    file_id,
    household_id,
    note_id,
    shopping_item_id,
    expense_id,
    message_id,
    created_at
  )
SELECT f.id,
  f.household_id,
  f.note_id,
  f.shopping_item_id,
  f.expense_id,
  f.message_id,
  f.created_at
FROM files f
WHERE num_nonnulls(
    f.note_id,
    f.shopping_item_id,
    f.expense_id,
    f.message_id
  ) = 1;
INSERT INTO file_refs (file_id, household_id, message_id)
SELECT file_id,
  household_id,
  message_id
FROM message_file_refs ON CONFLICT DO NOTHING;
UPDATE files
SET kept = TRUE
WHERE num_nonnulls(note_id, shopping_item_id, expense_id, message_id) = 0;
-- -----------------------------------------------------------------------------
DROP TABLE message_file_refs;
ALTER TABLE files DROP CONSTRAINT files_check,
  DROP COLUMN note_id,
  DROP COLUMN shopping_item_id,
  DROP COLUMN expense_id,
  DROP COLUMN message_id;
-- -----------------------------------------------------------------------------
-- a file only on private notes is for their authors and its uploader, a kept file is for
-- everyone wherever it is attached, NULL as the user asks whether everyone in the household sees it
CREATE FUNCTION file_visible(fid UUID, uid UUID) RETURNS BOOLEAN AS $$
SELECT EXISTS (
    SELECT 1
    FROM files f
    WHERE f.id = fid
      AND (
        f.kept
        OR f.uploader_id = uid
      )
  )
  OR NOT EXISTS (
    SELECT 1
    FROM file_refs r
    WHERE r.file_id = fid
  )
  OR EXISTS (
    SELECT 1
    FROM file_refs r
      LEFT JOIN notes n ON n.id = r.note_id
    WHERE r.file_id = fid
      AND (
        n.id IS NULL
        OR n.visibility = 'household'
        OR n.author_id = uid
      )
  );
$$ LANGUAGE sql STABLE;
-- also runs for the cascades, so a deleted note, item, expense or message takes its files along
CREATE FUNCTION delete_unreferenced_file() RETURNS TRIGGER AS $$ BEGIN
DELETE FROM files f
WHERE f.id = OLD.file_id
  AND NOT f.kept
  AND NOT EXISTS (
    SELECT 1
    FROM file_refs r
    WHERE r.file_id = f.id
  );
RETURN NULL;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER delete_unreferenced_file
AFTER DELETE ON file_refs FOR EACH ROW EXECUTE FUNCTION delete_unreferenced_file();
