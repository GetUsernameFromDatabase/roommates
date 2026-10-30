-- -----------------------------------------------------------------------------
-- target of the composite key below, so a reference never crosses households
ALTER TABLE files
ADD UNIQUE (id, household_id);
-- a chat message pointing at a file that already exists, unlike files.message_id
-- the file stays when the message goes, and the reference goes with the file
CREATE TABLE message_file_refs (
  message_id UUID NOT NULL,
  file_id UUID NOT NULL,
  household_id UUID NOT NULL,
  PRIMARY KEY (message_id, file_id),
  FOREIGN KEY (message_id, household_id) REFERENCES messages (id, household_id) ON DELETE CASCADE,
  FOREIGN KEY (file_id, household_id) REFERENCES files (id, household_id) ON DELETE CASCADE
);
CREATE INDEX idx_message_file_refs_file_id ON message_file_refs (file_id);
