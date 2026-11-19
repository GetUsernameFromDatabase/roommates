-- name: UpsertPushSubscription :exec
INSERT INTO push_subscriptions (endpoint, user_id, p256dh, auth, lang)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (endpoint) DO
UPDATE
SET user_id = EXCLUDED.user_id,
  p256dh = EXCLUDED.p256dh,
  auth = EXCLUDED.auth,
  lang = EXCLUDED.lang;
-- name: DeletePushSubscription :exec
DELETE FROM push_subscriptions
WHERE endpoint = $1
  AND user_id = $2;
-- the push service said the browser unsubscribed
-- name: DeleteGonePushSubscription :exec
DELETE FROM push_subscriptions
WHERE endpoint = $1;
-- leaves out the users who turned kind off in the household, see notification_mutes
-- name: ListPushSubscriptionsOfUsers :many
SELECT p.endpoint,
  p.user_id,
  p.p256dh,
  p.auth,
  p.lang
FROM push_subscriptions p
WHERE p.user_id = ANY(sqlc.arg(user_ids)::uuid [])
  AND NOT EXISTS (
    SELECT 1
    FROM notification_mutes n
    WHERE n.household_id = sqlc.arg(household_id)
      AND n.user_id = p.user_id
      AND n.kind = sqlc.arg(kind)
  );
