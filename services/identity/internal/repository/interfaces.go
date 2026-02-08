package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
)

// TenantRepository defines the interface for tenant data access
type TenantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	GetByCode(ctx context.Context, code string) (*domain.Tenant, error)
	Create(ctx context.Context, tenant *domain.Tenant) error
	Update(ctx context.Context, tenant *domain.Tenant) error
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*domain.User, error)
	GetByInvitationToken(ctx context.Context, token string) (*domain.User, error)
	List(ctx context.Context, filter domain.UserFilter) ([]domain.User, int, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error // Soft delete
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	SetPassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	ClearInvitation(ctx context.Context, id uuid.UUID) error
}

// CompanyRepository defines the interface for company data access
type CompanyRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error)
	GetBySAPNumber(ctx context.Context, tenantID uuid.UUID, sapNumber string) (*domain.Company, error)
	List(ctx context.Context, filter domain.CompanyFilter) ([]domain.Company, int, error)
	Create(ctx context.Context, company *domain.Company) error
	Update(ctx context.Context, company *domain.Company) error
	Delete(ctx context.Context, id uuid.UUID) error // Soft delete
}

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error)
	GetByName(ctx context.Context, tenantID uuid.UUID, companyID *uuid.UUID, name string) (*domain.Role, error)
	List(ctx context.Context, filter domain.RoleFilter) ([]domain.Role, int, error)
	ListSystemRoles(ctx context.Context, tenantID uuid.UUID) ([]domain.Role, error)
	Create(ctx context.Context, role *domain.Role) error
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// UserCompanyRepository defines the interface for user-company relationship data access
type UserCompanyRepository interface {
	GetByUserAndCompany(ctx context.Context, userID, companyID uuid.UUID) (*domain.UserCompany, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.UserCompany, error)
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.UserCompany, error)
	Create(ctx context.Context, uc *domain.UserCompany) error
	Update(ctx context.Context, uc *domain.UserCompany) error
	Delete(ctx context.Context, userID, companyID uuid.UUID) error
	CountByUser(ctx context.Context, userID uuid.UUID) (int, error)
}

// RefreshTokenRepository defines the interface for refresh token data access
type RefreshTokenRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	Create(ctx context.Context, token *domain.RefreshToken) error
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

// PasswordResetRepository defines the interface for password reset data access
type PasswordResetRepository interface {
	GetByTokenHash(ctx context.Context, tokenHash string) (*domain.PasswordReset, error)
	GetLatestByUser(ctx context.Context, userID uuid.UUID) (*domain.PasswordReset, error)
	Create(ctx context.Context, reset *domain.PasswordReset) error
	MarkUsed(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

// AuthLogRepository defines the interface for authentication log data access
type AuthLogRepository interface {
	Create(ctx context.Context, log *domain.AuthenticationLog) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.AuthenticationLog, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.AuthenticationLog, error)
}
