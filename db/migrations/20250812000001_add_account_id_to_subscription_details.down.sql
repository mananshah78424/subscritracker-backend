-- Remove foreign key constraint
ALTER TABLE subscription_details 
DROP CONSTRAINT IF EXISTS fk_subscription_details_account;

-- Remove unique constraint
ALTER TABLE subscription_details 
DROP CONSTRAINT IF EXISTS unique_account_subscription_channel;

-- Remove account_id column
ALTER TABLE subscription_details 
DROP COLUMN IF EXISTS account_id; 