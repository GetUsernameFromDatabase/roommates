-- -----------------------------------------------------------------------------
DROP TRIGGER mdt_users ON users;
DROP INDEX idx_users_username;
DROP INDEX idxh_users_email;
DROP TABLE users;
-- -----------------------------------------------------------------------------