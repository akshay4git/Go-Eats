DO $$
BEGIN 
    IF EXITS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='restaurants' AND column_name='photo'
    ) THEN
        ALTER TABLE restaurants RENAME COLUMN photo TO store_image;
    END IF;
END $$;