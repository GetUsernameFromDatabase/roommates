-- -----------------------------------------------------------------------------
-- one chat per household, a message has a body, files or both
CREATE TABLE messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  author_id UUID NOT NULL REFERENCES users (id),
  body TEXT NOT NULL CHECK (length(body) <= 4000),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  -- target of the composite key in files
  UNIQUE (id, household_id)
);
CREATE INDEX idx_messages_household_id ON messages (household_id, created_at DESC, id DESC);
-- -----------------------------------------------------------------------------
-- messages of others after it are unread, a new member starts with nothing unread
ALTER TABLE memberships
ADD COLUMN chat_read_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;
-- -----------------------------------------------------------------------------
-- photos and files sent in the chat
ALTER TABLE files
ADD COLUMN message_id UUID,
ADD FOREIGN KEY (message_id, household_id) REFERENCES messages (id, household_id) ON DELETE CASCADE,
DROP CONSTRAINT files_check,
ADD CONSTRAINT files_check CHECK (
    num_nonnulls(note_id, shopping_item_id, expense_id, message_id) <= 1
  );
CREATE INDEX idx_files_message_id ON files (message_id);
-- -----------------------------------------------------------------------------
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
CREATE TRIGGER notify_shopping_lists
AFTER INSERT OR UPDATE OR DELETE ON shopping_lists FOR EACH ROW EXECUTE FUNCTION notify_household('shopping');
CREATE TRIGGER notify_shopping_items
AFTER INSERT OR UPDATE OR DELETE ON shopping_items FOR EACH ROW EXECUTE FUNCTION notify_household('shopping');
CREATE TRIGGER notify_notes
AFTER INSERT OR UPDATE OR DELETE ON notes FOR EACH ROW EXECUTE FUNCTION notify_household('notes');
CREATE TRIGGER notify_files
AFTER INSERT OR UPDATE OR DELETE ON files FOR EACH ROW EXECUTE FUNCTION notify_household('files');
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
