-- expenses for the CSV export, a member only gets those they paid or have a share of,
-- user_id NULL is all of them
-- name: ListExportExpenses :many
SELECT e.id,
  e.category,
  e.title,
  e.amount_cents,
  e.currency,
  e.incurred_on,
  e.payer_id,
  payer.username AS payer_name
FROM expenses e
  JOIN users payer ON payer.id = e.payer_id
WHERE e.household_id = sqlc.arg(household_id)
  AND (
    sqlc.narg(user_id)::uuid IS NULL
    OR e.payer_id = sqlc.narg(user_id)::uuid
    OR EXISTS (
      SELECT 1
      FROM expense_shares s
      WHERE s.expense_id = e.id
        AND s.user_id = sqlc.narg(user_id)::uuid
    )
  )
ORDER BY e.incurred_on,
  e.created_at,
  e.id;
-- the shares of the expenses ListExportExpenses returns for the same arguments,
-- the payer's own part counts as paid
-- name: ListExportShares :many
SELECT s.expense_id,
  s.user_id,
  u.username,
  s.amount_cents,
  (
    s.paid_at IS NOT NULL
    OR s.user_id = e.payer_id
  )::boolean AS paid
FROM expense_shares s
  JOIN expenses e ON e.id = s.expense_id
  JOIN users u ON u.id = s.user_id
WHERE e.household_id = sqlc.arg(household_id)
  AND (
    sqlc.narg(user_id)::uuid IS NULL
    OR e.payer_id = sqlc.narg(user_id)::uuid
    OR EXISTS (
      SELECT 1
      FROM expense_shares own
      WHERE own.expense_id = e.id
        AND own.user_id = sqlc.narg(user_id)::uuid
    )
  )
ORDER BY u.username,
  s.user_id;
