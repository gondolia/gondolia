// Package noop provides a no-op implementation of the Auth provider.
package noop

import (
	"context"

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/auth"
)

func init() {
	provider.Register[auth.AuthProvider]("auth", "noop",
		provider.Metadata{
			Name:        "noop",
			DisplayName: "No-Op Auth Provider",
			Category:    "auth",
			Version:     "1.0.0",
			Description: "A no-operation auth provider for development and testing",
			ConfigSpec:  []provider.ConfigField{},
		},
		NewProvider,
	)
}

// Provider is a no-op Auth provider.
type Provider struct{}

// NewProvider creates a new no-op Auth provider.
func NewProvider(config map[string]any) (auth.AuthProvider, error) {
	return &Provider{}, nil
}

func (p *Provider) GetAuthURL(ctx context.Context, state string, redirectURL string) (string, error) {
	return redirectURL + "?code=noop-code&state=" + state, nil
}

func (p *Provider) HandleCallback(ctx context.Context, code string, state string) (*auth.SSOUser, error) {
	return &auth.SSOUser{
		ExternalID: "noop-user-001",
		Email:      "user@example.com",
		FirstName:  "Demo",
		LastName:   "User",
		Groups:     []string{},
		Attributes: make(map[string]string),
		RawClaims:  make(map[string]any),
	}, nil
}

func (p *Provider) ValidateToken(ctx context.Context, token string) (*auth.SSOUser, error) {
	return &auth.SSOUser{
		ExternalID: "noop-user-001",
		Email:      "user@example.com",
		FirstName:  "Demo",
		LastName:   "User",
		Groups:     []string{},
		Attributes: make(map[string]string),
		RawClaims:  make(map[string]any),
	}, nil
}

func (p *Provider) GetUserInfo(ctx context.Context, accessToken string) (*auth.SSOUser, error) {
	return &auth.SSOUser{
		ExternalID: "noop-user-001",
		Email:      "user@example.com",
		FirstName:  "Demo",
		LastName:   "User",
		Groups:     []string{},
		Attributes: make(map[string]string),
		RawClaims:  make(map[string]any),
	}, nil
}

func (p *Provider) Metadata() auth.Metadata {
	return auth.Metadata{
		Name:     "noop",
		Protocol: "none",
		Issuer:   "noop",
	}
}
