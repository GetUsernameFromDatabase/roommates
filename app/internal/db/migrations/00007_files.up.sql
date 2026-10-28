-- -----------------------------------------------------------------------------
-- targets of the composite keys in files, so a file never attaches across households
ALTER TABLE notes
ADD UNIQUE (id, household_id);
ALTER TABLE shopping_items
ADD UNIQUE (id, household_id);
-- -----------------------------------------------------------------------------
-- the bytes live in object storage under object_key, a file is attached to at most
-- one thing and is deleted with it; a file on a private note is seen only by its author
CREATE TABLE files (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  uploader_id UUID NOT NULL REFERENCES users (id),
  object_key TEXT NOT NULL UNIQUE,
  -- a small JPEG of an image, NULL for other files or when it could not be made
  thumb_key TEXT UNIQUE,
  name TEXT NOT NULL,
  mime TEXT NOT NULL,
  size_bytes BIGINT NOT NULL,
  note_id UUID,
  shopping_item_id UUID,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (note_id, household_id) REFERENCES notes (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (shopping_item_id, household_id) REFERENCES shopping_items (id, household_id) ON DELETE CASCADE,
  CHECK (num_nonnulls(note_id, shopping_item_id) <= 1)
);
CREATE INDEX idx_files_household_id ON files (household_id);
CREATE INDEX idx_files_note_id ON files (note_id);
CREATE INDEX idx_files_shopping_item_id ON files (shopping_item_id);
-- -----------------------------------------------------------------------------
-- objects of deleted files, cascades included, waiting to be removed from object storage
CREATE TABLE file_purges (
  object_key TEXT PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE FUNCTION queue_file_purge() RETURNS TRIGGER AS $$ BEGIN
INSERT INTO file_purges (object_key)
VALUES (OLD.object_key) ON CONFLICT DO NOTHING;
IF OLD.thumb_key IS NOT NULL THEN
INSERT INTO file_purges (object_key)
VALUES (OLD.thumb_key) ON CONFLICT DO NOTHING;
END IF;
RETURN OLD;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER queue_file_purge
AFTER DELETE ON files FOR EACH ROW EXECUTE FUNCTION queue_file_purge();
