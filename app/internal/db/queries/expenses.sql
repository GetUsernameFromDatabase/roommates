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
-- name: HouseholdHasExpenses :one
SELECT EXISTS (
    SELECT 1
    FROM expenses
    WHERE household_id = $1
  );
-- newest first, q is a part of the title in any case, open_only keeps those
-- someone else still owes a part of, a share of or to a former member nets to zero
-- and nobody can mark it paid, so it does not count
-- shortcut: capped at 200 rows, add paging when a household gets there
-- name: ListExpenses :many
SELECT e.id,
  e.category,
  e.title,
  e.amount_cents,
  e.incurred_on,
  e.payer_id,
  payer.username AS payer_name,
  -- someone else owed a part of it and nothing of that is open
  (
    EXISTS (
      SELECT 1
      FROM expense_shares s
      WHERE s.expense_id = e.id
        AND s.user_id <> e.payer_id
        AND s.amount_cents > 0
    )
    AND NOT EXISTS (
      SELECT 1
      FROM expense_shares s
      WHERE s.expense_id = e.id
        AND s.user_id <> e.payer_id
        AND s.amount_cents > 0
        AND s.paid_at IS NULL
    )
  )::boolean AS fully_paid
FROM expenses e
  JOIN users payer ON payer.id = e.payer_id
WHERE e.household_id = sqlc.arg(household_id)
  AND strpos(lower(e.title), lower(sqlc.arg(q)::text)) > 0
  AND (
    NOT sqlc.arg(open_only)::boolean
    OR EXISTS (
      SELECT 1
      FROM expense_shares s
        JOIN memberships sharer ON sharer.household_id = e.household_id
        AND sharer.user_id = s.user_id
        JOIN memberships payer_member ON payer_member.household_id = e.household_id
        AND payer_member.user_id = e.payer_id
      WHERE s.expense_id = e.id
        AND s.user_id <> e.payer_id
        AND s.amount_cents > 0
        AND s.paid_at IS NULL
        AND sharer.left_at IS NULL
        AND payer_member.left_at IS NULL
    )
  )
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
-- keeps what would go with the expense, a file with another reference outlives it anyway
-- name: KeepExpenseFiles :exec
UPDATE files f
SET kept = TRUE
FROM file_refs r
WHERE r.file_id = f.id
  AND r.expense_id = $1
  AND r.household_id = $2
  AND (
    SELECT count(*)
    FROM file_refs o
    WHERE o.file_id = f.id
  ) = 1;
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
-- locks the share while it is marked
-- name: GetShareToMark :one
SELECT s.amount_cents,
  s.paid_at,
  e.payer_id
FROM expense_shares s
  JOIN expenses e ON e.id = s.expense_id
WHERE s.expense_id = $1
  AND s.user_id = $2
  AND e.household_id = $3 FOR
UPDATE OF s;
-- only of an active member, marking an already paid share keeps when it was paid
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
-- confirmed settlement with money left (a settlement is a debt the other way, so to_user
-- paid it), only rows involving user_id when it is given, oldest first
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
    SELECT st.to_user,
      st.from_user,
      NULL::uuid,
      ''::text,
      st.confirmed_at::date,
      st.left_cents
    FROM (
        SELECT c.to_user,
          c.from_user,
          c.confirmed_at,
          c.amount_cents - (
            SELECT COALESCE(SUM(ss.cents), 0)::bigint
            FROM settlement_shares ss
            WHERE ss.settlement_id = c.id
          ) AS left_cents
        FROM settlements c
        WHERE c.household_id = sqlc.arg(household_id)
          AND c.confirmed_at IS NOT NULL
      ) st
    WHERE st.left_cents > 0
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
-- name: ConfirmSettlement :one
UPDATE settlements
SET confirmed_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND household_id = $2
  AND to_user = $3
  AND confirmed_at IS NULL
RETURNING from_user,
  to_user;
-- its own statement, so ListSettlementCredit then sees what others spent while this waited
-- name: LockSettlementCredit :exec
SELECT id
FROM settlements
WHERE household_id = $1
  AND from_user = $2
  AND to_user = $3
  AND confirmed_at IS NOT NULL FOR
UPDATE;
-- what is left of the confirmed settlements from_user paid to_user, oldest first
-- name: ListSettlementCredit :many
SELECT st.id,
  (st.amount_cents - (
    SELECT COALESCE(SUM(ss.cents), 0)::bigint
    FROM settlement_shares ss
    WHERE ss.settlement_id = st.id
  ))::bigint AS left_cents
FROM settlements st
WHERE st.household_id = $1
  AND st.from_user = $2
  AND st.to_user = $3
  AND st.confirmed_at IS NOT NULL
ORDER BY st.confirmed_at;
-- name: InsertSettlementShare :exec
INSERT INTO settlement_shares (settlement_id, expense_id, user_id, cents)
VALUES ($1, $2, $3, $4);
-- gives back what a share spent
-- name: DeleteSettlementShares :exec
DELETE FROM settlement_shares
WHERE expense_id = $1
  AND user_id = $2;
-- the sender cancels or the receiver rejects a pending settlement
-- name: DeleteSettlement :execrows
DELETE FROM settlements
WHERE id = sqlc.arg(id)
  AND household_id = sqlc.arg(household_id)
  AND sqlc.arg(user_id)::uuid IN (from_user, to_user)
  AND confirmed_at IS NULL;
