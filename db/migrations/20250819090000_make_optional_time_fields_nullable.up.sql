-- Make optional time fields nullable to allow NULL values instead of default time values
-- This will prevent fields like reminder_date, reminder_time, start_time, due_time, end_date
-- from showing as "0001-01-01T00:00:00Z" when they are not provided

-- First, let's check if these columns exist and their current state
-- Then make them properly nullable and remove any default values

-- Make end_date nullable (if it exists)
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'subscription_details' AND column_name = 'end_date') THEN
        ALTER TABLE subscription_details ALTER COLUMN end_date DROP NOT NULL;
        ALTER TABLE subscription_details ALTER COLUMN end_date DROP DEFAULT;
    END IF;
END $$;

-- Make start_time nullable (if it exists)
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'subscription_details' AND column_name = 'start_time') THEN
        ALTER TABLE subscription_details ALTER COLUMN start_time DROP NOT NULL;
        ALTER TABLE subscription_details ALTER COLUMN start_time DROP DEFAULT;
    END IF;
END $$;

-- Make due_time nullable (if it exists)
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'subscription_details' AND column_name = 'due_time') THEN
        ALTER TABLE subscription_details ALTER COLUMN due_time DROP NOT NULL;
        ALTER TABLE subscription_details ALTER COLUMN due_time DROP DEFAULT;
    END IF;
END $$;

-- Make reminder_date nullable (if it exists)
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'subscription_details' AND column_name = 'reminder_date') THEN
        ALTER TABLE subscription_details ALTER COLUMN reminder_date DROP NOT NULL;
        ALTER TABLE subscription_details ALTER COLUMN reminder_date DROP DEFAULT;
    END IF;
END $$;

-- Make reminder_time nullable (if it exists)
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'subscription_details' AND column_name = 'reminder_time') THEN
        ALTER TABLE subscription_details ALTER COLUMN reminder_time DROP NOT NULL;
        ALTER TABLE subscription_details ALTER COLUMN reminder_time DROP DEFAULT;
    END IF;
END $$;

-- Update any existing records that have the default time value to NULL
-- This will clean up existing data that has the problematic default values
UPDATE subscription_details 
SET start_time = NULL 
WHERE start_time = '0001-01-01 00:00:00' OR start_time = '0001-01-01 00:00:00+00';

UPDATE subscription_details 
SET due_time = NULL 
WHERE due_time = '0001-01-01 00:00:00' OR due_time = '0001-01-01 00:00:00+00';

UPDATE subscription_details 
SET reminder_time = NULL 
WHERE reminder_time = '0001-01-01 00:00:00' OR reminder_time = '0001-01-01 00:00:00+00';

UPDATE subscription_details 
SET end_date = NULL 
WHERE end_date = '0001-01-01' OR end_date = '0001-01-01 00:00:00+00';

UPDATE subscription_details 
SET reminder_date = NULL 
WHERE reminder_date = '0001-01-01' OR reminder_date = '0001-01-01 00:00:00+00'; 