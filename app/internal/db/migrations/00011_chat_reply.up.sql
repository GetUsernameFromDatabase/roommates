-- -----------------------------------------------------------------------------
-- a message can answer another message of the same household,
-- deleting the answered message keeps the reply without its quote
ALTER TABLE messages
ADD COLUMN reply_to_id UUID,
ADD FOREIGN KEY (reply_to_id, household_id) REFERENCES messages (id, household_id) ON DELETE SET NULL (reply_to_id);
