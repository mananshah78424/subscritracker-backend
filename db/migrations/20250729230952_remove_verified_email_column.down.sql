-- Restore the verified_email column
ALTER TABLE account ADD COLUMN verified_email BOOLEAN DEFAULT false; 