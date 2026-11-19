-- name: InsertActivity :exec
INSERT INTO activity (
    household_id,
    actor_id,
    grp,
    kind,
    subject_id,
    target_user_id,
    title,
    amount_cents,
    currency
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
-- one page of the household's activity, newest first, before_at (NULL for the first page) leaves out
-- what is not older; a private note's rows are for its author, a file's rows for those who see the file;
-- the *_exists columns say whether the subject is still there to link to
-- name: ListActivity :many
SELECT a.id,
  a.actor_id,
  actor.username AS actor_name,
  a.kind,
  a.subject_id,
  a.target_user_id,
  target.username AS target_name,
  a.title,
  a.amount_cents,
  a.currency,
  a.created_at,
  EXISTS (
    SELECT 1
    FROM expenses x
    WHERE x.id = a.subject_id
      AND x.household_id = a.household_id
  ) AS expense_exists,
  EXISTS (
    SELECT 1
    FROM shopping_lists l
    WHERE l.id = a.subject_id
      AND l.household_id = a.household_id
  ) AS list_exists,
  EXISTS (
    SELECT 1
    FROM notes n
    WHERE n.id = a.subject_id
      AND n.household_id = a.household_id
  ) AS note_exists,
  EXISTS (
    SELECT 1
    FROM files f
    WHERE f.id = a.subject_id
      AND f.household_id = a.household_id
  ) AS file_exists,
  EXISTS (
    SELECT 1
    FROM recurring_expenses r
    WHERE r.id = a.subject_id
      AND r.household_id = a.household_id
  ) AS recurring_exists,
  EXISTS (
    SELECT 1
    FROM chores c
    WHERE c.id = a.subject_id
      AND c.household_id = a.household_id
  ) AS chore_exists,
  EXISTS (
    SELECT 1
    FROM events ev
    WHERE ev.id = a.subject_id
      AND ev.household_id = a.household_id
  ) AS event_exists
FROM activity a
  LEFT JOIN users actor ON actor.id = a.actor_id
  LEFT JOIN users target ON target.id = a.target_user_id
WHERE a.household_id = sqlc.arg(household_id)
  AND (
    sqlc.narg(before_at)::timestamptz IS NULL
    OR a.created_at < sqlc.narg(before_at)::timestamptz
  )
  AND NOT EXISTS (
    SELECT 1
    FROM notes n
    WHERE n.id = a.subject_id
      AND n.household_id = a.household_id
      AND n.visibility = 'private'
      AND n.author_id IS DISTINCT FROM sqlc.arg(user_id)::uuid
  )
  AND (
    NOT EXISTS (
      SELECT 1
      FROM files f
      WHERE f.id = a.subject_id
        AND f.household_id = a.household_id
    )
    OR file_visible(a.subject_id, sqlc.arg(user_id)::uuid)
  )
ORDER BY a.created_at DESC,
  a.id DESC
LIMIT sqlc.arg(page_size);
-- name: PurgeActivity :exec
DELETE FROM activity
WHERE created_at < now() - make_interval(
    secs => CASE
      grp
      WHEN 'money' THEN sqlc.arg(money_secs)::float8
      WHEN 'shopping' THEN sqlc.arg(shopping_secs)::float8
      WHEN 'notes' THEN sqlc.arg(notes_secs)::float8
      WHEN 'chores' THEN sqlc.arg(chores_secs)::float8
      WHEN 'calendar' THEN sqlc.arg(calendar_secs)::float8
      ELSE sqlc.arg(members_secs)::float8
    END
  );
