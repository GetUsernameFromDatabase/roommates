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
-- name: ListPushSubscriptionsOfUsers :many
SELECT endpoint,
  user_id,
  p256dh,
  auth,
  lang
FROM push_subscriptions
WHERE user_id = ANY(sqlc.arg(user_ids)::uuid []);
