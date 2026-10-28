-- name: InsertShoppingList :one
INSERT INTO shopping_lists (household_id, name, created_by)
VALUES ($1, $2, $3)
RETURNING id;
-- name: ListShoppingLists :many
SELECT l.id,
  l.name,
  count(i.id) FILTER (
    WHERE i.checked_at IS NULL
  ) AS open_items
FROM shopping_lists l
  LEFT JOIN shopping_items i ON i.list_id = l.id
WHERE l.household_id = $1
GROUP BY l.id
ORDER BY l.created_at;
-- name: GetShoppingList :one
SELECT id,
  name,
  created_by
FROM shopping_lists
WHERE id = $1
  AND household_id = $2;
-- only the creator deletes a list
-- name: DeleteShoppingList :execrows
DELETE FROM shopping_lists
WHERE id = $1
  AND household_id = $2
  AND created_by = $3;
-- shopping items
-- name: InsertShoppingItem :exec
INSERT INTO shopping_items (list_id, household_id, text, qty, created_by)
VALUES ($1, $2, $3, $4, $5);
-- open items first in the order added, then the most recently checked
-- name: ListShoppingItems :many
SELECT i.id,
  i.text,
  i.qty,
  i.created_by,
  i.claimed_by,
  claimer.username AS claimed_by_name,
  i.checked_at,
  checker.username AS checked_by_name,
  (
    SELECT count(*)
    FROM files f
    WHERE f.shopping_item_id = i.id
  ) AS file_count
FROM shopping_items i
  LEFT JOIN users claimer ON claimer.id = i.claimed_by
  LEFT JOIN users checker ON checker.id = i.checked_by
WHERE i.list_id = $1
  AND i.household_id = $2
ORDER BY i.checked_at IS NOT NULL,
  i.checked_at DESC,
  i.created_at;
-- name: GetShoppingItem :one
SELECT id,
  text,
  qty
FROM shopping_items
WHERE id = $1
  AND list_id = $2
  AND household_id = $3;
-- claim true claims the item for user_id, false drops only user_id's own claim,
-- a stale page cannot undo what someone else did
-- name: SetShoppingItemClaim :execrows
UPDATE shopping_items
SET claimed_by = CASE
    WHEN sqlc.arg(claim)::boolean THEN sqlc.arg(user_id)::uuid
    ELSE NULLIF(claimed_by, sqlc.arg(user_id)::uuid)
  END
WHERE id = sqlc.arg(id)
  AND list_id = sqlc.arg(list_id)
  AND household_id = sqlc.arg(household_id);
-- checking an already checked item keeps who bought it and when
-- name: SetShoppingItemChecked :execrows
UPDATE shopping_items
SET checked_at = CASE
    WHEN sqlc.arg(checked)::boolean THEN COALESCE(checked_at, CURRENT_TIMESTAMP)
  END,
  checked_by = CASE
    WHEN sqlc.arg(checked)::boolean THEN COALESCE(checked_by, sqlc.arg(user_id)::uuid)
  END
WHERE id = sqlc.arg(id)
  AND list_id = sqlc.arg(list_id)
  AND household_id = sqlc.arg(household_id);
-- only the creator deletes an item
-- name: DeleteShoppingItem :execrows
DELETE FROM shopping_items
WHERE id = $1
  AND list_id = $2
  AND household_id = $3
  AND created_by = $4;
-- name: DeleteCheckedShoppingItems :exec
DELETE FROM shopping_items
WHERE list_id = $1
  AND household_id = $2
  AND checked_at IS NOT NULL;
