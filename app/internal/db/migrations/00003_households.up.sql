-- -----------------------------------------------------------------------------
CREATE TABLE households (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  -- ISO 4217, amounts are stored as integer minor units (cents)
  currency TEXT NOT NULL DEFAULT 'EUR',
  created_by UUID NOT NULL REFERENCES users (id),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER mdt_households BEFORE
UPDATE ON households FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
-- -----------------------------------------------------------------------------
-- a row is kept after leaving (left_at set) so history still shows former members,
-- rejoining clears left_at
CREATE TABLE memberships (
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users (id),
  role TEXT NOT NULL CHECK (role IN ('owner', 'member')),
  joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  left_at TIMESTAMP WITH TIME ZONE,
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
