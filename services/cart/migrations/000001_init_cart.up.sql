-- Carts table
CREATE TABLE IF NOT EXISTS carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    user_id UUID,
    session_id TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT check_cart_status CHECK (status IN ('active', 'merged', 'completed')),
    CONSTRAINT check_cart_owner CHECK (
        (user_id IS NOT NULL AND session_id IS NULL) OR
        (user_id IS NULL AND session_id IS NOT NULL)
    )
);

-- Indexes for carts
CREATE INDEX idx_carts_tenant ON carts(tenant_id);
CREATE INDEX idx_carts_user ON carts(tenant_id, user_id, status) WHERE user_id IS NOT NULL;
CREATE INDEX idx_carts_session ON carts(tenant_id, session_id, status) WHERE session_id IS NOT NULL;

-- Cart items table
CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    variant_id UUID,
    product_type VARCHAR(30) NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    configuration JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT check_product_type CHECK (product_type IN ('simple', 'variant', 'bundle', 'parametric'))
);

-- Indexes for cart items
CREATE INDEX idx_cart_items_cart ON cart_items(cart_id);
CREATE INDEX idx_cart_items_product ON cart_items(product_id);

-- Comments
COMMENT ON TABLE carts IS 'Shopping carts for users and guest sessions';
COMMENT ON TABLE cart_items IS 'Items in shopping carts';
COMMENT ON COLUMN carts.user_id IS 'User ID for authenticated users (NULL for guests)';
COMMENT ON COLUMN carts.session_id IS 'Session ID for guest users (NULL for authenticated)';
COMMENT ON COLUMN cart_items.configuration IS 'JSONB configuration for bundle components and parametric parameters';
