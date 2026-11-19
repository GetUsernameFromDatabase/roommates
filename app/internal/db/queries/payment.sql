-- name: GetPaymentDetails :one
SELECT payee_name,
  payee_iban
FROM users
WHERE id = $1;
-- both NULL clears
-- name: SetPaymentDetails :exec
UPDATE users
SET payee_name = sqlc.narg(name),
  payee_iban = sqlc.narg(iban)
WHERE id = sqlc.arg(user_id);
-- name: GetShowPaymentDetails :one
SELECT show_payment_details
FROM memberships
WHERE household_id = $1
  AND user_id = $2
  AND left_at IS NULL;
-- name: SetShowPaymentDetails :exec
UPDATE memberships
SET show_payment_details = sqlc.arg(show)
WHERE household_id = sqlc.arg(household_id)
  AND user_id = sqlc.arg(user_id)
  AND left_at IS NULL;
-- the active members who chose to show where to pay them
-- name: ListPayeeDetails :many
SELECT u.id,
  u.payee_name,
  u.payee_iban
FROM memberships m
  JOIN users u ON u.id = m.user_id
WHERE m.household_id = $1
  AND m.left_at IS NULL
  AND m.show_payment_details
  AND u.payee_iban IS NOT NULL
  AND u.payee_name IS NOT NULL;
