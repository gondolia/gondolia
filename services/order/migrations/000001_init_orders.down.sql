-- Drop tables in reverse order (respect foreign keys)
DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
-- Note: We don't drop tenants table as it might be shared
