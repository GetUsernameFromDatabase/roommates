-- name: InsertNote :one
INSERT INTO notes (
    household_id,
    author_id,
    title,
    body_md,
    visibility,
    edited_by
  )
VALUES (
    sqlc.arg(household_id),
    sqlc.arg(author_id),
    sqlc.arg(title),
    sqlc.arg(body_md),
    sqlc.arg(visibility),
    sqlc.arg(author_id)
  )
RETURNING id;
-- pinned first, then the most recently edited
-- name: ListNotes :many
SELECT n.id,
  n.title,
  n.visibility,
  n.pinned,
  n.edited_at,
  editor.username AS edited_by_name,
  -- blank lines at the top would fill the slip's three lines
  left(btrim(n.body_md, E' \t\r\n'), 160) AS excerpt
FROM notes n
  JOIN users editor ON editor.id = n.edited_by
WHERE n.household_id = sqlc.arg(household_id)
  AND (
    n.visibility = 'household'
    OR n.author_id = sqlc.arg(user_id)
  )
ORDER BY n.pinned DESC,
  n.edited_at DESC;
-- name: GetNote :one
SELECT n.id,
  n.author_id,
  n.title,
  n.body_md,
  n.visibility,
  n.pinned,
  n.edited_at,
  editor.username AS edited_by_name
FROM notes n
  JOIN users editor ON editor.id = n.edited_by
WHERE n.id = sqlc.arg(id)
  AND n.household_id = sqlc.arg(household_id)
  AND (
    n.visibility = 'household'
    OR n.author_id = sqlc.arg(user_id)
  );
-- only the author changes the visibility, so nobody hides a note from its author
-- name: UpdateNote :execrows
UPDATE notes
SET title = sqlc.arg(title),
  body_md = sqlc.arg(body_md),
  visibility = CASE
    WHEN author_id = sqlc.arg(user_id) THEN sqlc.arg(visibility)
    ELSE visibility
  END,
  edited_by = sqlc.arg(user_id),
  edited_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
  AND household_id = sqlc.arg(household_id)
  AND (
    visibility = 'household'
    OR author_id = sqlc.arg(user_id)
  );
-- name: SetNotePinned :execrows
UPDATE notes
SET pinned = sqlc.arg(pinned)
WHERE id = sqlc.arg(id)
  AND household_id = sqlc.arg(household_id)
  AND (
    visibility = 'household'
    OR author_id = sqlc.arg(user_id)
  );
-- only the author deletes a note
-- name: DeleteNote :execrows
DELETE FROM notes
WHERE id = $1
  AND household_id = $2
  AND author_id = $3;
