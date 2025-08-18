-- Change due_date to due_day for recurring subscriptions
-- This makes more sense since subscriptions recur monthly on the same day

-- First, add the new due_day column
ALTER TABLE subscription_details 
ADD COLUMN due_day INT CHECK (due_day >= 1 AND due_day <= 31);

-- Update existing due_date values to extract the day
UPDATE subscription_details 
SET due_day = EXTRACT(DAY FROM due_date) 
WHERE due_date IS NOT NULL;

-- Drop the old due_date column
ALTER TABLE subscription_details 
DROP COLUMN due_date;

-- Add index on due_day for better performance
CREATE INDEX idx_subscription_details_due_day ON subscription_details(due_day);

-- Update the existing index to use due_day instead
DROP INDEX IF EXISTS idx_subscription_details_due_date; 