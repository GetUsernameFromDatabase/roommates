-- name: UpsertMembership :exec
INSERT INTO memberships (household_id, user_id, role)
VALUES ($1, $2, $3) ON CONFLICT (household_id, user_id) DO
UPDATE
SET role = EXCLUDED.role,
  joined_at = CURRENT_TIMESTAMP,
  left_at = NULL;
-- name: GetActiveMembership :one
SELECT m.household_id,
  m.user_id,
  m.role,
  h.currency,
  (
    SELECT count(*)
    FROM messages msg
    WHERE msg.household_id = m.household_id
      AND msg.created_at > m.chat_read_at
      AND msg.author_id <> m.user_id
  ) AS unread_messages
FROM memberships m
  JOIN households h ON h.id = m.household_id
WHERE m.household_id = $1
  AND m.user_id = $2
  AND m.left_at IS NULL;
-- name: ListActiveMembers :many
SELECT m.user_id,
  u.username,
  u.avatar_key,
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
