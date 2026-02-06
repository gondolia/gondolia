// Package tax defines the interface for tax calculation integrations.
package tax

import "context"

// TaxProvider abstracts tax calculations (internal, Avalara, TaxJar, etc.).
type TaxProvider interface {
	// CalculateTax calculates taxes for a list of items.
	CalculateTax(ctx context.Context, req TaxRequest) (*TaxResult, error)

	// Metadata returns provider information.
	Metadata() Metadata
}

// TaxRequest contains data for tax calculation.
type TaxRequest struct {
	Country  string // ISO 3166-1 alpha-2
	Region   string
	Currency string
	Items    []TaxItem
}

// TaxItem represents an item for tax calculation.
type TaxItem struct {
	SKU       string
	Quantity  float64
	UnitPrice int64 // Cents
	TaxCode   string
}

// TaxResult is the result of tax calculation.
type TaxResult struct {
	Items    []TaxItemResult
	TotalTax int64
}

// TaxItemResult represents the tax calculation for an item.
type TaxItemResult struct {
	SKU       string
	TaxRate   float64 // e.g. 0.081 for 8.1%
	TaxAmount int64
}

// Metadata provides information about the tax provider.
type Metadata struct {
	Name string
}
