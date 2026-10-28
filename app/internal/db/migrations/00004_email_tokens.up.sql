-- -----------------------------------------------------------------------------
-- single use tokens sent by email, only the sha256 of the token is stored
-- so a database leak does not leak working links
CREATE TABLE email_tokens (
  token_hash BYTEA PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  kind TEXT NOT NULL CHECK (kind IN ('verify', 'reset')),
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_email_tokens_user_id ON email_tokens (user_id);
