// Package noop provides a no-op implementation of the CRM provider.
package noop

import (
	"context"

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/crm"
)

func init() {
	provider.Register[crm.CRMProvider]("crm", "noop",
		provider.Metadata{
			Name:        "noop",
			DisplayName: "No-Op CRM Provider",
			Category:    "crm",
			Version:     "1.0.0",
			Description: "A no-operation CRM provider for development and testing",
			ConfigSpec:  []provider.ConfigField{},
		},
		NewProvider,
	)
}

// Provider is a no-op CRM provider.
type Provider struct{}

// NewProvider creates a new no-op CRM provider.
func NewProvider(config map[string]any) (crm.CRMProvider, error) {
	return &Provider{}, nil
}

func (p *Provider) SyncContact(ctx context.Context, contact crm.Contact) (*crm.SyncResult, error) {
	return &crm.SyncResult{
		ExternalID: "noop-contact-001",
		Action:     "unchanged",
	}, nil
}

func (p *Provider) GetAccount(ctx context.Context, accountID string) (*crm.Account, error) {
	return &crm.Account{
		ExternalID:    accountID,
		Name:          "Demo Company",
		ERPCustomerID: "",
		Attributes:    make(map[string]any),
	}, nil
}

func (p *Provider) ListAccounts(ctx context.Context, filter crm.AccountFilter) ([]crm.Account, error) {
	return []crm.Account{}, nil
}

func (p *Provider) Metadata() crm.Metadata {
	return crm.Metadata{
		Name: "noop",
	}
}
