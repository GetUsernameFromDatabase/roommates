DROP TRIGGER IF EXISTS mdt_messages ON messages;
DROP TABLE IF EXISTS messages;
--
DROP TABLE IF EXISTS conversations;
DROP TYPE IF EXISTS conversation_recipient_type;
--
DROP TRIGGER IF EXISTS mdt_house_payments_files ON house_payments_files;
DROP INDEX IF EXISTS idxh_house_payments_files_payment;
DROP TABLE IF EXISTS house_payments_files;
--
DROP TABLE IF EXISTS house_payment_payers;
DROP TYPE IF EXISTS house_payment_status;
--
DROP TRIGGER IF EXISTS mdt_house_payments ON house_payments;
DROP INDEX IF EXISTS idxh_house_payments_requester;
DROP TABLE IF EXISTS house_payments;
--
DROP TRIGGER IF EXISTS mdt_house_notes ON house_notes;
DROP INDEX IF EXISTS idxh_house_notes_maker_id;
DROP INDEX IF EXISTS idxh_house_notes_house_id;
DROP TABLE IF EXISTS house_notes;
--
DROP TRIGGER IF EXISTS mdt_house_reminders ON house_reminders;
DROP TABLE IF EXISTS house_reminders;
DROP TYPE IF EXISTS house_reminder_status;
--
DROP TABLE IF EXISTS user_contact_information;
--
DROP TABLE IF EXISTS user_houses;
--
DROP TRIGGER IF EXISTS mdt_houses ON houses;
DROP INDEX IF EXISTS idxh_houses_maker_id;
DROP TABLE IF EXISTS houses;
--
DROP TRIGGER IF EXISTS mdt_users ON users;
DROP INDEX IF EXISTS idxh_users_email;
DROP TABLE IF EXISTS users;
