-- 000010: Add bundle products support
-- Bundles are composite products made of multiple components

BEGIN;

-- 1. Extend products table with bundle-specific fields
ALTER TABLE products ADD COLUMN bundle_mode VARCHAR(20);        -- 'fixed' | 'configurable'
ALTER TABLE products ADD COLUMN bundle_price_mode VARCHAR(20);  -- 'computed' | 'fixed'

-- Add CHECK constraints for bundle fields
ALTER TABLE products ADD CONSTRAINT check_bundle_mode
  CHECK (bundle_mode IS NULL OR bundle_mode IN ('fixed', 'configurable'));

ALTER TABLE products ADD CONSTRAINT check_bundle_price_mode
  CHECK (bundle_price_mode IS NULL OR bundle_price_mode IN ('computed', 'fixed'));

-- Ensure bundle fields are only set for bundle products
ALTER TABLE products ADD CONSTRAINT check_bundle_fields
  CHECK (
    (product_type = 'bundle' AND bundle_mode IS NOT NULL AND bundle_price_mode IS NOT NULL)
    OR (product_type != 'bundle' AND bundle_mode IS NULL AND bundle_price_mode IS NULL)
  );

-- 2. Update product_type constraint to include 'bundle'
ALTER TABLE products DROP CONSTRAINT check_product_type;
ALTER TABLE products ADD CONSTRAINT check_product_type CHECK (
  product_type IN ('simple', 'variant_parent', 'variant', 'parametric', 'bundle')
);

-- 3. Create bundle_components table
CREATE TABLE bundle_components (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL,
  bundle_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  component_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,

  -- Quantity settings
  quantity INT NOT NULL DEFAULT 1,           -- Fixed quantity (mode=fixed) or default (mode=configurable)
  min_quantity INT,                           -- Only for configurable mode
  max_quantity INT,                           -- Only for configurable mode

  -- Sorting
  sort_order INT NOT NULL DEFAULT 0,

  -- Parametric defaults (optional, JSON)
  -- If component is parametric, can store default axis values
  default_parameters JSONB,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  UNIQUE(tenant_id, bundle_product_id, component_product_id)
);

CREATE INDEX idx_bundle_components_bundle ON bundle_components(bundle_product_id);
CREATE INDEX idx_bundle_components_component ON bundle_components(component_product_id);
CREATE INDEX idx_bundle_components_tenant ON bundle_components(tenant_id);

-- 4. Validate component types via trigger (CHECK cannot use subqueries)
-- Components can be: simple, variant, parametric (NOT variant_parent or bundle)
CREATE OR REPLACE FUNCTION check_bundle_component_type()
RETURNS TRIGGER AS $$
DECLARE
  comp_type VARCHAR(20);
BEGIN
  SELECT product_type INTO comp_type FROM products WHERE id = NEW.component_product_id;
  IF comp_type NOT IN ('simple', 'variant', 'parametric') THEN
    RAISE EXCEPTION 'Invalid component type: %. Must be simple, variant, or parametric.', comp_type;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_bundle_component_type
  BEFORE INSERT OR UPDATE ON bundle_components
  FOR EACH ROW EXECUTE FUNCTION check_bundle_component_type();

-- 5. Seed example bundle products
-- Example: Tool Kit Bundle (fixed mode, computed price)
INSERT INTO products (
  id, tenant_id, product_type, sku, name, description,
  category_ids, status, images, bundle_mode, bundle_price_mode, created_at, updated_at
) VALUES (
  'b0000000-1000-0000-0000-000000000001',
  '00000000-0000-0000-0000-000000000001',
  'bundle',
  'BUNDLE-TOOLSET-01',
  '{"de": "Werkzeugset Basic", "en": "Tool Set Basic"}',
  '{"de": "Komplett-Set für Einsteiger mit allen wichtigen Werkzeugen", "en": "Complete starter set with all essential tools"}',
  ARRAY[]::UUID[],
  'active',
  '[{"url": "/images/products/bundle-toolset.svg", "alt_text": "Werkzeugset Basic Bundle", "sort_order": 0, "is_primary": true}]',
  'fixed',
  'computed',
  now(),
  now()
);

-- Example: Custom Furniture Bundle (configurable mode, fixed price)
INSERT INTO products (
  id, tenant_id, product_type, sku, name, description,
  category_ids, status, images, bundle_mode, bundle_price_mode, created_at, updated_at
) VALUES (
  'b0000000-1000-0000-0000-000000000002',
  '00000000-0000-0000-0000-000000000001',
  'bundle',
  'BUNDLE-FURNITURE-CUSTOM',
  '{"de": "Möbel-Paket Individuell", "en": "Custom Furniture Package"}',
  '{"de": "Stellen Sie Ihr individuelles Möbelpaket zusammen", "en": "Assemble your custom furniture package"}',
  ARRAY[]::UUID[],
  'active',
  '[{"url": "/images/products/bundle-furniture.svg", "alt_text": "Möbel-Paket Individuell", "sort_order": 0, "is_primary": true}]',
  'configurable',
  'fixed',
  now(),
  now()
);

-- Add a fixed price for the configurable bundle
-- (Note: This references prices table which should exist from earlier migrations)
INSERT INTO prices (
  id, tenant_id, product_id, min_quantity, price, currency,
  valid_from, created_at, updated_at
) VALUES (
  gen_random_uuid(),
  '00000000-0000-0000-0000-000000000001',
  'b0000000-1000-0000-0000-000000000002',
  1,
  999.99,
  'EUR',
  now(),
  now(),
  now()
);

COMMIT;
