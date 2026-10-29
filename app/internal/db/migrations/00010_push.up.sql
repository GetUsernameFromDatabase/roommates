-- -----------------------------------------------------------------------------
-- web push: one row per browser that turned notifications on,
-- a browser that signs in as another user and subscribes moves over to them
CREATE TABLE push_subscriptions (
  endpoint TEXT PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  p256dh TEXT NOT NULL,
  auth TEXT NOT NULL,
  -- language of the notifications, the one used when subscribing
  lang TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_push_subscriptions_user_id ON push_subscriptions (user_id);
