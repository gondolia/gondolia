-- ============================================================
-- Migration 000007 ROLLBACK: Revert variants back to simple
-- NOTE: IM3000 and SAFEGRIP-PRO are NOT touched.
-- ============================================================

-- All variant_axis_values for converted products are deleted via CASCADE
-- when we delete the variant_parent. But we need to restore product_type
-- first, then delete parents.

-- Step 1: Restore all variants to simple
UPDATE products SET product_type = 'simple', parent_id = NULL
WHERE product_type = 'variant'
  AND parent_id IN (
    SELECT id FROM products
    WHERE product_type = 'variant_parent'
      AND sku NOT IN ('IM3000', 'SAFEGRIP-PRO')
      AND tenant_id = '00000000-0000-0000-0000-000000000001'
  )
  AND tenant_id = '00000000-0000-0000-0000-000000000001';

-- Step 2: Delete variant_parents (CASCADE removes axes + axis_values)
DELETE FROM products
WHERE product_type = 'variant_parent'
  AND sku NOT IN ('IM3000', 'SAFEGRIP-PRO')
  AND tenant_id = '00000000-0000-0000-0000-000000000001';

-- Step 3: Remove new attribute_translations
DELETE FROM attribute_translations
WHERE tenant_id = '00000000-0000-0000-0000-000000000001'
  AND attribute_key IN (
    'power_kw', 'wattage', 'bearing_type', 'belt_type', 'chain_type',
    'conductors', 'cross_section', 'cable_category', 'characteristic',
    'poles', 'nennstrom', 'volume', 'viscosity', 'flow_rate', 'bore_stroke',
    'filtration', 'nominal_diam', 'connection', 'size', 'protection_lvl',
    'snr', 'tint', 'format', 'valve_type', 'hose_thread', 'pal_size'
  );
