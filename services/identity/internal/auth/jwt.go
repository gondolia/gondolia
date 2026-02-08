package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenConfig holds JWT configuration
type TokenConfig struct {
	AccessSecret        string
	RefreshSecret       string
	AccessTokenExpiry   time.Duration
	RefreshTokenExpiry  time.Duration
	Issuer              string
}

// DefaultTokenConfig returns default token configuration
func DefaultTokenConfig() TokenConfig {
	return TokenConfig{
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour, // 7 days
		Issuer:             "identity-service",
	}
}

// AccessTokenClaims represents the claims in an access token
type AccessTokenClaims struct {
	jwt.RegisteredClaims

	// User Info
	UserID    uuid.UUID `json:"user_id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`

	// Current Context
	CompanyID   uuid.UUID `json:"company_id"`
	CompanyName string    `json:"company_name"`
	RoleID      uuid.UUID `json:"role_id"`
	RoleName    string    `json:"role_name"`

	// Permissions (flattened for fast checks)
	Permissions []string `json:"permissions"`

	// Flags
	IsSalesMaster bool `json:"is_salesmaster"`
	SSOOnly       bool `json:"sso_only"`
}

// RefreshTokenClaims represents the claims in a refresh token
type RefreshTokenClaims struct {
	jwt.RegisteredClaims

	UserID   uuid.UUID `json:"user_id"`
	TenantID uuid.UUID `json:"tenant_id"`
	TokenID  uuid.UUID `json:"token_id"` // ID in database
}

// JWTManager handles JWT operations
type JWTManager struct {
	config TokenConfig
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(config TokenConfig) *JWTManager {
	return &JWTManager{config: config}
}

// GenerateAccessToken generates a new access token
func (m *JWTManager) GenerateAccessToken(claims AccessTokenClaims) (string, error) {
	now := time.Now()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    m.config.Issuer,
		Subject:   claims.UserID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.config.AccessTokenExpiry)),
		NotBefore: jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.AccessSecret))
}

// GenerateRefreshToken generates a new refresh token
func (m *JWTManager) GenerateRefreshToken(claims RefreshTokenClaims) (string, error) {
	now := time.Now()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    m.config.Issuer,
		Subject:   claims.UserID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.config.RefreshTokenExpiry)),
		NotBefore: jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.RefreshSecret))
}

// ValidateAccessToken validates and parses an access token
func (m *JWTManager) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.config.AccessSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// ValidateRefreshToken validates and parses a refresh token
func (m *JWTManager) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.config.RefreshSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// AccessTokenExpiry returns the access token expiry duration
func (m *JWTManager) AccessTokenExpiry() time.Duration {
	return m.config.AccessTokenExpiry
}

// RefreshTokenExpiry returns the refresh token expiry duration
func (m *JWTManager) RefreshTokenExpiry() time.Duration {
	return m.config.RefreshTokenExpiry
}
