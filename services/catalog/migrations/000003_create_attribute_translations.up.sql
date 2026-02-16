-- Create attribute translations table
CREATE TABLE attribute_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    attribute_key VARCHAR(100) NOT NULL,
    locale VARCHAR(2) NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    unit VARCHAR(50),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Unique constraint: one translation per attribute key and locale per tenant
    UNIQUE(tenant_id, attribute_key, locale)
);

-- Indexes
CREATE INDEX idx_attribute_translations_tenant_locale ON attribute_translations(tenant_id, locale);
CREATE INDEX idx_attribute_translations_key ON attribute_translations(attribute_key);

-- Comments
COMMENT ON TABLE attribute_translations IS 'Translations for product attribute keys (e.g., thickness_mm -> Dicke)';
COMMENT ON COLUMN attribute_translations.attribute_key IS 'Technical attribute key from product.attributes (e.g., thickness_mm, voltage)';
COMMENT ON COLUMN attribute_translations.locale IS 'ISO 639-1 language code (de, en, fr, it)';
COMMENT ON COLUMN attribute_translations.display_name IS 'Human-readable translated name';
COMMENT ON COLUMN attribute_translations.unit IS 'Optional unit (e.g., mm, V, kg)';
