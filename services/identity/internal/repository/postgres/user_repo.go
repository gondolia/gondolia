package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, is_active, is_imported, is_salesmaster, sso_only,
		       sap_user_id, sap_customer_number, email, password_hash, firstname, lastname,
		       phone, mobile, default_language, default_company_id, invitation_token,
		       invited_at, last_login_at, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	return r.scanUser(r.db.Pool.QueryRow(ctx, query, id))
}

func (r *UserRepository) GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, is_active, is_imported, is_salesmaster, sso_only,
		       sap_user_id, sap_customer_number, email, password_hash, firstname, lastname,
		       phone, mobile, default_language, default_company_id, invitation_token,
		       invited_at, last_login_at, created_at, updated_at, deleted_at
		FROM users
		WHERE tenant_id = $1 AND LOWER(email) = LOWER($2) AND deleted_at IS NULL
	`

	return r.scanUser(r.db.Pool.QueryRow(ctx, query, tenantID, email))
}

func (r *UserRepository) GetByInvitationToken(ctx context.Context, token string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, is_active, is_imported, is_salesmaster, sso_only,
		       sap_user_id, sap_customer_number, email, password_hash, firstname, lastname,
		       phone, mobile, default_language, default_company_id, invitation_token,
		       invited_at, last_login_at, created_at, updated_at, deleted_at
		FROM users
		WHERE invitation_token = $1 AND deleted_at IS NULL
	`

	return r.scanUser(r.db.Pool.QueryRow(ctx, query, token))
}

func (r *UserRepository) List(ctx context.Context, filter domain.UserFilter) ([]domain.User, int, error) {
	var conditions []string
	var args []any
	argNum := 1

	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argNum))
	args = append(args, filter.TenantID)
	argNum++

	conditions = append(conditions, "deleted_at IS NULL")

	if filter.CompanyID != nil {
		conditions = append(conditions, fmt.Sprintf(`
			id IN (SELECT user_id FROM user_companies WHERE company_id = $%d)
		`, argNum))
		args = append(args, *filter.CompanyID)
		argNum++
	}

	if filter.Email != nil {
		conditions = append(conditions, fmt.Sprintf("LOWER(email) = LOWER($%d)", argNum))
		args = append(args, *filter.Email)
		argNum++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argNum))
		args = append(args, *filter.IsActive)
		argNum++
	}

	if filter.Search != nil {
		conditions = append(conditions, fmt.Sprintf(`
			(LOWER(firstname) LIKE LOWER($%d) OR LOWER(lastname) LIKE LOWER($%d) OR LOWER(email) LIKE LOWER($%d))
		`, argNum, argNum, argNum))
		args = append(args, "%"+*filter.Search+"%")
		argNum++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s", whereClause)
	var total int
	if err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT id, tenant_id, is_active, is_imported, is_salesmaster, sso_only,
		       sap_user_id, sap_customer_number, email, password_hash, firstname, lastname,
		       phone, mobile, default_language, default_company_id, invitation_token,
		       invited_at, last_login_at, created_at, updated_at, deleted_at
		FROM users
		WHERE %s
		ORDER BY lastname, firstname
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		user, err := r.scanUserFromRows(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, *user)
	}

	return users, total, rows.Err()
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			id, tenant_id, is_active, is_imported, is_salesmaster, sso_only,
			sap_user_id, sap_customer_number, email, password_hash, firstname, lastname,
			phone, mobile, default_language, default_company_id, invitation_token,
			invited_at, last_login_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.TenantID, user.IsActive, user.IsImported, user.IsSalesMaster, user.SSOOnly,
		user.SAPUserID, user.SAPCustomerNumber, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Phone, user.Mobile, user.DefaultLanguage, user.DefaultCompanyID, user.InvitationToken,
		user.InvitedAt, user.LastLoginAt, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET
			is_active = $2, is_imported = $3, is_salesmaster = $4, sso_only = $5,
			sap_user_id = $6, sap_customer_number = $7, firstname = $8, lastname = $9,
			phone = $10, mobile = $11, default_language = $12, default_company_id = $13,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Pool.Exec(ctx, query,
		user.ID, user.IsActive, user.IsImported, user.IsSalesMaster, user.SSOOnly,
		user.SAPUserID, user.SAPCustomerNumber, user.FirstName, user.LastName,
		user.Phone, user.Mobile, user.DefaultLanguage, user.DefaultCompanyID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_login_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

func (r *UserRepository) SetPassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1`
	result, err := r.db.Pool.Exec(ctx, query, id, passwordHash)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) ClearInvitation(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET invitation_token = NULL, invited_at = NULL, updated_at = NOW() WHERE id = $1`
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) scanUser(row pgx.Row) (*domain.User, error) {
	var u domain.User
	err := row.Scan(
		&u.ID, &u.TenantID, &u.IsActive, &u.IsImported, &u.IsSalesMaster, &u.SSOOnly,
		&u.SAPUserID, &u.SAPCustomerNumber, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
		&u.Phone, &u.Mobile, &u.DefaultLanguage, &u.DefaultCompanyID, &u.InvitationToken,
		&u.InvitedAt, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) scanUserFromRows(rows pgx.Rows) (*domain.User, error) {
	var u domain.User
	err := rows.Scan(
		&u.ID, &u.TenantID, &u.IsActive, &u.IsImported, &u.IsSalesMaster, &u.SSOOnly,
		&u.SAPUserID, &u.SAPCustomerNumber, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
		&u.Phone, &u.Mobile, &u.DefaultLanguage, &u.DefaultCompanyID, &u.InvitationToken,
		&u.InvitedAt, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
