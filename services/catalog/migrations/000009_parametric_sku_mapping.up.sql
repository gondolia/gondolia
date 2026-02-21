-- 000009: Add SKU mapping for parametric products
-- Each combination of select axis values maps to a unique SKU with its own unit_price.
-- Range axes (length, width) only affect quantity calculation, not SKU.

BEGIN;

-- 1. SKU mapping table
CREATE TABLE parametric_sku_mapping (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  selections JSONB NOT NULL DEFAULT '{}'::jsonb,  -- {"surface": "galvanized", "thickness_mm": "2"}
  sku VARCHAR(100) NOT NULL,
  unit_price DOUBLE PRECISION NOT NULL,            -- price per unit (m², running meter, etc.)
  base_price DOUBLE PRECISION NOT NULL DEFAULT 0,  -- fixed surcharge for this combination
  stock INT,                                        -- optional: stock per SKU variant
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(product_id, sku)
);

CREATE INDEX idx_parametric_sku_product ON parametric_sku_mapping(product_id);
CREATE INDEX idx_parametric_sku_selections ON parametric_sku_mapping USING gin(selections);

-- 2. Convert thickness from range to select for steel sheet
UPDATE variant_axes SET input_type = 'select', min_value = NULL, max_value = NULL, step_value = NULL, unit = 'mm'
WHERE id = 'd0000000-9000-0000-0000-000000000004';

-- Add thickness options
INSERT INTO axis_options (axis_id, code, label, position) VALUES
  ('d0000000-9000-0000-0000-000000000004', '0.5', '{"de": "0,5 mm", "en": "0.5 mm"}', 0),
  ('d0000000-9000-0000-0000-000000000004', '0.75', '{"de": "0,75 mm", "en": "0.75 mm"}', 1),
  ('d0000000-9000-0000-0000-000000000004', '1', '{"de": "1 mm", "en": "1 mm"}', 2),
  ('d0000000-9000-0000-0000-000000000004', '1.5', '{"de": "1,5 mm", "en": "1.5 mm"}', 3),
  ('d0000000-9000-0000-0000-000000000004', '2', '{"de": "2 mm", "en": "2 mm"}', 4),
  ('d0000000-9000-0000-0000-000000000004', '3', '{"de": "3 mm", "en": "3 mm"}', 5);

-- 3. Convert thickness from range to select for aluminum plate
UPDATE variant_axes SET input_type = 'select', min_value = NULL, max_value = NULL, step_value = NULL, unit = 'mm'
WHERE id = 'd0000000-9000-0000-0000-000000000010';

INSERT INTO axis_options (axis_id, code, label, position) VALUES
  ('d0000000-9000-0000-0000-000000000010', '2', '{"de": "2 mm", "en": "2 mm"}', 0),
  ('d0000000-9000-0000-0000-000000000010', '3', '{"de": "3 mm", "en": "3 mm"}', 1),
  ('d0000000-9000-0000-0000-000000000010', '5', '{"de": "5 mm", "en": "5 mm"}', 2),
  ('d0000000-9000-0000-0000-000000000010', '8', '{"de": "8 mm", "en": "8 mm"}', 3),
  ('d0000000-9000-0000-0000-000000000010', '10', '{"de": "10 mm", "en": "10 mm"}', 4),
  ('d0000000-9000-0000-0000-000000000010', '15', '{"de": "15 mm", "en": "15 mm"}', 5),
  ('d0000000-9000-0000-0000-000000000010', '20', '{"de": "20 mm", "en": "20 mm"}', 6);

-- 4. Seed SKU mappings for Steel Sheet (surface × thickness = 18 SKUs)
INSERT INTO parametric_sku_mapping (product_id, selections, sku, unit_price, base_price) VALUES
  -- Verzinkt
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"galvanized","thickness_mm":"0.5"}', 'MTL-STEEL-VZ-05', 28.00, 5.00),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"galvanized","thickness_mm":"0.75"}', 'MTL-STEEL-VZ-075', 31.00, 5.50),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"galvanized","thickness_mm":"1"}', 'MTL-STEEL-VZ-10', 35.00, 6.00),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"galvanized","thickness_mm":"1.5"}', 'MTL-STEEL-VZ-15', 42.00, 7.00),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"galvanized","thickness_mm":"2"}', 'MTL-STEEL-VZ-20', 52.00, 8.50),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"galvanized","thickness_mm":"3"}', 'MTL-STEEL-VZ-30', 68.00, 10.00),
  -- RAL 9005 (schwarz)
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"ral9005","thickness_mm":"0.5"}', 'MTL-STEEL-PB9005-05', 38.00, 8.00),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"ral9005","thickness_mm":"0.75"}', 'MTL-STEEL-PB9005-075', 41.00, 8.50),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"ral9005","thickness_mm":"1"}', 'MTL-STEEL-PB9005-10', 45.00, 9.00),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"ral9005","thickness_mm":"1.5"}', 'MTL-STEEL-PB9005-15', 52.00, 10.00),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"ral9005","thickness_mm":"2"}', 'MTL-STEEL-PB9005-20', 62.00, 12.00),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"ral9005","thickness_mm":"3"}', 'MTL-STEEL-PB9005-30', 78.00, 14.00),
  -- RAL 7035 (lichtgrau)
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"ral7035","thickness_mm":"0.5"}', 'MTL-STEEL-PB7035-05', 38.00, 8.00),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"ral7035","thickness_mm":"0.75"}', 'MTL-STEEL-PB7035-075', 41.00, 8.50),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"ral7035","thickness_mm":"1"}', 'MTL-STEEL-PB7035-10', 45.00, 9.00),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"ral7035","thickness_mm":"1.5"}', 'MTL-STEEL-PB7035-15', 52.00, 10.00),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"ral7035","thickness_mm":"2"}', 'MTL-STEEL-PB7035-20', 62.00, 12.00),
  ('c0000000-9000-0000-0000-000000000001', '{"surface":"ral7035","thickness_mm":"3"}', 'MTL-STEEL-PB7035-30', 78.00, 14.00);

-- 5. Seed SKU mappings for Hydraulic Hose (inner_diameter = 5 SKUs)
INSERT INTO parametric_sku_mapping (product_id, selections, sku, unit_price, base_price) VALUES
  ('c0000000-9000-0000-0000-000000000002', '{"inner_diameter":"dn6"}', 'HYD-HOSE-DN6', 8.50, 3.00),
  ('c0000000-9000-0000-0000-000000000002', '{"inner_diameter":"dn10"}', 'HYD-HOSE-DN10', 12.50, 4.00),
  ('c0000000-9000-0000-0000-000000000002', '{"inner_diameter":"dn12"}', 'HYD-HOSE-DN12', 15.00, 4.50),
  ('c0000000-9000-0000-0000-000000000002', '{"inner_diameter":"dn16"}', 'HYD-HOSE-DN16', 19.50, 5.00),
  ('c0000000-9000-0000-0000-000000000002', '{"inner_diameter":"dn20"}', 'HYD-HOSE-DN20', 24.00, 6.00);

-- 6. Seed SKU mappings for Aluminum Plate (alloy × thickness = 14 SKUs)
INSERT INTO parametric_sku_mapping (product_id, selections, sku, unit_price, base_price) VALUES
  -- AlMg3
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almg3","thickness_mm":"2"}', 'MTL-ALU-MG3-02', 72.00, 8.00),
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almg3","thickness_mm":"3"}', 'MTL-ALU-MG3-03', 85.00, 10.00),
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almg3","thickness_mm":"5"}', 'MTL-ALU-MG3-05', 105.00, 12.00),
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almg3","thickness_mm":"8"}', 'MTL-ALU-MG3-08', 135.00, 15.00),
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almg3","thickness_mm":"10"}', 'MTL-ALU-MG3-10', 155.00, 18.00),
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almg3","thickness_mm":"15"}', 'MTL-ALU-MG3-15', 195.00, 22.00),
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almg3","thickness_mm":"20"}', 'MTL-ALU-MG3-20', 240.00, 28.00),
  -- AlMgSi1
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almgsi1","thickness_mm":"2"}', 'MTL-ALU-MGSI1-02', 82.00, 10.00),
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almgsi1","thickness_mm":"3"}', 'MTL-ALU-MGSI1-03', 95.00, 12.00),
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almgsi1","thickness_mm":"5"}', 'MTL-ALU-MGSI1-05', 118.00, 14.00),
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almgsi1","thickness_mm":"8"}', 'MTL-ALU-MGSI1-08', 150.00, 18.00),
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almgsi1","thickness_mm":"10"}', 'MTL-ALU-MGSI1-10', 175.00, 22.00),
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almgsi1","thickness_mm":"15"}', 'MTL-ALU-MGSI1-15', 220.00, 28.00),
  ('c0000000-9000-0000-0000-000000000003', '{"alloy":"almgsi1","thickness_mm":"20"}', 'MTL-ALU-MGSI1-20', 270.00, 35.00);

-- 7. Remove unit_price from parametric_pricing (now per-SKU in mapping table)
-- Keep base_price as fallback / formula_type / currency / min_order_value
-- Actually let's keep unit_price as default fallback if no SKU mapping found

COMMIT;
