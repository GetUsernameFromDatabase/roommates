-- name: InsertFile :one
INSERT INTO files (
    household_id,
    uploader_id,
    object_key,
    thumb_key,
    name,
    mime,
    size_bytes,
    note_id,
    shopping_item_id,
    expense_id,
    message_id
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id;
-- every file the user may see, newest first, with what it is attached to,
-- q, kind and place narrow it down unless they are empty,
-- shared_only leaves out the user's own private notes too
-- name: ListFiles :many
SELECT f.id,
  f.name,
  f.mime,
  f.size_bytes,
  f.created_at,
  f.uploader_id,
  uploader.username AS uploader_name,
  (f.thumb_key IS NOT NULL)::boolean AS has_thumb,
  f.note_id,
  n.title AS note_title,
  i.list_id AS item_list_id,
  f.shopping_item_id,
  i.text AS item_text,
  f.expense_id,
  e.title AS expense_title,
  f.message_id
FROM files f
  JOIN users uploader ON uploader.id = f.uploader_id
  LEFT JOIN notes n ON n.id = f.note_id
  LEFT JOIN shopping_items i ON i.id = f.shopping_item_id
  LEFT JOIN expenses e ON e.id = f.expense_id
WHERE f.household_id = sqlc.arg(household_id)
  AND (
    n.id IS NULL
    OR n.visibility = 'household'
    OR n.author_id = sqlc.arg(user_id)
  )
  AND (
    NOT sqlc.arg(shared_only)::boolean
    OR n.id IS NULL
    OR n.visibility = 'household'
  )
  AND (
    sqlc.arg(q)::text = ''
    OR strpos(lower(f.name), lower(sqlc.arg(q))) > 0
  )
  AND CASE
    sqlc.arg(kind)::text
    WHEN '' THEN TRUE
    WHEN 'image' THEN f.mime LIKE 'image/%'
    WHEN 'pdf' THEN f.mime = 'application/pdf'
    ELSE f.mime NOT LIKE 'image/%'
    AND f.mime <> 'application/pdf'
  END
  AND CASE
    sqlc.arg(place)::text
    WHEN '' THEN TRUE
    WHEN 'note' THEN f.note_id IS NOT NULL
    WHEN 'item' THEN f.shopping_item_id IS NOT NULL
    WHEN 'expense' THEN f.expense_id IS NOT NULL
    WHEN 'message' THEN f.message_id IS NOT NULL
    ELSE num_nonnulls(
      f.note_id,
      f.shopping_item_id,
      f.expense_id,
      f.message_id
    ) = 0
  END
ORDER BY f.created_at DESC;
-- the caller has already checked that the user sees the note or item
-- name: ListAttachedFiles :many
SELECT f.id,
  f.name,
  f.mime,
  f.size_bytes,
  f.created_at,
  f.uploader_id,
  uploader.username AS uploader_name,
  (f.thumb_key IS NOT NULL)::boolean AS has_thumb
FROM files f
  JOIN users uploader ON uploader.id = f.uploader_id
WHERE f.household_id = sqlc.arg(household_id)
  AND (
    f.note_id = sqlc.narg(note_id)
    OR f.shopping_item_id = sqlc.narg(shopping_item_id)
    OR f.expense_id = sqlc.narg(expense_id)
    OR f.message_id = sqlc.narg(message_id)
  )
ORDER BY f.created_at;
-- name: GetFile :one
SELECT f.id,
  f.object_key,
  f.thumb_key,
  f.name,
  f.mime,
  f.size_bytes,
  f.created_at,
  f.uploader_id,
  uploader.username AS uploader_name,
  f.note_id,
  n.title AS note_title,
  i.list_id AS item_list_id,
  f.shopping_item_id,
  i.text AS item_text,
  f.expense_id,
  e.title AS expense_title,
  f.message_id
FROM files f
  JOIN users uploader ON uploader.id = f.uploader_id
  LEFT JOIN notes n ON n.id = f.note_id
  LEFT JOIN shopping_items i ON i.id = f.shopping_item_id
  LEFT JOIN expenses e ON e.id = f.expense_id
WHERE f.id = sqlc.arg(id)
  AND f.household_id = sqlc.arg(household_id)
  AND (
    n.id IS NULL
    OR n.visibility = 'household'
    OR n.author_id = sqlc.arg(user_id)
  );
-- only the uploader deletes a file
-- name: DeleteFile :execrows
DELETE FROM files
WHERE id = $1
  AND household_id = $2
  AND uploader_id = $3;
-- name: ListFilePurges :many
SELECT object_key
FROM file_purges
ORDER BY created_at
LIMIT $1;
-- name: DeleteFilePurge :exec
DELETE FROM file_purges
WHERE object_key = $1;
