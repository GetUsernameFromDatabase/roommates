-- emails are stored lowercase so Foo@x.com and foo@x.com are the same account
-- fails if there are existing accounts that only differ by email case
UPDATE users SET email = lower(email) WHERE email <> lower(email);
ALTER TABLE users ADD CONSTRAINT chk_users_email_lowercase CHECK (email = lower(email));
-- UNIQUE (email) already creates an index, the hash index was redundant
DROP INDEX idxh_users_email;
-- -----------------------------------------------------------------------------
