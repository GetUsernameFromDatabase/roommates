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
    shopping_item_id
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;
-- every file the user may see, newest first, with what it is attached to
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
  i.text AS item_text
FROM files f
  JOIN users uploader ON uploader.id = f.uploader_id
  LEFT JOIN notes n ON n.id = f.note_id
  LEFT JOIN shopping_items i ON i.id = f.shopping_item_id
WHERE f.household_id = sqlc.arg(household_id)
  AND (
    n.id IS NULL
    OR n.visibility = 'household'
    OR n.author_id = sqlc.arg(user_id)
  )
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
  i.text AS item_text
FROM files f
  JOIN users uploader ON uploader.id = f.uploader_id
  LEFT JOIN notes n ON n.id = f.note_id
  LEFT JOIN shopping_items i ON i.id = f.shopping_item_id
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
