-- -----------------------------------------------------------------------------
CREATE TABLE shopping_lists (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  created_by UUID NOT NULL REFERENCES users (id),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  -- target of the composite key in shopping_items
  UNIQUE (id, household_id)
);
CREATE INDEX idx_shopping_lists_household_id ON shopping_lists (household_id);
CREATE TRIGGER mdt_shopping_lists BEFORE
UPDATE ON shopping_lists FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
-- -----------------------------------------------------------------------------
-- (list_id, household_id) references the list, so an item can never sit
-- in a list of another household
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
  FOREIGN KEY (list_id, household_id) REFERENCES shopping_lists (id, household_id) ON DELETE CASCADE
);
CREATE INDEX idx_shopping_items_list_id ON shopping_items (list_id);
CREATE TRIGGER mdt_shopping_items BEFORE
UPDATE ON shopping_items FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
