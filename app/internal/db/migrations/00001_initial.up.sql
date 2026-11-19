-- --- EXTENSIONS ---
CREATE EXTENSION IF NOT EXISTS moddatetime;
-- -----------------------------------------------------------------------------
-- emails are stored lowercase so Foo@x.com and foo@x.com are the same account,
-- avatar_key is the user's picture, a small JPEG in object storage
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT NOT NULL UNIQUE,
  email_validated BOOLEAN NOT NULL DEFAULT FALSE,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  avatar_key TEXT UNIQUE,
  -- NULL until the user picks one, then the browser's language is used
  language TEXT,
  -- where to pay the user, shown only in households whose membership has show_payment_details on
  payee_name TEXT,
  payee_iban TEXT,
  CONSTRAINT chk_users_email_lowercase CHECK (email = lower(email))
);
-- if in the future switch is made from like ex% to %ex% then this index won't do
CREATE INDEX idx_users_username ON users (username);
-- -----------------------------------------------------------------------------
CREATE TABLE households (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  -- ISO 4217, amounts are stored as integer minor units (cents)
  currency TEXT NOT NULL DEFAULT 'EUR',
  created_by UUID NOT NULL REFERENCES users (id),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  -- soft delete: read-only until purged g.HouseholdPurgeAfter later, an owner can restore it
  deleted_at TIMESTAMP WITH TIME ZONE,
  -- an IANA name from pg_timezone_names (the service checks it), the household's "today"
  -- for due dates is (now() AT TIME ZONE time_zone)::date
  time_zone TEXT NOT NULL DEFAULT 'Europe/Tallinn',
  -- local hour at which chore and recurring expense reminders are pushed
  reminder_hour SMALLINT NOT NULL DEFAULT 9 CHECK (reminder_hour BETWEEN 0 AND 23)
);
-- -----------------------------------------------------------------------------
-- a row is kept after leaving (left_at set) so history still shows former members,
-- rejoining clears left_at
CREATE TABLE memberships (
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users (id),
  role TEXT NOT NULL CHECK (role IN ('owner', 'member')),
  joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  left_at TIMESTAMP WITH TIME ZONE,
  -- messages of others after it are unread, a new member starts with nothing unread
  chat_read_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  show_payment_details BOOLEAN NOT NULL DEFAULT FALSE,
  PRIMARY KEY (household_id, user_id)
);
CREATE INDEX idx_memberships_user_id ON memberships (user_id);
-- -----------------------------------------------------------------------------
-- reusable until expires_at, the code works both in a link and typed by hand
CREATE TABLE invites (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  code TEXT NOT NULL UNIQUE,
  created_by UUID NOT NULL REFERENCES users (id),
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_invites_household_id ON invites (household_id);
-- -----------------------------------------------------------------------------
-- single use tokens sent by email, only the sha256 of the token is stored
-- so a database leak does not leak working links;
-- an email_change token carries the address it confirms, the swap happens when it is used
CREATE TABLE email_tokens (
  token_hash BYTEA PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  kind TEXT NOT NULL CHECK (kind IN ('verify', 'reset', 'email_change')),
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  new_email TEXT,
  CHECK ((kind = 'email_change') = (new_email IS NOT NULL))
);
CREATE INDEX idx_email_tokens_user_id ON email_tokens (user_id);
-- -----------------------------------------------------------------------------
-- the categories an expense or shopping list can have, config.ExpenseCategories must match
CREATE TABLE expense_categories (id TEXT PRIMARY KEY);
INSERT INTO expense_categories (id)
VALUES ('groceries'),
  ('utilities'),
  ('rent'),
  ('household'),
  ('other');
-- -----------------------------------------------------------------------------
CREATE TABLE shopping_lists (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  created_by UUID NOT NULL REFERENCES users (id),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  -- what a list's expenses are filed under
  category TEXT NOT NULL DEFAULT 'groceries' REFERENCES expense_categories (id) ON UPDATE CASCADE,
  -- target of the composite key in shopping_items
  UNIQUE (id, household_id)
);
CREATE INDEX idx_shopping_lists_household_id ON shopping_lists (household_id);
-- -----------------------------------------------------------------------------
-- whoever enters an expense paid it, so payer_id is also its creator
CREATE TABLE expenses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  payer_id UUID NOT NULL REFERENCES users (id),
  category TEXT NOT NULL REFERENCES expense_categories (id) ON UPDATE CASCADE,
  title TEXT NOT NULL,
  -- in currency's minor units, the household's currency when it was made; the household
  -- currency only changes while everyone is settled, so open balances never mix two
  amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
  currency TEXT NOT NULL,
  incurred_on DATE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  -- target of the composite keys in shopping_items and file_refs
  UNIQUE (id, household_id)
);
CREATE INDEX idx_expenses_household_id ON expenses (household_id, incurred_on DESC);
-- -----------------------------------------------------------------------------
-- (list_id, household_id) references the list, so an item can never sit
-- in a list of another household;
-- cleared_at hides a bought item from its list, it can be bought again with its
-- description and files, expense_id is the expense it was bought under, several items
-- can share one receipt and deleting the expense only drops the link
CREATE TABLE shopping_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  list_id UUID NOT NULL,
  household_id UUID NOT NULL,
  text TEXT NOT NULL,
  qty TEXT NOT NULL DEFAULT '',
  created_by UUID NOT NULL REFERENCES users (id),
  -- "I'm buying"
  claimed_by UUID REFERENCES users (id),
  checked_by UUID REFERENCES users (id),
  checked_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  description TEXT NOT NULL DEFAULT '',
  cleared_at TIMESTAMP WITH TIME ZONE,
  expense_id UUID,
  FOREIGN KEY (list_id, household_id) REFERENCES shopping_lists (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (expense_id, household_id) REFERENCES expenses (id, household_id) ON DELETE SET NULL (expense_id),
  -- target of the composite key in file_refs
  UNIQUE (id, household_id)
);
CREATE INDEX idx_shopping_items_list_id ON shopping_items (list_id);
CREATE INDEX idx_shopping_items_expense_id ON shopping_items (expense_id);
CREATE INDEX idx_shopping_items_cleared ON shopping_items (household_id, cleared_at DESC)
WHERE cleared_at IS NOT NULL;
-- -----------------------------------------------------------------------------
-- a household note is read and edited by every member, a private one only by its author
CREATE TABLE notes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  author_id UUID NOT NULL REFERENCES users (id),
  title TEXT NOT NULL,
  body_md TEXT NOT NULL DEFAULT '',
  visibility TEXT NOT NULL DEFAULT 'household' CHECK (visibility IN ('household', 'private')),
  pinned BOOLEAN NOT NULL DEFAULT FALSE,
  -- last write wins, these show who wrote it; updated_at also moves on a pin
  edited_by UUID NOT NULL REFERENCES users (id),
  edited_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  -- target of the composite key in file_refs
  UNIQUE (id, household_id)
);
CREATE INDEX idx_notes_household_id ON notes (household_id);
-- -----------------------------------------------------------------------------
-- the shares of an expense sum to its amount, paid_at is set by the payer
CREATE TABLE expense_shares (
  expense_id UUID NOT NULL REFERENCES expenses (id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users (id),
  amount_cents BIGINT NOT NULL CHECK (amount_cents >= 0),
  paid_at TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY (expense_id, user_id)
);
CREATE INDEX idx_expense_shares_user_id ON expense_shares (user_id);
-- -----------------------------------------------------------------------------
-- "from_user paid to_user", counts towards balances once to_user confirms it;
-- method is how from_user paid, NULL only on one an expense edit made (from_edit), that one
-- records a share the payer had marked paid outside the app, expense_id is that expense
CREATE TABLE settlements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  from_user UUID NOT NULL REFERENCES users (id),
  to_user UUID NOT NULL REFERENCES users (id),
  amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
  currency TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  confirmed_at TIMESTAMP WITH TIME ZONE,
  method TEXT CHECK (method IN ('transfer', 'cash', 'card', 'other')),
  note TEXT NOT NULL DEFAULT '',
  from_edit BOOLEAN NOT NULL DEFAULT FALSE,
  expense_id UUID,
  CHECK (from_user <> to_user),
  FOREIGN KEY (expense_id, household_id) REFERENCES expenses (id, household_id) ON DELETE
  SET NULL (expense_id)
);
CREATE INDEX idx_settlements_household_id ON settlements (household_id);
CREATE INDEX idx_settlements_expense_id ON settlements (expense_id);
-- -----------------------------------------------------------------------------
-- how much of a confirmed settlement went to which share of an expense, marking a share
-- paid spends the settlement, only what is left of it still lowers the debt
CREATE TABLE settlement_shares (
  settlement_id UUID NOT NULL REFERENCES settlements (id) ON DELETE CASCADE,
  expense_id UUID NOT NULL,
  user_id UUID NOT NULL,
  cents BIGINT NOT NULL CHECK (cents > 0),
  PRIMARY KEY (settlement_id, expense_id, user_id),
  FOREIGN KEY (expense_id, user_id) REFERENCES expense_shares (expense_id, user_id) ON DELETE CASCADE
);
CREATE INDEX idx_settlement_shares_share ON settlement_shares (expense_id, user_id);
-- -----------------------------------------------------------------------------
-- one chat per household, a message has a body, files or both;
-- a message can answer another message of the same household,
-- deleting the answered message keeps the reply without its quote
CREATE TABLE messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  author_id UUID NOT NULL REFERENCES users (id),
  body TEXT NOT NULL CHECK (length(body) <= 4000),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  reply_to_id UUID,
  -- target of the composite keys here and in file_refs
  UNIQUE (id, household_id),
  FOREIGN KEY (reply_to_id, household_id) REFERENCES messages (id, household_id) ON DELETE SET NULL (reply_to_id)
);
CREATE INDEX idx_messages_household_id ON messages (household_id, created_at DESC, id DESC);
-- -----------------------------------------------------------------------------
-- the bytes live in object storage under object_key; every attachment is a reference
-- in file_refs, a file lasts as long as something refers to it,
-- kept files, those of the Files page, last until they are deleted
CREATE TABLE files (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  uploader_id UUID NOT NULL REFERENCES users (id),
  object_key TEXT NOT NULL UNIQUE,
  -- a small JPEG of an image, NULL for other files or when it could not be made
  thumb_key TEXT UNIQUE,
  name TEXT NOT NULL,
  mime TEXT NOT NULL,
  size_bytes BIGINT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  kept BOOLEAN NOT NULL DEFAULT FALSE,
  -- target of the composite key in file_refs, so a reference never crosses households
  UNIQUE (id, household_id)
);
CREATE INDEX idx_files_household_id ON files (household_id);
-- -----------------------------------------------------------------------------
-- objects of deleted files and replaced or removed avatars, cascades included,
-- waiting to be removed from object storage
CREATE TABLE file_purges (
  object_key TEXT PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- -----------------------------------------------------------------------------
-- where a file is attached, exactly one of the four
CREATE TABLE file_refs (
  file_id UUID NOT NULL,
  household_id UUID NOT NULL,
  note_id UUID,
  shopping_item_id UUID,
  expense_id UUID,
  message_id UUID,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (file_id, household_id) REFERENCES files (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (note_id, household_id) REFERENCES notes (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (shopping_item_id, household_id) REFERENCES shopping_items (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (expense_id, household_id) REFERENCES expenses (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (message_id, household_id) REFERENCES messages (id, household_id) ON DELETE CASCADE,
  CHECK (
    num_nonnulls(note_id, shopping_item_id, expense_id, message_id) = 1
  ),
  UNIQUE (file_id, note_id),
  UNIQUE (file_id, shopping_item_id),
  UNIQUE (file_id, expense_id),
  UNIQUE (file_id, message_id)
);
CREATE INDEX idx_file_refs_note_id ON file_refs (note_id);
CREATE INDEX idx_file_refs_shopping_item_id ON file_refs (shopping_item_id);
CREATE INDEX idx_file_refs_expense_id ON file_refs (expense_id);
CREATE INDEX idx_file_refs_message_id ON file_refs (message_id);
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
-- -----------------------------------------------------------------------------
-- a row turns that kind of push off for the member in that household, no row means on
CREATE TABLE notification_mutes (
  household_id UUID NOT NULL,
  user_id UUID NOT NULL,
  kind TEXT NOT NULL CHECK (kind IN ('chat', 'money', 'recurring', 'household', 'chores')),
  PRIMARY KEY (household_id, user_id, kind),
  FOREIGN KEY (household_id, user_id) REFERENCES memberships (household_id, user_id) ON DELETE CASCADE
);
-- -----------------------------------------------------------------------------
-- a schedule makes a draft on next_due that the payer adds (through CreateExpense) or skips,
-- either moves next_due one interval on from anchor; ended by a trigger when the payer leaves
CREATE TABLE recurring_expenses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  payer_id UUID NOT NULL,
  category TEXT NOT NULL REFERENCES expense_categories (id) ON UPDATE CASCADE,
  title TEXT NOT NULL,
  -- in the household's currency when the draft is added
  amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
  custom BOOLEAN NOT NULL,
  interval TEXT NOT NULL CHECK (interval IN ('week', 'month', 'year')),
  anchor DATE NOT NULL,
  next_due DATE NOT NULL,
  -- the next_due the payer was pushed about
  notified_on DATE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (id, household_id),
  FOREIGN KEY (household_id, payer_id) REFERENCES memberships (household_id, user_id)
);
CREATE INDEX idx_recurring_expenses_due ON recurring_expenses (next_due);
CREATE TABLE recurring_expense_shares (
  recurring_id UUID NOT NULL,
  household_id UUID NOT NULL,
  user_id UUID NOT NULL,
  -- NULL when the schedule splits equally
  amount_cents BIGINT CHECK (amount_cents >= 0),
  PRIMARY KEY (recurring_id, user_id),
  FOREIGN KEY (recurring_id, household_id) REFERENCES recurring_expenses (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (household_id, user_id) REFERENCES memberships (household_id, user_id)
);
-- -----------------------------------------------------------------------------
-- a chore is due on next_due for assignee_id (NULL: nobody's turn), either on weekdays
-- (bit 0 Monday .. bit 6 Sunday) or every `every` every_unit counted from when it was done;
-- chore_next_assignee says whose turn comes next
CREATE TABLE chores (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  rotation TEXT NOT NULL CHECK (rotation IN ('fixed', 'least_recent', 'manual')),
  every SMALLINT CHECK (every BETWEEN 1 AND 365),
  every_unit TEXT CHECK (every_unit IN ('day', 'week', 'month')),
  weekdays SMALLINT CHECK (weekdays BETWEEN 1 AND 127),
  next_due DATE NOT NULL,
  assignee_id UUID,
  -- the next_due the assignee was pushed about, on the day and the day after
  notified_on DATE,
  nudged_on DATE,
  created_by UUID NOT NULL REFERENCES users (id),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (id, household_id),
  CHECK ((every IS NULL) = (every_unit IS NULL)),
  CHECK ((every IS NULL) <> (weekdays IS NULL)),
  FOREIGN KEY (household_id, assignee_id) REFERENCES memberships (household_id, user_id)
);
CREATE INDEX idx_chores_household_id ON chores (household_id, next_due);
CREATE INDEX idx_chores_due ON chores (next_due);
-- the rota of a chore in position order, a leaver is taken out by end_recurring_of_leaver
CREATE TABLE chore_members (
  chore_id UUID NOT NULL,
  household_id UUID NOT NULL,
  user_id UUID NOT NULL,
  position INT NOT NULL,
  PRIMARY KEY (chore_id, user_id),
  FOREIGN KEY (chore_id, household_id) REFERENCES chores (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (household_id, user_id) REFERENCES memberships (household_id, user_id)
);
-- what happened to each due date: user_id did, skipped or passed it on assignee_id's turn;
-- the last done per member orders a least_recent rota
--
-- shortcut: never purged, rows go with their chore; purge old ones in RunCleanup if it grows
CREATE TABLE chore_log (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  chore_id UUID NOT NULL,
  household_id UUID NOT NULL,
  kind TEXT NOT NULL CHECK (kind IN ('done', 'skipped', 'passed')),
  user_id UUID NOT NULL REFERENCES users (id),
  assignee_id UUID REFERENCES users (id),
  due DATE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (chore_id, household_id) REFERENCES chores (id, household_id) ON DELETE CASCADE
);
CREATE INDEX idx_chore_log_chore_id ON chore_log (chore_id, created_at DESC);
-- -----------------------------------------------------------------------------
-- who did what, written by the services in the same transaction, kept per grp for
-- g.ActivityRetention; title and amount are copies so the row outlives its subject
CREATE TABLE activity (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  actor_id UUID REFERENCES users (id) ON DELETE SET NULL,
  grp TEXT NOT NULL CHECK (grp IN ('money', 'shopping', 'notes', 'members', 'chores')),
  kind TEXT NOT NULL,
  subject_id UUID,
  target_user_id UUID REFERENCES users (id) ON DELETE SET NULL,
  title TEXT NOT NULL DEFAULT '',
  amount_cents BIGINT,
  currency TEXT,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_activity_household ON activity (household_id, created_at DESC);
CREATE INDEX idx_activity_created ON activity (created_at);
-- -----------------------------------------------------------------------------
-- http_sessions belongs to the session store (pgstore), which creates it on start, after the
-- migrations, so it is created here the same way, only user_id is ours: sessions.Login sets it,
-- a password reset deletes the user's sessions by it
CREATE TABLE IF NOT EXISTS http_sessions (
  id BIGSERIAL PRIMARY KEY,
  key BYTEA,
  data BYTEA,
  created_on TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  modified_on TIMESTAMPTZ,
  expires_on TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS http_sessions_expiry_idx ON http_sessions (expires_on);
CREATE INDEX IF NOT EXISTS http_sessions_key_idx ON http_sessions (key);
ALTER TABLE http_sessions
ADD COLUMN user_id uuid REFERENCES users (id) ON DELETE CASCADE;
CREATE INDEX http_sessions_user_id_idx ON http_sessions (user_id);
-- -----------------------------------------------------------------------------
-- --- FUNCTIONS ---
CREATE FUNCTION queue_file_purge() RETURNS TRIGGER AS $$ BEGIN
INSERT INTO file_purges (object_key)
VALUES (OLD.object_key) ON CONFLICT DO NOTHING;
IF OLD.thumb_key IS NOT NULL THEN
INSERT INTO file_purges (object_key)
VALUES (OLD.thumb_key) ON CONFLICT DO NOTHING;
END IF;
RETURN OLD;
END;
$$ LANGUAGE plpgsql;
CREATE FUNCTION queue_avatar_purge() RETURNS TRIGGER AS $$ BEGIN
INSERT INTO file_purges (object_key)
VALUES (OLD.avatar_key) ON CONFLICT DO NOTHING;
RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- whose turn a chore is after after's: fixed is the next in the rota (wrapping round, the first
-- when after is not in it), least_recent who did it longest ago (never first, then rota order),
-- manual stays with after; exclude_after leaves after out (a pass, a leaver), NULL when nobody is left
CREATE FUNCTION chore_next_assignee(chore UUID, after UUID, exclude_after BOOLEAN) RETURNS UUID AS $$
SELECT CASE
    c.rotation
    WHEN 'manual' THEN CASE
      WHEN exclude_after THEN NULL
      ELSE after
    END
    WHEN 'fixed' THEN (
      SELECT m.user_id
      FROM chore_members m
      WHERE m.chore_id = c.id
        AND NOT (
          exclude_after
          AND m.user_id IS NOT DISTINCT FROM after
        )
      ORDER BY m.position <= COALESCE(
          (
            SELECT p.position
            FROM chore_members p
            WHERE p.chore_id = c.id
              AND p.user_id = after
          ),
          -1
        ),
        m.position
      LIMIT 1
    )
    ELSE (
      SELECT m.user_id
      FROM chore_members m
      WHERE m.chore_id = c.id
        AND NOT (
          exclude_after
          AND m.user_id IS NOT DISTINCT FROM after
        )
      ORDER BY (
          SELECT max(l.created_at)
          FROM chore_log l
          WHERE l.chore_id = c.id
            AND l.user_id = m.user_id
            AND l.kind = 'done'
        ) NULLS FIRST,
        m.position
      LIMIT 1
    )
  END
FROM chores c
WHERE c.id = chore;
$$ LANGUAGE sql STABLE;
-- a leaver's schedules, and their shares in the others' schedules, end with the membership;
-- their chore turns move on and they leave every rota
CREATE FUNCTION end_recurring_of_leaver() RETURNS TRIGGER AS $$ BEGIN
DELETE FROM recurring_expenses
WHERE household_id = NEW.household_id
  AND payer_id = NEW.user_id;
DELETE FROM recurring_expense_shares
WHERE household_id = NEW.household_id
  AND user_id = NEW.user_id;
UPDATE chores
SET assignee_id = chore_next_assignee(id, NEW.user_id, TRUE),
  notified_on = NULL,
  nudged_on = NULL
WHERE household_id = NEW.household_id
  AND assignee_id = NEW.user_id;
DELETE FROM chore_members
WHERE household_id = NEW.household_id
  AND user_id = NEW.user_id;
RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- live updates: '<household_id> <topic>' on household_events, sent on commit,
-- the server tells open pages of that household to refetch
CREATE FUNCTION notify_household() RETURNS TRIGGER AS $$
DECLARE r RECORD;
BEGIN
IF TG_OP = 'DELETE' THEN r := OLD;
ELSE r := NEW;
END IF;
PERFORM pg_notify(
  'household_events',
  r.household_id::text || ' ' || TG_ARGV [0]
);
RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- expense_shares has no household_id, in a cascade the expense is already gone
-- and its own delete notifies
CREATE FUNCTION notify_household_expense_share() RETURNS TRIGGER AS $$
DECLARE hid UUID;
BEGIN
SELECT household_id INTO hid
FROM expenses
WHERE id = CASE
    WHEN TG_OP = 'DELETE' THEN OLD.expense_id
    ELSE NEW.expense_id
  END;
IF hid IS NOT NULL THEN PERFORM pg_notify('household_events', hid::text || ' expenses');
END IF;
RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- a file only on private notes is for their authors and its uploader, a kept file is for
-- everyone wherever it is attached, NULL as the user asks whether everyone in the household sees it
CREATE FUNCTION file_visible(fid UUID, uid UUID) RETURNS BOOLEAN AS $$
SELECT EXISTS (
    SELECT 1
    FROM files f
    WHERE f.id = fid
      AND (
        f.kept
        OR f.uploader_id = uid
      )
  )
  OR NOT EXISTS (
    SELECT 1
    FROM file_refs r
    WHERE r.file_id = fid
  )
  OR EXISTS (
    SELECT 1
    FROM file_refs r
      LEFT JOIN notes n ON n.id = r.note_id
    WHERE r.file_id = fid
      AND (
        n.id IS NULL
        OR n.visibility = 'household'
        OR n.author_id = uid
      )
  );
$$ LANGUAGE sql STABLE;
-- also runs for the cascades, so a deleted note, item, expense or message takes its files along
CREATE FUNCTION delete_unreferenced_file() RETURNS TRIGGER AS $$ BEGIN
DELETE FROM files f
WHERE f.id = OLD.file_id
  AND NOT f.kept
  AND NOT EXISTS (
    SELECT 1
    FROM file_refs r
    WHERE r.file_id = f.id
  );
RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- -----------------------------------------------------------------------------
-- --- TRIGGERS ---
CREATE TRIGGER mdt_users BEFORE
UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
CREATE TRIGGER mdt_households BEFORE
UPDATE ON households FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
CREATE TRIGGER mdt_shopping_lists BEFORE
UPDATE ON shopping_lists FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
CREATE TRIGGER mdt_shopping_items BEFORE
UPDATE ON shopping_items FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
CREATE TRIGGER mdt_notes BEFORE
UPDATE ON notes FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
CREATE TRIGGER queue_file_purge
AFTER DELETE ON files FOR EACH ROW EXECUTE FUNCTION queue_file_purge();
CREATE TRIGGER queue_avatar_purge
AFTER
UPDATE OF avatar_key ON users FOR EACH ROW
  WHEN (
    OLD.avatar_key IS NOT NULL
    AND OLD.avatar_key IS DISTINCT FROM NEW.avatar_key
  ) EXECUTE FUNCTION queue_avatar_purge();
CREATE TRIGGER queue_avatar_purge_delete
AFTER DELETE ON users FOR EACH ROW
  WHEN (OLD.avatar_key IS NOT NULL) EXECUTE FUNCTION queue_avatar_purge();
CREATE TRIGGER delete_unreferenced_file
AFTER DELETE ON file_refs FOR EACH ROW EXECUTE FUNCTION delete_unreferenced_file();
CREATE TRIGGER notify_shopping_lists
AFTER INSERT OR UPDATE OR DELETE ON shopping_lists FOR EACH ROW EXECUTE FUNCTION notify_household('shopping');
CREATE TRIGGER notify_shopping_items
AFTER INSERT OR UPDATE OR DELETE ON shopping_items FOR EACH ROW EXECUTE FUNCTION notify_household('shopping');
CREATE TRIGGER notify_notes
AFTER INSERT OR UPDATE OR DELETE ON notes FOR EACH ROW EXECUTE FUNCTION notify_household('notes');
CREATE TRIGGER notify_files
AFTER INSERT OR UPDATE OR DELETE ON files FOR EACH ROW EXECUTE FUNCTION notify_household('files');
CREATE TRIGGER notify_file_refs
AFTER INSERT OR UPDATE OR DELETE ON file_refs FOR EACH ROW EXECUTE FUNCTION notify_household('files');
CREATE TRIGGER notify_expenses
AFTER INSERT OR UPDATE OR DELETE ON expenses FOR EACH ROW EXECUTE FUNCTION notify_household('expenses');
CREATE TRIGGER notify_expense_shares
AFTER INSERT OR UPDATE OR DELETE ON expense_shares FOR EACH ROW EXECUTE FUNCTION notify_household_expense_share();
CREATE TRIGGER notify_settlements
AFTER INSERT OR UPDATE OR DELETE ON settlements FOR EACH ROW EXECUTE FUNCTION notify_household('expenses');
-- reading the chat moves chat_read_at, which is no news for others
CREATE TRIGGER notify_memberships
AFTER INSERT OR DELETE OR UPDATE OF role, left_at ON memberships FOR EACH ROW EXECUTE FUNCTION notify_household('people');
CREATE TRIGGER notify_invites
AFTER INSERT OR UPDATE OR DELETE ON invites FOR EACH ROW EXECUTE FUNCTION notify_household('people');
CREATE TRIGGER notify_messages
AFTER INSERT OR UPDATE OR DELETE ON messages FOR EACH ROW EXECUTE FUNCTION notify_household('chat');
CREATE TRIGGER notify_recurring_expenses
AFTER INSERT OR UPDATE OR DELETE ON recurring_expenses FOR EACH ROW EXECUTE FUNCTION notify_household('expenses');
CREATE TRIGGER notify_chores
AFTER INSERT OR UPDATE OR DELETE ON chores FOR EACH ROW EXECUTE FUNCTION notify_household('chores');
CREATE TRIGGER notify_chore_members
AFTER INSERT OR UPDATE OR DELETE ON chore_members FOR EACH ROW EXECUTE FUNCTION notify_household('chores');
CREATE TRIGGER end_recurring_of_leaver
AFTER UPDATE OF left_at ON memberships FOR EACH ROW
  WHEN (NEW.left_at IS NOT NULL) EXECUTE FUNCTION end_recurring_of_leaver();
