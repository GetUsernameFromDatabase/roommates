DROP TABLE files;
DROP FUNCTION queue_file_purge;
DROP TABLE file_purges;
ALTER TABLE shopping_items DROP CONSTRAINT shopping_items_id_household_id_key;
ALTER TABLE notes DROP CONSTRAINT notes_id_household_id_key;
