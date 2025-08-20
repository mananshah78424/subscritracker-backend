-- Reverse the nullable changes for optional time fields
-- This will make the fields NOT NULL again (though this might cause issues with existing data)

-- Make end_date NOT NULL
ALTER TABLE subscription_details 
ALTER COLUMN end_date SET NOT NULL;

-- Make start_time NOT NULL
ALTER TABLE subscription_details 
ALTER COLUMN start_time SET NOT NULL;

-- Make due_time NOT NULL
ALTER TABLE subscription_details 
ALTER COLUMN due_time SET NOT NULL;

-- Make reminder_date NOT NULL
ALTER TABLE subscription_details 
ALTER COLUMN reminder_date SET NOT NULL;

-- Make reminder_time NOT NULL
ALTER TABLE subscription_details 
ALTER COLUMN reminder_time SET NOT NULL; 