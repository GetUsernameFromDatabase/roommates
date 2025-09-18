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
WHERE email = $1
LIMIT 1;
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
  h.name
FROM houses h
WHERE id IN (
    SELECT house_id
    FROM user_houses uh
    WHERE uh.user_id = $1
  );
-- name: UserLike :many
SELECT id,
  username
FROM users
WHERE id ILIKE @username::text || '%'
LIMIT 10;