-- Restore NOT NULL constraint on google_id
ALTER TABLE account ALTER COLUMN google_id SET NOT NULL; 