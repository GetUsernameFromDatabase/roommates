-- name: InsertHousehold :one
INSERT INTO households (name, created_by)
VALUES ($1, $2)
RETURNING id;
-- name: ListHouseholdsOfUser :many
SELECT h.id,
  h.name,
  m.role
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
