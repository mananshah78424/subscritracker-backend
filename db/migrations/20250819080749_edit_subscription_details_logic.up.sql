-- # Change due_date to next_due_date that is calcuated based on start_date, due_day, and due_type
-- # Add column due_type and column due_day_of_month
-- # Add column end_date that is calculated based on start_date, due_day, and due_type

ALTER TABLE subscription_details
ADD COLUMN next_due_date DATE;

ALTER TABLE subscription_details
ADD COLUMN due_type VARCHAR(255) DEFAULT 'monthly' CHECK (due_type IN ('monthly', 'yearly', 'weekly', 'daily'));

ALTER TABLE subscription_details
ADD COLUMN due_day_of_month INT;

ALTER TABLE subscription_details
ADD COLUMN end_date DATE;

-- # Drop due_date column
ALTER TABLE subscription_details
DROP COLUMN due_date;



