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
