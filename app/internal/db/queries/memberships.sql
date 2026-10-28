-- name: UpsertMembership :exec
INSERT INTO memberships (household_id, user_id, role)
VALUES ($1, $2, $3) ON CONFLICT (household_id, user_id) DO
UPDATE
SET role = EXCLUDED.role,
  joined_at = CURRENT_TIMESTAMP,
  left_at = NULL;
-- name: GetActiveMembership :one
SELECT household_id,
  user_id,
  role
FROM memberships
WHERE household_id = $1
  AND user_id = $2
  AND left_at IS NULL;
-- name: ListActiveMembers :many
SELECT m.user_id,
  u.username,
  m.role,
  m.joined_at
FROM memberships m
  JOIN users u ON u.id = m.user_id
WHERE m.household_id = $1
  AND m.left_at IS NULL
ORDER BY m.joined_at;
-- name: CountActiveMembers :one
SELECT count(*) FILTER (
    WHERE role = 'owner'
  ) AS owners,
  count(*) AS members
FROM memberships
WHERE household_id = $1
  AND left_at IS NULL;
-- name: EndMembership :exec
UPDATE memberships
SET left_at = CURRENT_TIMESTAMP
WHERE household_id = $1
  AND user_id = $2
  AND left_at IS NULL;
-- name: SetMembershipRole :exec
UPDATE memberships
SET role = $3
WHERE household_id = $1
  AND user_id = $2
  AND left_at IS NULL;
