-- -----------------------------------------------------------------------------
-- the categories an expense or shopping list can have, config.ExpenseCategories must match
CREATE TABLE expense_categories (id TEXT PRIMARY KEY);
INSERT INTO expense_categories (id)
VALUES ('groceries'),
  ('utilities'),
  ('rent'),
  ('household'),
  ('other');
ALTER TABLE expenses DROP CONSTRAINT expenses_category_check,
  ADD FOREIGN KEY (category) REFERENCES expense_categories (id) ON UPDATE CASCADE;
ALTER TABLE shopping_lists DROP CONSTRAINT shopping_lists_category_check,
  ADD FOREIGN KEY (category) REFERENCES expense_categories (id) ON UPDATE CASCADE;
