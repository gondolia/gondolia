-- Seed example variant products based on architecture document

-- First, let's assume we have a tenant (you may need to adjust tenant_id)
-- And some existing products to convert

-- 1. Create variant parent for Industriemotor
INSERT INTO products (
    id, tenant_id, product_type, parent_id, sku, name, description, 
    category_ids, attributes, status, images, created_at, updated_at
) VALUES (
    gen_random_uuid(),
    (SELECT id FROM tenants LIMIT 1), -- Use existing tenant
    'variant_parent',
    NULL,
    'IM3000',
    '{"de": "TurboTech Industriemotor IM3000", "en": "TurboTech Industrial Motor IM3000"}',
    '{"de": "Leistungsstarker Drehstrommotor für industrielle Anwendungen", "en": "Powerful three-phase motor for industrial applications"}',
    '{}',
    '[]',
    'active',
    '[]',
    NOW(),
    NOW()
) ON CONFLICT DO NOTHING;

-- Define variant axes for Industriemotor
INSERT INTO variant_axes (id, product_id, attribute_code, position)
SELECT 
    gen_random_uuid(),
    p.id,
    'power_rating',
    0
FROM products p
WHERE p.sku = 'IM3000' AND p.product_type = 'variant_parent'
ON CONFLICT DO NOTHING;

INSERT INTO variant_axes (id, product_id, attribute_code, position)
SELECT 
    gen_random_uuid(),
    p.id,
    'voltage',
    1
FROM products p
WHERE p.sku = 'IM3000' AND p.product_type = 'variant_parent'
ON CONFLICT DO NOTHING;

INSERT INTO variant_axes (id, product_id, attribute_code, position)
SELECT 
    gen_random_uuid(),
    p.id,
    'mounting',
    2
FROM products p
WHERE p.sku = 'IM3000' AND p.product_type = 'variant_parent'
ON CONFLICT DO NOTHING;

-- Create variants for Industriemotor
-- Variant 1: 1.5 kW, 230V, B3
INSERT INTO products (
    id, tenant_id, product_type, parent_id, sku, name, description,
    category_ids, attributes, status, images, created_at, updated_at
) 
SELECT 
    gen_random_uuid(),
    (SELECT id FROM tenants LIMIT 1),
    'variant',
    p.id,
    'IM3000-15-230-B3',
    '{"de": "TurboTech IM3000 - 1.5 kW, 230V, B3"}',
    '{}',
    '{}',
    '[{"key": "weight", "type": "number", "value": 12.5}]',
    'active',
    '[]',
    NOW(),
    NOW()
FROM products p
WHERE p.sku = 'IM3000' AND p.product_type = 'variant_parent'
ON CONFLICT DO NOTHING;

-- Axis values for Variant 1
INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT 
    v.id,
    va.id,
    '1_5kw'
FROM products v
JOIN products p ON v.parent_id = p.id
JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'IM3000-15-230-B3' 
  AND va.attribute_code = 'power_rating'
ON CONFLICT DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT 
    v.id,
    va.id,
    '230v'
FROM products v
JOIN products p ON v.parent_id = p.id
JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'IM3000-15-230-B3' 
  AND va.attribute_code = 'voltage'
ON CONFLICT DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT 
    v.id,
    va.id,
    'b3'
FROM products v
JOIN products p ON v.parent_id = p.id
JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'IM3000-15-230-B3' 
  AND va.attribute_code = 'mounting'
ON CONFLICT DO NOTHING;

-- Variant 2: 1.5 kW, 230V, B5
INSERT INTO products (
    id, tenant_id, product_type, parent_id, sku, name, description,
    category_ids, attributes, status, images, created_at, updated_at
) 
SELECT 
    gen_random_uuid(),
    (SELECT id FROM tenants LIMIT 1),
    'variant',
    p.id,
    'IM3000-15-230-B5',
    '{"de": "TurboTech IM3000 - 1.5 kW, 230V, B5"}',
    '{}',
    '{}',
    '[{"key": "weight", "type": "number", "value": 13.0}]',
    'active',
    '[]',
    NOW(),
    NOW()
FROM products p
WHERE p.sku = 'IM3000' AND p.product_type = 'variant_parent'
ON CONFLICT DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '1_5kw'
FROM products v
JOIN products p ON v.parent_id = p.id
JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'IM3000-15-230-B5' AND va.attribute_code = 'power_rating'
ON CONFLICT DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '230v'
FROM products v
JOIN products p ON v.parent_id = p.id
JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'IM3000-15-230-B5' AND va.attribute_code = 'voltage'
ON CONFLICT DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'b5'
FROM products v
JOIN products p ON v.parent_id = p.id
JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'IM3000-15-230-B5' AND va.attribute_code = 'mounting'
ON CONFLICT DO NOTHING;

-- Variant 3: 4.0 kW, 400V, B3
INSERT INTO products (
    id, tenant_id, product_type, parent_id, sku, name, description,
    category_ids, attributes, status, images, created_at, updated_at
) 
SELECT 
    gen_random_uuid(),
    (SELECT id FROM tenants LIMIT 1),
    'variant',
    p.id,
    'IM3000-40-400-B3',
    '{"de": "TurboTech IM3000 - 4.0 kW, 400V, B3"}',
    '{}',
    '{}',
    '[{"key": "weight", "type": "number", "value": 38.0}]',
    'active',
    '[]',
    NOW(),
    NOW()
FROM products p
WHERE p.sku = 'IM3000' AND p.product_type = 'variant_parent'
ON CONFLICT DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '4_0kw'
FROM products v
JOIN products p ON v.parent_id = p.id
JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'IM3000-40-400-B3' AND va.attribute_code = 'power_rating'
ON CONFLICT DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, '400v'
FROM products v
JOIN products p ON v.parent_id = p.id
JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'IM3000-40-400-B3' AND va.attribute_code = 'voltage'
ON CONFLICT DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'b3'
FROM products v
JOIN products p ON v.parent_id = p.id
JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'IM3000-40-400-B3' AND va.attribute_code = 'mounting'
ON CONFLICT DO NOTHING;

-- 2. Create variant parent for Sicherheitshandschuhe
INSERT INTO products (
    id, tenant_id, product_type, parent_id, sku, name, description,
    category_ids, attributes, status, images, created_at, updated_at
) VALUES (
    gen_random_uuid(),
    (SELECT id FROM tenants LIMIT 1),
    'variant_parent',
    NULL,
    'SAFEGRIP-PRO',
    '{"de": "SafeGrip Pro Sicherheitshandschuhe", "en": "SafeGrip Pro Safety Gloves"}',
    '{"de": "Hochwertige Schutzhandschuhe für professionellen Einsatz", "en": "High-quality protective gloves for professional use"}',
    '{}',
    '[]',
    'active',
    '[]',
    NOW(),
    NOW()
) ON CONFLICT DO NOTHING;

-- Define variant axis for Handschuhe (size)
INSERT INTO variant_axes (id, product_id, attribute_code, position)
SELECT 
    gen_random_uuid(),
    p.id,
    'size',
    0
FROM products p
WHERE p.sku = 'SAFEGRIP-PRO' AND p.product_type = 'variant_parent'
ON CONFLICT DO NOTHING;

-- Create glove variants (S, M, L, XL)
INSERT INTO products (id, tenant_id, product_type, parent_id, sku, name, description, category_ids, attributes, status, images, created_at, updated_at)
SELECT gen_random_uuid(), (SELECT id FROM tenants LIMIT 1), 'variant', p.id, 'SAFEGRIP-PRO-S', '{"de": "SafeGrip Pro - Größe S"}', '{}', '{}', '[]', 'active', '[]', NOW(), NOW()
FROM products p WHERE p.sku = 'SAFEGRIP-PRO' AND p.product_type = 'variant_parent' ON CONFLICT DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 's' FROM products v JOIN products p ON v.parent_id = p.id JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'SAFEGRIP-PRO-S' AND va.attribute_code = 'size' ON CONFLICT DO NOTHING;

INSERT INTO products (id, tenant_id, product_type, parent_id, sku, name, description, category_ids, attributes, status, images, created_at, updated_at)
SELECT gen_random_uuid(), (SELECT id FROM tenants LIMIT 1), 'variant', p.id, 'SAFEGRIP-PRO-M', '{"de": "SafeGrip Pro - Größe M"}', '{}', '{}', '[]', 'active', '[]', NOW(), NOW()
FROM products p WHERE p.sku = 'SAFEGRIP-PRO' AND p.product_type = 'variant_parent' ON CONFLICT DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'm' FROM products v JOIN products p ON v.parent_id = p.id JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'SAFEGRIP-PRO-M' AND va.attribute_code = 'size' ON CONFLICT DO NOTHING;

INSERT INTO products (id, tenant_id, product_type, parent_id, sku, name, description, category_ids, attributes, status, images, created_at, updated_at)
SELECT gen_random_uuid(), (SELECT id FROM tenants LIMIT 1), 'variant', p.id, 'SAFEGRIP-PRO-L', '{"de": "SafeGrip Pro - Größe L"}', '{}', '{}', '[]', 'active', '[]', NOW(), NOW()
FROM products p WHERE p.sku = 'SAFEGRIP-PRO' AND p.product_type = 'variant_parent' ON CONFLICT DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'l' FROM products v JOIN products p ON v.parent_id = p.id JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'SAFEGRIP-PRO-L' AND va.attribute_code = 'size' ON CONFLICT DO NOTHING;

INSERT INTO products (id, tenant_id, product_type, parent_id, sku, name, description, category_ids, attributes, status, images, created_at, updated_at)
SELECT gen_random_uuid(), (SELECT id FROM tenants LIMIT 1), 'variant', p.id, 'SAFEGRIP-PRO-XL', '{"de": "SafeGrip Pro - Größe XL"}', '{}', '{}', '[]', 'active', '[]', NOW(), NOW()
FROM products p WHERE p.sku = 'SAFEGRIP-PRO' AND p.product_type = 'variant_parent' ON CONFLICT DO NOTHING;

INSERT INTO variant_axis_values (variant_id, axis_id, option_code)
SELECT v.id, va.id, 'xl' FROM products v JOIN products p ON v.parent_id = p.id JOIN variant_axes va ON va.product_id = p.id
WHERE v.sku = 'SAFEGRIP-PRO-XL' AND va.attribute_code = 'size' ON CONFLICT DO NOTHING;

COMMENT ON TABLE products IS 'Catalog products - includes simple, variant_parent, and variant types';
COMMENT ON TABLE variant_axes IS 'Defines variant axes for parent products';
COMMENT ON TABLE variant_axis_values IS 'Stores axis values for each variant';
