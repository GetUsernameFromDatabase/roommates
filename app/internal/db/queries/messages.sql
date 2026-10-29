-- name: InsertMessage :one
INSERT INTO messages (household_id, author_id, body)
VALUES ($1, $2, $3)
RETURNING id;
-- the newest messages, older than the before message when it is given
-- name: ListMessages :many
SELECT m.id,
  m.author_id,
  author.username AS author_name,
  m.body,
  m.created_at
FROM messages m
  JOIN users author ON author.id = m.author_id
WHERE m.household_id = sqlc.arg(household_id)
  AND (
    sqlc.narg(before)::uuid IS NULL
    OR (m.created_at, m.id) < (
      SELECT b.created_at,
        b.id
      FROM messages b
      WHERE b.id = sqlc.narg(before)
        AND b.household_id = sqlc.arg(household_id)
    )
  )
ORDER BY m.created_at DESC,
  m.id DESC
LIMIT sqlc.arg(max_count);
-- name: ListMessageFiles :many
SELECT f.id,
  f.message_id,
  f.name,
  f.mime,
  f.size_bytes,
  (f.thumb_key IS NOT NULL)::boolean AS has_thumb
FROM files f
WHERE f.household_id = sqlc.arg(household_id)
  AND f.message_id = ANY(sqlc.arg(message_ids)::uuid [])
ORDER BY f.created_at;
-- name: GetMessageAuthor :one
SELECT author_id
FROM messages
WHERE id = $1
  AND household_id = $2;
-- only the author deletes a message, its files go with it
-- name: DeleteMessage :execrows
DELETE FROM messages
WHERE id = $1
  AND household_id = $2
  AND author_id = $3;
-- only the nav badges listen to it, the chat itself would mark read again and loop
-- name: NotifyChatRead :exec
SELECT pg_notify(
    'household_events',
    sqlc.arg(household_id)::uuid::text || ' read'
  );
-- name: MarkChatRead :exec
UPDATE memberships
SET chat_read_at = CURRENT_TIMESTAMP
WHERE household_id = $1
  AND user_id = $2;
