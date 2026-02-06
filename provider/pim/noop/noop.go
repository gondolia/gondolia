// Package noop provides a no-op implementation of the PIM provider.
package noop

import (
	"context"
	"io"

	"github.com/gondolia/gondolia/provider"
	"github.com/gondolia/gondolia/provider/pim"
)

func init() {
	provider.Register[pim.PIMProvider]("pim", "noop",
		provider.Metadata{
			Name:        "noop",
			DisplayName: "No-Op PIM Provider",
			Category:    "pim",
			Version:     "1.0.0",
			Description: "A no-operation PIM provider for development and testing",
			ConfigSpec:  []provider.ConfigField{},
		},
		NewProvider,
	)
}

// Provider is a no-op PIM provider.
type Provider struct{}

// NewProvider creates a new no-op PIM provider.
func NewProvider(config map[string]any) (pim.PIMProvider, error) {
	return &Provider{}, nil
}

func (p *Provider) FetchProducts(ctx context.Context, filter pim.ProductFilter) (*pim.ProductPage, error) {
	return &pim.ProductPage{
		Products:   []pim.Product{},
		NextCursor: "",
		TotalCount: 0,
	}, nil
}

func (p *Provider) FetchProduct(ctx context.Context, identifier string) (*pim.Product, error) {
	return &pim.Product{
		Identifier: identifier,
		Family:     "",
		Categories: []string{},
		Enabled:    false,
		Values:     make(map[string][]pim.AttributeValue),
	}, nil
}

func (p *Provider) FetchCategories(ctx context.Context) ([]pim.Category, error) {
	return []pim.Category{}, nil
}

func (p *Provider) FetchAttributes(ctx context.Context) ([]pim.Attribute, error) {
	return []pim.Attribute{}, nil
}

func (p *Provider) DownloadAsset(ctx context.Context, assetCode string) (io.ReadCloser, string, error) {
	return io.NopCloser(nil), "", nil
}

func (p *Provider) Metadata() pim.Metadata {
	return pim.Metadata{
		Name:    "noop",
		Version: "1.0.0",
	}
}
