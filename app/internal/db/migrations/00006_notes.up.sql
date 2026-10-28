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
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_notes_household_id ON notes (household_id);
CREATE TRIGGER mdt_notes BEFORE
UPDATE ON notes FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
