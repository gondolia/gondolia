-- Roles table for RBAC
CREATE TABLE roles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    company_id      UUID REFERENCES companies(id) ON DELETE CASCADE,  -- NULL = System Role

    name            VARCHAR(100) NOT NULL,
    permissions     JSONB NOT NULL DEFAULT '{}',
    is_system       BOOLEAN DEFAULT false,

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT uq_roles_tenant_company_name UNIQUE (tenant_id, company_id, name)
);

CREATE INDEX idx_roles_tenant ON roles(tenant_id);
CREATE INDEX idx_roles_company ON roles(company_id);
CREATE INDEX idx_roles_system ON roles(tenant_id, is_system) WHERE is_system = true;

CREATE TRIGGER update_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
