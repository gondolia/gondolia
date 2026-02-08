package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTManager_GenerateAndValidateAccessToken(t *testing.T) {
	config := TokenConfig{
		AccessSecret:       "test-access-secret-min-32-characters",
		RefreshSecret:      "test-refresh-secret-min-32-characters",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		Issuer:             "identity-service",
	}

	manager := NewJWTManager(config)

	userID := uuid.MustParse("00000000-0000-0000-0000-000000000005")
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	companyID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	roleID := uuid.MustParse("00000000-0000-0000-0000-000000000003")

	claims := AccessTokenClaims{
		UserID:        userID,
		TenantID:      tenantID,
		Email:         "admin@demo.local",
		Name:          "Admin User",
		CompanyID:     companyID,
		CompanyName:   "Demo Company GmbH",
		RoleID:        roleID,
		RoleName:      "Administrator",
		Permissions:   []string{"company.manage", "sales.create-order"},
		IsSalesMaster: true,
		SSOOnly:       false,
	}

	// Generate token
	token, err := manager.GenerateAccessToken(claims)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	if token == "" {
		t.Fatal("GenerateAccessToken() returned empty token")
	}

	// Validate token
	parsed, err := manager.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	// Verify claims
	if parsed.UserID != userID {
		t.Errorf("UserID = %v, want %v", parsed.UserID, userID)
	}
	if parsed.TenantID != tenantID {
		t.Errorf("TenantID = %v, want %v", parsed.TenantID, tenantID)
	}
	if parsed.Email != "admin@demo.local" {
		t.Errorf("Email = %v, want admin@demo.local", parsed.Email)
	}
	if parsed.CompanyID != companyID {
		t.Errorf("CompanyID = %v, want %v", parsed.CompanyID, companyID)
	}
	if parsed.RoleID != roleID {
		t.Errorf("RoleID = %v, want %v", parsed.RoleID, roleID)
	}
	if !parsed.IsSalesMaster {
		t.Error("IsSalesMaster should be true")
	}
	if len(parsed.Permissions) != 2 {
		t.Errorf("Permissions length = %d, want 2", len(parsed.Permissions))
	}
}

func TestJWTManager_GenerateAndValidateRefreshToken(t *testing.T) {
	config := TokenConfig{
		AccessSecret:       "test-access-secret-min-32-characters",
		RefreshSecret:      "test-refresh-secret-min-32-characters",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		Issuer:             "identity-service",
	}

	manager := NewJWTManager(config)

	userID := uuid.MustParse("00000000-0000-0000-0000-000000000005")
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	tokenID := uuid.New()

	claims := RefreshTokenClaims{
		UserID:   userID,
		TenantID: tenantID,
		TokenID:  tokenID,
	}

	// Generate token
	token, err := manager.GenerateRefreshToken(claims)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	if token == "" {
		t.Fatal("GenerateRefreshToken() returned empty token")
	}

	// Validate token
	parsed, err := manager.ValidateRefreshToken(token)
	if err != nil {
		t.Fatalf("ValidateRefreshToken() error = %v", err)
	}

	// Verify claims
	if parsed.UserID != userID {
		t.Errorf("UserID = %v, want %v", parsed.UserID, userID)
	}
	if parsed.TenantID != tenantID {
		t.Errorf("TenantID = %v, want %v", parsed.TenantID, tenantID)
	}
	if parsed.TokenID != tokenID {
		t.Errorf("TokenID = %v, want %v", parsed.TokenID, tokenID)
	}
}

func TestJWTManager_InvalidToken(t *testing.T) {
	config := TokenConfig{
		AccessSecret:       "test-access-secret-min-32-characters",
		RefreshSecret:      "test-refresh-secret-min-32-characters",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		Issuer:             "identity-service",
	}

	manager := NewJWTManager(config)

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "invalid format",
			token: "not-a-valid-jwt",
		},
		{
			name:  "tampered token",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.tampered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+" access", func(t *testing.T) {
			_, err := manager.ValidateAccessToken(tt.token)
			if err == nil {
				t.Error("ValidateAccessToken() should return error for invalid token")
			}
		})

		t.Run(tt.name+" refresh", func(t *testing.T) {
			_, err := manager.ValidateRefreshToken(tt.token)
			if err == nil {
				t.Error("ValidateRefreshToken() should return error for invalid token")
			}
		})
	}
}

func TestJWTManager_WrongSecret(t *testing.T) {
	config1 := TokenConfig{
		AccessSecret:       "secret-one-min-32-characters-here",
		RefreshSecret:      "refresh-one-min-32-characters-here",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		Issuer:             "identity-service",
	}

	config2 := TokenConfig{
		AccessSecret:       "secret-two-min-32-characters-here",
		RefreshSecret:      "refresh-two-min-32-characters-here",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		Issuer:             "identity-service",
	}

	manager1 := NewJWTManager(config1)
	manager2 := NewJWTManager(config2)

	userID := uuid.New()
	tenantID := uuid.New()

	// Generate with manager1
	accessToken, _ := manager1.GenerateAccessToken(AccessTokenClaims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    "test@test.com",
	})

	refreshToken, _ := manager1.GenerateRefreshToken(RefreshTokenClaims{
		UserID:   userID,
		TenantID: tenantID,
		TokenID:  uuid.New(),
	})

	// Try to validate with manager2 (different secret)
	_, err := manager2.ValidateAccessToken(accessToken)
	if err == nil {
		t.Error("ValidateAccessToken() should fail with wrong secret")
	}

	_, err = manager2.ValidateRefreshToken(refreshToken)
	if err == nil {
		t.Error("ValidateRefreshToken() should fail with wrong secret")
	}
}

func TestJWTManager_TokenExpiry(t *testing.T) {
	config := TokenConfig{
		AccessSecret:       "test-access-secret-min-32-characters",
		RefreshSecret:      "test-refresh-secret-min-32-characters",
		AccessTokenExpiry:  -1 * time.Hour, // Already expired
		RefreshTokenExpiry: -1 * time.Hour, // Already expired
		Issuer:             "identity-service",
	}

	manager := NewJWTManager(config)

	userID := uuid.New()
	tenantID := uuid.New()

	// Generate expired access token
	accessToken, err := manager.GenerateAccessToken(AccessTokenClaims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    "test@test.com",
	})
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	// Validation should fail due to expiry
	_, err = manager.ValidateAccessToken(accessToken)
	if err == nil {
		t.Error("ValidateAccessToken() should fail for expired token")
	}

	// Generate expired refresh token
	refreshToken, err := manager.GenerateRefreshToken(RefreshTokenClaims{
		UserID:   userID,
		TenantID: tenantID,
		TokenID:  uuid.New(),
	})
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	// Validation should fail due to expiry
	_, err = manager.ValidateRefreshToken(refreshToken)
	if err == nil {
		t.Error("ValidateRefreshToken() should fail for expired token")
	}
}

func TestDefaultTokenConfig(t *testing.T) {
	config := DefaultTokenConfig()

	if config.AccessTokenExpiry != 15*time.Minute {
		t.Errorf("AccessTokenExpiry = %v, want 15m", config.AccessTokenExpiry)
	}
	if config.RefreshTokenExpiry != 7*24*time.Hour {
		t.Errorf("RefreshTokenExpiry = %v, want 7d", config.RefreshTokenExpiry)
	}
	if config.Issuer != "identity-service" {
		t.Errorf("Issuer = %v, want identity-service", config.Issuer)
	}
}
