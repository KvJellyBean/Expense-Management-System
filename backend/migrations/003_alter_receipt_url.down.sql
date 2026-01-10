-- Revert receipt_url column to VARCHAR(512)
ALTER TABLE expenses ALTER COLUMN receipt_url TYPE VARCHAR(512);
