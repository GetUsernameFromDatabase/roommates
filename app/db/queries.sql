-- This file uses SQLC -- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#schema-and-queries
-- name: GetUserCredentials :one
SELECT id,
  email,
  username,
  password
FROM users
WHERE email = $1;
-- name: InsertUser :one
INSERT INTO users (email, username, password)
VALUES ($1, $2, $3)
RETURNING id;