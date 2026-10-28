-- emails are not changed back to their original case, that information is lost
CREATE INDEX idxh_users_email ON users USING HASH (email);
ALTER TABLE users DROP CONSTRAINT chk_users_email_lowercase;
-- -----------------------------------------------------------------------------
