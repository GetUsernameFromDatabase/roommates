-- name: InsertShoppingList :one
INSERT INTO shopping_lists (household_id, name, category, created_by)
VALUES ($1, $2, $3, $4)
RETURNING id;
-- name: ListShoppingLists :many
SELECT l.id,
  l.name,
  count(i.id) FILTER (
    WHERE i.checked_at IS NULL
      AND i.cleared_at IS NULL
  ) AS open_items
FROM shopping_lists l
  LEFT JOIN shopping_items i ON i.list_id = l.id
WHERE l.household_id = $1
GROUP BY l.id
ORDER BY l.created_at;
-- name: GetShoppingList :one
SELECT id,
  name,
  category,
  created_by
FROM shopping_lists
WHERE id = $1
  AND household_id = $2;
-- any member renames a list
-- name: SetShoppingListName :execrows
UPDATE shopping_lists
SET name = $3
WHERE id = $1
  AND household_id = $2;
-- any member files a list under another category
-- name: SetShoppingListCategory :execrows
UPDATE shopping_lists
SET category = $3
WHERE id = $1
  AND household_id = $2;
-- the files that go with the list, cleared items hold on to their files too
-- name: CountListOnlyFiles :one
SELECT count(*)
FROM files f
WHERE f.household_id = sqlc.arg(household_id)
  AND NOT f.kept
  AND EXISTS (
    SELECT 1
    FROM file_refs r
      JOIN shopping_items i ON i.id = r.shopping_item_id
    WHERE r.file_id = f.id
      AND i.list_id = sqlc.arg(list_id)
  )
  AND NOT EXISTS (
    SELECT 1
    FROM file_refs r
      LEFT JOIN shopping_items i ON i.id = r.shopping_item_id
    WHERE r.file_id = f.id
      AND i.list_id IS DISTINCT FROM sqlc.arg(list_id)
  );
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
  i.description,
  i.created_by,
  i.claimed_by,
  claimer.username AS claimed_by_name,
  claimer.avatar_key AS claimed_by_avatar_key,
  i.checked_at,
  checker.username AS checked_by_name,
  i.expense_id,
  -- shortcut: counts a file only on a private note too, the item page hides it,
  -- use file_visible here if the two ever need to agree
  (
    SELECT count(*)
    FROM file_refs r
    WHERE r.shopping_item_id = i.id
  ) AS file_count
FROM shopping_items i
  LEFT JOIN users claimer ON claimer.id = i.claimed_by
  LEFT JOIN users checker ON checker.id = i.checked_by
WHERE i.list_id = $1
  AND i.household_id = $2
  AND i.cleared_at IS NULL
ORDER BY i.checked_at IS NOT NULL,
  i.checked_at DESC,
  i.created_at;
-- name: GetShoppingItem :one
SELECT id,
  text,
  qty,
  description,
  created_by,
  expense_id
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
-- checking an already checked item keeps who bought it and when, unchecking drops the
-- expense, buying it later is a new purchase
-- name: SetShoppingItemChecked :execrows
UPDATE shopping_items
SET checked_at = CASE
    WHEN sqlc.arg(checked)::boolean THEN COALESCE(checked_at, CURRENT_TIMESTAMP)
  END,
  checked_by = CASE
    WHEN sqlc.arg(checked)::boolean THEN COALESCE(checked_by, sqlc.arg(user_id)::uuid)
  END,
  expense_id = CASE
    WHEN sqlc.arg(checked)::boolean THEN expense_id
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
-- any member describes an item
-- name: SetShoppingItemDescription :execrows
UPDATE shopping_items
SET description = sqlc.arg(description)
WHERE id = sqlc.arg(id)
  AND list_id = sqlc.arg(list_id)
  AND household_id = sqlc.arg(household_id);
-- removing bought items hides them, they can be bought again
-- name: HideCheckedShoppingItems :exec
UPDATE shopping_items
SET cleared_at = CURRENT_TIMESTAMP,
  claimed_by = NULL
WHERE list_id = $1
  AND household_id = $2
  AND checked_at IS NOT NULL
  AND cleared_at IS NULL;
-- buying again brings back the household's latest hidden item of that name with its
-- description and files, on list_id at the end of it, the old expense stays but this
-- purchase is a new one and must not show its receipt
-- name: ReviveShoppingItem :execrows
UPDATE shopping_items
SET list_id = sqlc.arg(list_id),
  qty = sqlc.arg(qty),
  created_by = sqlc.arg(created_by),
  created_at = CURRENT_TIMESTAMP,
  cleared_at = NULL,
  checked_at = NULL,
  checked_by = NULL,
  claimed_by = NULL,
  expense_id = NULL
WHERE id = (
    SELECT h.id
    FROM shopping_items h
    WHERE h.household_id = sqlc.arg(household_id)
      AND h.cleared_at IS NOT NULL
      AND lower(h.text) = lower(sqlc.arg(text))
    ORDER BY h.cleared_at DESC
    LIMIT 1
  );
-- what can be bought again on list_id, newest first, each name once with its last
-- quantity, not what is on the list already
-- shortcut: capped at 100 rows, a household with more only sees its latest
-- name: ListBuyAgain :many
SELECT h.text,
  h.qty
FROM (
    SELECT DISTINCT ON (lower(i.text)) i.text,
      i.qty,
      i.cleared_at
    FROM shopping_items i
    WHERE i.household_id = sqlc.arg(household_id)
      AND i.cleared_at IS NOT NULL
      AND NOT EXISTS (
        SELECT 1
        FROM shopping_items o
        WHERE o.list_id = sqlc.arg(list_id)
          AND o.cleared_at IS NULL
          AND lower(o.text) = lower(i.text)
      )
    ORDER BY lower(i.text),
      i.cleared_at DESC
  ) h
ORDER BY h.cleared_at DESC
LIMIT 100;
-- the bought items of a list a new expense paid for, one item_id or all not linked yet,
-- anything that is not this household's bought item simply links nothing
-- name: LinkShoppingItems :exec
UPDATE shopping_items
SET expense_id = sqlc.arg(expense_id)
WHERE household_id = sqlc.arg(household_id)
  AND list_id = sqlc.arg(list_id)
  AND checked_at IS NOT NULL
  AND cleared_at IS NULL
  AND expense_id IS NULL
  AND (
    sqlc.narg(item_id)::uuid IS NULL
    OR id = sqlc.narg(item_id)::uuid
  );
