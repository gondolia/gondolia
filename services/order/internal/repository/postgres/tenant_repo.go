package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gondolia/gondolia/services/order/internal/domain"
)

type TenantRepository struct {
	db *DB
}

func NewTenantRepository(db *DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	query := `
		SELECT id, code, name, config, is_active, created_at, updated_at
		FROM tenants
		WHERE id = $1
	`

	return r.scanTenant(r.db.Pool.QueryRow(ctx, query, id))
}

func (r *TenantRepository) GetByCode(ctx context.Context, code string) (*domain.Tenant, error) {
	query := `
		SELECT id, code, name, config, is_active, created_at, updated_at
		FROM tenants
		WHERE code = $1
	`

	return r.scanTenant(r.db.Pool.QueryRow(ctx, query, code))
}

func (r *TenantRepository) scanTenant(row pgx.Row) (*domain.Tenant, error) {
	var tenant domain.Tenant
	err := row.Scan(
		&tenant.ID,
		&tenant.Code,
		&tenant.Name,
		&tenant.Config,
		&tenant.IsActive,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrTenantNotFound
		}
		return nil, err
	}
	return &tenant, nil
}
