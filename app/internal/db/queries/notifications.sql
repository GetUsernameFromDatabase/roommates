-- the kinds the member turned off in the household, see notification_mutes
-- name: ListNotificationMutes :many
SELECT kind
FROM notification_mutes
WHERE household_id = $1
  AND user_id = $2;
-- name: MuteNotification :exec
INSERT INTO notification_mutes (household_id, user_id, kind)
VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;
-- name: UnmuteNotification :exec
DELETE FROM notification_mutes
WHERE household_id = $1
  AND user_id = $2
  AND kind = $3;
