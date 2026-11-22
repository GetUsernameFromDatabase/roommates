-- a new chore starts with nobody's turn, SetChoreTurn sets it once the rota is in
-- name: InsertChore :one
INSERT INTO chores (
    household_id,
    title,
    rotation,
    every,
    every_unit,
    weekdays,
    next_due,
    created_by
  )
VALUES (
    @household_id,
    @title,
    @rotation,
    sqlc.narg(every)::smallint,
    sqlc.narg(every_unit)::text,
    sqlc.narg(weekdays)::smallint,
    @next_due,
    @created_by
  )
RETURNING id;
-- a moved due date is pushed about again
-- name: UpdateChore :exec
UPDATE chores
SET title = @title,
  rotation = @rotation,
  every = sqlc.narg(every)::smallint,
  every_unit = sqlc.narg(every_unit)::text,
  weekdays = sqlc.narg(weekdays)::smallint,
  notified_on = CASE
    WHEN next_due = sqlc.arg(next_due)::date THEN notified_on
  END,
  nudged_on = CASE
    WHEN next_due = sqlc.arg(next_due)::date THEN nudged_on
  END,
  next_due = sqlc.arg(next_due)::date
WHERE id = @id
  AND household_id = @household_id;
-- sets whose turn it is on which day, the new turn is pushed about again; a handed turn (a pass
-- or an assignment) was just pushed to its new assignee, which stands for the due push and, when
-- it is already overdue, the nudge
-- name: SetChoreTurn :exec
UPDATE chores c
SET next_due = @next_due,
  assignee_id = sqlc.narg(assignee_id)::uuid,
  notified_on = CASE
    WHEN @handed::boolean
    AND @next_due <= (now() AT TIME ZONE h.time_zone)::date THEN @next_due
  END,
  nudged_on = CASE
    WHEN @handed::boolean
    AND @next_due < (now() AT TIME ZONE h.time_zone)::date THEN @next_due
  END
FROM households h
WHERE c.id = @id
  AND c.household_id = @household_id
  AND h.id = c.household_id;
-- name: DeleteChore :execrows
DELETE FROM chores
WHERE id = @id
  AND household_id = @household_id;
-- locks a chore while a turn moves on or it is edited
-- name: LockChore :one
SELECT *
FROM chores
WHERE id = @id
  AND household_id = @household_id FOR
UPDATE;
-- whose turn it is after after's, see chore_next_assignee
-- name: ChoreNextAssignee :one
SELECT chore_next_assignee(
    @chore_id::uuid,
    sqlc.narg(after)::uuid,
    @exclude_after::boolean
  )::uuid AS next;
-- name: InsertChoreMember :exec
INSERT INTO chore_members (chore_id, household_id, user_id, position)
VALUES (@chore_id, @household_id, @user_id, @position);
-- name: DeleteChoreMembers :exec
DELETE FROM chore_members
WHERE chore_id = @chore_id
  AND household_id = @household_id;
-- name: SetChoreMemberPosition :exec
UPDATE chore_members
SET position = @position
WHERE chore_id = @chore_id
  AND household_id = @household_id
  AND user_id = @user_id;
-- the household's chores, only the one id when it is set, with the household's local date
-- name: ListChores :many
SELECT c.id,
  c.title,
  c.rotation,
  c.every,
  c.every_unit,
  c.weekdays,
  c.next_due,
  c.assignee_id,
  COALESCE(u.username, '')::text AS assignee_name,
  (now() AT TIME ZONE h.time_zone)::date AS today
FROM chores c
  JOIN households h ON h.id = c.household_id
  LEFT JOIN users u ON u.id = c.assignee_id
WHERE c.household_id = @household_id
  AND (
    sqlc.narg(id)::uuid IS NULL
    OR c.id = sqlc.narg(id)::uuid
  )
ORDER BY c.next_due,
  c.title,
  c.id;
-- the rotas of the household's chores in order, only chore_id's when it is set, with when
-- each member last did it
-- name: ListChoreMembers :many
SELECT cm.chore_id,
  cm.user_id,
  cm.position,
  u.username,
  (
    SELECT max(l.created_at)
    FROM chore_log l
    WHERE l.chore_id = cm.chore_id
      AND l.user_id = cm.user_id
      AND l.kind = 'done'
  )::timestamptz AS last_done
FROM chore_members cm
  JOIN users u ON u.id = cm.user_id
WHERE cm.household_id = @household_id
  AND (
    sqlc.narg(chore_id)::uuid IS NULL
    OR cm.chore_id = sqlc.narg(chore_id)::uuid
  )
ORDER BY cm.chore_id,
  cm.position;
-- name: InsertChoreLog :exec
INSERT INTO chore_log (chore_id, household_id, kind, user_id, assignee_id, due)
VALUES (
    @chore_id,
    @household_id,
    @kind,
    @user_id,
    sqlc.narg(assignee_id)::uuid,
    @due
  );
-- name: ListChoreLog :many
SELECT l.kind,
  u.username AS user_name,
  COALESCE(a.username, '')::text AS assignee_name,
  -- by id, an erased assignee's name is '' like no assignee
  COALESCE(l.assignee_id <> l.user_id, FALSE)::boolean AS for_other,
  l.due,
  l.created_at
FROM chore_log l
  JOIN users u ON u.id = l.user_id
  LEFT JOIN users a ON a.id = l.assignee_id
WHERE l.chore_id = @chore_id
  AND l.household_id = @household_id
ORDER BY l.created_at DESC,
  l.id
LIMIT @lim;
-- chores of live households whose active assignee is due a push: due is the push about the due
-- date (on the day from the reminder hour on, at once when it is already past), the nudge the
-- one the day after from the reminder hour on
-- name: ListChoresToNotify :many
SELECT c.id,
  c.household_id,
  c.assignee_id,
  c.title,
  c.next_due,
  (c.notified_on IS DISTINCT FROM c.next_due)::boolean AS due,
  (c.next_due < (now() AT TIME ZONE h.time_zone)::date)::boolean AS overdue
FROM chores c
  JOIN households h ON h.id = c.household_id
  JOIN memberships m ON m.household_id = c.household_id
  AND m.user_id = c.assignee_id
WHERE h.deleted_at IS NULL
  AND m.left_at IS NULL
  AND (
    (
      c.notified_on IS DISTINCT FROM c.next_due
      AND (
        c.next_due < (now() AT TIME ZONE h.time_zone)::date
        OR (
          c.next_due = (now() AT TIME ZONE h.time_zone)::date
          AND extract(
            hour
            FROM now() AT TIME ZONE h.time_zone
          ) >= h.reminder_hour
        )
      )
    )
    OR (
      c.nudged_on IS DISTINCT FROM c.next_due
      AND c.next_due < (now() AT TIME ZONE h.time_zone)::date
      AND extract(
        hour
        FROM now() AT TIME ZONE h.time_zone
      ) >= h.reminder_hour
    )
  )
ORDER BY c.next_due;
-- claims the due push, and the nudge too when it is already overdue so one run sends one push;
-- 0 rows when another run did or the turn moved on
-- name: ClaimChoreNotified :execrows
UPDATE chores
SET notified_on = next_due,
  nudged_on = CASE
    WHEN @overdue::boolean THEN next_due
    ELSE nudged_on
  END
WHERE id = @id
  AND next_due = @next_due
  AND assignee_id = @assignee_id
  AND notified_on IS DISTINCT FROM next_due;
-- 0 rows too when the due push was cleared meanwhile, the next run sends that one first
-- name: ClaimChoreNudged :execrows
UPDATE chores
SET nudged_on = next_due
WHERE id = @id
  AND next_due = @next_due
  AND assignee_id = @assignee_id
  AND notified_on = next_due
  AND nudged_on IS DISTINCT FROM next_due;
