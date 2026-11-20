-- name: MarkAccountDeleted :exec
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1;
-- no row when the erase got there first, a sign-in racing it must not revive the erased row
-- name: RestoreAccount :execrows
UPDATE users
SET deleted_at = NULL
WHERE id = $1
  AND username <> '';
-- name: DeletePushSubscriptionsOfUser :exec
DELETE FROM push_subscriptions
WHERE user_id = $1;
-- deleted accounts past their grace period, erased ones have username ''
-- name: ListAccountsToErase :many
SELECT id
FROM users
WHERE deleted_at < sqlc.arg(deleted_before)
  AND username <> '';
-- the row again under lock, no row when it was restored or erased meanwhile
-- name: LockAccountToErase :one
SELECT email,
  username,
  language
FROM users
WHERE id = sqlc.arg(id)
  AND deleted_at < sqlc.arg(deleted_before)
  AND username <> '' FOR
UPDATE;
-- name: EndMembershipsOfUser :exec
UPDATE memberships
SET left_at = CURRENT_TIMESTAMP
WHERE user_id = $1
  AND left_at IS NULL;
-- file refs cascade, delete_unreferenced_file removes the files left without one
-- name: DeletePrivateNotesOfUser :exec
DELETE FROM notes
WHERE author_id = $1
  AND visibility = 'private';
-- name: DeleteAllEmailTokensOfUser :exec
DELETE FROM email_tokens
WHERE user_id = $1;
-- queue_export_purge queues the archives
-- name: DeleteDataExportsOfUser :exec
DELETE FROM data_exports
WHERE user_id = $1;
-- the row stays for what others still show (expenses, household notes) as "Former member",
-- clearing avatar_key queues the picture's purge
-- name: EraseUser :exec
UPDATE users
SET username = '',
  email = id::text || '@erased.invalid',
  password = '',
  avatar_key = NULL,
  payee_name = NULL,
  payee_iban = NULL,
  language = NULL,
  email_validated = FALSE
WHERE id = $1;
