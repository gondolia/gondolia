-- ============================================================
-- Migration 000007: Convert simple products to variant groups
-- Tenant: 00000000-0000-0000-0000-000000000001
-- NOTE: IM3000 and SAFEGRIP-PRO variants are NOT touched.
-- ============================================================

-- ============================================================
-- ATTRIBUTE TRANSLATIONS for new axis codes
-- ============================================================
INSERT INTO attribute_translations (tenant_id, attribute_key, locale, display_name) VALUES
  ('00000000-0000-0000-0000-000000000001', 'power_kw',       'de', 'Leistung'),
  ('00000000-0000-0000-0000-000000000001', 'power_kw',       'en', 'Power'),
  ('00000000-0000-0000-0000-000000000001', 'wattage',        'de', 'Leistung'),
  ('00000000-0000-0000-0000-000000000001', 'wattage',        'en', 'Wattage'),
  ('00000000-0000-0000-0000-000000000001', 'bearing_type',   'de', 'Lagertyp'),
  ('00000000-0000-0000-0000-000000000001', 'bearing_type',   'en', 'Bearing Type'),
  ('00000000-0000-0000-0000-000000000001', 'belt_type',      'de', 'Riementyp'),
  ('00000000-0000-0000-0000-000000000001', 'belt_type',      'en', 'Belt Type'),
  ('00000000-0000-0000-0000-000000000001', 'chain_type',     'de', 'Kettentyp'),
  ('00000000-0000-0000-0000-000000000001', 'chain_type',     'en', 'Chain Type'),
  ('00000000-0000-0000-0000-000000000001', 'conductors',     'de', 'Anzahl Adern'),
  ('00000000-0000-0000-0000-000000000001', 'conductors',     'en', 'Conductors'),
  ('00000000-0000-0000-0000-000000000001', 'cross_section',  'de', 'Querschnitt'),
  ('00000000-0000-0000-0000-000000000001', 'cross_section',  'en', 'Cross Section'),
  ('00000000-0000-0000-0000-000000000001', 'cable_category', 'de', 'Netzwerkkategorie'),
  ('00000000-0000-0000-0000-000000000001', 'cable_category', 'en', 'Network Category'),
  ('00000000-0000-0000-0000-000000000001', 'characteristic', 'de', 'Charakteristik'),
  ('00000000-0000-0000-0000-000000000001', 'characteristic', 'en', 'Characteristic'),
  ('00000000-0000-0000-0000-000000000001', 'poles',          'de', 'Polanzahl'),
  ('00000000-0000-0000-0000-000000000001', 'poles',          'en', 'Poles'),
  ('00000000-0000-0000-0000-000000000001', 'nennstrom',      'de', 'Nennstrom'),
  ('00000000-0000-0000-0000-000000000001', 'nennstrom',      'en', 'Rated Current'),
  ('00000000-0000-0000-0000-000000000001', 'volume',         'de', 'Volumen'),
  ('00000000-0000-0000-0000-000000000001', 'volume',         'en', 'Volume'),
  ('00000000-0000-0000-0000-000000000001', 'viscosity',      'de', 'Viskositaetsklasse'),
  ('00000000-0000-0000-0000-000000000001', 'viscosity',      'en', 'Viscosity Class'),
  ('00000000-0000-0000-0000-000000000001', 'flow_rate',      'de', 'Foerdermenge'),
  ('00000000-0000-0000-0000-000000000001', 'flow_rate',      'en', 'Flow Rate'),
  ('00000000-0000-0000-0000-000000000001', 'bore_stroke',    'de', 'Kolbendurchmesser / Hub'),
  ('00000000-0000-0000-0000-000000000001', 'bore_stroke',    'en', 'Bore / Stroke'),
  ('00000000-0000-0000-0000-000000000001', 'filtration',     'de', 'Filterfeinheit'),
  ('00000000-0000-0000-0000-000000000001', 'filtration',     'en', 'Filtration Grade'),
  ('00000000-0000-0000-0000-000000000001', 'nominal_diam',   'de', 'Nennweite'),
  ('00000000-0000-0000-0000-000000000001', 'nominal_diam',   'en', 'Nominal Diameter'),
  ('00000000-0000-0000-0000-000000000001', 'connection',     'de', 'Anschluss'),
  ('00000000-0000-0000-0000-000000000001', 'connection',     'en', 'Connection'),
  ('00000000-0000-0000-0000-000000000001', 'size',           'de', 'Groesse'),
  ('00000000-0000-0000-0000-000000000001', 'size',           'en', 'Size'),
  ('00000000-0000-0000-0000-000000000001', 'protection_lvl', 'de', 'Schutzstufe'),
  ('00000000-0000-0000-0000-000000000001', 'protection_lvl', 'en', 'Protection Level'),
  ('00000000-0000-0000-0000-000000000001', 'snr',            'de', 'SNR-Wert'),
  ('00000000-0000-0000-0000-000000000001', 'snr',            'en', 'SNR Value'),
  ('00000000-0000-0000-0000-000000000001', 'tint',           'de', 'Toenung'),
  ('00000000-0000-0000-0000-000000000001', 'tint',           'en', 'Lens Tint'),
  ('00000000-0000-0000-0000-000000000001', 'format',         'de', 'Format'),
  ('00000000-0000-0000-0000-000000000001', 'format',         'en', 'Format'),
  ('00000000-0000-0000-0000-000000000001', 'valve_type',     'de', 'Ventiltyp'),
  ('00000000-0000-0000-0000-000000000001', 'valve_type',     'en', 'Valve Type'),
  ('00000000-0000-0000-0000-000000000001', 'hose_thread',    'de', 'Schlauch / Gewinde'),
  ('00000000-0000-0000-0000-000000000001', 'hose_thread',    'en', 'Hose / Thread'),
  ('00000000-0000-0000-0000-000000000001', 'pal_size',       'de', 'Palettentyp'),
  ('00000000-0000-0000-0000-000000000001', 'pal_size',       'en', 'Pallet Type')
ON CONFLICT (tenant_id, attribute_key, locale) DO NOTHING;

-- ============================================================
-- HELPER: reusable subquery pattern
-- For each group:
--   1. INSERT variant_parent (ON CONFLICT DO NOTHING)
--   2. INSERT variant_axes  (ON CONFLICT DO NOTHING)
--   3. UPDATE variants      (only if still 'simple')
--   4. INSERT variant_axis_values (ON CONFLICT DO NOTHING)
-- ============================================================


-- ============================================================
-- GROUP 1: DRV-MOT — Elektromotoren
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'DRV-MOT',
  '{"de": "Elektromotor", "en": "Electric Motor"}',
  '{"de": "Drehstrommotoren in verschiedenen Leistungsstufen fuer Industrieanwendungen", "en": "Three-phase motors in various power ratings for industrial applications"}',
  ARRAY['a0000000-1300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'power_kw', 0 FROM products WHERE sku = 'DRV-MOT' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'DRV-MOT' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('DRV-MOT-0.37KW','DRV-MOT-0.75KW','DRV-MOT-1.5KW','DRV-MOT-2.2KW')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '0.37_kw' FROM products v
  JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-MOT-0.37KW' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'power_kw'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '0.75_kw' FROM products v
  JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-MOT-0.75KW' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'power_kw'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1.5_kw' FROM products v
  JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-MOT-1.5KW' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'power_kw'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '2.2_kw' FROM products v
  JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-MOT-2.2KW' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'power_kw'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 2: DRV-VFD — Frequenzumrichter
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'DRV-VFD',
  '{"de": "Frequenzumrichter", "en": "Variable Frequency Drive"}',
  '{"de": "Frequenzumrichter zur stufenlosen Drehzahlregelung", "en": "Variable frequency drives for stepless speed control"}',
  ARRAY['a0000000-1300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'power_kw', 0 FROM products WHERE sku = 'DRV-VFD' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'DRV-VFD' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('DRV-VFD-0.75KW','DRV-VFD-1.5KW','DRV-VFD-2.2KW')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '0.75_kw' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-VFD-0.75KW' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'power_kw'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1.5_kw' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-VFD-1.5KW' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'power_kw'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '2.2_kw' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-VFD-2.2KW' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'power_kw'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 3: DRV-BRG — Kugellager
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'DRV-BRG',
  '{"de": "Kugellager 2RS", "en": "Ball Bearing 2RS"}',
  '{"de": "Rillenkugellager mit beidseitiger Dichtung in verschiedenen Groessen", "en": "Deep groove ball bearings with double seal in various sizes"}',
  ARRAY['a0000000-1300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'bearing_type', 0 FROM products WHERE sku = 'DRV-BRG' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'DRV-BRG' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('DRV-BRG-6001-2RS','DRV-BRG-6204-2RS','DRV-BRG-6306-2RS')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '6001-2rs' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-BRG-6001-2RS' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'bearing_type'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '6204-2rs' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-BRG-6204-2RS' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'bearing_type'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '6306-2rs' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-BRG-6306-2RS' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'bearing_type'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 4: DRV-BELT — Zahnriemen
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'DRV-BELT',
  '{"de": "Zahnriemen", "en": "Timing Belt"}',
  '{"de": "Zahnriemen in verschiedenen Profilen und Laengen", "en": "Timing belts in various profiles and lengths"}',
  ARRAY['a0000000-1300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'belt_type', 0 FROM products WHERE sku = 'DRV-BELT' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'DRV-BELT' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('DRV-BELT-AT10-1200','DRV-BELT-T10-800','DRV-BELT-T5-500')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'at10-1200' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-BELT-AT10-1200' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'belt_type'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 't10-800' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-BELT-T10-800' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'belt_type'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 't5-500' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-BELT-T5-500' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'belt_type'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 5: DRV-CHAIN — Rollenkette
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'DRV-CHAIN',
  '{"de": "Rollenkette 5m", "en": "Roller Chain 5m"}',
  '{"de": "Einfachrollenketten nach DIN 8187 in verschiedenen Teilungen", "en": "Single-strand roller chains DIN 8187 in various pitches"}',
  ARRAY['a0000000-1300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'chain_type', 0 FROM products WHERE sku = 'DRV-CHAIN' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'DRV-CHAIN' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('DRV-CHAIN-08B','DRV-CHAIN-10B','DRV-CHAIN-12B')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '08b' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-CHAIN-08B' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'chain_type'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '10b' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-CHAIN-10B' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'chain_type'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '12b' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'DRV-CHAIN-12B' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'chain_type'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 6: ELC-CAB-NYM — NYM-J Installationskabel
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-CAB-NYM',
  '{"de": "NYM-J Installationskabel", "en": "NYM-J Installation Cable"}',
  '{"de": "Installationskabel NYM-J fuer feste Verlegung in verschiedenen Querschnitten", "en": "NYM-J installation cable for fixed wiring in various cross-sections"}',
  ARRAY['a0000000-2100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'conductors', 0 FROM products WHERE sku = 'ELC-CAB-NYM' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'cross_section', 1 FROM products WHERE sku = 'ELC-CAB-NYM' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-CAB-NYM' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-CAB-NYM-3x1.5','ELC-CAB-NYM-3x2.5','ELC-CAB-NYM-5x1.5','ELC-CAB-NYM-5x2.5')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

-- NYM-3x1.5: 3 Adern, 1.5mm²
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '3' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-NYM-3x1.5' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'conductors'
ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1.5' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-NYM-3x1.5' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cross_section'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

-- NYM-3x2.5: 3 Adern, 2.5mm²
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '3' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-NYM-3x2.5' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'conductors'
ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '2.5' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-NYM-3x2.5' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cross_section'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

-- NYM-5x1.5: 5 Adern, 1.5mm²
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '5' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-NYM-5x1.5' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'conductors'
ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1.5' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-NYM-5x1.5' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cross_section'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

-- NYM-5x2.5: 5 Adern, 2.5mm²
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '5' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-NYM-5x2.5' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'conductors'
ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '2.5' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-NYM-5x2.5' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cross_section'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 7: ELC-CAB-CAT — Netzwerkkabel
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-CAB-CAT',
  '{"de": "Netzwerkkabel 305m", "en": "Network Cable 305m"}',
  '{"de": "Netzwerkkabel verschiedener Kategorien auf 305m-Trommel", "en": "Network cables of various categories on 305m drum"}',
  ARRAY['a0000000-2100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'cable_category', 0 FROM products WHERE sku = 'ELC-CAB-CAT' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-CAB-CAT' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-CAB-CAT6-UTP','ELC-CAB-CAT6A-SFTP','ELC-CAB-CAT7-SFTP')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'cat6-utp' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-CAT6-UTP' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cable_category'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'cat6a-sftp' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-CAT6A-SFTP' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cable_category'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'cat7-sftp' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-CAT7-SFTP' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cable_category'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 8: ELC-CAB-H07V — H07V-K Aderleitung
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-CAB-H07V',
  '{"de": "H07V-K Aderleitung 100m", "en": "H07V-K Single Core Cable 100m"}',
  '{"de": "Flexible Aderleitung H07V-K in verschiedenen Querschnitten und Farben, 100m-Trommel", "en": "Flexible single core cable H07V-K in various cross-sections and colors, 100m drum"}',
  ARRAY['a0000000-2100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'cross_section', 0 FROM products WHERE sku = 'ELC-CAB-H07V' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'color', 1 FROM products WHERE sku = 'ELC-CAB-H07V' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-CAB-H07V' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-CAB-H07V-1.5-BLK','ELC-CAB-H07V-1.5-BLU','ELC-CAB-H07V-1.5-YEL','ELC-CAB-H07V-2.5-BLK')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

-- H07V-1.5-BLK
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1.5' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-H07V-1.5-BLK' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cross_section'
ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'schwarz' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-H07V-1.5-BLK' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'color'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

-- H07V-1.5-BLU
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1.5' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-H07V-1.5-BLU' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cross_section'
ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'blau' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-H07V-1.5-BLU' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'color'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

-- H07V-1.5-YEL
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1.5' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-H07V-1.5-YEL' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cross_section'
ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'gelb-gruen' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-H07V-1.5-YEL' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'color'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

-- H07V-2.5-BLK
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '2.5' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-H07V-2.5-BLK' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cross_section'
ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'schwarz' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CAB-H07V-2.5-BLK' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'color'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 9: ELC-BREAKER-B — LS-Schalter B-Char. 1-polig
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-BREAKER-B',
  '{"de": "LS-Schalter B-Charakteristik 1-polig", "en": "MCB B-Characteristic 1-Pole"}',
  '{"de": "Leitungsschutzschalter B-Charakteristik 1-polig in verschiedenen Nennstroemen", "en": "Miniature circuit breaker B-characteristic 1-pole in various rated currents"}',
  ARRAY['a0000000-2200-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'nennstrom', 0 FROM products WHERE sku = 'ELC-BREAKER-B' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-BREAKER-B' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-BREAKER-B16','ELC-BREAKER-B20')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '16a' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-BREAKER-B16' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nennstrom'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '20a' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-BREAKER-B20' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nennstrom'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 10: ELC-BREAKER-C3P — LS-Schalter C-Char. 3-polig
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-BREAKER-C3P',
  '{"de": "LS-Schalter C-Charakteristik 3-polig", "en": "MCB C-Characteristic 3-Pole"}',
  '{"de": "Leitungsschutzschalter C-Charakteristik 3-polig in verschiedenen Nennstroemen", "en": "Miniature circuit breaker C-characteristic 3-pole in various rated currents"}',
  ARRAY['a0000000-2200-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'nennstrom', 0 FROM products WHERE sku = 'ELC-BREAKER-C3P' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-BREAKER-C3P' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-BREAKER-C16-3P','ELC-BREAKER-C20-3P')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '16a' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-BREAKER-C16-3P' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nennstrom'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '20a' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-BREAKER-C20-3P' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nennstrom'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 11: ELC-CON-CEE — CEE-Steckdose 5-polig
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-CON-CEE',
  '{"de": "CEE-Steckdose 5-polig", "en": "CEE Socket 5-Pin"}',
  '{"de": "CEE-Kraftsteckdosen 5-polig in verschiedenen Nennstroemen", "en": "CEE power sockets 5-pin in various rated currents"}',
  ARRAY['a0000000-2200-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'nennstrom', 0 FROM products WHERE sku = 'ELC-CON-CEE' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-CON-CEE' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-CON-CEE-16A-5P','ELC-CON-CEE-32A-5P','ELC-CON-CEE-63A-5P')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '16a' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CON-CEE-16A-5P' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nennstrom'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '32a' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CON-CEE-32A-5P' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nennstrom'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '63a' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CON-CEE-63A-5P' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nennstrom'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 12: ELC-CON-RJ45 — RJ45-Stecker
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-CON-RJ45',
  '{"de": "RJ45-Stecker (100 Stk)", "en": "RJ45 Connector (100 pcs)"}',
  '{"de": "RJ45-Netzwerkstecker in verschiedenen Kategorien, 100 Stueck", "en": "RJ45 network connectors in various categories, 100 pieces"}',
  ARRAY['a0000000-2200-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'cable_category', 0 FROM products WHERE sku = 'ELC-CON-RJ45' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-CON-RJ45' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-CON-RJ45-CAT6','ELC-CON-RJ45-CAT6A')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'cat6' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CON-RJ45-CAT6' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cable_category'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'cat6a' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CON-RJ45-CAT6A' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cable_category'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 13: ELC-CON-SCHUKO — Schuko-Steckdose
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-CON-SCHUKO',
  '{"de": "Schuko-Steckdose", "en": "Schuko Socket"}',
  '{"de": "Schuko-Aufbausteckdosen in verschiedenen Nennstroemen", "en": "Schuko surface-mounted sockets in various rated currents"}',
  ARRAY['a0000000-2200-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'nennstrom', 0 FROM products WHERE sku = 'ELC-CON-SCHUKO' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-CON-SCHUKO' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-CON-SCHUKO-10A','ELC-CON-SCHUKO-16A')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '10a' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CON-SCHUKO-10A' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nennstrom'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '16a' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-CON-SCHUKO-16A' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nennstrom'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 14: ELC-LED-FLOOD — LED-Fluter
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-LED-FLOOD',
  '{"de": "LED-Fluter 4000K IP65", "en": "LED Floodlight 4000K IP65"}',
  '{"de": "LED-Aussenfluter 4000K IP65 in verschiedenen Leistungsstufen", "en": "LED outdoor floodlight 4000K IP65 in various power ratings"}',
  ARRAY['a0000000-2300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'wattage', 0 FROM products WHERE sku = 'ELC-LED-FLOOD' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-LED-FLOOD' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-LED-FLOOD-50W','ELC-LED-FLOOD-100W')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '50w' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-LED-FLOOD-50W' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'wattage'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '100w' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-LED-FLOOD-100W' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'wattage'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 15: ELC-LED-HIGH — LED-Hallenleuchte
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-LED-HIGH',
  '{"de": "LED-Hallenleuchte", "en": "LED High Bay Light"}',
  '{"de": "LED-Hallenleuchten fuer Industrie und Lager in verschiedenen Leistungsstufen", "en": "LED high bay lights for industry and warehouses in various power ratings"}',
  ARRAY['a0000000-2300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'wattage', 0 FROM products WHERE sku = 'ELC-LED-HIGH' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-LED-HIGH' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-LED-HIGH-100W','ELC-LED-HIGH-150W','ELC-LED-HIGH-200W')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '100w' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-LED-HIGH-100W' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'wattage'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '150w' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-LED-HIGH-150W' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'wattage'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '200w' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-LED-HIGH-200W' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'wattage'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 16: ELC-LED-PANEL — LED-Panel
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-LED-PANEL',
  '{"de": "LED-Panel 4000K", "en": "LED Panel 4000K"}',
  '{"de": "LED-Einlegepanele 4000K fuer Rasterdecken in verschiedenen Leistungsstufen", "en": "LED recessed panels 4000K for grid ceilings in various power ratings"}',
  ARRAY['a0000000-2300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'wattage', 0 FROM products WHERE sku = 'ELC-LED-PANEL' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-LED-PANEL' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-LED-PANEL-40W','ELC-LED-PANEL-48W')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '40w' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-LED-PANEL-40W' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'wattage'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '48w' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-LED-PANEL-48W' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'wattage'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 17: ELC-LED-SPOT — LED-Einbaustrahler
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-LED-SPOT',
  '{"de": "LED-Einbaustrahler 4000K", "en": "LED Recessed Spotlight 4000K"}',
  '{"de": "LED-Einbaustrahler 4000K in verschiedenen Leistungsstufen", "en": "LED recessed spotlights 4000K in various power ratings"}',
  ARRAY['a0000000-2300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'wattage', 0 FROM products WHERE sku = 'ELC-LED-SPOT' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-LED-SPOT' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-LED-SPOT-9W','ELC-LED-SPOT-12W')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '9w' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-LED-SPOT-9W' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'wattage'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '12w' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-LED-SPOT-12W' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'wattage'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 18: ELC-LED-T8 — LED-Roehre T8
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-LED-T8',
  '{"de": "LED-Roehre T8", "en": "LED Tube T8"}',
  '{"de": "LED-Roehren T8 als direkter Ersatz fuer Leuchtstoffroehren", "en": "LED tubes T8 as direct replacement for fluorescent tubes"}',
  ARRAY['a0000000-2300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'wattage', 0 FROM products WHERE sku = 'ELC-LED-T8' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-LED-T8' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-LED-T8-18W-120','ELC-LED-T8-22W-150')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '18w-120cm' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-LED-T8-18W-120' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'wattage'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '22w-150cm' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-LED-T8-22W-150' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'wattage'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 19: ELC-RCD — FI-Schutzschalter 30mA
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-RCD',
  '{"de": "FI-Schutzschalter 30mA", "en": "RCCB 30mA"}',
  '{"de": "Fehlerstromschutzschalter 30mA Typ A in verschiedenen Nennstroemen", "en": "Residual current circuit breaker 30mA type A in various rated currents"}',
  ARRAY['a0000000-2200-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'nennstrom', 0 FROM products WHERE sku = 'ELC-RCD' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-RCD' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-RCD-30MA-40A','ELC-RCD-30MA-63A')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '40a' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-RCD-30MA-40A' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nennstrom'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '63a' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-RCD-30MA-63A' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nennstrom'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 20: ELC-TERM — Adernendhuelsen
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'ELC-TERM',
  '{"de": "Adernendhuelsen (100 Stk)", "en": "Wire Ferrules (100 pcs)"}',
  '{"de": "Adernendhuelsen in verschiedenen Querschnitten und Farben, 100 Stueck", "en": "Wire ferrules in various cross-sections and colors, 100 pieces"}',
  ARRAY['a0000000-2200-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'cross_section', 0 FROM products WHERE sku = 'ELC-TERM' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'ELC-TERM' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('ELC-TERM-1.5-BLUE','ELC-TERM-2.5-GRAY','ELC-TERM-4-RED')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1.5mm2' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-TERM-1.5-BLUE' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cross_section'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '2.5mm2' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-TERM-2.5-GRAY' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cross_section'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '4mm2' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'ELC-TERM-4-RED' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'cross_section'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 21: HYD-ACCU — Hydrospeicher 210 bar
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'HYD-ACCU',
  '{"de": "Hydrospeicher 210 bar", "en": "Hydraulic Accumulator 210 bar"}',
  '{"de": "Blasenspeicher 210 bar in verschiedenen Volumina", "en": "Bladder accumulators 210 bar in various volumes"}',
  ARRAY['a0000000-1100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'HYD-ACCU' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'HYD-ACCU' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('HYD-ACCU-0.5L','HYD-ACCU-1.0L','HYD-ACCU-2.5L')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '0.5l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-ACCU-0.5L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1.0l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-ACCU-1.0L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '2.5l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-ACCU-2.5L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 22: HYD-FILTER — Hydraulikfilter
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'HYD-FILTER',
  '{"de": "Hydraulikfilter", "en": "Hydraulic Filter"}',
  '{"de": "Hydraulikfilter in verschiedenen Filterfeinheiten", "en": "Hydraulic filters in various filtration grades"}',
  ARRAY['a0000000-1100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'filtration', 0 FROM products WHERE sku = 'HYD-FILTER' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'HYD-FILTER' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('HYD-FILTER-10MIC','HYD-FILTER-25MIC')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '10um' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-FILTER-10MIC' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'filtration'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '25um' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-FILTER-25MIC' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'filtration'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 23: HYD-HSE — Hydraulikschlauch
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'HYD-HSE',
  '{"de": "Hydraulikschlauch", "en": "Hydraulic Hose"}',
  '{"de": "Hochdruckhydraulikschlaeuche in verschiedenen Nennweiten", "en": "High-pressure hydraulic hoses in various nominal diameters"}',
  ARRAY['a0000000-1100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'nominal_diam', 0 FROM products WHERE sku = 'HYD-HSE' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'HYD-HSE' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('HYD-HSE-6-1SN','HYD-HSE-10-1SN','HYD-HSE-16-2SN')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'dn6' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-HSE-6-1SN' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nominal_diam'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'dn10' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-HSE-10-1SN' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nominal_diam'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'dn16' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-HSE-16-2SN' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'nominal_diam'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 24: HYD-OIL-HLP — Hydraulikoel HLP 20L
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'HYD-OIL-HLP',
  '{"de": "Hydraulikoel HLP 20L", "en": "Hydraulic Oil HLP 20L"}',
  '{"de": "Hydraulikoel HLP nach DIN 51524-2 in verschiedenen Viskositaetsklassen, 20L Kanne", "en": "Hydraulic oil HLP to DIN 51524-2 in various viscosity classes, 20L can"}',
  ARRAY['a0000000-1100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'viscosity', 0 FROM products WHERE sku = 'HYD-OIL-HLP' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'HYD-OIL-HLP' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('HYD-OIL-HLP32-20L','HYD-OIL-HLP46-20L','HYD-OIL-HLP68-20L')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'hlp32' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-OIL-HLP32-20L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'viscosity'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'hlp46' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-OIL-HLP46-20L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'viscosity'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'hlp68' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-OIL-HLP68-20L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'viscosity'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 25: HYD-PMP — Hydraulikpumpe
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'HYD-PMP',
  '{"de": "Hydraulikpumpe", "en": "Hydraulic Pump"}',
  '{"de": "Zahnradpumpen in verschiedenen Foerdermengen", "en": "Gear pumps in various flow rates"}',
  ARRAY['a0000000-1100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'flow_rate', 0 FROM products WHERE sku = 'HYD-PMP' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'HYD-PMP' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('HYD-PMP-10','HYD-PMP-16','HYD-PMP-25')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '10lmin' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-PMP-10' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'flow_rate'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '16lmin' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-PMP-16' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'flow_rate'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '25lmin' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-PMP-25' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'flow_rate'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 26: HYD-ZYL — Hydraulikzylinder
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'HYD-ZYL',
  '{"de": "Hydraulikzylinder", "en": "Hydraulic Cylinder"}',
  '{"de": "Doppelwirkende Hydraulikzylinder in verschiedenen Kolbendurchmessern und Hueben", "en": "Double-acting hydraulic cylinders in various bore diameters and strokes"}',
  ARRAY['a0000000-1100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'bore_stroke', 0 FROM products WHERE sku = 'HYD-ZYL' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'HYD-ZYL' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('HYD-ZYL-50-200','HYD-ZYL-63-300','HYD-ZYL-80-400','HYD-ZYL-100-500')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '50x200' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-ZYL-50-200' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'bore_stroke'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '63x300' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-ZYL-63-300' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'bore_stroke'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '80x400' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-ZYL-80-400' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'bore_stroke'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '100x500' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-ZYL-100-500' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'bore_stroke'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 27: HYD-VLV-4WAY — 4-Wege-Hydraulikventil
-- (HYD-VLV-DR-10 Druckbegrenzungsventil bleibt simple)
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'HYD-VLV-4WAY',
  '{"de": "4-Wege-Hydraulikventil", "en": "4-Way Hydraulic Valve"}',
  '{"de": "4-Wege-Wegeventile in verschiedenen Schaltstellungen", "en": "4-way directional control valves in various switching positions"}',
  ARRAY['a0000000-1100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'valve_type', 0 FROM products WHERE sku = 'HYD-VLV-4WAY' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'HYD-VLV-4WAY' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('HYD-VLV-4-2','HYD-VLV-4-3')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '4-2' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-VLV-4-2' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'valve_type'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '4-3' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'HYD-VLV-4-3' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'valve_type'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 28: LAB-CHEM-ACET — Aceton 99.5%
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-CHEM-ACET',
  '{"de": "Aceton 99,5%", "en": "Acetone 99.5%"}',
  '{"de": "Aceton 99,5% reinst in verschiedenen Gebindegrossen", "en": "Acetone 99.5% pure in various container sizes"}',
  ARRAY['a0000000-5100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-CHEM-ACET' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-CHEM-ACET' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-CHEM-ACET-1L','LAB-CHEM-ACET-5L')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CHEM-ACET-1L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '5l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CHEM-ACET-5L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 29: LAB-CHEM-ETHANOL — Ethanol 96%
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-CHEM-ETHANOL',
  '{"de": "Ethanol 96%", "en": "Ethanol 96%"}',
  '{"de": "Ethanol 96% vergaellt in verschiedenen Gebindegrossen", "en": "Ethanol 96% denatured in various container sizes"}',
  ARRAY['a0000000-5100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-CHEM-ETHANOL' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-CHEM-ETHANOL' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-CHEM-ETHANOL-1L','LAB-CHEM-ETHANOL-5L')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CHEM-ETHANOL-1L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '5l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CHEM-ETHANOL-5L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 30: LAB-CHEM-ISOPROP — Isopropanol 99.9%
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-CHEM-ISOPROP',
  '{"de": "Isopropanol 99,9%", "en": "Isopropanol 99.9%"}',
  '{"de": "Isopropanol 99,9% reinst in verschiedenen Gebindegrossen", "en": "Isopropanol 99.9% pure in various container sizes"}',
  ARRAY['a0000000-5100-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-CHEM-ISOPROP' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-CHEM-ISOPROP' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-CHEM-ISOPROP-1L','LAB-CHEM-ISOPROP-5L')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CHEM-ISOPROP-1L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '5l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CHEM-ISOPROP-5L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 31: LAB-CLEAN-ALL — Universalreiniger
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-CLEAN-ALL',
  '{"de": "Universalreiniger", "en": "Universal Cleaner"}',
  '{"de": "Universalreiniger fuer Buero und Betrieb in verschiedenen Gebindegrossen", "en": "Universal cleaner for office and industrial use in various container sizes"}',
  ARRAY['a0000000-5300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-CLEAN-ALL' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-CLEAN-ALL' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-CLEAN-ALL-5L','LAB-CLEAN-ALL-10L')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '5l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CLEAN-ALL-5L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '10l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CLEAN-ALL-10L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 32: LAB-CLEAN-DEGREASE — Industrieentfetter
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-CLEAN-DEGREASE',
  '{"de": "Industrieentfetter", "en": "Industrial Degreaser"}',
  '{"de": "Industrieentfetter fuer Metall und Maschinen in verschiedenen Gebindegrossen", "en": "Industrial degreaser for metal and machinery in various container sizes"}',
  ARRAY['a0000000-5300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-CLEAN-DEGREASE' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-CLEAN-DEGREASE' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-CLEAN-DEGREASE-5L','LAB-CLEAN-DEGREASE-10L')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '5l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CLEAN-DEGREASE-5L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '10l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CLEAN-DEGREASE-10L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 33: LAB-CLEAN-DISINFECT — Flaechendesinfektion
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-CLEAN-DISINFECT',
  '{"de": "Flaechendesinfektion", "en": "Surface Disinfectant"}',
  '{"de": "Alkoholische Flaechendesinfektion in verschiedenen Gebindegrossen", "en": "Alcoholic surface disinfectant in various container sizes"}',
  ARRAY['a0000000-5300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-CLEAN-DISINFECT' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-CLEAN-DISINFECT' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-CLEAN-DISINFECT-1L','LAB-CLEAN-DISINFECT-5L')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CLEAN-DISINFECT-1L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '5l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CLEAN-DISINFECT-5L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 34: LAB-CLEAN-GLASS — Glasreiniger
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-CLEAN-GLASS',
  '{"de": "Glasreiniger", "en": "Glass Cleaner"}',
  '{"de": "Streifenfreier Glasreiniger in verschiedenen Gebindegrossen", "en": "Streak-free glass cleaner in various container sizes"}',
  ARRAY['a0000000-5300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-CLEAN-GLASS' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-CLEAN-GLASS' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-CLEAN-GLASS-5L','LAB-CLEAN-GLASS-10L')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '5l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CLEAN-GLASS-5L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '10l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CLEAN-GLASS-10L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 35: LAB-CLEAN-HANDDISINFECT — Haendedesinfektion
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-CLEAN-HANDDISINFECT',
  '{"de": "Haendedesinfektion", "en": "Hand Disinfectant"}',
  '{"de": "Alkoholisches Haendedesinfektionsmittel in verschiedenen Gebindegrossen", "en": "Alcoholic hand disinfectant in various container sizes"}',
  ARRAY['a0000000-5300-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-CLEAN-HANDDISINFECT' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-CLEAN-HANDDISINFECT' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-CLEAN-HANDDISINFECT-500ML','LAB-CLEAN-HANDDISINFECT-1L')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '500ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CLEAN-HANDDISINFECT-500ML' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1l' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id
WHERE v.sku = 'LAB-CLEAN-HANDDISINFECT-1L' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume'
ON CONFLICT (variant_id, axis_id) DO NOTHING;


-- ============================================================
-- GROUP 36: LAB-GLASS-BEAKER — Becherglas
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-GLASS-BEAKER',
  '{"de": "Becherglas", "en": "Beaker"}',
  '{"de": "Borosilikatglas-Becherglaeser mit Ausguss in verschiedenen Volumina", "en": "Borosilicate glass beakers with spout in various volumes"}',
  ARRAY['a0000000-5200-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-GLASS-BEAKER' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-GLASS-BEAKER' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-GLASS-BEAKER-50','LAB-GLASS-BEAKER-100','LAB-GLASS-BEAKER-250','LAB-GLASS-BEAKER-500','LAB-GLASS-BEAKER-1000')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '50ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-BEAKER-50' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '100ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-BEAKER-100' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '250ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-BEAKER-250' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '500ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-BEAKER-500' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '1000ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-BEAKER-1000' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;

-- ============================================================
-- GROUP 37: LAB-GLASS-CYL — Messzylinder
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-GLASS-CYL',
  '{"de": "Messzylinder Klasse A", "en": "Graduated Cylinder Class A"}',
  '{"de": "Geeichte Messzylinder Klasse A in verschiedenen Volumina", "en": "Calibrated graduated cylinders Class A in various volumes"}',
  ARRAY['a0000000-5200-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-GLASS-CYL' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-GLASS-CYL' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-GLASS-CYL-10','LAB-GLASS-CYL-25','LAB-GLASS-CYL-100','LAB-GLASS-CYL-250','LAB-GLASS-CYL-500','LAB-GLASS-CYL-1000')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '10ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-CYL-10' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '25ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-CYL-25' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '100ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-CYL-100' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '250ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-CYL-250' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '500ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-CYL-500' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '1000ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-CYL-1000' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;

-- ============================================================
-- GROUP 38: LAB-GLASS-ERLEN — Erlenmeyerkolben
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-GLASS-ERLEN',
  '{"de": "Erlenmeyerkolben", "en": "Erlenmeyer Flask"}',
  '{"de": "Erlenmeyerkolben aus Borosilikatglas in verschiedenen Volumina", "en": "Erlenmeyer flasks of borosilicate glass in various volumes"}',
  ARRAY['a0000000-5200-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-GLASS-ERLEN' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-GLASS-ERLEN' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-GLASS-ERLEN-100','LAB-GLASS-ERLEN-250','LAB-GLASS-ERLEN-500','LAB-GLASS-ERLEN-1000')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '100ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-ERLEN-100' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '250ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-ERLEN-250' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '500ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-ERLEN-500' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '1000ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-ERLEN-1000' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;

-- ============================================================
-- GROUP 39: LAB-GLASS-FLASK — Messkolben Klasse A
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-GLASS-FLASK',
  '{"de": "Messkolben Klasse A", "en": "Volumetric Flask Class A"}',
  '{"de": "Geeichte Messkolben Klasse A in verschiedenen Volumina", "en": "Calibrated volumetric flasks Class A in various volumes"}',
  ARRAY['a0000000-5200-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-GLASS-FLASK' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-GLASS-FLASK' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-GLASS-FLASK-100ML','LAB-GLASS-FLASK-250ML','LAB-GLASS-FLASK-500ML','LAB-GLASS-FLASK-1000ML')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '100ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-FLASK-100ML' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '250ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-FLASK-250ML' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '500ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-FLASK-500ML' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '1000ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-FLASK-1000ML' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;

-- ============================================================
-- GROUP 40: LAB-GLASS-PIPETTE — Vollpipette Klasse AS
-- ============================================================
INSERT INTO products (id, tenant_id, sku, name, description, category_ids, attributes, status, images, product_type, created_at, updated_at)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'LAB-GLASS-PIPETTE',
  '{"de": "Vollpipette Klasse AS", "en": "Volumetric Pipette Class AS"}',
  '{"de": "Geeichte Vollpipetten Klasse AS in verschiedenen Volumina", "en": "Calibrated volumetric pipettes Class AS in various volumes"}',
  ARRAY['a0000000-5200-0000-0000-000000000001']::uuid[], '[]'::jsonb, 'active', '[]'::jsonb, 'variant_parent', NOW(), NOW())
ON CONFLICT (tenant_id, sku) DO NOTHING;

INSERT INTO variant_axes (product_id, attribute_code, position)
SELECT id, 'volume', 0 FROM products WHERE sku = 'LAB-GLASS-PIPETTE' AND tenant_id = '00000000-0000-0000-0000-000000000001'
ON CONFLICT (product_id, attribute_code) DO NOTHING;

UPDATE products SET product_type = 'variant',
  parent_id = (SELECT id FROM products WHERE sku = 'LAB-GLASS-PIPETTE' AND tenant_id = '00000000-0000-0000-0000-000000000001')
WHERE sku IN ('LAB-GLASS-PIPETTE-5ML','LAB-GLASS-PIPETTE-10ML','LAB-GLASS-PIPETTE-25ML')
  AND tenant_id = '00000000-0000-0000-0000-000000000001' AND product_type = 'simple';

INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '5ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-PIPETTE-5ML' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '10ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-PIPETTE-10ML' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id, axis_id, option_code) SELECT v.id, va.id, '25ml' FROM products v JOIN variant_axes va ON va.product_id = v.parent_id WHERE v.sku = 'LAB-GLASS-PIPETTE-25ML' AND v.tenant_id = '00000000-0000-0000-0000-000000000001' AND va.attribute_code = 'volume' ON CONFLICT (variant_id, axis_id) DO NOTHING;

-- ============================================================
-- GROUPS 41-72: PKG + PNE + WWR (kompakt)
-- ============================================================

-- GROUP 41: PKG-BOX
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PKG-BOX','{"de":"Faltkarton","en":"Folding Carton"}','{"de":"Faltkartons in verschiedenen Grossen","en":"Folding cartons in various sizes"}',ARRAY['a0000000-3100-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='PKG-BOX' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PKG-BOX' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PKG-BOX-200x150x100','PKG-BOX-300x200x150','PKG-BOX-400x300x200','PKG-BOX-600x400x400','PKG-BOX-800x600x600') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'200x150x100' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-BOX-200x150x100' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'300x200x150' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-BOX-300x200x150' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'400x300x200' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-BOX-400x300x200' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'600x400x400' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-BOX-600x400x400' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'800x600x600' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-BOX-800x600x600' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 42: PKG-CORNER
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PKG-CORNER','{"de":"Kantenschutz L-Profil 1m","en":"Edge Protector L-Profile 1m"}','{"de":"Karton-Kantenschutzwinkel in verschiedenen Profilbreiten","en":"Cardboard edge protector in various profile widths"}',ARRAY['a0000000-3100-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='PKG-CORNER' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PKG-CORNER' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PKG-CORNER-L50','PKG-CORNER-L70') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'50x50mm' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-CORNER-L50' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'70x70mm' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-CORNER-L70' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 43: PKG-DISPENSER
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PKG-DISPENSER','{"de":"Paketbandabroller","en":"Tape Dispenser"}','{"de":"Ergonomischer Paketbandabroller in verschiedenen Breiten","en":"Ergonomic tape dispenser in various widths"}',ARRAY['a0000000-3300-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='PKG-DISPENSER' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PKG-DISPENSER' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PKG-DISPENSER-50','PKG-DISPENSER-75') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'50mm' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-DISPENSER-50' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'75mm' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-DISPENSER-75' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 44: PKG-FILM-BUBBLE
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PKG-FILM-BUBBLE','{"de":"Luftpolsterfolie","en":"Bubble Wrap"}','{"de":"Luftpolsterfolie in verschiedenen Breiten","en":"Bubble wrap in various widths"}',ARRAY['a0000000-3200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='PKG-FILM-BUBBLE' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PKG-FILM-BUBBLE' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PKG-FILM-BUBBLE-100','PKG-FILM-BUBBLE-150') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'100cm' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-FILM-BUBBLE-100' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'150cm' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-FILM-BUBBLE-150' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 45: PKG-FILM-SHRINK
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PKG-FILM-SHRINK','{"de":"Schrumpffolie","en":"Shrink Film"}','{"de":"Schrumpffolie in verschiedenen Ausfuehrungen","en":"Shrink film in various configurations"}',ARRAY['a0000000-3200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='PKG-FILM-SHRINK' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PKG-FILM-SHRINK' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PKG-FILM-SHRINK-400-15','PKG-FILM-SHRINK-500-19') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'400mm-15my' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-FILM-SHRINK-400-15' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'500mm-19my' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-FILM-SHRINK-500-19' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 46: PKG-FILM-STRETCH
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PKG-FILM-STRETCH','{"de":"Stretchfolie 500mm","en":"Stretch Film 500mm"}','{"de":"Handstretchfolie 500mm in verschiedenen Staerken","en":"Hand stretch film 500mm in various thicknesses"}',ARRAY['a0000000-3200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='PKG-FILM-STRETCH' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PKG-FILM-STRETCH' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PKG-FILM-STRETCH-500-17','PKG-FILM-STRETCH-500-23') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'17my' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-FILM-STRETCH-500-17' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'23my' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-FILM-STRETCH-500-23' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 47: PKG-LABEL
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PKG-LABEL','{"de":"Versandetiketten","en":"Shipping Labels"}','{"de":"Selbstklebende Versandetiketten in verschiedenen Formaten","en":"Self-adhesive shipping labels in various formats"}',ARRAY['a0000000-3100-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'format',0 FROM products WHERE sku='PKG-LABEL' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PKG-LABEL' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PKG-LABEL-A6-TH','PKG-LABEL-100x150-TT') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'a6-thermodirekt' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-LABEL-A6-TH' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='format' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'100x150-thermotransfer' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-LABEL-100x150-TT' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='format' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 48: PKG-PAL
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PKG-PAL','{"de":"Holzpalette NEU","en":"Wooden Pallet NEW"}','{"de":"Neue Einwegholzpaletten in verschiedenen Grossen","en":"New single-use wooden pallets in various sizes"}',ARRAY['a0000000-3100-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'pal_size',0 FROM products WHERE sku='PKG-PAL' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PKG-PAL' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PKG-PAL-EUR-NEU','PKG-PAL-HALF-NEU') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'euro-1200x800' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-PAL-EUR-NEU' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='pal_size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'halb-800x600' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-PAL-HALF-NEU' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='pal_size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 49: PKG-TAPE-MASK
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PKG-TAPE-MASK','{"de":"Abdeckband (Kreppband)","en":"Masking Tape"}','{"de":"Abdeckband in verschiedenen Breiten","en":"Masking tape in various widths"}',ARRAY['a0000000-3300-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='PKG-TAPE-MASK' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PKG-TAPE-MASK' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PKG-TAPE-MASK-30','PKG-TAPE-MASK-50') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'30mm' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-TAPE-MASK-30' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'50mm' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-TAPE-MASK-50' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 50: PKG-TAPE-PP
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PKG-TAPE-PP','{"de":"Packband PP 50mm","en":"PP Packaging Tape 50mm"}','{"de":"Polypropylen-Packband 50mm in verschiedenen Farben","en":"Polypropylene packaging tape 50mm in various colors"}',ARRAY['a0000000-3300-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'color',0 FROM products WHERE sku='PKG-TAPE-PP' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PKG-TAPE-PP' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PKG-TAPE-PP-50-BRN','PKG-TAPE-PP-50-CLR') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'braun' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-TAPE-PP-50-BRN' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='color' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'transparent' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-TAPE-PP-50-CLR' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='color' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 51: PKG-TAPE-PVC
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PKG-TAPE-PVC','{"de":"Packband PVC 50mm","en":"PVC Packaging Tape 50mm"}','{"de":"PVC-Packband 50mm in verschiedenen Farben","en":"PVC packaging tape 50mm in various colors"}',ARRAY['a0000000-3300-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'color',0 FROM products WHERE sku='PKG-TAPE-PVC' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PKG-TAPE-PVC' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PKG-TAPE-PVC-50-BRN','PKG-TAPE-PVC-50-CLR') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'braun' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-TAPE-PVC-50-BRN' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='color' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'transparent' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PKG-TAPE-PVC-50-CLR' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='color' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 52: PNE-CYL
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PNE-CYL','{"de":"Pneumatikzylinder","en":"Pneumatic Cylinder"}','{"de":"Doppelwirkende Pneumatikzylinder in verschiedenen Bohrungen und Hueben","en":"Double-acting pneumatic cylinders in various bores and strokes"}',ARRAY['a0000000-1200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'bore_stroke',0 FROM products WHERE sku='PNE-CYL' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PNE-CYL' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PNE-CYL-32-50','PNE-CYL-40-100','PNE-CYL-50-150','PNE-CYL-63-200') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'oe32-50' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-CYL-32-50' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='bore_stroke' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'oe40-100' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-CYL-40-100' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='bore_stroke' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'oe50-150' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-CYL-50-150' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='bore_stroke' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'oe63-200' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-CYL-63-200' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='bore_stroke' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 53: PNE-FRL
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PNE-FRL','{"de":"Wartungseinheit FRL","en":"FRL Maintenance Unit"}','{"de":"Druckluft-Wartungseinheit in verschiedenen Anschlussgroessen","en":"Compressed air maintenance unit in various connection sizes"}',ARRAY['a0000000-1200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'connection',0 FROM products WHERE sku='PNE-FRL' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PNE-FRL' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PNE-FRL-G14','PNE-FRL-G38','PNE-FRL-G12') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'g1-4' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-FRL-G14' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='connection' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'g3-8' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-FRL-G38' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='connection' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'g1-2' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-FRL-G12' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='connection' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 54: PNE-PUSH
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PNE-PUSH','{"de":"Steckverschraubung","en":"Push-In Fitting"}','{"de":"Pneumatik-Steckverschraubungen in verschiedenen Groessen","en":"Pneumatic push-in fittings in various sizes"}',ARRAY['a0000000-1200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'hose_thread',0 FROM products WHERE sku='PNE-PUSH' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PNE-PUSH' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PNE-PUSH-4-M5','PNE-PUSH-6-G14','PNE-PUSH-8-G38') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'4mm-m5' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-PUSH-4-M5' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='hose_thread' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'6mm-g1-4' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-PUSH-6-G14' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='hose_thread' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'8mm-g3-8' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-PUSH-8-G38' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='hose_thread' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 55: PNE-REG
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PNE-REG','{"de":"Druckregler","en":"Pressure Regulator"}','{"de":"Pneumatik-Druckregler in verschiedenen Anschlussgroessen","en":"Pneumatic pressure regulators in various connection sizes"}',ARRAY['a0000000-1200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'connection',0 FROM products WHERE sku='PNE-REG' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PNE-REG' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PNE-REG-G14','PNE-REG-G38') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'g1-4' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-REG-G14' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='connection' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'g3-8' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-REG-G38' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='connection' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 56: PNE-TUB
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PNE-TUB','{"de":"Pneumatikschlauch PU","en":"Pneumatic Tubing PU"}','{"de":"Polyurethan-Pneumatikschlaeuche in verschiedenen Durchmessern","en":"Polyurethane pneumatic tubing in various diameters"}',ARRAY['a0000000-1200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='PNE-TUB' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PNE-TUB' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PNE-TUB-4-BLU','PNE-TUB-6-BLK','PNE-TUB-8-BLU','PNE-TUB-10-BLK') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'4mm-blau' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-TUB-4-BLU' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'6mm-schwarz' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-TUB-6-BLK' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'8mm-blau' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-TUB-8-BLU' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'10mm-schwarz' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-TUB-10-BLK' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 57: PNE-VLV
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','PNE-VLV','{"de":"Pneumatik-Wegeventil","en":"Pneumatic Directional Valve"}','{"de":"Elektromagnetisch betaetigte Pneumatik-Wegeventile","en":"Solenoid-operated pneumatic directional valves"}',ARRAY['a0000000-1200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'valve_type',0 FROM products WHERE sku='PNE-VLV' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='PNE-VLV' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('PNE-VLV-3-2','PNE-VLV-5-2') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'3-2' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-VLV-3-2' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='valve_type' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'5-2' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='PNE-VLV-5-2' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='valve_type' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 58: WWR-BIB
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-BIB','{"de":"Latzhose marine","en":"Bib Overalls Navy"}','{"de":"Latzhose in Marine in verschiedenen Groessen","en":"Bib overalls in navy in various sizes"}',ARRAY['a0000000-4100-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='WWR-BIB' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-BIB' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-BIB-50-NVY','WWR-BIB-52-NVY') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'50' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-BIB-50-NVY' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'52' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-BIB-52-NVY' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 59: WWR-EAR-MUFF
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-EAR-MUFF','{"de":"Kapselgehoerschutz","en":"Ear Muff"}','{"de":"Kapselgehoerschutz in verschiedenen Daemmwerten","en":"Ear muffs with various SNR values"}',ARRAY['a0000000-4200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'snr',0 FROM products WHERE sku='WWR-EAR-MUFF' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-EAR-MUFF' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-EAR-MUFF-SNR27','WWR-EAR-MUFF-SNR32') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'snr27' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-EAR-MUFF-SNR27' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='snr' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'snr32' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-EAR-MUFF-SNR32' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='snr' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 60: WWR-GLASS
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-GLASS','{"de":"Schutzbrille","en":"Safety Goggles"}','{"de":"Schutzbrille in verschiedenen Glastoeningen","en":"Safety goggles in various lens tints"}',ARRAY['a0000000-4200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'tint',0 FROM products WHERE sku='WWR-GLASS' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-GLASS' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-GLASS-CLR','WWR-GLASS-DARK') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'klar' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLASS-CLR' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='tint' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'getoeint' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLASS-DARK' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='tint' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 61: WWR-GLOVE-CUT5
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-GLOVE-CUT5','{"de":"Schnittschutzhandschuhe Level 5","en":"Cut-Resistant Gloves Level 5"}','{"de":"Schnittschutzhandschuhe Stufe 5 in verschiedenen Groessen","en":"Cut-resistant gloves level 5 in various sizes"}',ARRAY['a0000000-4300-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='WWR-GLOVE-CUT5' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-GLOVE-CUT5' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-GLOVE-CUT5-9','WWR-GLOVE-CUT5-10') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'9' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLOVE-CUT5-9' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'10' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLOVE-CUT5-10' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 62: WWR-GLOVE-LATEX
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-GLOVE-LATEX','{"de":"Latexhandschuhe (100 Stk)","en":"Latex Gloves (100 pcs)"}','{"de":"Latexhandschuhe ungepudert in verschiedenen Groessen","en":"Latex gloves unpowdered in various sizes"}',ARRAY['a0000000-4300-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='WWR-GLOVE-LATEX' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-GLOVE-LATEX' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-GLOVE-LATEX-8','WWR-GLOVE-LATEX-9') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'8' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLOVE-LATEX-8' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'9' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLOVE-LATEX-9' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 63: WWR-GLOVE-LEATHER
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-GLOVE-LEATHER','{"de":"Lederhandschuhe","en":"Leather Work Gloves"}','{"de":"Rindleder-Arbeitshandschuhe in verschiedenen Groessen","en":"Cowhide leather work gloves in various sizes"}',ARRAY['a0000000-4300-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='WWR-GLOVE-LEATHER' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-GLOVE-LEATHER' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-GLOVE-LEATHER-9','WWR-GLOVE-LEATHER-10') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'9' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLOVE-LEATHER-9' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'10' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLOVE-LEATHER-10' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 64: WWR-GLOVE-MONT
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-GLOVE-MONT','{"de":"Montagehandschuhe","en":"Assembly Gloves"}','{"de":"Mehrzweck-Montagehandschuhe in verschiedenen Groessen","en":"Multi-purpose assembly gloves in various sizes"}',ARRAY['a0000000-4300-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='WWR-GLOVE-MONT' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-GLOVE-MONT' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-GLOVE-MONT-8','WWR-GLOVE-MONT-9','WWR-GLOVE-MONT-10') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'8' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLOVE-MONT-8' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'9' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLOVE-MONT-9' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'10' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLOVE-MONT-10' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 65: WWR-GLOVE-NITR
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-GLOVE-NITR','{"de":"Nitrilhandschuhe (100 Stk)","en":"Nitrile Gloves (100 pcs)"}','{"de":"Nitrilhandschuhe puderfrei in verschiedenen Groessen","en":"Nitrile gloves powder-free in various sizes"}',ARRAY['a0000000-4300-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='WWR-GLOVE-NITR' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-GLOVE-NITR' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-GLOVE-NITR-8','WWR-GLOVE-NITR-9','WWR-GLOVE-NITR-10') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'8' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLOVE-NITR-8' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'9' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLOVE-NITR-9' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'10' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-GLOVE-NITR-10' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 66: WWR-HELM
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-HELM','{"de":"Schutzhelm","en":"Hard Hat"}','{"de":"Industrieschutzhelm EN 397 in verschiedenen Farben","en":"Industrial hard hat EN 397 in various colors"}',ARRAY['a0000000-4200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'color',0 FROM products WHERE sku='WWR-HELM' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-HELM' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-HELM-WHT','WWR-HELM-YEL','WWR-HELM-RED') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'weiss' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-HELM-WHT' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='color' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'gelb' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-HELM-YEL' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='color' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'rot' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-HELM-RED' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='color' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 67: WWR-JACKET
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-JACKET','{"de":"Arbeitsjacke schwarz","en":"Work Jacket Black"}','{"de":"Robuste Arbeitsjacke Schwarz in verschiedenen Groessen","en":"Robust work jacket black in various sizes"}',ARRAY['a0000000-4100-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='WWR-JACKET' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-JACKET' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-JACKET-L-BLK','WWR-JACKET-XL-BLK','WWR-JACKET-XXL-BLK') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'l' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-JACKET-L-BLK' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'xl' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-JACKET-XL-BLK' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'xxl' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-JACKET-XXL-BLK' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 68: WWR-PANT
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-PANT','{"de":"Arbeitshose schwarz","en":"Work Trousers Black"}','{"de":"Robuste Arbeitshose Schwarz in verschiedenen Groessen","en":"Robust work trousers black in various sizes"}',ARRAY['a0000000-4100-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='WWR-PANT' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-PANT' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-PANT-50-BLK','WWR-PANT-52-BLK','WWR-PANT-54-BLK') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'50' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-PANT-50-BLK' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'52' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-PANT-52-BLK' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'54' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-PANT-54-BLK' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 69: WWR-RESP
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-RESP','{"de":"Atemschutzmaske NR","en":"Respirator NR"}','{"de":"Partikelfiltrierende Halbmasken in verschiedenen Schutzstufen","en":"Particle filtering half masks in various protection levels"}',ARRAY['a0000000-4200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'protection_lvl',0 FROM products WHERE sku='WWR-RESP' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-RESP' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-RESP-FFP2-NR','WWR-RESP-FFP3-NR') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'ffp2-nr' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-RESP-FFP2-NR' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='protection_lvl' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'ffp3-nr' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-RESP-FFP3-NR' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='protection_lvl' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 70: WWR-SHIRT
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-SHIRT','{"de":"Arbeitshemd Langarm blau","en":"Work Shirt Long Sleeve Blue"}','{"de":"Langarm-Arbeitshemd Blau in verschiedenen Groessen","en":"Long-sleeve work shirt blue in various sizes"}',ARRAY['a0000000-4100-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='WWR-SHIRT' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-SHIRT' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-SHIRT-L-BLU','WWR-SHIRT-XL-BLU') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'l' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-SHIRT-L-BLU' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'xl' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-SHIRT-XL-BLU' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 71: WWR-SHOE
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-SHOE','{"de":"Sicherheitsschuhe S3","en":"Safety Shoes S3"}','{"de":"Sicherheitsschuhe S3 mit Stahlkappe in verschiedenen Groessen","en":"Safety shoes S3 with steel toe cap in various sizes"}',ARRAY['a0000000-4200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='WWR-SHOE' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-SHOE' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-SHOE-42-S3','WWR-SHOE-43-S3','WWR-SHOE-44-S3') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'42' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-SHOE-42-S3' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'43' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-SHOE-43-S3' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'44' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-SHOE-44-S3' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- GROUP 72: WWR-VEST
INSERT INTO products (id,tenant_id,sku,name,description,category_ids,attributes,status,images,product_type,created_at,updated_at) VALUES (gen_random_uuid(),'00000000-0000-0000-0000-000000000001','WWR-VEST','{"de":"Warnweste EN ISO 20471","en":"High-Vis Vest EN ISO 20471"}','{"de":"Warnschutzweste EN ISO 20471 in verschiedenen Groessen","en":"High-visibility vest EN ISO 20471 in various sizes"}',ARRAY['a0000000-4200-0000-0000-000000000001']::uuid[],'[]'::jsonb,'active','[]'::jsonb,'variant_parent',NOW(),NOW()) ON CONFLICT (tenant_id,sku) DO NOTHING;
INSERT INTO variant_axes (product_id,attribute_code,position) SELECT id,'size',0 FROM products WHERE sku='WWR-VEST' AND tenant_id='00000000-0000-0000-0000-000000000001' ON CONFLICT (product_id,attribute_code) DO NOTHING;
UPDATE products SET product_type='variant', parent_id=(SELECT id FROM products WHERE sku='WWR-VEST' AND tenant_id='00000000-0000-0000-0000-000000000001') WHERE sku IN ('WWR-VEST-L-YEL','WWR-VEST-XL-YEL','WWR-VEST-XXL-ORG') AND tenant_id='00000000-0000-0000-0000-000000000001' AND product_type='simple';
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'l' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-VEST-L-YEL' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'xl' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-VEST-XL-YEL' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;
INSERT INTO variant_axis_values (variant_id,axis_id,option_code) SELECT v.id,va.id,'xxl' FROM products v JOIN variant_axes va ON va.product_id=v.parent_id WHERE v.sku='WWR-VEST-XXL-ORG' AND v.tenant_id='00000000-0000-0000-0000-000000000001' AND va.attribute_code='size' ON CONFLICT (variant_id,axis_id) DO NOTHING;

-- ============================================================
-- END OF MIGRATION 000007
-- Products remaining as simple (unique, no sensible group):
--   HYD-VLV-DR-10, LAB-CHEM-H2SO4-1L, LAB-CHEM-H3PO4-1L,
--   LAB-CHEM-HCL-1L, LAB-CHEM-KOH-1L, LAB-CHEM-NA2SO4-1KG,
--   LAB-CHEM-NACL-1KG, LAB-CHEM-NAOH-1L,
--   LAB-CLEAN-SOAP-5L, LAB-CLEAN-WIPES-BOX,
--   WWR-EAR-PLUG-FOAM, PKG-TAPE-DOUBLE-19
-- ============================================================
