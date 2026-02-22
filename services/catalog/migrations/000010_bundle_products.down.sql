BEGIN;

-- Remove seed data
DELETE FROM prices WHERE product_id IN ('b0000000-1000-0000-0000-000000000001', 'b0000000-1000-0000-0000-000000000002');
DELETE FROM products WHERE product_type = 'bundle';

-- Drop trigger and function
DROP TRIGGER IF EXISTS trg_check_bundle_component_type ON bundle_components;
DROP FUNCTION IF EXISTS check_bundle_component_type();

-- Drop bundle_components table
DROP TABLE IF EXISTS bundle_components;

-- Remove bundle fields from products table
ALTER TABLE products DROP CONSTRAINT IF EXISTS check_bundle_fields;
ALTER TABLE products DROP CONSTRAINT IF EXISTS check_bundle_price_mode;
ALTER TABLE products DROP CONSTRAINT IF EXISTS check_bundle_mode;
ALTER TABLE products DROP COLUMN IF EXISTS bundle_price_mode;
ALTER TABLE products DROP COLUMN IF EXISTS bundle_mode;

-- Restore product_type constraint (remove 'bundle')
ALTER TABLE products DROP CONSTRAINT check_product_type;
ALTER TABLE products ADD CONSTRAINT check_product_type CHECK (
  product_type IN ('simple', 'variant_parent', 'variant', 'parametric')
);

COMMIT;
