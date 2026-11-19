-- name: InsertHousehold :one
INSERT INTO households (name, created_by)
VALUES ($1, $2)
RETURNING id;
-- name: ListHouseholdsOfUser :many
SELECT h.id,
  h.name,
  m.role,
  h.deleted_at
FROM households h
  JOIN memberships m ON m.household_id = h.id
WHERE m.user_id = $1
  AND m.left_at IS NULL
ORDER BY h.name;
-- name: GetHousehold :one
SELECT id,
  name,
  currency
FROM households
WHERE id = $1;
-- locks the household row so membership changes in it run one at a time
-- name: LockHousehold :exec
SELECT id
FROM households
WHERE id = $1 FOR
UPDATE;
-- the owners' settings, the caller checks that the currency may change
-- name: UpdateHousehold :exec
UPDATE households
SET name = sqlc.arg(name),
  currency = sqlc.arg(currency)
WHERE id = sqlc.arg(id);
-- a deleted household is read-only until PurgeDeletedHouseholds removes it
-- name: SoftDeleteHousehold :exec
UPDATE households
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND deleted_at IS NULL;
-- name: RestoreHousehold :exec
UPDATE households
SET deleted_at = NULL
WHERE id = $1;
-- every household table cascades, the files trigger queues their storage keys in file_purges
-- name: PurgeDeletedHouseholds :execrows
DELETE FROM households
WHERE deleted_at < sqlc.arg(deleted_before);
