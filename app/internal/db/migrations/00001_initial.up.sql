-- --- EXTENSIONS ---
CREATE EXTENSION IF NOT EXISTS moddatetime;
-- -----------------------------------------------------------------------------
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT NOT NULL UNIQUE,
  email_validated BOOLEAN NOT NULL DEFAULT FALSE,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- if in the future switch is made from like ex% to %ex% then this index won't do
CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idxh_users_email ON users USING HASH (email);
CREATE TRIGGER mdt_users BEFORE
UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
-- -----------------------------------------------------------------------------
