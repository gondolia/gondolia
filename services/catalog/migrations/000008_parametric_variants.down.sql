BEGIN;

-- Remove seed data
DELETE FROM products WHERE sku IN ('MTL-SHEET-STEEL', 'HYD-HOSE-METER', 'MTL-PLATE-ALU');
DELETE FROM categories WHERE id = 'a0000000-3000-0000-0000-000000000001';

-- Drop axis_options table
DROP TABLE IF EXISTS axis_options;

-- Drop parametric_pricing table
DROP TABLE IF EXISTS parametric_pricing;

-- Remove parametric columns from variant_axes
ALTER TABLE variant_axes DROP CONSTRAINT IF EXISTS check_axis_input_type;
ALTER TABLE variant_axes DROP COLUMN IF EXISTS unit;
ALTER TABLE variant_axes DROP COLUMN IF EXISTS step_value;
ALTER TABLE variant_axes DROP COLUMN IF EXISTS max_value;
ALTER TABLE variant_axes DROP COLUMN IF EXISTS min_value;
ALTER TABLE variant_axes DROP COLUMN IF EXISTS input_type;

-- Restore original product_type constraint
ALTER TABLE products DROP CONSTRAINT check_product_type;
ALTER TABLE products ADD CONSTRAINT check_product_type CHECK (
  product_type IN ('simple', 'variant_parent', 'variant')
);

COMMIT;
