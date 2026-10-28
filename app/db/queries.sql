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
RETURNING id;