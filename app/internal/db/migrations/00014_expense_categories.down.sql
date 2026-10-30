ALTER TABLE shopping_lists DROP CONSTRAINT shopping_lists_category_fkey,
  ADD CONSTRAINT shopping_lists_category_check CHECK (
    category IN ('groceries', 'utilities', 'rent', 'household', 'other')
  );
ALTER TABLE expenses DROP CONSTRAINT expenses_category_fkey,
  ADD CONSTRAINT expenses_category_check CHECK (
    category IN ('groceries', 'utilities', 'rent', 'household', 'other')
  );
DROP TABLE expense_categories;
