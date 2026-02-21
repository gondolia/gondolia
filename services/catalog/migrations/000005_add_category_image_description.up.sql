-- Add image and description fields to categories table

-- Add description column (JSONB for i18n support, like name field)
ALTER TABLE categories ADD COLUMN description JSONB DEFAULT '{}'::jsonb;

-- Add image column (URL to category image)
ALTER TABLE categories ADD COLUMN image VARCHAR(500);

-- Create index on description for better search performance
CREATE INDEX idx_categories_description ON categories USING gin (description);

COMMENT ON COLUMN categories.description IS 'Multi-language description (locale -> text)';
COMMENT ON COLUMN categories.image IS 'URL to category image';
