-- the table itself stays, it belongs to the session store
ALTER TABLE http_sessions DROP COLUMN user_id;
-- -----------------------------------------------------------------------------
-- triggers go with their tables
DROP TABLE activity,
recurring_expense_shares,
recurring_expenses,
notification_mutes,
push_subscriptions,
file_refs,
file_purges,
files,
messages,
settlement_shares,
settlements,
expense_shares,
notes,
shopping_items,
expenses,
shopping_lists,
expense_categories,
email_tokens,
invites,
memberships,
households,
users;
-- -----------------------------------------------------------------------------
DROP FUNCTION end_recurring_of_leaver(),
delete_unreferenced_file(),
file_visible(UUID, UUID),
notify_household_expense_share(),
notify_household(),
queue_avatar_purge(),
queue_file_purge();
