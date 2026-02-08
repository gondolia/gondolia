package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
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

	var t domain.Tenant
	var configJSON []byte

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Code, &t.Name, &configJSON, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTenantNotFound
		}
		return nil, err
	}

	if configJSON != nil {
		if err := json.Unmarshal(configJSON, &t.Config); err != nil {
			return nil, err
		}
	}

	return &t, nil
}

func (r *TenantRepository) GetByCode(ctx context.Context, code string) (*domain.Tenant, error) {
	query := `
		SELECT id, code, name, config, is_active, created_at, updated_at
		FROM tenants
		WHERE code = $1
	`

	var t domain.Tenant
	var configJSON []byte

	err := r.db.Pool.QueryRow(ctx, query, code).Scan(
		&t.ID, &t.Code, &t.Name, &configJSON, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTenantNotFound
		}
		return nil, err
	}

	if configJSON != nil {
		if err := json.Unmarshal(configJSON, &t.Config); err != nil {
			return nil, err
		}
	}

	return &t, nil
}

func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		INSERT INTO tenants (id, code, name, config, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	configJSON, err := json.Marshal(tenant.Config)
	if err != nil {
		return err
	}

	_, err = r.db.Pool.Exec(ctx, query,
		tenant.ID, tenant.Code, tenant.Name, configJSON, tenant.IsActive, tenant.CreatedAt, tenant.UpdatedAt,
	)
	return err
}

func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		UPDATE tenants
		SET name = $2, config = $3, is_active = $4, updated_at = NOW()
		WHERE id = $1
	`

	configJSON, err := json.Marshal(tenant.Config)
	if err != nil {
		return err
	}

	result, err := r.db.Pool.Exec(ctx, query, tenant.ID, tenant.Name, configJSON, tenant.IsActive)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTenantNotFound
	}

	return nil
}
