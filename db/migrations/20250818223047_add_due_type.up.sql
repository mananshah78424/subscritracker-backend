-- Add due_type column to support different subscription frequencies
ALTER TABLE subscription_details 
ADD COLUMN due_type VARCHAR(255) NOT NULL DEFAULT 'monthly' CHECK (due_type IN ('monthly', 'yearly', 'custom', 'weekly', 'daily'));



