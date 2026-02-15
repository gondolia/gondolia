-- Drop variant tables
DROP TABLE IF EXISTS variant_axis_values;
DROP TABLE IF EXISTS variant_axes;

-- Remove variant columns from products
DROP INDEX IF EXISTS idx_products_type;
DROP INDEX IF EXISTS idx_products_parent;

ALTER TABLE products DROP CONSTRAINT IF EXISTS check_variant_parent;
ALTER TABLE products DROP CONSTRAINT IF EXISTS check_product_type;
ALTER TABLE products DROP COLUMN IF EXISTS parent_id;
ALTER TABLE products DROP COLUMN IF EXISTS product_type;
