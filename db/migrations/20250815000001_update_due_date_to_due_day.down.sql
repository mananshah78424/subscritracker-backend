-- Revert due_day back to due_date

-- Add the old due_date column back
ALTER TABLE subscription_details 
ADD COLUMN due_date DATE;

-- Convert due_day back to due_date (using current month as example)
UPDATE subscription_details 
SET due_date = CURRENT_DATE + (due_day - EXTRACT(DAY FROM CURRENT_DATE)) * INTERVAL '1 day'
WHERE due_day IS NOT NULL;

-- Drop the due_day column
ALTER TABLE subscription_details 
DROP COLUMN due_day;

-- Recreate the old index
CREATE INDEX idx_subscription_details_due_date ON subscription_details(due_date);

-- Drop the due_day index
DROP INDEX IF EXISTS idx_subscription_details_due_day; 