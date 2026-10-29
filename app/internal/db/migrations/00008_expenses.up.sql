-- -----------------------------------------------------------------------------
-- whoever enters an expense paid it, so payer_id is also its creator
CREATE TABLE expenses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  payer_id UUID NOT NULL REFERENCES users (id),
  category TEXT NOT NULL CHECK (
    category IN ('groceries', 'utilities', 'rent', 'household', 'other')
  ),
  title TEXT NOT NULL,
  -- in the household currency's minor units
  amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
  incurred_on DATE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  -- target of the composite key in files
  UNIQUE (id, household_id)
);
CREATE INDEX idx_expenses_household_id ON expenses (household_id, incurred_on DESC);
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
-- "from_user paid to_user", counts towards balances once to_user confirms it
CREATE TABLE settlements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
  from_user UUID NOT NULL REFERENCES users (id),
  to_user UUID NOT NULL REFERENCES users (id),
  amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  confirmed_at TIMESTAMP WITH TIME ZONE,
  CHECK (from_user <> to_user)
);
CREATE INDEX idx_settlements_household_id ON settlements (household_id);
-- -----------------------------------------------------------------------------
-- receipts and bills attach to an expense
ALTER TABLE files
ADD COLUMN expense_id UUID,
ADD FOREIGN KEY (expense_id, household_id) REFERENCES expenses (id, household_id) ON DELETE CASCADE,
DROP CONSTRAINT files_check,
ADD CONSTRAINT files_check CHECK (
    num_nonnulls(note_id, shopping_item_id, expense_id) <= 1
  );
CREATE INDEX idx_files_expense_id ON files (expense_id);
