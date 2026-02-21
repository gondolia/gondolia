-- Remove image and description fields from categories table

DROP INDEX IF EXISTS idx_categories_description;
ALTER TABLE categories DROP COLUMN IF EXISTS description;
ALTER TABLE categories DROP COLUMN IF EXISTS image;
