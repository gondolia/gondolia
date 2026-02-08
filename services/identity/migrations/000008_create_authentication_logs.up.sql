-- Authentication event logs
CREATE TABLE authentication_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,

    event_type      VARCHAR(50) NOT NULL,
    ip_address      VARCHAR(45),
    user_agent      VARCHAR(500),

    metadata        JSONB DEFAULT '{}',

    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Event types: login, logout, failed_login, password_reset, password_changed, invitation_sent, invitation_accepted
CREATE INDEX idx_auth_logs_tenant ON authentication_logs(tenant_id);
CREATE INDEX idx_auth_logs_user ON authentication_logs(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_auth_logs_created ON authentication_logs(created_at);
CREATE INDEX idx_auth_logs_event ON authentication_logs(event_type, created_at);

-- Partitioning hint: Consider partitioning by created_at for large deployments
COMMENT ON TABLE authentication_logs IS 'Audit log for authentication events. Consider partitioning by created_at monthly.';
