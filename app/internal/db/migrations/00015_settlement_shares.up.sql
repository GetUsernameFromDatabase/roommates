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
