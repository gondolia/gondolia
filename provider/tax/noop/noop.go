// Package noop provides a no-op implementation of the Tax provider.
package noop

import (
	"context"

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/tax"
)

func init() {
	provider.Register[tax.TaxProvider]("tax", "noop",
		provider.Metadata{
			Name:        "noop",
			DisplayName: "No-Op Tax Provider",
			Category:    "tax",
			Version:     "1.0.0",
			Description: "A no-operation tax provider for development and testing",
			ConfigSpec:  []provider.ConfigField{},
		},
		NewProvider,
	)
}

// Provider is a no-op Tax provider.
type Provider struct{}

// NewProvider creates a new no-op Tax provider.
func NewProvider(config map[string]any) (tax.TaxProvider, error) {
	return &Provider{}, nil
}

func (p *Provider) CalculateTax(ctx context.Context, req tax.TaxRequest) (*tax.TaxResult, error) {
	items := make([]tax.TaxItemResult, len(req.Items))
	var totalTax int64
	for i, item := range req.Items {
		items[i] = tax.TaxItemResult{
			SKU:       item.SKU,
			TaxRate:   0,
			TaxAmount: 0,
		}
	}
	return &tax.TaxResult{
		Items:    items,
		TotalTax: totalTax,
	}, nil
}

func (p *Provider) Metadata() tax.Metadata {
	return tax.Metadata{
		Name: "noop",
	}
}
