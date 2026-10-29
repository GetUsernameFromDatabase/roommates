DROP TRIGGER queue_avatar_purge_delete ON users;
DROP TRIGGER queue_avatar_purge ON users;
DROP FUNCTION queue_avatar_purge;
ALTER TABLE users DROP COLUMN avatar_key;
