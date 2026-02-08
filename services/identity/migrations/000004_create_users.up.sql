-- Users table (formerly customers)
CREATE TABLE users (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,

    -- Status
    is_active               BOOLEAN DEFAULT false,
    is_imported             BOOLEAN DEFAULT false,
    is_salesmaster          BOOLEAN DEFAULT false,
    sso_only                BOOLEAN DEFAULT false,

    -- SAP Mapping
    sap_user_id             VARCHAR(50),
    sap_customer_number     VARCHAR(50),

    -- Profile
    email                   VARCHAR(255) NOT NULL,
    password_hash           VARCHAR(255),
    firstname               VARCHAR(100) NOT NULL,
    lastname                VARCHAR(100) NOT NULL,
    phone                   VARCHAR(50),
    mobile                  VARCHAR(50),
    default_language        VARCHAR(10) DEFAULT 'de',

    -- Company Context
    default_company_id      UUID REFERENCES companies(id) ON DELETE SET NULL,

    -- Invitation
    invitation_token        VARCHAR(100) UNIQUE,
    invited_at              TIMESTAMPTZ,

    -- Tracking
    last_login_at           TIMESTAMPTZ,
    created_at              TIMESTAMPTZ DEFAULT NOW(),
    updated_at              TIMESTAMPTZ DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ,

    CONSTRAINT uq_users_tenant_email UNIQUE (tenant_id, email)
);

CREATE INDEX idx_users_tenant ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_sap_customer ON users(sap_customer_number) WHERE sap_customer_number IS NOT NULL;
CREATE INDEX idx_users_invitation ON users(invitation_token) WHERE invitation_token IS NOT NULL;
CREATE INDEX idx_users_active ON users(tenant_id, is_active) WHERE is_active = true AND deleted_at IS NULL;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
