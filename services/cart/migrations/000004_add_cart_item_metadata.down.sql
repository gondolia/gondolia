-- Remove product metadata fields from cart_items table
DROP INDEX IF EXISTS idx_cart_items_sku;

ALTER TABLE cart_items
  DROP COLUMN IF EXISTS product_name,
  DROP COLUMN IF EXISTS sku,
  DROP COLUMN IF EXISTS image_url,
  DROP COLUMN IF EXISTS total_price;
