package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
)

type RoleRepository struct {
	db *DB
}

func NewRoleRepository(db *DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	query := `
		SELECT id, tenant_id, company_id, name, permissions, is_system, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	return r.scanRole(r.db.Pool.QueryRow(ctx, query, id))
}

func (r *RoleRepository) GetByName(ctx context.Context, tenantID uuid.UUID, companyID *uuid.UUID, name string) (*domain.Role, error) {
	var query string
	var args []any

	if companyID != nil {
		query = `
			SELECT id, tenant_id, company_id, name, permissions, is_system, created_at, updated_at
			FROM roles
			WHERE tenant_id = $1 AND company_id = $2 AND LOWER(name) = LOWER($3)
		`
		args = []any{tenantID, *companyID, name}
	} else {
		query = `
			SELECT id, tenant_id, company_id, name, permissions, is_system, created_at, updated_at
			FROM roles
			WHERE tenant_id = $1 AND company_id IS NULL AND LOWER(name) = LOWER($2)
		`
		args = []any{tenantID, name}
	}

	return r.scanRole(r.db.Pool.QueryRow(ctx, query, args...))
}

func (r *RoleRepository) List(ctx context.Context, filter domain.RoleFilter) ([]domain.Role, int, error) {
	var conditions []string
	var args []any
	argNum := 1

	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argNum))
	args = append(args, filter.TenantID)
	argNum++

	if filter.CompanyID != nil {
		// Include both company-specific and system roles
		conditions = append(conditions, fmt.Sprintf("(company_id = $%d OR company_id IS NULL)", argNum))
		args = append(args, *filter.CompanyID)
		argNum++
	} else {
		// Only system roles
		conditions = append(conditions, "company_id IS NULL")
	}

	if filter.IsSystem != nil {
		conditions = append(conditions, fmt.Sprintf("is_system = $%d", argNum))
		args = append(args, *filter.IsSystem)
		argNum++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM roles WHERE %s", whereClause)
	var total int
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT id, tenant_id, company_id, name, permissions, is_system, created_at, updated_at
		FROM roles
		WHERE %s
		ORDER BY is_system DESC, name
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		role, err := r.scanRoleFromRows(rows)
		if err != nil {
			return nil, 0, err
		}
		roles = append(roles, *role)
	}

	return roles, total, rows.Err()
}

func (r *RoleRepository) ListSystemRoles(ctx context.Context, tenantID uuid.UUID) ([]domain.Role, error) {
	query := `
		SELECT id, tenant_id, company_id, name, permissions, is_system, created_at, updated_at
		FROM roles
		WHERE tenant_id = $1 AND is_system = true
		ORDER BY name
	`

	rows, err := r.db.Pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		role, err := r.scanRoleFromRows(rows)
		if err != nil {
			return nil, err
		}
		roles = append(roles, *role)
	}

	return roles, rows.Err()
}

func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	query := `
		INSERT INTO roles (id, tenant_id, company_id, name, permissions, is_system, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	permJSON, err := json.Marshal(role.Permissions)
	if err != nil {
		return err
	}

	_, err = r.db.Pool.Exec(ctx, query,
		role.ID, role.TenantID, role.CompanyID, role.Name, permJSON, role.IsSystem,
		role.CreatedAt, role.UpdatedAt,
	)
	return err
}

func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	query := `
		UPDATE roles SET
			name = $2, permissions = $3, updated_at = NOW()
		WHERE id = $1 AND is_system = false
	`

	permJSON, err := json.Marshal(role.Permissions)
	if err != nil {
		return err
	}

	result, err := r.db.Pool.Exec(ctx, query, role.ID, role.Name, permJSON)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		// Check if it's a system role
		var isSystem bool
		checkQuery := `SELECT is_system FROM roles WHERE id = $1`
		if err := r.db.Pool.QueryRow(ctx, checkQuery, role.ID).Scan(&isSystem); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrRoleNotFound
			}
			return err
		}
		if isSystem {
			return domain.ErrRoleIsSystem
		}
		return domain.ErrRoleNotFound
	}

	return nil
}

func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Don't allow deleting system roles
	query := `DELETE FROM roles WHERE id = $1 AND is_system = false`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		// Check if it's a system role
		var isSystem bool
		checkQuery := `SELECT is_system FROM roles WHERE id = $1`
		if err := r.db.Pool.QueryRow(ctx, checkQuery, id).Scan(&isSystem); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrRoleNotFound
			}
			return err
		}
		if isSystem {
			return domain.ErrRoleIsSystem
		}
		return domain.ErrRoleNotFound
	}

	return nil
}

func (r *RoleRepository) scanRole(row pgx.Row) (*domain.Role, error) {
	var role domain.Role
	var permJSON []byte

	err := row.Scan(
		&role.ID, &role.TenantID, &role.CompanyID, &role.Name,
		&permJSON, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRoleNotFound
		}
		return nil, err
	}

	if permJSON != nil {
		if err := json.Unmarshal(permJSON, &role.Permissions); err != nil {
			return nil, err
		}
	}

	return &role, nil
}

func (r *RoleRepository) scanRoleFromRows(rows pgx.Rows) (*domain.Role, error) {
	var role domain.Role
	var permJSON []byte

	err := rows.Scan(
		&role.ID, &role.TenantID, &role.CompanyID, &role.Name,
		&permJSON, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if permJSON != nil {
		if err := json.Unmarshal(permJSON, &role.Permissions); err != nil {
			return nil, err
		}
	}

	return &role, nil
}
