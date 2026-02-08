-- User-Company pivot table
-- User types: 0 = Admin, 1 = SalesRep, 2 = Invited
CREATE TABLE user_companies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id      UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    role_id         UUID REFERENCES roles(id) ON DELETE SET NULL,

    user_type       SMALLINT NOT NULL DEFAULT 2,

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),

    CONSTRAINT uq_user_company UNIQUE (user_id, company_id),
    CONSTRAINT chk_user_type CHECK (user_type IN (0, 1, 2))
);

CREATE INDEX idx_user_companies_user ON user_companies(user_id);
CREATE INDEX idx_user_companies_company ON user_companies(company_id);
CREATE INDEX idx_user_companies_type ON user_companies(user_type);
CREATE INDEX idx_user_companies_role ON user_companies(role_id) WHERE role_id IS NOT NULL;

CREATE TRIGGER update_user_companies_updated_at
    BEFORE UPDATE ON user_companies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
