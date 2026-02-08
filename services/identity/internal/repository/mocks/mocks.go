package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uuid.UUID]*domain.User),
	}
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, user := range m.users {
		if user.TenantID == tenantID && user.Email == email {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) GetByInvitationToken(ctx context.Context, token string) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, user := range m.users {
		if user.InvitationToken != nil && *user.InvitationToken == token {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) List(ctx context.Context, filter domain.UserFilter) ([]domain.User, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []domain.User
	for _, user := range m.users {
		if user.TenantID == filter.TenantID {
			result = append(result, *user)
		}
	}
	return result, len(result), nil
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if user, ok := m.users[id]; ok {
		now := time.Now()
		user.LastLoginAt = &now
	}
	return nil
}

func (m *MockUserRepository) SetPassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if user, ok := m.users[id]; ok {
		user.PasswordHash = &passwordHash
	}
	return nil
}

func (m *MockUserRepository) ClearInvitation(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if user, ok := m.users[id]; ok {
		user.InvitationToken = nil
		user.InvitedAt = nil
	}
	return nil
}

// AddUser adds a user to the mock repository
func (m *MockUserRepository) AddUser(user *domain.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.ID] = user
}

// MockCompanyRepository is a mock implementation of CompanyRepository
type MockCompanyRepository struct {
	mu        sync.RWMutex
	companies map[uuid.UUID]*domain.Company
}

func NewMockCompanyRepository() *MockCompanyRepository {
	return &MockCompanyRepository{
		companies: make(map[uuid.UUID]*domain.Company),
	}
}

func (m *MockCompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if company, ok := m.companies[id]; ok {
		return company, nil
	}
	return nil, domain.ErrCompanyNotFound
}

func (m *MockCompanyRepository) GetBySAPNumber(ctx context.Context, tenantID uuid.UUID, sapNumber string) (*domain.Company, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, company := range m.companies {
		if company.TenantID == tenantID && company.SAPCompanyNumber == sapNumber {
			return company, nil
		}
	}
	return nil, domain.ErrCompanyNotFound
}

func (m *MockCompanyRepository) List(ctx context.Context, filter domain.CompanyFilter) ([]domain.Company, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []domain.Company
	for _, company := range m.companies {
		if company.TenantID == filter.TenantID {
			result = append(result, *company)
		}
	}
	return result, len(result), nil
}

func (m *MockCompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.companies[company.ID] = company
	return nil
}

func (m *MockCompanyRepository) Update(ctx context.Context, company *domain.Company) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.companies[company.ID] = company
	return nil
}

func (m *MockCompanyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.companies, id)
	return nil
}

// AddCompany adds a company to the mock repository
func (m *MockCompanyRepository) AddCompany(company *domain.Company) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.companies[company.ID] = company
}

// MockRoleRepository is a mock implementation of RoleRepository
type MockRoleRepository struct {
	mu    sync.RWMutex
	roles map[uuid.UUID]*domain.Role
}

func NewMockRoleRepository() *MockRoleRepository {
	return &MockRoleRepository{
		roles: make(map[uuid.UUID]*domain.Role),
	}
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if role, ok := m.roles[id]; ok {
		return role, nil
	}
	return nil, domain.ErrRoleNotFound
}

func (m *MockRoleRepository) GetByName(ctx context.Context, tenantID uuid.UUID, companyID *uuid.UUID, name string) (*domain.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, role := range m.roles {
		if role.TenantID == tenantID && role.Name == name {
			return role, nil
		}
	}
	return nil, domain.ErrRoleNotFound
}

func (m *MockRoleRepository) List(ctx context.Context, filter domain.RoleFilter) ([]domain.Role, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []domain.Role
	for _, role := range m.roles {
		if role.TenantID == filter.TenantID {
			result = append(result, *role)
		}
	}
	return result, len(result), nil
}

func (m *MockRoleRepository) ListSystemRoles(ctx context.Context, tenantID uuid.UUID) ([]domain.Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []domain.Role
	for _, role := range m.roles {
		if role.TenantID == tenantID && role.IsSystem {
			result = append(result, *role)
		}
	}
	return result, nil
}

func (m *MockRoleRepository) Create(ctx context.Context, role *domain.Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.roles[role.ID] = role
	return nil
}

func (m *MockRoleRepository) Update(ctx context.Context, role *domain.Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.roles[role.ID] = role
	return nil
}

func (m *MockRoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.roles, id)
	return nil
}

// AddRole adds a role to the mock repository
func (m *MockRoleRepository) AddRole(role *domain.Role) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.roles[role.ID] = role
}

// MockUserCompanyRepository is a mock implementation of UserCompanyRepository
type MockUserCompanyRepository struct {
	mu            sync.RWMutex
	userCompanies []domain.UserCompany
}

func NewMockUserCompanyRepository() *MockUserCompanyRepository {
	return &MockUserCompanyRepository{
		userCompanies: []domain.UserCompany{},
	}
}

func (m *MockUserCompanyRepository) GetByUserAndCompany(ctx context.Context, userID, companyID uuid.UUID) (*domain.UserCompany, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, uc := range m.userCompanies {
		if uc.UserID == userID && uc.CompanyID == companyID {
			return &uc, nil
		}
	}
	return nil, domain.ErrUserNotInCompany
}

func (m *MockUserCompanyRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.UserCompany, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []domain.UserCompany
	for _, uc := range m.userCompanies {
		if uc.UserID == userID {
			result = append(result, uc)
		}
	}
	return result, nil
}

func (m *MockUserCompanyRepository) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.UserCompany, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []domain.UserCompany
	for _, uc := range m.userCompanies {
		if uc.CompanyID == companyID {
			result = append(result, uc)
		}
	}
	return result, nil
}

func (m *MockUserCompanyRepository) Create(ctx context.Context, uc *domain.UserCompany) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.userCompanies = append(m.userCompanies, *uc)
	return nil
}

func (m *MockUserCompanyRepository) Update(ctx context.Context, uc *domain.UserCompany) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, existing := range m.userCompanies {
		if existing.UserID == uc.UserID && existing.CompanyID == uc.CompanyID {
			m.userCompanies[i] = *uc
			return nil
		}
	}
	return nil
}

func (m *MockUserCompanyRepository) Delete(ctx context.Context, userID, companyID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, uc := range m.userCompanies {
		if uc.UserID == userID && uc.CompanyID == companyID {
			m.userCompanies = append(m.userCompanies[:i], m.userCompanies[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockUserCompanyRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	count := 0
	for _, uc := range m.userCompanies {
		if uc.UserID == userID {
			count++
		}
	}
	return count, nil
}

// AddUserCompany adds a user-company relationship to the mock repository
func (m *MockUserCompanyRepository) AddUserCompany(uc *domain.UserCompany) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.userCompanies = append(m.userCompanies, *uc)
}

// MockRefreshTokenRepository is a mock implementation of RefreshTokenRepository
type MockRefreshTokenRepository struct {
	mu     sync.RWMutex
	tokens map[uuid.UUID]*domain.RefreshToken
}

func NewMockRefreshTokenRepository() *MockRefreshTokenRepository {
	return &MockRefreshTokenRepository{
		tokens: make(map[uuid.UUID]*domain.RefreshToken),
	}
}

func (m *MockRefreshTokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if token, ok := m.tokens[id]; ok {
		return token, nil
	}
	return nil, domain.ErrRefreshTokenNotFound
}

func (m *MockRefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, token := range m.tokens {
		if token.TokenHash == tokenHash {
			return token, nil
		}
	}
	return nil, domain.ErrRefreshTokenNotFound
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokens[token.ID] = token
	return nil
}

func (m *MockRefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if token, ok := m.tokens[id]; ok {
		now := time.Now()
		token.RevokedAt = &now
	}
	return nil
}

func (m *MockRefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for _, token := range m.tokens {
		if token.UserID == userID {
			token.RevokedAt = &now
		}
	}
	return nil
}

func (m *MockRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return nil
}

// MockPasswordResetRepository is a mock implementation of PasswordResetRepository
type MockPasswordResetRepository struct {
	mu     sync.RWMutex
	resets map[uuid.UUID]*domain.PasswordReset
}

func NewMockPasswordResetRepository() *MockPasswordResetRepository {
	return &MockPasswordResetRepository{
		resets: make(map[uuid.UUID]*domain.PasswordReset),
	}
}

func (m *MockPasswordResetRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.PasswordReset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, reset := range m.resets {
		if reset.TokenHash == tokenHash {
			return reset, nil
		}
	}
	return nil, domain.ErrPasswordResetNotFound
}

func (m *MockPasswordResetRepository) GetLatestByUser(ctx context.Context, userID uuid.UUID) (*domain.PasswordReset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var latest *domain.PasswordReset
	for _, reset := range m.resets {
		if reset.UserID == userID {
			if latest == nil || reset.CreatedAt.After(latest.CreatedAt) {
				latest = reset
			}
		}
	}
	if latest == nil {
		return nil, domain.ErrPasswordResetNotFound
	}
	return latest, nil
}

func (m *MockPasswordResetRepository) Create(ctx context.Context, reset *domain.PasswordReset) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resets[reset.ID] = reset
	return nil
}

func (m *MockPasswordResetRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if reset, ok := m.resets[id]; ok {
		now := time.Now()
		reset.UsedAt = &now
	}
	return nil
}

func (m *MockPasswordResetRepository) DeleteExpired(ctx context.Context) error {
	return nil
}

// MockAuthLogRepository is a mock implementation of AuthLogRepository
type MockAuthLogRepository struct {
	mu   sync.RWMutex
	logs []domain.AuthenticationLog
}

func NewMockAuthLogRepository() *MockAuthLogRepository {
	return &MockAuthLogRepository{
		logs: []domain.AuthenticationLog{},
	}
}

func (m *MockAuthLogRepository) Create(ctx context.Context, log *domain.AuthenticationLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, *log)
	return nil
}

func (m *MockAuthLogRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.AuthenticationLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []domain.AuthenticationLog
	for _, log := range m.logs {
		if log.UserID != nil && *log.UserID == userID {
			result = append(result, log)
		}
	}
	return result, nil
}

func (m *MockAuthLogRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.AuthenticationLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []domain.AuthenticationLog
	for _, log := range m.logs {
		if log.TenantID == tenantID {
			result = append(result, log)
		}
	}
	return result, nil
}

// GetLogs returns all logged events (for testing)
func (m *MockAuthLogRepository) GetLogs() []domain.AuthenticationLog {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.logs
}

// MockTenantRepository is a mock implementation of TenantRepository
type MockTenantRepository struct {
	mu      sync.RWMutex
	tenants map[uuid.UUID]*domain.Tenant
}

func NewMockTenantRepository() *MockTenantRepository {
	return &MockTenantRepository{
		tenants: make(map[uuid.UUID]*domain.Tenant),
	}
}

func (m *MockTenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if tenant, ok := m.tenants[id]; ok {
		return tenant, nil
	}
	return nil, domain.ErrTenantNotFound
}

func (m *MockTenantRepository) GetByCode(ctx context.Context, code string) (*domain.Tenant, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, tenant := range m.tenants {
		if tenant.Code == code {
			return tenant, nil
		}
	}
	return nil, domain.ErrTenantNotFound
}

func (m *MockTenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tenants[tenant.ID] = tenant
	return nil
}

func (m *MockTenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tenants[tenant.ID] = tenant
	return nil
}

// AddTenant adds a tenant to the mock repository
func (m *MockTenantRepository) AddTenant(tenant *domain.Tenant) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tenants[tenant.ID] = tenant
}
