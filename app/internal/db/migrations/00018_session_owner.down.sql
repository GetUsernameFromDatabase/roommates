-- the table itself stays, it belongs to the session store
ALTER TABLE http_sessions DROP COLUMN user_id;
