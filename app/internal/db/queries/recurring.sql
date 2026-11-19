-- name: InsertRecurring :one
INSERT INTO recurring_expenses (
    household_id,
    payer_id,
    category,
    title,
    amount_cents,
    custom,
    interval,
    anchor,
    next_due
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;
-- amount_cents is NULL when the schedule splits equally
-- name: InsertRecurringShare :exec
INSERT INTO recurring_expense_shares (recurring_id, household_id, user_id, amount_cents)
VALUES ($1, $2, $3, $4);
-- name: DeleteRecurringShares :exec
DELETE FROM recurring_expense_shares
WHERE recurring_id = $1
  AND household_id = $2;
-- the household's schedules, only the one id when it is set, only those due_for pays that are due
-- (on the household's local date) when that is set
-- name: ListRecurring :many
SELECT r.id,
  r.payer_id,
  u.username AS payer_name,
  r.category,
  r.title,
  r.amount_cents,
  r.custom,
  r.interval,
  r.anchor,
  r.next_due
FROM recurring_expenses r
  JOIN users u ON u.id = r.payer_id
  JOIN households h ON h.id = r.household_id
WHERE r.household_id = @household_id
  AND (
    sqlc.narg(id)::uuid IS NULL
    OR r.id = sqlc.narg(id)::uuid
  )
  AND (
    sqlc.narg(due_for)::uuid IS NULL
    OR (
      r.payer_id = sqlc.narg(due_for)::uuid
      AND r.next_due <= (now() AT TIME ZONE h.time_zone)::date
    )
  )
ORDER BY r.next_due,
  r.title,
  r.id;
-- name: ListRecurringShares :many
SELECT user_id,
  amount_cents
FROM recurring_expense_shares
WHERE recurring_id = $1
  AND household_id = $2;
-- locks a schedule while a draft is added or skipped or it is edited
-- name: LockRecurring :one
SELECT *
FROM recurring_expenses
WHERE id = $1
  AND household_id = $2 FOR
UPDATE;
-- name: SetRecurringNextDue :exec
UPDATE recurring_expenses
SET next_due = $3
WHERE id = $1
  AND household_id = $2;
-- name: UpdateRecurring :exec
UPDATE recurring_expenses
SET category = $3,
  title = $4,
  amount_cents = $5,
  custom = $6,
  interval = $7,
  anchor = $8,
  next_due = $9
WHERE id = $1
  AND household_id = $2;
-- name: DeleteRecurring :execrows
DELETE FROM recurring_expenses
WHERE id = $1
  AND household_id = $2;
-- due schedules whose payer was not pushed about this due date yet, in live households;
-- on the due day itself only from the household's reminder hour on
-- name: ListRecurringToNotify :many
SELECT r.id,
  r.household_id,
  r.payer_id,
  r.title,
  r.next_due
FROM recurring_expenses r
  JOIN households h ON h.id = r.household_id
  JOIN memberships m ON m.household_id = r.household_id
  AND m.user_id = r.payer_id
WHERE (
    r.next_due < (now() AT TIME ZONE h.time_zone)::date
    OR (
      r.next_due = (now() AT TIME ZONE h.time_zone)::date
      AND extract(
        hour
        FROM now() AT TIME ZONE h.time_zone
      ) >= h.reminder_hour
    )
  )
  AND r.notified_on IS DISTINCT FROM r.next_due
  AND h.deleted_at IS NULL
  AND m.left_at IS NULL
ORDER BY r.next_due;
-- claims the push about due, 0 rows when another run did or the schedule moved on
-- name: SetRecurringNotified :execrows
UPDATE recurring_expenses
SET notified_on = next_due
WHERE id = $1
  AND next_due = $2
  AND notified_on IS DISTINCT FROM next_due;
