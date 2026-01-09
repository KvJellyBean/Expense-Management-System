-- Alter receipt_url column to support data URLs (base64)
ALTER TABLE expenses ALTER COLUMN receipt_url TYPE TEXT;
