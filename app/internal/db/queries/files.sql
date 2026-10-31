-- name: InsertFile :one
INSERT INTO files (
    household_id,
    uploader_id,
    object_key,
    thumb_key,
    name,
    mime,
    size_bytes,
    kept
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;
-- every file the user may see, newest first, with the first place it is attached to that
-- the user may see, q, kind and place narrow it down unless they are empty,
-- shared_only leaves out what only the user's own private notes hold too
-- name: ListFiles :many
SELECT f.id,
  f.name,
  f.mime,
  f.size_bytes,
  f.created_at,
  f.uploader_id,
  uploader.username AS uploader_name,
  (f.thumb_key IS NOT NULL)::boolean AS has_thumb,
  (
    SELECT count(*)
    FROM file_refs c
    WHERE c.file_id = f.id
  ) AS ref_count,
  r.note_id,
  n.title AS note_title,
  i.list_id AS item_list_id,
  r.shopping_item_id,
  i.text AS item_text,
  r.expense_id,
  e.title AS expense_title,
  r.message_id
FROM files f
  JOIN users uploader ON uploader.id = f.uploader_id
  LEFT JOIN LATERAL (
    SELECT p.note_id,
      p.shopping_item_id,
      p.expense_id,
      p.message_id
    FROM file_refs p
      LEFT JOIN notes pn ON pn.id = p.note_id
    WHERE p.file_id = f.id
      AND (
        pn.id IS NULL
        OR pn.visibility = 'household'
        OR pn.author_id = sqlc.arg(user_id)
      )
    ORDER BY p.created_at
    LIMIT 1
  ) r ON TRUE
  LEFT JOIN notes n ON n.id = r.note_id
  LEFT JOIN shopping_items i ON i.id = r.shopping_item_id
  LEFT JOIN expenses e ON e.id = r.expense_id
WHERE f.household_id = sqlc.arg(household_id)
  AND file_visible(f.id, sqlc.arg(user_id))
  AND (
    NOT sqlc.arg(shared_only)::boolean
    OR file_visible(f.id, NULL)
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
    WHEN 'note' THEN EXISTS (
      SELECT 1
      FROM file_refs x
      WHERE x.file_id = f.id
        AND x.note_id IS NOT NULL
    )
    WHEN 'item' THEN EXISTS (
      SELECT 1
      FROM file_refs x
      WHERE x.file_id = f.id
        AND x.shopping_item_id IS NOT NULL
    )
    WHEN 'expense' THEN EXISTS (
      SELECT 1
      FROM file_refs x
      WHERE x.file_id = f.id
        AND x.expense_id IS NOT NULL
    )
    WHEN 'message' THEN EXISTS (
      SELECT 1
      FROM file_refs x
      WHERE x.file_id = f.id
        AND x.message_id IS NOT NULL
    )
    ELSE NOT EXISTS (
      SELECT 1
      FROM file_refs x
      WHERE x.file_id = f.id
    )
  END
ORDER BY f.created_at DESC;
-- the caller has already checked that the user sees the note or item, a file the user may
-- not see is left out, its note can have turned private since, last_ref says that taking the
-- file off here deletes it
-- name: ListAttachedFiles :many
SELECT f.id,
  f.name,
  f.mime,
  f.size_bytes,
  f.created_at,
  f.uploader_id,
  uploader.username AS uploader_name,
  (f.thumb_key IS NOT NULL)::boolean AS has_thumb,
  (
    NOT f.kept
    AND (
      SELECT count(*)
      FROM file_refs o
      WHERE o.file_id = f.id
    ) = 1
  )::boolean AS last_ref
FROM file_refs r
  JOIN files f ON f.id = r.file_id
  JOIN users uploader ON uploader.id = f.uploader_id
WHERE r.household_id = sqlc.arg(household_id)
  AND (
    r.note_id = sqlc.narg(note_id)
    OR r.shopping_item_id = sqlc.narg(shopping_item_id)
    OR r.expense_id = sqlc.narg(expense_id)
    OR r.message_id = sqlc.narg(message_id)
  )
  AND file_visible(f.id, sqlc.arg(user_id))
ORDER BY r.created_at;
-- attaches a file everyone in the household can see, a file that is not one, or is already
-- there, is skipped
-- name: InsertFileRef :exec
INSERT INTO file_refs (
    file_id,
    household_id,
    note_id,
    shopping_item_id,
    expense_id,
    message_id
  )
SELECT f.id,
  f.household_id,
  sqlc.narg(note_id)::uuid,
  sqlc.narg(shopping_item_id)::uuid,
  sqlc.narg(expense_id)::uuid,
  sqlc.narg(message_id)::uuid
FROM files f
WHERE f.id = sqlc.arg(file_id)
  AND f.household_id = sqlc.arg(household_id)
  AND file_visible(f.id, NULL) ON CONFLICT DO NOTHING;
-- a trigger deletes the file with its last reference unless it is kept
-- name: DeleteFileRef :exec
DELETE FROM file_refs
WHERE file_id = sqlc.arg(file_id)
  AND household_id = sqlc.arg(household_id)
  AND (
    note_id = sqlc.narg(note_id)
    OR shopping_item_id = sqlc.narg(shopping_item_id)
    OR expense_id = sqlc.narg(expense_id)
  );
-- a kept file stays after its last reference, there is no way back
-- name: KeepFile :execrows
UPDATE files
SET kept = TRUE
WHERE id = $1
  AND household_id = $2;
-- name: GetFile :one
SELECT f.id,
  f.object_key,
  f.thumb_key,
  f.name,
  f.mime,
  f.size_bytes,
  f.created_at,
  f.kept,
  f.uploader_id,
  uploader.username AS uploader_name,
  (
    SELECT count(*)
    FROM file_refs c
    WHERE c.file_id = f.id
  ) AS ref_count,
  r.note_id,
  n.title AS note_title,
  i.list_id AS item_list_id,
  r.shopping_item_id,
  i.text AS item_text,
  r.expense_id,
  e.title AS expense_title,
  r.message_id
FROM files f
  JOIN users uploader ON uploader.id = f.uploader_id
  LEFT JOIN LATERAL (
    SELECT p.note_id,
      p.shopping_item_id,
      p.expense_id,
      p.message_id
    FROM file_refs p
      LEFT JOIN notes pn ON pn.id = p.note_id
    WHERE p.file_id = f.id
      AND (
        pn.id IS NULL
        OR pn.visibility = 'household'
        OR pn.author_id = sqlc.arg(user_id)
      )
    ORDER BY p.created_at
    LIMIT 1
  ) r ON TRUE
  LEFT JOIN notes n ON n.id = r.note_id
  LEFT JOIN shopping_items i ON i.id = r.shopping_item_id
  LEFT JOIN expenses e ON e.id = r.expense_id
WHERE f.id = sqlc.arg(id)
  AND f.household_id = sqlc.arg(household_id)
  AND file_visible(f.id, sqlc.arg(user_id));
-- only the uploader renames a file
-- name: RenameFile :execrows
UPDATE files
SET name = $1
WHERE id = $2
  AND household_id = $3
  AND uploader_id = $4;
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
