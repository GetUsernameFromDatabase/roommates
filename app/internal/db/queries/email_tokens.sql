-- shortcut: expired rows are only removed by DeleteEmailTokensOfUser, add a periodic cleanup if the table grows
-- name: InsertEmailToken :exec
INSERT INTO email_tokens (token_hash, user_id, kind, expires_at)
VALUES ($1, $2, $3, $4);
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
-- name: DeleteEmailTokensOfUser :exec
DELETE FROM email_tokens
WHERE user_id = $1
  AND kind = $2;
