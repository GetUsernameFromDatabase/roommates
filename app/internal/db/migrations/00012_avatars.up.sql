-- -----------------------------------------------------------------------------
-- a user's picture, a small JPEG in object storage; a replaced or removed picture
-- is queued in file_purges like a deleted file
ALTER TABLE users
ADD COLUMN avatar_key TEXT UNIQUE;
CREATE FUNCTION queue_avatar_purge() RETURNS TRIGGER AS $$ BEGIN
INSERT INTO file_purges (object_key)
VALUES (OLD.avatar_key) ON CONFLICT DO NOTHING;
RETURN OLD;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER queue_avatar_purge
AFTER
UPDATE OF avatar_key ON users FOR EACH ROW
  WHEN (
    OLD.avatar_key IS NOT NULL
    AND OLD.avatar_key IS DISTINCT FROM NEW.avatar_key
  ) EXECUTE FUNCTION queue_avatar_purge();
CREATE TRIGGER queue_avatar_purge_delete
AFTER DELETE ON users FOR EACH ROW
  WHEN (OLD.avatar_key IS NOT NULL) EXECUTE FUNCTION queue_avatar_purge();
