-- This file uses SQLC -- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#schema-and-queries
-- useful links
-- - annotations (like :one): https://docs.sqlc.dev/en/latest/reference/query-annotations.html#many
-- - named parameters: https://docs.sqlc.dev/en/latest/howto/named_parameters.html
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
    password,
    full_name,
    is_full_name_public
  )
VALUES ($1, $2, $3, $4, $5)
RETURNING id;
-- name: UserHouses :many
SELECT h.id,
  h.name,
  ARRAY_AGG(uh.user_id)::UUID [] as user_ids
FROM houses h
  LEFT JOIN user_houses uh ON uh.house_id = h.id
WHERE h.id IN (
    SELECT house_id
    FROM user_houses uh
    WHERE uh.user_id = $1
  )
GROUP BY h.id;
-- name: UsersLikeExcludingExisting :many
SELECT id,
  username
FROM users
WHERE username ILIKE @username::text || '%'
  AND username NOT IN (
    SELECT UNNEST(@existing_users::text [])
  )
LIMIT 10;
-- name: InsertHouse :one
INSERT INTO houses (name)
VALUES ($1)
RETURNING id;
-- name: UpdateHouse :exec
UPDATE houses
SET name = $1
WHERE id = $2;
-- name: InsertUserIntoHouse :exec
INSERT INTO user_houses (user_id, house_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING;
-- name: DeleteHouse :exec
DELETE FROM houses
WHERE id = $1;
-- name: DeleteHouseUsers :exec
DELETE FROM user_houses
WHERE house_id = $1;
-- name: DeleteUserFromHouse :exec
DELETE FROM user_houses
WHERE user_id = $1
  AND house_id = $2;
-- name: SelectUsername :one
SELECT username
FROM users
WHERE id = $1;
-- name: SelectHouseRoommates :many
SELECT u.id,
  u.username
FROM users u
WHERE u.id IN (
    SELECT user_id
    FROM user_houses
    WHERE house_id = $1
  );
-- name: SelectHouse :one
SELECT name
FROM houses
WHERE id = $1;