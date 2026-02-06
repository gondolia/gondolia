// Package auth defines the interface for SSO/Identity Provider integrations.
package auth

import "context"

// AuthProvider abstracts SSO/Identity Providers (Azure AD, Keycloak, Auth0, etc.).
type AuthProvider interface {
	// GetAuthURL returns the URL for SSO login.
	GetAuthURL(ctx context.Context, state string, redirectURL string) (string, error)

	// HandleCallback processes the SSO callback and returns user information.
	HandleCallback(ctx context.Context, code string, state string) (*SSOUser, error)

	// ValidateToken validates an SSO token (for API-based flows).
	ValidateToken(ctx context.Context, token string) (*SSOUser, error)

	// GetUserInfo retrieves user information from the provider.
	GetUserInfo(ctx context.Context, accessToken string) (*SSOUser, error)

	// Metadata returns provider information.
	Metadata() Metadata
}

// SSOUser represents a user from an SSO provider.
type SSOUser struct {
	ExternalID string
	Email      string
	FirstName  string
	LastName   string
	Groups     []string
	Attributes map[string]string
	RawClaims  map[string]any
}

// Metadata provides information about the auth provider.
type Metadata struct {
	Name     string
	Protocol string // "saml", "oidc", "oauth2"
	Issuer   string
}
