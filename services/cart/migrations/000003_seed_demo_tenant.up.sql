-- Seed demo tenant (idempotent)
INSERT INTO tenants (id, code, name, config, is_active, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001'::uuid,
    'demo',
    'Demo Tenant',
    '{}'::jsonb,
    true,
    NOW(),
    NOW()
)
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    is_active = EXCLUDED.is_active,
    updated_at = NOW();
