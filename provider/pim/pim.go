// Package pim defines the interface for PIM (Product Information Management) system integrations.
package pim

import (
	"context"
	"io"
	"time"
)

// PIMProvider abstracts the communication with a PIM system (Akeneo, Pimcore, etc.).
type PIMProvider interface {
	// --- Products ---

	// FetchProducts retrieves products with cursor-based pagination.
	FetchProducts(ctx context.Context, filter ProductFilter) (*ProductPage, error)

	// FetchProduct retrieves a single product by identifier.
	FetchProduct(ctx context.Context, identifier string) (*Product, error)

	// --- Categories ---

	// FetchCategories retrieves all categories.
	FetchCategories(ctx context.Context) ([]Category, error)

	// --- Attributes ---

	// FetchAttributes retrieves attribute definitions.
	FetchAttributes(ctx context.Context) ([]Attribute, error)

	// --- Assets/Media ---

	// DownloadAsset downloads an asset file.
	DownloadAsset(ctx context.Context, assetCode string) (io.ReadCloser, string, error)

	// --- Metadata ---

	// Metadata returns information about this provider.
	Metadata() Metadata
}

// ProductFilter contains filter criteria for fetching products.
type ProductFilter struct {
	UpdatedSince *time.Time
	Families     []string
	Categories   []string
	Cursor       string // For pagination
	Limit        int
}

// ProductPage represents a page of products.
type ProductPage struct {
	Products   []Product
	NextCursor string
	TotalCount int
}

// Product represents a product from the PIM system.
type Product struct {
	Identifier string
	Family     string
	Categories []string
	Enabled    bool
	Values     map[string][]AttributeValue // Attribute name -> localized values
	Created    time.Time
	Updated    time.Time
}

// AttributeValue represents a localized/scoped attribute value.
type AttributeValue struct {
	Locale string
	Scope  string
	Data   any
}

// Category represents a product category.
type Category struct {
	Code   string
	Parent string
	Labels map[string]string // Locale -> Label
}

// Attribute represents an attribute definition.
type Attribute struct {
	Code        string
	Type        string // "text", "number", "select", "media", etc.
	Group       string
	Localizable bool
	Scopable    bool
	Labels      map[string]string
}

// Metadata provides information about the PIM provider implementation.
type Metadata struct {
	Name    string
	Version string
}
