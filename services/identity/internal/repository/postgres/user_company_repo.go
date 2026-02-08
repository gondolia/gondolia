package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
)

type UserCompanyRepository struct {
	db *DB
}

func NewUserCompanyRepository(db *DB) *UserCompanyRepository {
	return &UserCompanyRepository{db: db}
}

func (r *UserCompanyRepository) GetByUserAndCompany(ctx context.Context, userID, companyID uuid.UUID) (*domain.UserCompany, error) {
	query := `
		SELECT uc.id, uc.user_id, uc.company_id, uc.role_id, uc.user_type, uc.created_at, uc.updated_at
		FROM user_companies uc
		WHERE uc.user_id = $1 AND uc.company_id = $2
	`

	var uc domain.UserCompany
	err := r.db.Pool.QueryRow(ctx, query, userID, companyID).Scan(
		&uc.ID, &uc.UserID, &uc.CompanyID, &uc.RoleID, &uc.UserType, &uc.CreatedAt, &uc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotInCompany
		}
		return nil, err
	}

	return &uc, nil
}

func (r *UserCompanyRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.UserCompany, error) {
	query := `
		SELECT uc.id, uc.user_id, uc.company_id, uc.role_id, uc.user_type, uc.created_at, uc.updated_at,
		       c.id, c.tenant_id, c.sap_company_number, c.name, c.is_active,
		       r.id, r.tenant_id, r.company_id, r.name, r.permissions, r.is_system
		FROM user_companies uc
		JOIN companies c ON c.id = uc.company_id AND c.deleted_at IS NULL
		LEFT JOIN roles r ON r.id = uc.role_id
		WHERE uc.user_id = $1 AND c.is_active = true
		ORDER BY c.name
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.UserCompany
	for rows.Next() {
		var uc domain.UserCompany
		var company domain.Company
		var role domain.Role
		var roleID, roleTenantID, roleCompanyID *uuid.UUID
		var roleName *string
		var rolePermJSON []byte
		var roleIsSystem *bool

		err := rows.Scan(
			&uc.ID, &uc.UserID, &uc.CompanyID, &uc.RoleID, &uc.UserType, &uc.CreatedAt, &uc.UpdatedAt,
			&company.ID, &company.TenantID, &company.SAPCompanyNumber, &company.Name, &company.IsActive,
			&roleID, &roleTenantID, &roleCompanyID, &roleName, &rolePermJSON, &roleIsSystem,
		)
		if err != nil {
			return nil, err
		}

		uc.Company = &company

		if roleID != nil {
			role.ID = *roleID
			role.TenantID = *roleTenantID
			role.CompanyID = roleCompanyID
			role.Name = *roleName
			role.IsSystem = *roleIsSystem
			if rolePermJSON != nil {
				_ = json.Unmarshal(rolePermJSON, &role.Permissions)
			}
			uc.Role = &role
		}

		results = append(results, uc)
	}

	return results, rows.Err()
}

func (r *UserCompanyRepository) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.UserCompany, error) {
	query := `
		SELECT uc.id, uc.user_id, uc.company_id, uc.role_id, uc.user_type, uc.created_at, uc.updated_at,
		       u.id, u.tenant_id, u.email, u.firstname, u.lastname, u.is_active,
		       r.id, r.tenant_id, r.company_id, r.name, r.permissions, r.is_system
		FROM user_companies uc
		JOIN users u ON u.id = uc.user_id AND u.deleted_at IS NULL
		LEFT JOIN roles r ON r.id = uc.role_id
		WHERE uc.company_id = $1
		ORDER BY u.lastname, u.firstname
	`

	rows, err := r.db.Pool.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.UserCompany
	for rows.Next() {
		var uc domain.UserCompany
		var user domain.User
		var role domain.Role
		var roleID, roleTenantID, roleCompanyID *uuid.UUID
		var roleName *string
		var rolePermJSON []byte
		var roleIsSystem *bool

		err := rows.Scan(
			&uc.ID, &uc.UserID, &uc.CompanyID, &uc.RoleID, &uc.UserType, &uc.CreatedAt, &uc.UpdatedAt,
			&user.ID, &user.TenantID, &user.Email, &user.FirstName, &user.LastName, &user.IsActive,
			&roleID, &roleTenantID, &roleCompanyID, &roleName, &rolePermJSON, &roleIsSystem,
		)
		if err != nil {
			return nil, err
		}

		uc.User = &user

		if roleID != nil {
			role.ID = *roleID
			role.TenantID = *roleTenantID
			role.CompanyID = roleCompanyID
			role.Name = *roleName
			role.IsSystem = *roleIsSystem
			if rolePermJSON != nil {
				_ = json.Unmarshal(rolePermJSON, &role.Permissions)
			}
			uc.Role = &role
		}

		results = append(results, uc)
	}

	return results, rows.Err()
}

func (r *UserCompanyRepository) Create(ctx context.Context, uc *domain.UserCompany) error {
	query := `
		INSERT INTO user_companies (id, user_id, company_id, role_id, user_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		uc.ID, uc.UserID, uc.CompanyID, uc.RoleID, uc.UserType, uc.CreatedAt, uc.UpdatedAt,
	)
	return err
}

func (r *UserCompanyRepository) Update(ctx context.Context, uc *domain.UserCompany) error {
	query := `
		UPDATE user_companies SET
			role_id = $3, user_type = $4, updated_at = NOW()
		WHERE user_id = $1 AND company_id = $2
	`

	result, err := r.db.Pool.Exec(ctx, query, uc.UserID, uc.CompanyID, uc.RoleID, uc.UserType)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotInCompany
	}

	return nil
}

func (r *UserCompanyRepository) Delete(ctx context.Context, userID, companyID uuid.UUID) error {
	query := `DELETE FROM user_companies WHERE user_id = $1 AND company_id = $2`

	result, err := r.db.Pool.Exec(ctx, query, userID, companyID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotInCompany
	}

	return nil
}

func (r *UserCompanyRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM user_companies uc
		JOIN companies c ON c.id = uc.company_id AND c.deleted_at IS NULL AND c.is_active = true
		WHERE uc.user_id = $1
	`

	var count int
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}
