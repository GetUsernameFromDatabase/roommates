-- name: InsertHousehold :one
INSERT INTO households (name, created_by)
VALUES ($1, $2)
RETURNING id;
-- the id breaks name ties, so account deletes lock shared households in the same order
-- name: ListHouseholdsOfUser :many
SELECT h.id,
  h.name,
  m.role,
  h.deleted_at
FROM households h
  JOIN memberships m ON m.household_id = h.id
WHERE m.user_id = $1
  AND m.left_at IS NULL
ORDER BY h.name,
  h.id;
-- name: GetHousehold :one
SELECT id,
  name,
  currency,
  time_zone,
  reminder_hour
FROM households
WHERE id = $1;
-- the zones a household can pick: the IANA region names and UTC, no posix/ or Etc/ aliases
-- name: ListTimeZones :many
SELECT name::text
FROM pg_catalog.pg_timezone_names
WHERE name ~ '^(Africa|America|Antarctica|Asia|Atlantic|Australia|Europe|Indian|Pacific)/'
  OR name = 'UTC'
ORDER BY name;
-- name: IsTimeZone :one
SELECT EXISTS (
    SELECT 1
    FROM pg_catalog.pg_timezone_names
    WHERE name = sqlc.arg(name)::text
      AND (
        name ~ '^(Africa|America|Antarctica|Asia|Atlantic|Australia|Europe|Indian|Pacific)/'
        OR name = 'UTC'
      )
  );
-- name: SetHouseholdTimeZone :exec
UPDATE households
SET time_zone = $2
WHERE id = $1;
-- the household's local date, what "today" is for its due dates
-- name: HouseholdToday :one
SELECT (now() AT TIME ZONE time_zone)::date
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
  currency = sqlc.arg(currency),
  time_zone = sqlc.arg(time_zone),
  reminder_hour = sqlc.arg(reminder_hour)
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
