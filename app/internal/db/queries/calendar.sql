-- name: InsertEvent :one
INSERT INTO events (
    household_id,
    title,
    location,
    notes,
    start_date,
    start_time,
    end_date,
    end_time,
    every,
    every_unit,
    weekdays,
    repeat_until,
    remind_minutes,
    detached_from,
    detached_on,
    created_by
  )
VALUES (
    @household_id,
    @title,
    @location,
    @notes,
    @start_date,
    sqlc.narg(start_time)::time,
    @end_date,
    sqlc.narg(end_time)::time,
    sqlc.narg(every)::smallint,
    sqlc.narg(every_unit)::text,
    sqlc.narg(weekdays)::smallint,
    sqlc.narg(repeat_until)::date,
    sqlc.narg(remind_minutes)::smallint,
    sqlc.narg(detached_from)::uuid,
    sqlc.narg(detached_on)::date,
    @created_by
  )
RETURNING id;
-- name: UpdateEvent :exec
UPDATE events
SET title = @title,
  location = @location,
  notes = @notes,
  start_date = @start_date,
  start_time = sqlc.narg(start_time)::time,
  end_date = @end_date,
  end_time = sqlc.narg(end_time)::time,
  every = sqlc.narg(every)::smallint,
  every_unit = sqlc.narg(every_unit)::text,
  weekdays = sqlc.narg(weekdays)::smallint,
  repeat_until = sqlc.narg(repeat_until)::date,
  remind_minutes = sqlc.narg(remind_minutes)::smallint
WHERE id = @id
  AND household_id = @household_id;
-- name: GetEvent :one
SELECT *
FROM events
WHERE id = @id
  AND household_id = @household_id;
-- name: LockEvent :one
SELECT *
FROM events
WHERE id = @id
  AND household_id = @household_id FOR
UPDATE;
-- detached one-offs cascade
-- name: DeleteEvent :execrows
DELETE FROM events
WHERE id = @id
  AND household_id = @household_id;
-- ends a series before a day, see the "following" scope
-- name: SetEventRepeatUntil :exec
UPDATE events
SET repeat_until = @repeat_until
WHERE id = @id
  AND household_id = @household_id;
-- name: InsertEventException :exec
INSERT INTO event_exceptions (event_id, household_id, occurs_on)
VALUES (@event_id, @household_id, @occurs_on) ON CONFLICT DO NOTHING;
-- the exceptions of a series, only those from from_date on when it is set
-- name: DeleteEventExceptions :exec
DELETE FROM event_exceptions
WHERE event_id = @event_id
  AND household_id = @household_id
  AND (
    sqlc.narg(from_date)::date IS NULL
    OR occurs_on >= sqlc.narg(from_date)::date
  );
-- the one-offs detached from a series, only those from from_date on when it is set
-- name: DeleteDetachedEvents :exec
DELETE FROM events
WHERE detached_from = @event_id
  AND household_id = @household_id
  AND (
    sqlc.narg(from_date)::date IS NULL
    OR detached_on >= sqlc.narg(from_date)::date
  );
-- hands a split series' exceptions from from_date on to the series that continues it
-- name: MoveEventExceptions :exec
UPDATE event_exceptions
SET event_id = @to_id
WHERE event_id = @event_id
  AND household_id = @household_id
  AND occurs_on >= @from_date;
-- name: MoveDetachedEvents :exec
UPDATE events
SET detached_from = @to_id
WHERE detached_from = @event_id
  AND household_id = @household_id
  AND detached_on >= @from_date;
-- the household's events that can have an occurrence overlapping the days [from_date, to_date)
-- name: ListEventsInWindow :many
SELECT *
FROM events
WHERE household_id = @household_id
  AND start_date < @to_date::date
  AND CASE
    WHEN every IS NULL
    AND weekdays IS NULL THEN end_date >= @from_date::date
    ELSE repeat_until IS NULL
    OR repeat_until + (end_date - start_date) >= @from_date::date
  END
ORDER BY start_date,
  start_time NULLS FIRST,
  title,
  id;
-- the days a series skips: cancelled ones and those edited alone
-- name: ListEventSkips :many
SELECT x.event_id,
  x.occurs_on
FROM event_exceptions x
WHERE x.household_id = ANY(@household_ids::uuid [])
  AND x.event_id = ANY(@event_ids::uuid [])
UNION ALL
SELECT d.detached_from,
  d.detached_on
FROM events d
WHERE d.household_id = ANY(@household_ids::uuid [])
  AND d.detached_from = ANY(@event_ids::uuid []);
-- name: ListEventMembers :many
SELECT em.event_id,
  em.user_id,
  u.username
FROM event_members em
  JOIN users u ON u.id = em.user_id
WHERE em.household_id = @household_id
  AND em.event_id = ANY(@event_ids::uuid [])
ORDER BY u.username,
  em.user_id;
-- name: InsertEventMember :exec
INSERT INTO event_members (event_id, household_id, user_id)
VALUES (@event_id, @household_id, @user_id);
-- name: DeleteEventMembers :exec
DELETE FROM event_members
WHERE event_id = @event_id
  AND household_id = @household_id;
-- name: CopyEventFileRefs :exec
INSERT INTO file_refs (file_id, household_id, event_id)
SELECT r.file_id,
  r.household_id,
  @to_id::uuid
FROM file_refs r
WHERE r.event_id = @event_id
  AND r.household_id = @household_id;
-- events whose last occurrence is today or later, detached one-offs are part of their series
-- name: CountUpcomingEvents :one
SELECT count(*)
FROM events
WHERE household_id = @household_id
  AND detached_from IS NULL
  AND CASE
    WHEN every IS NULL
    AND weekdays IS NULL THEN end_date >= @today::date
    ELSE repeat_until IS NULL
    OR repeat_until >= @today::date
  END;
-- events of live households with a reminder that can have an occurrence in [from_date, to_date],
-- the caller widens the days by the time zones and picks the reminders that are due
-- name: ListEventsToRemind :many
SELECT e.*
FROM events e
  JOIN households h ON h.id = e.household_id
WHERE h.deleted_at IS NULL
  AND e.remind_minutes IS NOT NULL
  AND e.start_date <= @to_date::date
  AND CASE
    WHEN e.every IS NULL
    AND e.weekdays IS NULL THEN e.end_date >= @from_date::date
    ELSE e.repeat_until IS NULL
    OR e.repeat_until >= @from_date::date
  END;
-- claimed before the push, 0 rows when another run did
-- name: ClaimEventReminder :execrows
INSERT INTO event_reminders (event_id, household_id, occurs_on)
VALUES (@event_id, @household_id, @occurs_on) ON CONFLICT DO NOTHING;
-- an edit that moves the start lets the event remind again
-- name: DeleteEventReminders :exec
DELETE FROM event_reminders
WHERE event_id = @event_id
  AND household_id = @household_id;
-- a claim only has to outlive its reminder, which comes at the latest on the occurrence's day
-- name: PurgeEventReminders :exec
DELETE FROM event_reminders
WHERE occurs_on < CURRENT_DATE - 2;
