-- Add email authentication fields
ALTER TABLE account ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);
ALTER TABLE account ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT false;
ALTER TABLE account ADD COLUMN IF NOT EXISTS verification_token VARCHAR(255);
ALTER TABLE account ADD COLUMN IF NOT EXISTS reset_token VARCHAR(255);
ALTER TABLE account ADD COLUMN IF NOT EXISTS reset_token_expires TIMESTAMP;

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_account_email_verified ON account(email_verified);
CREATE INDEX IF NOT EXISTS idx_account_verification_token ON account(verification_token);
CREATE INDEX IF NOT EXISTS idx_account_reset_token ON account(reset_token);