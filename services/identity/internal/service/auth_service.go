package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/identity/internal/auth"
	"github.com/gondolia/gondolia/services/identity/internal/domain"
	"github.com/gondolia/gondolia/services/identity/internal/repository"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo         repository.UserRepository
	companyRepo      repository.CompanyRepository
	roleRepo         repository.RoleRepository
	userCompanyRepo  repository.UserCompanyRepository
	refreshTokenRepo repository.RefreshTokenRepository
	passwordResetRepo repository.PasswordResetRepository
	authLogRepo      repository.AuthLogRepository
	jwtManager       *auth.JWTManager
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepository,
	companyRepo repository.CompanyRepository,
	roleRepo repository.RoleRepository,
	userCompanyRepo repository.UserCompanyRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	passwordResetRepo repository.PasswordResetRepository,
	authLogRepo repository.AuthLogRepository,
	jwtManager *auth.JWTManager,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		companyRepo:       companyRepo,
		roleRepo:          roleRepo,
		userCompanyRepo:   userCompanyRepo,
		refreshTokenRepo:  refreshTokenRepo,
		passwordResetRepo: passwordResetRepo,
		authLogRepo:       authLogRepo,
		jwtManager:        jwtManager,
	}
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, tenantID uuid.UUID, req domain.LoginRequest, ipAddress, userAgent string) (*domain.TokenPair, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, tenantID, req.Email)
	if err != nil {
		s.logAuthEvent(ctx, tenantID, nil, domain.AuthEventFailedLogin, ipAddress, userAgent, map[string]any{
			"email": req.Email,
			"error": "user not found",
		})
		return nil, domain.ErrInvalidCredentials
	}

	// Check if user can login
	if err := user.CanLogin(); err != nil {
		s.logAuthEvent(ctx, tenantID, &user.ID, domain.AuthEventFailedLogin, ipAddress, userAgent, map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}

	// Verify password
	if user.PasswordHash == nil || !auth.VerifyPassword(req.Password, *user.PasswordHash) {
		s.logAuthEvent(ctx, tenantID, &user.ID, domain.AuthEventFailedLogin, ipAddress, userAgent, map[string]any{
			"error": "invalid password",
		})
		return nil, domain.ErrInvalidCredentials
	}

	// Get user's companies and determine initial company
	company, role, err := s.determineInitialCompany(ctx, user)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	tokens, err := s.generateTokens(ctx, user, company, role, ipAddress, userAgent)
	if err != nil {
		return nil, err
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Log successful login
	s.logAuthEvent(ctx, tenantID, &user.ID, domain.AuthEventLogin, ipAddress, userAgent, map[string]any{
		"company_id": company.ID,
	})

	return tokens, nil
}

// Logout revokes the refresh token
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return domain.ErrTokenInvalid
	}

	// Revoke the token
	if err := s.refreshTokenRepo.Revoke(ctx, claims.TokenID); err != nil {
		return err
	}

	// Log logout
	s.logAuthEvent(ctx, claims.TenantID, &claims.UserID, domain.AuthEventLogout, "", "", nil)

	return nil
}

// RefreshToken generates new tokens from a valid refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string, ipAddress, userAgent string) (*domain.TokenPair, error) {
	// Validate refresh token JWT
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	// Check if token exists and is valid in database
	storedToken, err := s.refreshTokenRepo.GetByID(ctx, claims.TokenID)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	if !storedToken.IsValid() {
		if storedToken.IsRevoked() {
			return nil, domain.ErrTokenRevoked
		}
		return nil, domain.ErrTokenExpired
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Check if user can still login
	if err := user.CanLogin(); err != nil {
		return nil, err
	}

	// Get current company and role
	company, role, err := s.determineInitialCompany(ctx, user)
	if err != nil {
		return nil, err
	}

	// Revoke old token
	_ = s.refreshTokenRepo.Revoke(ctx, claims.TokenID)

	// Generate new tokens
	tokens, err := s.generateTokens(ctx, user, company, role, ipAddress, userAgent)
	if err != nil {
		return nil, err
	}

	// Log token refresh
	s.logAuthEvent(ctx, claims.TenantID, &claims.UserID, domain.AuthEventTokenRefresh, ipAddress, userAgent, nil)

	return tokens, nil
}

// SwitchCompany switches the user's current company context and returns new tokens
func (s *AuthService) SwitchCompany(ctx context.Context, userID, companyID uuid.UUID, ipAddress, userAgent string) (*domain.TokenPair, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check user is assigned to company
	uc, err := s.userCompanyRepo.GetByUserAndCompany(ctx, userID, companyID)
	if err != nil {
		return nil, domain.ErrUserNotInCompany
	}

	// Get company
	company, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	if !company.IsActive {
		return nil, domain.ErrCompanyNotActive
	}

	// Get role
	var role *domain.Role
	if uc.RoleID != nil {
		role, err = s.roleRepo.GetByID(ctx, *uc.RoleID)
		if err != nil {
			return nil, err
		}
	}

	// Update user's default company
	user.DefaultCompanyID = &companyID
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Generate new tokens
	return s.generateTokens(ctx, user, company, role, ipAddress, userAgent)
}

// ForgotPassword initiates password reset
func (s *AuthService) ForgotPassword(ctx context.Context, tenantID uuid.UUID, email string) (string, error) {
	// Get user
	user, err := s.userRepo.GetByEmail(ctx, tenantID, email)
	if err != nil {
		// Don't reveal if user exists
		return "", nil
	}

	// Check rate limiting (last reset < 60 seconds ago)
	lastReset, err := s.passwordResetRepo.GetLatestByUser(ctx, user.ID)
	if err == nil && time.Since(lastReset.CreatedAt) < 60*time.Second {
		// Silently ignore, don't reveal rate limiting
		return "", nil
	}

	// Generate token
	token, err := auth.GeneratePasswordResetToken()
	if err != nil {
		return "", err
	}

	// Store hashed token
	reset := &domain.PasswordReset{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: auth.HashToken(token),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(60 * time.Minute),
	}

	if err := s.passwordResetRepo.Create(ctx, reset); err != nil {
		return "", err
	}

	// Log event
	s.logAuthEvent(ctx, tenantID, &user.ID, domain.AuthEventPasswordReset, "", "", nil)

	return token, nil
}

// ResetPassword resets password with token
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate password strength
	if err := auth.ValidatePasswordStrength(newPassword); err != nil {
		return domain.ErrPasswordTooWeak
	}

	// Find reset by token hash
	tokenHash := auth.HashToken(token)
	reset, err := s.passwordResetRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return domain.ErrPasswordResetNotFound
	}

	// Check validity
	if reset.IsUsed() {
		return domain.ErrPasswordResetUsed
	}
	if reset.IsExpired() {
		return domain.ErrPasswordResetExpired
	}

	// Hash new password
	passwordHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	if err := s.userRepo.SetPassword(ctx, reset.UserID, passwordHash); err != nil {
		return err
	}

	// Mark reset as used
	if err := s.passwordResetRepo.MarkUsed(ctx, reset.ID); err != nil {
		return err
	}

	// Revoke all refresh tokens for user (force re-login)
	_ = s.refreshTokenRepo.RevokeAllForUser(ctx, reset.UserID)

	// Log event
	user, _ := s.userRepo.GetByID(ctx, reset.UserID)
	if user != nil {
		s.logAuthEvent(ctx, user.TenantID, &user.ID, domain.AuthEventPasswordChanged, "", "", map[string]any{
			"method": "forgot_password",
		})
	}

	return nil
}

// GetCurrentUser returns the current user with company context
func (s *AuthService) GetCurrentUser(ctx context.Context, userID, companyID uuid.UUID) (*domain.UserWithContext, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	company, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	uc, err := s.userCompanyRepo.GetByUserAndCompany(ctx, userID, companyID)
	if err != nil {
		return nil, err
	}

	var role *domain.Role
	var permissions []string

	if uc.RoleID != nil {
		role, err = s.roleRepo.GetByID(ctx, *uc.RoleID)
		if err != nil {
			return nil, err
		}
		permissions = role.PermissionStrings()
	}

	// Load all user's companies
	companies, err := s.userCompanyRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Companies = companies

	return &domain.UserWithContext{
		User:        user,
		Company:     company,
		Role:        role,
		Permissions: permissions,
	}, nil
}

// ValidateAccessToken validates an access token and returns claims
func (s *AuthService) ValidateAccessToken(token string) (*auth.AccessTokenClaims, error) {
	return s.jwtManager.ValidateAccessToken(token)
}

// determineInitialCompany determines the initial company for a user
func (s *AuthService) determineInitialCompany(ctx context.Context, user *domain.User) (*domain.Company, *domain.Role, error) {
	// Check if user has default company
	if user.DefaultCompanyID != nil {
		uc, err := s.userCompanyRepo.GetByUserAndCompany(ctx, user.ID, *user.DefaultCompanyID)
		if err == nil {
			company, err := s.companyRepo.GetByID(ctx, *user.DefaultCompanyID)
			if err == nil && company.IsActive {
				var role *domain.Role
				if uc.RoleID != nil {
					role, _ = s.roleRepo.GetByID(ctx, *uc.RoleID)
				}
				return company, role, nil
			}
		}
	}

	// Get first available company
	companies, err := s.userCompanyRepo.ListByUser(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}

	if len(companies) == 0 {
		return nil, nil, domain.ErrNoCompanyAssigned
	}

	// Return first active company
	for _, uc := range companies {
		if uc.Company != nil && uc.Company.IsActive {
			return uc.Company, uc.Role, nil
		}
	}

	return nil, nil, domain.ErrNoCompanyAssigned
}

// generateTokens generates access and refresh tokens
func (s *AuthService) generateTokens(ctx context.Context, user *domain.User, company *domain.Company, role *domain.Role, ipAddress, userAgent string) (*domain.TokenPair, error) {
	// Build permissions list
	var permissions []string
	var roleID uuid.UUID
	var roleName string

	if role != nil {
		roleID = role.ID
		roleName = role.Name
		permissions = role.PermissionStrings()
	}

	// Generate access token
	accessClaims := auth.AccessTokenClaims{
		UserID:        user.ID,
		TenantID:      user.TenantID,
		Email:         user.Email,
		Name:          user.Name(),
		CompanyID:     company.ID,
		CompanyName:   company.Name,
		RoleID:        roleID,
		RoleName:      roleName,
		Permissions:   permissions,
		IsSalesMaster: user.IsSalesMaster,
		SSOOnly:       user.SSOOnly,
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(accessClaims)
	if err != nil {
		return nil, err
	}

	// Generate refresh token ID
	refreshTokenID := uuid.New()

	refreshClaims := auth.RefreshTokenClaims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		TokenID:  refreshTokenID,
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(refreshClaims)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	storedToken := &domain.RefreshToken{
		ID:        refreshTokenID,
		UserID:    user.ID,
		TokenHash: auth.HashToken(refreshToken),
		UserAgent: nilIfEmpty(userAgent),
		IPAddress: nilIfEmpty(ipAddress),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.jwtManager.RefreshTokenExpiry()),
	}

	if err := s.refreshTokenRepo.Create(ctx, storedToken); err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtManager.AccessTokenExpiry().Seconds()),
	}, nil
}

// logAuthEvent logs an authentication event
func (s *AuthService) logAuthEvent(ctx context.Context, tenantID uuid.UUID, userID *uuid.UUID, eventType, ipAddress, userAgent string, metadata map[string]any) {
	log := &domain.AuthenticationLog{
		ID:        uuid.New(),
		TenantID:  tenantID,
		UserID:    userID,
		EventType: eventType,
		IPAddress: nilIfEmpty(ipAddress),
		UserAgent: nilIfEmpty(userAgent),
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}
	_ = s.authLogRepo.Create(ctx, log)
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
