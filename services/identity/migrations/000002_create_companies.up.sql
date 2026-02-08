-- Companies table
CREATE TABLE companies (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,

    -- SAP Mapping
    sap_company_number      VARCHAR(50) NOT NULL,
    sap_customer_group      VARCHAR(50),
    sap_shipping_plant      VARCHAR(50),
    sap_office              VARCHAR(50),
    sap_payment_type        VARCHAR(50),
    sap_price_group         VARCHAR(50),

    -- Profile
    name                    VARCHAR(255) NOT NULL,
    description             TEXT,
    email                   VARCHAR(255),
    currency                VARCHAR(3) DEFAULT 'CHF',

    -- Address
    street                  VARCHAR(255),
    house_number            VARCHAR(20),
    zip                     VARCHAR(20),
    city                    VARCHAR(100),
    country                 VARCHAR(2) DEFAULT 'CH',

    -- Contact
    phone                   VARCHAR(50),
    fax                     VARCHAR(50),
    url                     VARCHAR(255),

    -- Config
    config                  JSONB DEFAULT '{}',
    desired_delivery_days   TEXT[],
    default_shipping_note   TEXT,
    disable_order_feature   BOOLEAN DEFAULT false,

    -- Branding
    custom_primary_color    VARCHAR(7),
    custom_secondary_color  VARCHAR(7),

    -- Status
    is_active               BOOLEAN DEFAULT false,
    created_at              TIMESTAMPTZ DEFAULT NOW(),
    updated_at              TIMESTAMPTZ DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ,

    CONSTRAINT uq_companies_tenant_sap UNIQUE (tenant_id, sap_company_number)
);

CREATE INDEX idx_companies_tenant ON companies(tenant_id);
CREATE INDEX idx_companies_sap ON companies(sap_company_number);
CREATE INDEX idx_companies_active ON companies(tenant_id, is_active) WHERE is_active = true AND deleted_at IS NULL;

CREATE TRIGGER update_companies_updated_at
    BEFORE UPDATE ON companies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
