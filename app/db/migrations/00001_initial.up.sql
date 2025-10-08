-- --- EXTENSIONS ---
CREATE EXTENSION IF NOT EXISTS moddatetime;
-- --- the core ---
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT NOT NULL UNIQUE,
  email_validated BOOLEAN NOT NULL DEFAULT FALSE,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  full_name TEXT,
  is_full_name_public BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER mdt_users BEFORE
UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
-- if in the future switch is made from like ex% to %ex% then this index won't do
CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idxh_users_email ON users USING HASH (email);
--
CREATE TABLE houses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  maker_id UUID REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idxh_houses_maker_id ON houses USING HASH (maker_id);
CREATE TRIGGER mdt_houses BEFORE
UPDATE ON houses FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
-- --- user section ---
CREATE TABLE user_houses (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  house_id UUID NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, house_id)
);
--
CREATE TABLE user_contact_information (
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  contact_information JSONB
);
-- --- house section ---
CREATE TYPE house_reminder_status AS ENUM ('canceled', 'in-progress', 'complete');
CREATE TABLE house_reminders (
  id SERIAL PRIMARY KEY,
  content JSONB NOT NULL,
  reminder_status house_reminder_status NOT NULL DEFAULT 'in-progress',
  house_id UUID NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
  maker_id UUID REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER mdt_house_reminders BEFORE
UPDATE ON house_reminders FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
--
CREATE TABLE house_notes (
  id SERIAL PRIMARY KEY,
  title text NOT NULL,
  content TEXT NOT NULL,
  house_id UUID NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
  maker_id UUID REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idxh_house_notes_maker_id ON house_notes USING HASH (maker_id);
CREATE INDEX idxh_house_notes_house_id ON house_notes USING HASH (house_id);
CREATE TRIGGER mdt_house_notes BEFORE
UPDATE ON house_notes FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
--
CREATE TABLE house_payments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  payment_name TEXT NOT NULL,
  -- euros (€) -- expected use case
  amount MONEY,
  requester_id UUID REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idxh_house_payments_requester ON house_payments USING HASH (requester_id);
CREATE TRIGGER mdt_house_payments BEFORE
UPDATE ON house_payments FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
--
CREATE TYPE house_payment_status AS ENUM ('done', 'incomplete');
CREATE TABLE house_payment_payers (
  payment_id UUID NOT NULL REFERENCES house_payments(id) ON DELETE CASCADE,
  payer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  payment_status house_payment_status NOT NULL DEFAULT 'incomplete',
  PRIMARY KEY (payment_id, payer_id)
);
--
-- TODO: nonono to storing files in the database minio
CREATE TABLE house_payments_files (
  payment_id UUID NOT NULL REFERENCES house_payments(id) ON DELETE CASCADE,
  file_name TEXT NOT NULL,
  file_data BYTEA NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idxh_house_payments_files_payment ON house_payments_files USING HASH (payment_id);
CREATE TRIGGER mdt_house_payments_files BEFORE
UPDATE ON house_payments_files FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);
-- --- messaging section ---
CREATE TYPE conversation_recipient_type AS ENUM ('house', 'direct', 'group');
CREATE TABLE conversations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT,
  conversation_image BYTEA,
  recipient_ids TEXT [] NOT NULL,
  recipient_type conversation_recipient_type NOT NULL
);
--
CREATE TABLE messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  content TEXT NOT NULL,
  conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  sender_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER mdt_messages BEFORE
UPDATE ON messages FOR EACH ROW EXECUTE PROCEDURE moddatetime (updated_at);