-- Remove seed variant products
DELETE FROM variant_axis_values WHERE variant_id IN (
    SELECT id FROM products WHERE parent_id IN (
        SELECT id FROM products WHERE sku IN ('IM3000', 'SAFEGRIP-PRO')
    )
);

DELETE FROM variant_axes WHERE product_id IN (
    SELECT id FROM products WHERE sku IN ('IM3000', 'SAFEGRIP-PRO')
);

DELETE FROM products WHERE parent_id IN (
    SELECT id FROM products WHERE sku IN ('IM3000', 'SAFEGRIP-PRO')
);

DELETE FROM products WHERE sku IN ('IM3000', 'SAFEGRIP-PRO');
