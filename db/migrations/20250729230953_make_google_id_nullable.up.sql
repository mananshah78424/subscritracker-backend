-- Make google_id nullable to support email signup without Google OAuth
ALTER TABLE account ALTER COLUMN google_id DROP NOT NULL; 