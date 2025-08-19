
ALTER TABLE subscription_details
ADD COLUMN due_date DATE;

ALTER TABLE subscription_details
DROP COLUMN next_due_date;

ALTER TABLE subscription_details
DROP COLUMN due_type;

ALTER TABLE subscription_details
DROP COLUMN due_day_of_month;

ALTER TABLE subscription_details
DROP COLUMN end_date;