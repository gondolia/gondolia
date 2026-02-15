-- Add variant support to products table
ALTER TABLE products ADD COLUMN product_type VARCHAR(30) NOT NULL DEFAULT 'simple';
-- Constraint: product_type must be one of 'simple', 'variant_parent', 'variant'
ALTER TABLE products ADD CONSTRAINT check_product_type 
    CHECK (product_type IN ('simple', 'variant_parent', 'variant'));

-- Add parent_id for variant products
ALTER TABLE products ADD COLUMN parent_id UUID REFERENCES products(id) ON DELETE CASCADE;
-- Constraint: only variant type can have parent_id
ALTER TABLE products ADD CONSTRAINT check_variant_parent
    CHECK (
        (product_type = 'variant' AND parent_id IS NOT NULL) OR
        (product_type != 'variant' AND parent_id IS NULL)
    );

-- Indexes for performance
CREATE INDEX idx_products_parent ON products(parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX idx_products_type ON products(tenant_id, product_type);

-- Variant axes definition table
-- Defines which attributes are variant axes for a parent product
CREATE TABLE variant_axes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    attribute_code VARCHAR(100) NOT NULL,
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(product_id, attribute_code),
    CONSTRAINT fk_variant_axes_product FOREIGN KEY (product_id) 
        REFERENCES products(id) ON DELETE CASCADE
);

-- Constraint: product_id must reference a variant_parent
CREATE INDEX idx_variant_axes_product ON variant_axes(product_id);

-- Variant axis values table
-- Stores the specific axis value for each variant
CREATE TABLE variant_axis_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    variant_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    axis_id UUID NOT NULL REFERENCES variant_axes(id) ON DELETE CASCADE,
    option_code VARCHAR(200) NOT NULL,

    UNIQUE(variant_id, axis_id),
    CONSTRAINT fk_variant_axis_values_variant FOREIGN KEY (variant_id)
        REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT fk_variant_axis_values_axis FOREIGN KEY (axis_id)
        REFERENCES variant_axes(id) ON DELETE CASCADE
);

CREATE INDEX idx_variant_axis_values_variant ON variant_axis_values(variant_id);
CREATE INDEX idx_variant_axis_values_axis ON variant_axis_values(axis_id, option_code);

-- Set existing products to 'simple' type (already done by DEFAULT, but explicit for clarity)
UPDATE products SET product_type = 'simple' WHERE product_type IS NULL;

COMMENT ON COLUMN products.product_type IS 'Type: simple (standalone), variant_parent (has variants), variant (child of parent)';
COMMENT ON COLUMN products.parent_id IS 'For variant type: reference to parent product';
COMMENT ON TABLE variant_axes IS 'Defines which attributes form the variant axes for a parent product';
COMMENT ON TABLE variant_axis_values IS 'Stores the axis values for each variant product';
