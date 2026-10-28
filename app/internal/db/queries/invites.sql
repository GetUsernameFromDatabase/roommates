-- name: InsertInvite :one
INSERT INTO invites (household_id, code, created_by, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id;
-- name: ListActiveInvites :many
SELECT id,
  code,
  expires_at
FROM invites
WHERE household_id = $1
  AND expires_at > CURRENT_TIMESTAMP
ORDER BY created_at DESC;
-- name: GetInviteByCode :one
SELECT i.id,
  i.household_id,
  h.name AS household_name
FROM invites i
  JOIN households h ON h.id = i.household_id
WHERE i.code = $1
  AND i.expires_at > CURRENT_TIMESTAMP;
-- household_id in the filter so an owner can only delete invites of their household
-- name: DeleteInvite :exec
DELETE FROM invites
WHERE id = $1
  AND household_id = $2;
-- name: DeleteInvitesOfHousehold :exec
DELETE FROM invites
WHERE household_id = $1;
