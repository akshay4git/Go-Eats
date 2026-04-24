-- Fix restaurants table (already correct from migration 1)
-- Fix any data inconsistency from old store_image column if it exists
DO $$ 
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns 
    WHERE table_name='restaurants' AND column_name='store_image'
  ) THEN
    ALTER TABLE restaurants RENAME COLUMN store_image TO photo;
  END IF;
END $$;