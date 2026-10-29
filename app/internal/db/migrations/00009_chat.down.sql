DROP TRIGGER notify_messages ON messages;
DROP TRIGGER notify_invites ON invites;
DROP TRIGGER notify_memberships ON memberships;
DROP TRIGGER notify_settlements ON settlements;
DROP TRIGGER notify_expense_shares ON expense_shares;
DROP TRIGGER notify_expenses ON expenses;
DROP TRIGGER notify_files ON files;
DROP TRIGGER notify_notes ON notes;
DROP TRIGGER notify_shopping_items ON shopping_items;
DROP TRIGGER notify_shopping_lists ON shopping_lists;
DROP FUNCTION notify_household_expense_share();
DROP FUNCTION notify_household();
DELETE FROM files
WHERE message_id IS NOT NULL;
ALTER TABLE files DROP CONSTRAINT files_check,
  DROP COLUMN message_id,
  ADD CONSTRAINT files_check CHECK (
    num_nonnulls(note_id, shopping_item_id, expense_id) <= 1
  );
ALTER TABLE memberships DROP COLUMN chat_read_at;
DROP TABLE messages;
