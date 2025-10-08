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
  h.maker_id
FROM houses h
WHERE h.id IN (
    SELECT house_id
    FROM user_houses uh
    WHERE uh.user_id = $1
  )
  OR h.maker_id = $1
ORDER BY h.name;
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
INSERT INTO houses (name, maker_id)
VALUES ($1, $2)
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
-- name: SelectHouseRoommates :many
SELECT u.id,
  u.username
FROM users u
WHERE u.id IN (
    SELECT user_id
    FROM user_houses
    WHERE house_id = $1
  )
ORDER BY u.username;
-- name: SelectHouse :one
SELECT *
FROM houses
WHERE id = $1;
-- name: SelectUserHousesWithNotes :many
SELECT h.id house_id,
  h.name house_name,
  COALESCE(
    ARRAY_AGG(hn.id) FILTER (
      WHERE hn.id IS NOT NULL
    ),
    '{}'
  )::int [] note_ids
FROM houses h
  LEFT JOIN house_notes hn ON h.id = hn.house_id
WHERE h.id IN (
    SELECT house_id
    FROM user_houses uh
    WHERE uh.user_id = $1
  )
GROUP BY h.id
ORDER BY h.name;
-- name: SelectNote :one
SELECT hn.id note_id,
  hn.title,
  hn.content,
  hn.maker_id,
  h.id house_id,
  h.name house_name
FROM house_notes hn
  INNER JOIN houses h ON hn.house_id = h.id
WHERE hn.id = $1
ORDER BY hn.updated_at;
-- name: InsertNote :one
INSERT INTO house_notes (title, content, house_id, maker_id)
VALUES ($1, $2, $3, $4)
RETURNING id;
-- name: UpdateNote :exec
UPDATE house_notes
SET title = $2,
  content = $3
WHERE id = $1;
-- name: DeleteNote :exec
DELETE FROM house_notes
WHERE id = $1;
-- name: IsUserHouseMaker :one
SELECT EXISTS (
    SELECT 1
    FROM houses
    WHERE id = @house_id
      AND maker_id = @user_id
  );
-- name: IsUserNoteMaker :one
SELECT EXISTS (
    SELECT 1
    FROM house_notes
    WHERE id = @note_id
      AND maker_id = @user_id
  );