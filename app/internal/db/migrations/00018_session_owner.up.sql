-- http_sessions belongs to the session store (pgstore), which creates it on start, after the
-- migrations, so it is created here the same way, only user_id is ours: sessions.Login sets it,
-- a password reset deletes the user's sessions by it
CREATE TABLE IF NOT EXISTS http_sessions (
  id BIGSERIAL PRIMARY KEY,
  key BYTEA,
  data BYTEA,
  created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  modified_on TIMESTAMPTZ,
  expires_on TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS http_sessions_expiry_idx ON http_sessions (expires_on);
CREATE INDEX IF NOT EXISTS http_sessions_key_idx ON http_sessions (key);
-- -----------------------------------------------------------------------------
ALTER TABLE http_sessions
ADD COLUMN user_id uuid REFERENCES users (id) ON DELETE CASCADE;
CREATE INDEX http_sessions_user_id_idx ON http_sessions (user_id);
