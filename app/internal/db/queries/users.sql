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
  email_validated,
  avatar_key
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
-- name: SetUserAvatar :exec
UPDATE users
SET avatar_key = $2
WHERE id = $1;
-- the picture of a user, only for themselves and the people they live with
-- name: GetUserAvatarKey :one
SELECT u.avatar_key::text
FROM users u
WHERE u.id = sqlc.arg(user_id)
  AND u.avatar_key IS NOT NULL
  AND (
    u.id = sqlc.arg(viewer_id)
    OR EXISTS (
      SELECT 1
      FROM memberships mine
        JOIN memberships theirs ON theirs.household_id = mine.household_id
      WHERE mine.user_id = sqlc.arg(viewer_id)
        AND theirs.user_id = u.id
        AND mine.left_at IS NULL
        AND theirs.left_at IS NULL
    )
  );
