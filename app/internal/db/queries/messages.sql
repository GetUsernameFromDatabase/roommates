-- name: InsertMessage :one
INSERT INTO messages (household_id, author_id, body, reply_to_id)
VALUES ($1, $2, $3, $4)
RETURNING id;
-- the newest messages, older than the before message when it is given
-- name: ListMessages :many
SELECT m.id,
  m.author_id,
  author.username AS author_name,
  author.avatar_key AS author_avatar_key,
  m.body,
  m.created_at,
  m.reply_to_id,
  reply_author.username AS reply_author_name,
  coalesce(left(reply.body, 200), '')::text AS reply_body
FROM messages m
  JOIN users author ON author.id = m.author_id
  LEFT JOIN messages reply ON reply.id = m.reply_to_id
  LEFT JOIN users reply_author ON reply_author.id = reply.author_id
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
-- the quote of a message being answered
-- name: GetReplyPreview :one
SELECT m.id,
  author.username AS author_name,
  left(m.body, 200)::text AS body
FROM messages m
  JOIN users author ON author.id = m.author_id
WHERE m.id = $1
  AND m.household_id = $2;
-- the files sent with the messages and those they refer to, a referred file the user may
-- not see is left out, its note can have turned private since
-- name: ListMessageFiles :many
SELECT mf.id,
  mf.message_id,
  mf.name,
  mf.mime,
  mf.size_bytes,
  mf.has_thumb
FROM (
    SELECT f.id,
      f.message_id,
      f.name,
      f.mime,
      f.size_bytes,
      (f.thumb_key IS NOT NULL)::boolean AS has_thumb,
      f.created_at
    FROM files f
    WHERE f.household_id = sqlc.arg(household_id)
      AND f.message_id = ANY(sqlc.arg(message_ids)::uuid [])
    UNION ALL
    SELECT f.id,
      r.message_id,
      f.name,
      f.mime,
      f.size_bytes,
      (f.thumb_key IS NOT NULL)::boolean AS has_thumb,
      f.created_at
    FROM message_file_refs r
      JOIN files f ON f.id = r.file_id
      LEFT JOIN notes n ON n.id = f.note_id
    WHERE r.household_id = sqlc.arg(household_id)
      AND r.message_id = ANY(sqlc.arg(message_ids)::uuid [])
      AND (
        n.id IS NULL
        OR n.visibility = 'household'
        OR n.author_id = sqlc.arg(user_id)
      )
  ) mf
ORDER BY mf.created_at;
-- refers to the files everyone in the household can see, the others are skipped
-- name: InsertMessageFileRefs :execrows
INSERT INTO message_file_refs (message_id, file_id, household_id)
SELECT sqlc.arg(message_id)::uuid,
  f.id,
  f.household_id
FROM files f
  LEFT JOIN notes n ON n.id = f.note_id
WHERE f.id = ANY(sqlc.arg(file_ids)::uuid [])
  AND f.household_id = sqlc.arg(household_id)
  AND (
    n.id IS NULL
    OR n.visibility = 'household'
  );
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
