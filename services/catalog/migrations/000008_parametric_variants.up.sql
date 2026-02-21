-- 000008: Add parametric variant product type
-- Parametric products have both discrete axes (dropdowns) and range axes (numeric inputs with min/max/step).
-- Price is calculated based on parameters (per m², per running meter, etc.)

BEGIN;

-- 1. Extend product_type to include 'parametric'
ALTER TABLE products DROP CONSTRAINT check_product_type;
ALTER TABLE products ADD CONSTRAINT check_product_type CHECK (
  product_type IN ('simple', 'variant_parent', 'variant', 'parametric')
);

-- 2. Extend variant_axes with parametric fields
ALTER TABLE variant_axes ADD COLUMN input_type VARCHAR(20) NOT NULL DEFAULT 'select';
ALTER TABLE variant_axes ADD COLUMN min_value DOUBLE PRECISION;
ALTER TABLE variant_axes ADD COLUMN max_value DOUBLE PRECISION;
ALTER TABLE variant_axes ADD COLUMN step_value DOUBLE PRECISION;
ALTER TABLE variant_axes ADD COLUMN unit VARCHAR(20);
ALTER TABLE variant_axes ADD CONSTRAINT check_axis_input_type CHECK (input_type IN ('select', 'range'));

-- 3. Parametric pricing table
CREATE TABLE parametric_pricing (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  formula_type VARCHAR(30) NOT NULL DEFAULT 'per_unit',
  base_price DOUBLE PRECISION NOT NULL DEFAULT 0,
  unit_price DOUBLE PRECISION,
  currency VARCHAR(3) NOT NULL DEFAULT 'CHF',
  min_order_value DOUBLE PRECISION,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(product_id)
);

CREATE INDEX idx_parametric_pricing_product ON parametric_pricing(product_id);

-- 4. Axis options table for parametric select axes
-- (For variant_parent, options are derived from actual variant axis_values.
--  For parametric products there are no variants, so options are stored directly.)
CREATE TABLE IF NOT EXISTS axis_options (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  axis_id UUID NOT NULL REFERENCES variant_axes(id) ON DELETE CASCADE,
  code VARCHAR(100) NOT NULL,
  label JSONB NOT NULL DEFAULT '{}'::jsonb,
  position INT NOT NULL DEFAULT 0,
  UNIQUE(axis_id, code)
);

CREATE INDEX idx_axis_options_axis ON axis_options(axis_id);

-- 5. Seed: Category "Metallbearbeitung"
INSERT INTO categories (id, tenant_id, code, name, description, parent_id, sort_order, active)
VALUES (
  'a0000000-3000-0000-0000-000000000001',
  'b0000000-0000-0000-0000-000000000001',
  'metallbearbeitung',
  '{"de": "Metallbearbeitung", "en": "Metalworking"}'::jsonb,
  '{"de": "Bleche, Platten und Zuschnitte aus Metall", "en": "Metal sheets, plates and custom cuts"}'::jsonb,
  'a0000000-1000-0000-0000-000000000001',
  16,
  true
) ON CONFLICT (id) DO NOTHING;

-- 6. Seed: Parametric products

-- Product 1: Stahlblech verzinkt
INSERT INTO products (id, tenant_id, product_type, sku, name, description, category_ids, attributes, status, images)
VALUES (
  'c0000000-9000-0000-0000-000000000001',
  'b0000000-0000-0000-0000-000000000001',
  'parametric',
  'MTL-SHEET-STEEL',
  '{"de": "Stahlblech verzinkt", "en": "Galvanized Steel Sheet"}'::jsonb,
  '{"de": "Stahlblech nach Mass. Wählen Sie Oberfläche, Länge, Breite und Stärke.", "en": "Custom cut steel sheet. Choose surface, length, width and thickness."}'::jsonb,
  ARRAY['a0000000-3000-0000-0000-000000000001']::uuid[],
  '[{"key": "material", "type": "text", "value": "Stahl DC01"},
    {"key": "norm", "type": "text", "value": "EN 10130"},
    {"key": "density_kg_m3", "type": "number", "value": 7850}]'::jsonb,
  'active',
  '[{"url": "/images/products/MTL-SHEET-STEEL.png", "is_primary": true, "sort_order": 0}]'::jsonb
);

INSERT INTO variant_axes (id, product_id, attribute_code, position, input_type, min_value, max_value, step_value, unit) VALUES
  ('d0000000-9000-0000-0000-000000000001', 'c0000000-9000-0000-0000-000000000001', 'surface', 0, 'select', NULL, NULL, NULL, NULL),
  ('d0000000-9000-0000-0000-000000000002', 'c0000000-9000-0000-0000-000000000001', 'length_mm', 1, 'range', 100, 3000, 1, 'mm'),
  ('d0000000-9000-0000-0000-000000000003', 'c0000000-9000-0000-0000-000000000001', 'width_mm', 2, 'range', 100, 1500, 1, 'mm'),
  ('d0000000-9000-0000-0000-000000000004', 'c0000000-9000-0000-0000-000000000001', 'thickness_mm', 3, 'range', 0.5, 5, 0.5, 'mm');

INSERT INTO axis_options (axis_id, code, label, position) VALUES
  ('d0000000-9000-0000-0000-000000000001', 'galvanized', '{"de": "Verzinkt", "en": "Galvanized"}', 0),
  ('d0000000-9000-0000-0000-000000000001', 'ral9005', '{"de": "Pulverbeschichtet RAL 9005 (schwarz)", "en": "Powder Coated RAL 9005 (black)"}', 1),
  ('d0000000-9000-0000-0000-000000000001', 'ral7035', '{"de": "Pulverbeschichtet RAL 7035 (lichtgrau)", "en": "Powder Coated RAL 7035 (light gray)"}', 2);

INSERT INTO parametric_pricing (product_id, formula_type, base_price, unit_price, currency)
VALUES ('c0000000-9000-0000-0000-000000000001', 'per_m2', 8.50, 35.00, 'CHF');

-- Product 2: Hydraulikschlauch Meterware
INSERT INTO products (id, tenant_id, product_type, sku, name, description, category_ids, attributes, status, images)
VALUES (
  'c0000000-9000-0000-0000-000000000002',
  'b0000000-0000-0000-0000-000000000001',
  'parametric',
  'HYD-HOSE-METER',
  '{"de": "Hydraulikschlauch Meterware", "en": "Hydraulic Hose by the Meter"}'::jsonb,
  '{"de": "Hydraulikschlauch als Meterware. Wählen Sie Innendurchmesser und Länge.", "en": "Hydraulic hose sold by the meter. Choose inner diameter and length."}'::jsonb,
  ARRAY['a0000000-1100-0000-0000-000000000001']::uuid[],
  '[{"key": "pressure_bar", "type": "number", "value": 350},
    {"key": "norm", "type": "text", "value": "SAE 100R2AT / EN 853 2SN"},
    {"key": "temperature_range", "type": "text", "value": "-40°C bis +100°C"}]'::jsonb,
  'active',
  '[{"url": "/images/products/HYD-HOSE-METER.png", "is_primary": true, "sort_order": 0}]'::jsonb
);

INSERT INTO variant_axes (id, product_id, attribute_code, position, input_type, min_value, max_value, step_value, unit) VALUES
  ('d0000000-9000-0000-0000-000000000005', 'c0000000-9000-0000-0000-000000000002', 'inner_diameter', 0, 'select', NULL, NULL, NULL, NULL),
  ('d0000000-9000-0000-0000-000000000006', 'c0000000-9000-0000-0000-000000000002', 'length_m', 1, 'range', 0.5, 50, 0.1, 'm');

INSERT INTO axis_options (axis_id, code, label, position) VALUES
  ('d0000000-9000-0000-0000-000000000005', 'dn6', '{"de": "DN 6 (1/4\")", "en": "DN 6 (1/4\")"}', 0),
  ('d0000000-9000-0000-0000-000000000005', 'dn10', '{"de": "DN 10 (3/8\")", "en": "DN 10 (3/8\")"}', 1),
  ('d0000000-9000-0000-0000-000000000005', 'dn12', '{"de": "DN 12 (1/2\")", "en": "DN 12 (1/2\")"}', 2),
  ('d0000000-9000-0000-0000-000000000005', 'dn16', '{"de": "DN 16 (5/8\")", "en": "DN 16 (5/8\")"}', 3),
  ('d0000000-9000-0000-0000-000000000005', 'dn20', '{"de": "DN 20 (3/4\")", "en": "DN 20 (3/4\")"}', 4);

INSERT INTO parametric_pricing (product_id, formula_type, base_price, unit_price, currency)
VALUES ('c0000000-9000-0000-0000-000000000002', 'per_running_meter', 5.00, 12.50, 'CHF');

-- Product 3: Aluminium-Platte Zuschnitt
INSERT INTO products (id, tenant_id, product_type, sku, name, description, category_ids, attributes, status, images)
VALUES (
  'c0000000-9000-0000-0000-000000000003',
  'b0000000-0000-0000-0000-000000000001',
  'parametric',
  'MTL-PLATE-ALU',
  '{"de": "Aluminium-Platte Zuschnitt", "en": "Custom Cut Aluminum Plate"}'::jsonb,
  '{"de": "Aluminium-Platte nach Mass. Wählen Sie Legierung, Länge, Breite und Stärke.", "en": "Custom cut aluminum plate. Choose alloy, length, width and thickness."}'::jsonb,
  ARRAY['a0000000-3000-0000-0000-000000000001']::uuid[],
  '[{"key": "material", "type": "text", "value": "Aluminium"},
    {"key": "density_kg_m3", "type": "number", "value": 2700}]'::jsonb,
  'active',
  '[{"url": "/images/products/MTL-PLATE-ALU.png", "is_primary": true, "sort_order": 0}]'::jsonb
);

INSERT INTO variant_axes (id, product_id, attribute_code, position, input_type, min_value, max_value, step_value, unit) VALUES
  ('d0000000-9000-0000-0000-000000000007', 'c0000000-9000-0000-0000-000000000003', 'alloy', 0, 'select', NULL, NULL, NULL, NULL),
  ('d0000000-9000-0000-0000-000000000008', 'c0000000-9000-0000-0000-000000000003', 'length_mm', 1, 'range', 50, 2000, 1, 'mm'),
  ('d0000000-9000-0000-0000-000000000009', 'c0000000-9000-0000-0000-000000000003', 'width_mm', 2, 'range', 50, 1000, 1, 'mm'),
  ('d0000000-9000-0000-0000-000000000010', 'c0000000-9000-0000-0000-000000000003', 'thickness_mm', 3, 'range', 2, 30, 1, 'mm');

INSERT INTO axis_options (axis_id, code, label, position) VALUES
  ('d0000000-9000-0000-0000-000000000007', 'almg3', '{"de": "AlMg3 (EN AW-5754)", "en": "AlMg3 (EN AW-5754)"}', 0),
  ('d0000000-9000-0000-0000-000000000007', 'almgsi1', '{"de": "AlMgSi1 (EN AW-6082)", "en": "AlMgSi1 (EN AW-6082)"}', 1);

INSERT INTO parametric_pricing (product_id, formula_type, base_price, unit_price, currency)
VALUES ('c0000000-9000-0000-0000-000000000003', 'per_m2', 12.00, 85.00, 'CHF');

-- 7. Seed: Attribute translations for parametric axes
INSERT INTO attribute_translations (tenant_id, attribute_key, locale, display_name) VALUES
  ('b0000000-0000-0000-0000-000000000001', 'surface', 'de', 'Oberfläche'),
  ('b0000000-0000-0000-0000-000000000001', 'surface', 'en', 'Surface'),
  ('b0000000-0000-0000-0000-000000000001', 'length_mm', 'de', 'Länge'),
  ('b0000000-0000-0000-0000-000000000001', 'length_mm', 'en', 'Length'),
  ('b0000000-0000-0000-0000-000000000001', 'width_mm', 'de', 'Breite'),
  ('b0000000-0000-0000-0000-000000000001', 'width_mm', 'en', 'Width'),
  ('b0000000-0000-0000-0000-000000000001', 'thickness_mm', 'de', 'Stärke'),
  ('b0000000-0000-0000-0000-000000000001', 'thickness_mm', 'en', 'Thickness'),
  ('b0000000-0000-0000-0000-000000000001', 'inner_diameter', 'de', 'Innendurchmesser'),
  ('b0000000-0000-0000-0000-000000000001', 'inner_diameter', 'en', 'Inner Diameter'),
  ('b0000000-0000-0000-0000-000000000001', 'length_m', 'de', 'Länge'),
  ('b0000000-0000-0000-0000-000000000001', 'length_m', 'en', 'Length'),
  ('b0000000-0000-0000-0000-000000000001', 'alloy', 'de', 'Legierung'),
  ('b0000000-0000-0000-0000-000000000001', 'alloy', 'en', 'Alloy')
ON CONFLICT (tenant_id, attribute_key, locale) DO UPDATE SET display_name = EXCLUDED.display_name;

COMMIT;
