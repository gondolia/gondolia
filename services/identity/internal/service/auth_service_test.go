package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/identity/internal/auth"
	"github.com/gondolia/gondolia/services/identity/internal/domain"
	"github.com/gondolia/gondolia/services/identity/internal/repository/mocks"
)

// Test fixtures
var (
	testTenantID  = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	testCompanyID = uuid.MustParse("00000000-0000-0000-0000-000000000002")
	testRoleID    = uuid.MustParse("00000000-0000-0000-0000-000000000003")
	testUserID    = uuid.MustParse("00000000-0000-0000-0000-000000000005")
)

func setupAuthServiceTest() (*AuthService, *mocks.MockUserRepository, *mocks.MockCompanyRepository, *mocks.MockRoleRepository, *mocks.MockUserCompanyRepository) {
	userRepo := mocks.NewMockUserRepository()
	companyRepo := mocks.NewMockCompanyRepository()
	roleRepo := mocks.NewMockRoleRepository()
	userCompanyRepo := mocks.NewMockUserCompanyRepository()
	refreshTokenRepo := mocks.NewMockRefreshTokenRepository()
	passwordResetRepo := mocks.NewMockPasswordResetRepository()
	authLogRepo := mocks.NewMockAuthLogRepository()

	jwtConfig := auth.TokenConfig{
		AccessSecret:       "test-access-secret-min-32-characters",
		RefreshSecret:      "test-refresh-secret-min-32-characters",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		Issuer:             "identity-service",
	}
	jwtManager := auth.NewJWTManager(jwtConfig)

	service := NewAuthService(
		userRepo,
		companyRepo,
		roleRepo,
		userCompanyRepo,
		refreshTokenRepo,
		passwordResetRepo,
		authLogRepo,
		jwtManager,
	)

	return service, userRepo, companyRepo, roleRepo, userCompanyRepo
}

func createTestUser(passwordHash string) *domain.User {
	return &domain.User{
		ID:               testUserID,
		TenantID:         testTenantID,
		IsActive:         true,
		IsImported:       false,
		IsSalesMaster:    true,
		SSOOnly:          false,
		Email:            "admin@demo.local",
		PasswordHash:     &passwordHash,
		FirstName:        "Admin",
		LastName:         "User",
		DefaultLanguage:  "de",
		DefaultCompanyID: &testCompanyID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func createTestCompany() *domain.Company {
	return &domain.Company{
		ID:               testCompanyID,
		TenantID:         testTenantID,
		SAPCompanyNumber: "1000",
		Name:             "Demo Company GmbH",
		Currency:         "EUR",
		Country:          "DE",
		IsActive:         true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func createTestRole() *domain.Role {
	return &domain.Role{
		ID:       testRoleID,
		TenantID: testTenantID,
		Name:     "Administrator",
		Permissions: map[domain.Permission]bool{
			domain.PermManageCompany: true,
			domain.PermCreateOrder:   true,
		},
		IsSystem:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createTestUserCompany(company *domain.Company, role *domain.Role) *domain.UserCompany {
	return &domain.UserCompany{
		ID:        uuid.New(),
		UserID:    testUserID,
		CompanyID: testCompanyID,
		RoleID:    &testRoleID,
		UserType:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Company:   company,
		Role:      role,
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	service, userRepo, companyRepo, roleRepo, userCompanyRepo := setupAuthServiceTest()
	ctx := context.Background()

	// Setup test data
	passwordHash, _ := auth.HashPassword("admin123")
	user := createTestUser(passwordHash)
	company := createTestCompany()
	role := createTestRole()
	userCompany := createTestUserCompany(company, role)

	userRepo.AddUser(user)
	companyRepo.AddCompany(company)
	roleRepo.AddRole(role)
	userCompanyRepo.AddUserCompany(userCompany)

	// Test login
	req := domain.LoginRequest{
		Email:    "admin@demo.local",
		Password: "admin123",
	}

	tokens, err := service.Login(ctx, testTenantID, req, "127.0.0.1", "TestAgent/1.0")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if tokens.AccessToken == "" {
		t.Error("Login() returned empty access token")
	}
	if tokens.RefreshToken == "" {
		t.Error("Login() returned empty refresh token")
	}
	if tokens.TokenType != "Bearer" {
		t.Errorf("Login() token type = %s, want Bearer", tokens.TokenType)
	}
	if tokens.ExpiresIn <= 0 {
		t.Error("Login() returned non-positive expiry")
	}

	// Validate access token contains correct claims
	claims, err := service.ValidateAccessToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	if claims.UserID != testUserID {
		t.Errorf("Claims UserID = %v, want %v", claims.UserID, testUserID)
	}
	if claims.Email != "admin@demo.local" {
		t.Errorf("Claims Email = %v, want admin@demo.local", claims.Email)
	}
	if claims.CompanyID != testCompanyID {
		t.Errorf("Claims CompanyID = %v, want %v", claims.CompanyID, testCompanyID)
	}
	if claims.RoleID != testRoleID {
		t.Errorf("Claims RoleID = %v, want %v", claims.RoleID, testRoleID)
	}
	if !claims.IsSalesMaster {
		t.Error("Claims IsSalesMaster should be true")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	service, userRepo, companyRepo, roleRepo, userCompanyRepo := setupAuthServiceTest()
	ctx := context.Background()

	// Setup test data
	passwordHash, _ := auth.HashPassword("admin123")
	user := createTestUser(passwordHash)
	company := createTestCompany()
	role := createTestRole()
	userCompany := createTestUserCompany(company, role)

	userRepo.AddUser(user)
	companyRepo.AddCompany(company)
	roleRepo.AddRole(role)
	userCompanyRepo.AddUserCompany(userCompany)

	// Test login with wrong password
	req := domain.LoginRequest{
		Email:    "admin@demo.local",
		Password: "wrongpassword",
	}

	_, err := service.Login(ctx, testTenantID, req, "127.0.0.1", "TestAgent/1.0")
	if err == nil {
		t.Fatal("Login() should return error for wrong password")
	}
	if err != domain.ErrInvalidCredentials {
		t.Errorf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	service, _, _, _, _ := setupAuthServiceTest()
	ctx := context.Background()

	req := domain.LoginRequest{
		Email:    "nonexistent@demo.local",
		Password: "anypassword",
	}

	_, err := service.Login(ctx, testTenantID, req, "127.0.0.1", "TestAgent/1.0")
	if err == nil {
		t.Fatal("Login() should return error for non-existent user")
	}
	if err != domain.ErrInvalidCredentials {
		t.Errorf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	service, userRepo, companyRepo, roleRepo, userCompanyRepo := setupAuthServiceTest()
	ctx := context.Background()

	// Setup inactive user
	passwordHash, _ := auth.HashPassword("admin123")
	user := createTestUser(passwordHash)
	user.IsActive = false

	company := createTestCompany()
	role := createTestRole()
	userCompany := createTestUserCompany(company, role)

	userRepo.AddUser(user)
	companyRepo.AddCompany(company)
	roleRepo.AddRole(role)
	userCompanyRepo.AddUserCompany(userCompany)

	req := domain.LoginRequest{
		Email:    "admin@demo.local",
		Password: "admin123",
	}

	_, err := service.Login(ctx, testTenantID, req, "127.0.0.1", "TestAgent/1.0")
	if err == nil {
		t.Fatal("Login() should return error for inactive user")
	}
}

func TestAuthService_Login_SSOOnlyUser(t *testing.T) {
	service, userRepo, companyRepo, roleRepo, userCompanyRepo := setupAuthServiceTest()
	ctx := context.Background()

	// Setup SSO-only user
	passwordHash, _ := auth.HashPassword("admin123")
	user := createTestUser(passwordHash)
	user.SSOOnly = true

	company := createTestCompany()
	role := createTestRole()
	userCompany := createTestUserCompany(company, role)

	userRepo.AddUser(user)
	companyRepo.AddCompany(company)
	roleRepo.AddRole(role)
	userCompanyRepo.AddUserCompany(userCompany)

	req := domain.LoginRequest{
		Email:    "admin@demo.local",
		Password: "admin123",
	}

	_, err := service.Login(ctx, testTenantID, req, "127.0.0.1", "TestAgent/1.0")
	if err == nil {
		t.Fatal("Login() should return error for SSO-only user")
	}
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	service, userRepo, companyRepo, roleRepo, userCompanyRepo := setupAuthServiceTest()
	ctx := context.Background()

	// Setup test data
	passwordHash, _ := auth.HashPassword("admin123")
	user := createTestUser(passwordHash)
	company := createTestCompany()
	role := createTestRole()
	userCompany := createTestUserCompany(company, role)

	userRepo.AddUser(user)
	companyRepo.AddCompany(company)
	roleRepo.AddRole(role)
	userCompanyRepo.AddUserCompany(userCompany)

	// First, login to get tokens
	req := domain.LoginRequest{
		Email:    "admin@demo.local",
		Password: "admin123",
	}

	tokens, err := service.Login(ctx, testTenantID, req, "127.0.0.1", "TestAgent/1.0")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	// Now refresh the token
	newTokens, err := service.RefreshToken(ctx, tokens.RefreshToken, "127.0.0.1", "TestAgent/1.0")
	if err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}

	if newTokens.AccessToken == "" {
		t.Error("RefreshToken() returned empty access token")
	}
	if newTokens.RefreshToken == "" {
		t.Error("RefreshToken() returned empty refresh token")
	}
	if newTokens.TokenType != "Bearer" {
		t.Errorf("RefreshToken() token type = %s, want Bearer", newTokens.TokenType)
	}
	// Refresh tokens should be different (new token ID)
	if newTokens.RefreshToken == tokens.RefreshToken {
		t.Error("RefreshToken() should return new refresh token")
	}
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	service, _, _, _, _ := setupAuthServiceTest()
	ctx := context.Background()

	_, err := service.RefreshToken(ctx, "invalid-token", "127.0.0.1", "TestAgent/1.0")
	if err == nil {
		t.Fatal("RefreshToken() should return error for invalid token")
	}
}

func TestAuthService_Logout_Success(t *testing.T) {
	service, userRepo, companyRepo, roleRepo, userCompanyRepo := setupAuthServiceTest()
	ctx := context.Background()

	// Setup test data
	passwordHash, _ := auth.HashPassword("admin123")
	user := createTestUser(passwordHash)
	company := createTestCompany()
	role := createTestRole()
	userCompany := createTestUserCompany(company, role)

	userRepo.AddUser(user)
	companyRepo.AddCompany(company)
	roleRepo.AddRole(role)
	userCompanyRepo.AddUserCompany(userCompany)

	// Login first
	req := domain.LoginRequest{
		Email:    "admin@demo.local",
		Password: "admin123",
	}

	tokens, err := service.Login(ctx, testTenantID, req, "127.0.0.1", "TestAgent/1.0")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	// Logout
	err = service.Logout(ctx, tokens.RefreshToken)
	if err != nil {
		t.Fatalf("Logout() error = %v", err)
	}

	// Refresh should fail after logout
	_, err = service.RefreshToken(ctx, tokens.RefreshToken, "127.0.0.1", "TestAgent/1.0")
	if err == nil {
		t.Error("RefreshToken() should fail after logout")
	}
}

func TestAuthService_GetCurrentUser_Success(t *testing.T) {
	service, userRepo, companyRepo, roleRepo, userCompanyRepo := setupAuthServiceTest()
	ctx := context.Background()

	// Setup test data
	passwordHash, _ := auth.HashPassword("admin123")
	user := createTestUser(passwordHash)
	company := createTestCompany()
	role := createTestRole()
	userCompany := createTestUserCompany(company, role)

	userRepo.AddUser(user)
	companyRepo.AddCompany(company)
	roleRepo.AddRole(role)
	userCompanyRepo.AddUserCompany(userCompany)

	// Get current user
	userWithContext, err := service.GetCurrentUser(ctx, testUserID, testCompanyID)
	if err != nil {
		t.Fatalf("GetCurrentUser() error = %v", err)
	}

	if userWithContext.User.Email != "admin@demo.local" {
		t.Errorf("User.Email = %v, want admin@demo.local", userWithContext.User.Email)
	}
	if userWithContext.Company.Name != "Demo Company GmbH" {
		t.Errorf("Company.Name = %v, want Demo Company GmbH", userWithContext.Company.Name)
	}
	if userWithContext.Role.Name != "Administrator" {
		t.Errorf("Role.Name = %v, want Administrator", userWithContext.Role.Name)
	}
	if len(userWithContext.Permissions) == 0 {
		t.Error("Permissions should not be empty")
	}
}

func TestAuthService_ForgotPassword_Success(t *testing.T) {
	service, userRepo, _, _, _ := setupAuthServiceTest()
	ctx := context.Background()

	// Setup test data
	passwordHash, _ := auth.HashPassword("admin123")
	user := createTestUser(passwordHash)
	userRepo.AddUser(user)

	// Request password reset
	token, err := service.ForgotPassword(ctx, testTenantID, "admin@demo.local")
	if err != nil {
		t.Fatalf("ForgotPassword() error = %v", err)
	}

	if token == "" {
		t.Error("ForgotPassword() returned empty token")
	}
}

func TestAuthService_ForgotPassword_NonexistentUser(t *testing.T) {
	service, _, _, _, _ := setupAuthServiceTest()
	ctx := context.Background()

	// Request password reset for non-existent user
	// Should not return error (security - don't reveal if user exists)
	token, err := service.ForgotPassword(ctx, testTenantID, "nonexistent@demo.local")
	if err != nil {
		t.Fatalf("ForgotPassword() error = %v", err)
	}

	// Token should be empty for non-existent user
	if token != "" {
		t.Error("ForgotPassword() should return empty token for non-existent user")
	}
}

func TestAuthService_ResetPassword_Success(t *testing.T) {
	service, userRepo, _, _, _ := setupAuthServiceTest()
	ctx := context.Background()

	// Setup test data
	passwordHash, _ := auth.HashPassword("admin123")
	user := createTestUser(passwordHash)
	userRepo.AddUser(user)

	// First, request password reset
	token, err := service.ForgotPassword(ctx, testTenantID, "admin@demo.local")
	if err != nil {
		t.Fatalf("ForgotPassword() error = %v", err)
	}

	// Reset password
	err = service.ResetPassword(ctx, token, "newpassword123")
	if err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}

	// Verify new password works (by checking the hash was updated)
	updatedUser, _ := userRepo.GetByID(ctx, testUserID)
	if updatedUser.PasswordHash == nil {
		t.Fatal("Password hash should be set")
	}
	if !auth.VerifyPassword("newpassword123", *updatedUser.PasswordHash) {
		t.Error("New password should be valid")
	}
}

func TestAuthService_ResetPassword_WeakPassword(t *testing.T) {
	service, userRepo, _, _, _ := setupAuthServiceTest()
	ctx := context.Background()

	// Setup test data
	passwordHash, _ := auth.HashPassword("admin123")
	user := createTestUser(passwordHash)
	userRepo.AddUser(user)

	// Request password reset
	token, _ := service.ForgotPassword(ctx, testTenantID, "admin@demo.local")

	// Try to reset with weak password
	err := service.ResetPassword(ctx, token, "weak")
	if err == nil {
		t.Fatal("ResetPassword() should fail for weak password")
	}
}

func TestAuthService_ResetPassword_InvalidToken(t *testing.T) {
	service, _, _, _, _ := setupAuthServiceTest()
	ctx := context.Background()

	err := service.ResetPassword(ctx, "invalid-token", "newpassword123")
	if err == nil {
		t.Fatal("ResetPassword() should fail for invalid token")
	}
}
