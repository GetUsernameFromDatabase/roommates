-- name: InsertExpense :one
INSERT INTO expenses (
    household_id,
    payer_id,
    category,
    title,
    amount_cents,
    incurred_on
  )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;
-- name: InsertExpenseShare :exec
INSERT INTO expense_shares (expense_id, user_id, amount_cents)
VALUES ($1, $2, $3);
-- newest first
-- shortcut: capped at 200 rows, add paging when a household gets there
-- name: ListExpenses :many
SELECT e.id,
  e.category,
  e.title,
  e.amount_cents,
  e.incurred_on,
  e.payer_id,
  payer.username AS payer_name
FROM expenses e
  JOIN users payer ON payer.id = e.payer_id
WHERE e.household_id = $1
ORDER BY e.incurred_on DESC,
  e.created_at DESC
LIMIT 200;
-- name: GetExpense :one
SELECT e.id,
  e.category,
  e.title,
  e.amount_cents,
  e.incurred_on,
  e.payer_id,
  payer.username AS payer_name
FROM expenses e
  JOIN users payer ON payer.id = e.payer_id
WHERE e.id = $1
  AND e.household_id = $2;
-- the caller has already checked the expense is in the household
-- name: ListExpenseShares :many
SELECT s.user_id,
  u.username,
  s.amount_cents,
  s.paid_at
FROM expense_shares s
  JOIN users u ON u.id = s.user_id
WHERE s.expense_id = $1
ORDER BY u.username;
-- only the payer deletes an expense, and not once someone sharing it has left,
-- their settled balance would change with nobody left to settle it
-- name: DeleteExpense :execrows
DELETE FROM expenses e
WHERE e.id = $1
  AND e.household_id = $2
  AND e.payer_id = $3
  AND NOT EXISTS (
    SELECT 1
    FROM expense_shares s
      LEFT JOIN memberships m ON m.household_id = e.household_id
      AND m.user_id = s.user_id
      AND m.left_at IS NULL
    WHERE s.expense_id = e.id
      AND m.user_id IS NULL
  );
-- only the payer, who got the money, marks a share paid, and only of an active member,
-- marking an already paid share keeps when it was paid
-- name: SetExpenseSharePaid :execrows
UPDATE expense_shares s
SET paid_at = CASE
    WHEN sqlc.arg(paid)::boolean THEN COALESCE(s.paid_at, CURRENT_TIMESTAMP)
  END
FROM expenses e
WHERE e.id = s.expense_id
  AND s.expense_id = sqlc.arg(expense_id)
  AND s.user_id = sqlc.arg(user_id)
  AND e.household_id = sqlc.arg(household_id)
  AND e.payer_id = sqlc.arg(payer_id)
  AND EXISTS (
    SELECT 1
    FROM memberships m
    WHERE m.household_id = e.household_id
      AND m.user_id = s.user_id
      AND m.left_at IS NULL
  );
-- what from_user owes to_user, one row per unpaid share of others' expenses and per
-- confirmed settlement (a settlement is a debt the other way, so to_user paid it),
-- only rows involving user_id when it is given, oldest first
-- name: ListDebtItems :many
SELECT f.from_user,
  fu.username AS from_name,
  f.to_user,
  tu.username AS to_name,
  f.expense_id,
  f.title,
  f.on_date,
  f.amount_cents
FROM (
    SELECT s.user_id AS from_user,
      e.payer_id AS to_user,
      e.id AS expense_id,
      e.title,
      e.incurred_on AS on_date,
      s.amount_cents
    FROM expense_shares s
      JOIN expenses e ON e.id = s.expense_id
    WHERE e.household_id = sqlc.arg(household_id)
      AND s.paid_at IS NULL
      AND s.user_id <> e.payer_id
      AND s.amount_cents > 0
    UNION ALL
    SELECT to_user,
      from_user,
      NULL::uuid,
      ''::text,
      confirmed_at::date,
      amount_cents
    FROM settlements
    WHERE household_id = sqlc.arg(household_id)
      AND confirmed_at IS NOT NULL
  ) f
  JOIN users fu ON fu.id = f.from_user
  JOIN users tu ON tu.id = f.to_user
WHERE sqlc.narg(user_id)::uuid IS NULL
  OR sqlc.narg(user_id)::uuid IN (f.from_user, f.to_user)
ORDER BY f.on_date;
-- settlements
-- name: InsertSettlement :exec
INSERT INTO settlements (household_id, from_user, to_user, amount_cents)
VALUES ($1, $2, $3, $4);
-- unconfirmed settlements the user sent or has to confirm
-- name: ListPendingSettlements :many
SELECT s.id,
  s.from_user,
  fu.username AS from_name,
  s.to_user,
  tu.username AS to_name,
  s.amount_cents,
  s.created_at
FROM settlements s
  JOIN users fu ON fu.id = s.from_user
  JOIN users tu ON tu.id = s.to_user
WHERE s.household_id = sqlc.arg(household_id)
  AND s.confirmed_at IS NULL
  AND sqlc.arg(user_id)::uuid IN (s.from_user, s.to_user)
ORDER BY s.created_at;
-- only the receiver confirms
-- name: ConfirmSettlement :execrows
UPDATE settlements
SET confirmed_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND household_id = $2
  AND to_user = $3
  AND confirmed_at IS NULL;
-- the sender cancels or the receiver rejects a pending settlement
-- name: DeleteSettlement :execrows
DELETE FROM settlements
WHERE id = sqlc.arg(id)
  AND household_id = sqlc.arg(household_id)
  AND sqlc.arg(user_id)::uuid IN (from_user, to_user)
  AND confirmed_at IS NULL;
