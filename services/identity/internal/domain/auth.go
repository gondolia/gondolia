package domain

import (
	"time"

	"github.com/google/uuid"
)

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // seconds until access token expires
}

// RefreshToken represents a stored refresh token
type RefreshToken struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	TokenHash string     `json:"-"`
	UserAgent *string    `json:"user_agent,omitempty"`
	IPAddress *string    `json:"ip_address,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

// IsExpired checks if token is expired
func (t *RefreshToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsRevoked checks if token is revoked
func (t *RefreshToken) IsRevoked() bool {
	return t.RevokedAt != nil
}

// IsValid checks if token is valid (not expired, not revoked)
func (t *RefreshToken) IsValid() bool {
	return !t.IsExpired() && !t.IsRevoked()
}

// PasswordReset represents a password reset request
type PasswordReset struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	TokenHash string     `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
}

// IsExpired checks if reset token is expired
func (r *PasswordReset) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// IsUsed checks if reset token has been used
func (r *PasswordReset) IsUsed() bool {
	return r.UsedAt != nil
}

// IsValid checks if reset token is valid
func (r *PasswordReset) IsValid() bool {
	return !r.IsExpired() && !r.IsUsed()
}

// AuthenticationLog represents an authentication event
type AuthenticationLog struct {
	ID        uuid.UUID      `json:"id"`
	TenantID  uuid.UUID      `json:"tenant_id"`
	UserID    *uuid.UUID     `json:"user_id,omitempty"`
	EventType string         `json:"event_type"`
	IPAddress *string        `json:"ip_address,omitempty"`
	UserAgent *string        `json:"user_agent,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// Auth event types
const (
	AuthEventLogin              = "login"
	AuthEventLogout             = "logout"
	AuthEventFailedLogin        = "failed_login"
	AuthEventPasswordReset      = "password_reset"
	AuthEventPasswordChanged    = "password_changed"
	AuthEventInvitationSent     = "invitation_sent"
	AuthEventInvitationAccepted = "invitation_accepted"
	AuthEventTokenRefresh       = "token_refresh"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents a password reset request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// SwitchCompanyRequest represents a request to switch company context
type SwitchCompanyRequest struct {
	CompanyID uuid.UUID `json:"company_id" binding:"required"`
}

// UserWithContext represents user data with current company context
type UserWithContext struct {
	User        *User    `json:"user"`
	Company     *Company `json:"company"`
	Role        *Role    `json:"role"`
	Permissions []string `json:"permissions"`
}

// AcceptInvitationRequest represents a request to accept an invitation
type AcceptInvitationRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}
