-- Add product metadata fields to cart_items table
ALTER TABLE cart_items
  ADD COLUMN product_name TEXT NOT NULL DEFAULT '',
  ADD COLUMN sku TEXT NOT NULL DEFAULT '',
  ADD COLUMN image_url TEXT NOT NULL DEFAULT '',
  ADD COLUMN total_price DECIMAL(10,2) NOT NULL DEFAULT 0;

-- Remove defaults after adding columns (for new inserts, values must be explicit)
ALTER TABLE cart_items
  ALTER COLUMN product_name DROP DEFAULT,
  ALTER COLUMN sku DROP DEFAULT,
  ALTER COLUMN image_url DROP DEFAULT,
  ALTER COLUMN total_price DROP DEFAULT;

-- Create index on SKU for faster lookups
CREATE INDEX idx_cart_items_sku ON cart_items(sku);
