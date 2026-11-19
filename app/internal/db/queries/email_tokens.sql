-- name: InsertEmailToken :exec
INSERT INTO email_tokens (token_hash, user_id, kind, expires_at, new_email)
VALUES ($1, $2, $3, $4, $5);
-- name: GetEmailToken :one
SELECT user_id
FROM email_tokens
WHERE token_hash = $1
  AND kind = $2
  AND expires_at > CURRENT_TIMESTAMP;
-- single use: the row is deleted while reading it
-- name: ConsumeEmailToken :one
DELETE FROM email_tokens
WHERE token_hash = $1
  AND kind = $2
  AND expires_at > CURRENT_TIMESTAMP
RETURNING user_id;
-- name: ConsumeEmailChangeToken :one
DELETE FROM email_tokens
WHERE token_hash = $1
  AND kind = 'email_change'
  AND expires_at > CURRENT_TIMESTAMP
RETURNING user_id,
  new_email::text;
-- name: DeleteEmailTokensOfUser :exec
DELETE FROM email_tokens
WHERE user_id = $1
  AND kind = $2;
-- name: DeleteExpiredEmailTokens :exec
DELETE FROM email_tokens
WHERE expires_at <= CURRENT_TIMESTAMP;
