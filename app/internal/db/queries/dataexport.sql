-- queues an export unless the user has one pending, then their older ones are deleted
-- (the trigger queues their objects for the purge)
-- name: RequestDataExport :exec
WITH pending AS (
  SELECT 1
  FROM data_exports
  WHERE user_id = sqlc.arg(user_id)
    AND ready_at IS NULL
),
gone AS (
  DELETE FROM data_exports
  WHERE user_id = sqlc.arg(user_id)
    AND NOT EXISTS (
      SELECT 1
      FROM pending
    )
)
INSERT INTO data_exports (user_id, lang)
SELECT sqlc.arg(user_id),
  sqlc.arg(lang)
WHERE NOT EXISTS (
    SELECT 1
    FROM pending
  );
-- the user's newest export
-- name: GetLatestDataExport :one
SELECT id,
  ready_at
FROM data_exports
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 1;
-- the pending export tried longest ago (never first, then the oldest), stamped as tried in its
-- own statement so a failed build still moves it to the back of the queue
-- name: ClaimPendingDataExport :one
UPDATE data_exports
SET attempted_at = CURRENT_TIMESTAMP
WHERE id = (
    SELECT id
    FROM data_exports
    WHERE ready_at IS NULL
    ORDER BY attempted_at NULLS FIRST,
      created_at
    LIMIT 1 FOR UPDATE SKIP LOCKED
  )
RETURNING id,
  user_id,
  lang;
-- name: SetDataExportReady :one
UPDATE data_exports
SET object_key = sqlc.arg(object_key),
  size_bytes = sqlc.arg(size_bytes),
  ready_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
  AND ready_at IS NULL
RETURNING ready_at;
-- the user's own export, built after ready_after
-- name: GetReadyDataExport :one
SELECT object_key,
  size_bytes,
  ready_at
FROM data_exports
WHERE id = sqlc.arg(id)
  AND user_id = sqlc.arg(user_id)
  AND ready_at > sqlc.arg(ready_after);
-- exports built before ready_before, and pending ones made before pending_before
-- name: DeleteExpiredDataExports :exec
DELETE FROM data_exports
WHERE ready_at < sqlc.arg(ready_before)
  OR (
    ready_at IS NULL
    AND created_at < sqlc.arg(pending_before)
  );
-- name: GetDataExportAccount :one
SELECT username,
  email,
  email_validated,
  language,
  created_at,
  payee_name,
  payee_iban
FROM users
WHERE id = $1;
-- every household the user is or was a member of
-- name: ListDataExportMemberships :many
SELECT h.id AS household_id,
  h.name AS household,
  m.role,
  m.joined_at,
  m.left_at
FROM memberships m
  JOIN households h ON h.id = m.household_id
WHERE m.user_id = $1
ORDER BY m.joined_at;
-- expenses the user paid or has a share of, with their own share
-- name: ListDataExportExpenses :many
SELECT e.incurred_on,
  e.title,
  e.category,
  e.amount_cents,
  e.currency,
  payer.username AS payer_name,
  own.amount_cents AS my_share_cents,
  COALESCE(
    own.paid_at IS NOT NULL
    OR own.user_id = e.payer_id,
    FALSE
  )::boolean AS my_share_paid
FROM expenses e
  JOIN users payer ON payer.id = e.payer_id
  LEFT JOIN expense_shares own ON own.expense_id = e.id
  AND own.user_id = sqlc.arg(user_id)
WHERE e.household_id = sqlc.arg(household_id)
  AND (
    e.payer_id = sqlc.arg(user_id)
    OR own.user_id IS NOT NULL
  )
ORDER BY e.incurred_on,
  e.created_at;
-- settlements from or to the user
-- name: ListDataExportSettlements :many
SELECT (
    CASE
      WHEN s.from_user = sqlc.arg(user_id) THEN 'paid'
      ELSE 'received'
    END
  )::text AS direction,
  other.username AS counterparty_name,
  s.amount_cents,
  s.currency,
  s.method,
  s.note,
  s.created_at,
  s.confirmed_at
FROM settlements s
  JOIN users other ON other.id = (
    CASE
      WHEN s.from_user = sqlc.arg(user_id) THEN s.to_user
      ELSE s.from_user
    END
  )
WHERE s.household_id = sqlc.arg(household_id)
  AND sqlc.arg(user_id) IN (s.from_user, s.to_user)
ORDER BY s.created_at;
-- name: ListDataExportMessages :many
SELECT m.created_at,
  m.body,
  COALESCE(
    array_agg(
      f.name
      ORDER BY f.name
    ) FILTER (
      WHERE f.id IS NOT NULL
    ),
    '{}'
  )::text [] AS file_names
FROM messages m
  LEFT JOIN file_refs r ON r.message_id = m.id
  LEFT JOIN files f ON f.id = r.file_id
WHERE m.household_id = sqlc.arg(household_id)
  AND m.author_id = sqlc.arg(user_id)
GROUP BY m.id
ORDER BY m.created_at;
-- notes the user wrote, private ones included
-- name: ListDataExportNotes :many
SELECT title,
  body_md,
  visibility,
  created_at,
  edited_at
FROM notes
WHERE household_id = sqlc.arg(household_id)
  AND author_id = sqlc.arg(user_id)
ORDER BY created_at;
-- name: ListDataExportShoppingItems :many
SELECT l.name AS list,
  i.text,
  i.qty,
  i.description,
  i.created_at,
  i.checked_at
FROM shopping_items i
  JOIN shopping_lists l ON l.id = i.list_id
WHERE i.household_id = sqlc.arg(household_id)
  AND i.created_by = sqlc.arg(user_id)
ORDER BY i.created_at;
-- name: ListDataExportChoreLog :many
SELECT c.title AS chore,
  l.kind,
  l.due,
  l.created_at
FROM chore_log l
  JOIN chores c ON c.id = l.chore_id
WHERE l.household_id = sqlc.arg(household_id)
  AND l.user_id = sqlc.arg(user_id)
ORDER BY l.created_at;
-- times as HH:MM, '' when the event has none (all-day)
-- name: ListDataExportEvents :many
SELECT title,
  start_date,
  COALESCE(to_char(start_time, 'HH24:MI'), '')::text AS start_time,
  end_date,
  COALESCE(to_char(end_time, 'HH24:MI'), '')::text AS end_time,
  location,
  notes,
  created_at
FROM events
WHERE household_id = sqlc.arg(household_id)
  AND created_by = sqlc.arg(user_id)
ORDER BY start_date;
-- name: ListDataExportFiles :many
SELECT name,
  mime,
  size_bytes,
  created_at,
  object_key
FROM files
WHERE household_id = sqlc.arg(household_id)
  AND uploader_id = sqlc.arg(user_id)
ORDER BY created_at;
