-- This file uses [SQLC](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#schema-and-queries)
-- useful links
-- - annotations (like :one): https://docs.sqlc.dev/en/latest/reference/query-annotations.html#many
-- - named parameters: https://docs.sqlc.dev/en/latest/howto/named_parameters.html
-- -----------------------------------------------------------------------------
-- name: NeverUseThisThisIsJustSoThatFileCommentSticksToThis :one
SELECT 1
WHERE 1 = 0;
-- name: GetUserCredentials :one
SELECT id,
  email,
  username,
  password
FROM users
WHERE email = $1;
-- name: InsertUser :one
INSERT INTO users (
    email,
    username,
    -- better not store password raaw
    password
  )
VALUES ($1, $2, $3)
RETURNING id;-- -----------------------------------------------------------------------------
-- households
-- -----------------------------------------------------------------------------
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
-- -----------------------------------------------------------------------------
-- memberships
-- -----------------------------------------------------------------------------
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
-- -----------------------------------------------------------------------------
-- invites
-- -----------------------------------------------------------------------------
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
-- -----------------------------------------------------------------------------
-- users
-- -----------------------------------------------------------------------------
-- name: GetUser :one
SELECT id,
  email,
  username,
  email_validated
FROM users
WHERE id = $1;
-- name: SetEmailValidated :exec
UPDATE users
SET email_validated = TRUE
WHERE id = $1;
-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1;
-- -----------------------------------------------------------------------------
-- email tokens
-- -----------------------------------------------------------------------------
-- shortcut: expired rows are only removed by DeleteEmailTokensOfUser, add a periodic cleanup if the table grows
-- name: InsertEmailToken :exec
INSERT INTO email_tokens (token_hash, user_id, kind, expires_at)
VALUES ($1, $2, $3, $4);
-- name: GetEmailToken :one
SELECT user_id
FROM email_tokens
WHERE token_hash = $1
  AND kind = $2
  AND expires_at > CURRENT_TIMESTAMP;
-- single use: the row is deleted while reading it
-- name: ConsumeEmailToken :one
DELETE FROM email_tokens
WHERE token_hash = $1
  AND kind = $2
  AND expires_at > CURRENT_TIMESTAMP
RETURNING user_id;
-- name: DeleteEmailTokensOfUser :exec
DELETE FROM email_tokens
WHERE user_id = $1
  AND kind = $2;
