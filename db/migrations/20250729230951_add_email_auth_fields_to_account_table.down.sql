ALTER TABLE account DROP COLUMN IF EXISTS password_hash;
ALTER TABLE account DROP COLUMN IF EXISTS email_verified;
ALTER TABLE account DROP COLUMN IF EXISTS verification_token;
ALTER TABLE account DROP COLUMN IF EXISTS reset_token;
ALTER TABLE account DROP COLUMN IF EXISTS reset_token_expires;
DROP INDEX IF EXISTS idx_account_email_verified;
DROP INDEX IF EXISTS idx_account_verification_token;
DROP INDEX IF EXISTS idx_account_reset_token;