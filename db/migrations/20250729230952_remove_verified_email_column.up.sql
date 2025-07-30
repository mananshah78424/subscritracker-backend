-- Remove the redundant verified_email column since we're using email_verified
ALTER TABLE account DROP COLUMN IF EXISTS verified_email; 