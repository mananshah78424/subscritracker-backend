-- Add account_id column to subscription_details table
ALTER TABLE subscription_details 
ADD COLUMN account_id INT NOT NULL DEFAULT 1;

-- Add foreign key constraint
ALTER TABLE subscription_details 
ADD CONSTRAINT fk_subscription_details_account 
FOREIGN KEY (account_id) REFERENCES account(id);

-- Update the unique constraint to include account_id
-- This allows one subscription per channel per account
ALTER TABLE subscription_details 
DROP CONSTRAINT IF EXISTS unique_subscription_channel;

ALTER TABLE subscription_details 
ADD CONSTRAINT unique_account_subscription_channel 
UNIQUE (account_id, subscription_channel_id); 