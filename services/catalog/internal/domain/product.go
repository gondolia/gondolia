package domain

import (
	"time"

	"github.com/google/uuid"
)

// ProductStatus represents the status of a product
type ProductStatus string

const (
	ProductStatusDraft    ProductStatus = "draft"
	ProductStatusActive   ProductStatus = "active"
	ProductStatusArchived ProductStatus = "archived"
)

// ProductType represents the type of a product
type ProductType string

const (
	ProductTypeSimple        ProductType = "simple"
	ProductTypeVariantParent ProductType = "variant_parent"
	ProductTypeVariant       ProductType = "variant"
)

// AttributeType represents the data type of a product attribute
type AttributeType string

const (
	AttributeTypeText    AttributeType = "text"
	AttributeTypeNumber  AttributeType = "number"
	AttributeTypeBoolean AttributeType = "boolean"
	AttributeTypeDate    AttributeType = "date"
)

// Product represents a product in the catalog
type Product struct {
	ID          uuid.UUID            `json:"id"`
	TenantID    uuid.UUID            `json:"tenant_id"`
	ProductType ProductType          `json:"product_type"`
	ParentID    *uuid.UUID           `json:"parent_id,omitempty"`
	SKU         string               `json:"sku"`
	Name        map[string]string    `json:"name"`        // locale -> name
	Description map[string]string    `json:"description"` // locale -> description
	CategoryIDs []uuid.UUID          `json:"category_ids"`
	Attributes  []ProductAttribute   `json:"attributes"`
	Status      ProductStatus        `json:"status"`
	Images      []ProductImage       `json:"images"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	DeletedAt   *time.Time           `json:"deleted_at,omitempty"`

	// PIM Integration
	PIMIdentifier *string    `json:"pim_identifier,omitempty"`
	LastSyncedAt  *time.Time `json:"last_synced_at,omitempty"`

	// Variant-specific fields (populated based on product_type)
	VariantAxes      []VariantAxis      `json:"variant_axes,omitempty"`       // Only for variant_parent
	Variants         []ProductVariant   `json:"variants,omitempty"`           // Only for variant_parent
	AxisValues       []AxisValueEntry   `json:"axis_values,omitempty"`        // Only for variant
}

// ProductAttribute represents a flexible product attribute
type ProductAttribute struct {
	Key   string        `json:"key"`
	Type  AttributeType `json:"type"`
	Value any           `json:"value"`
}

// ProductImage represents a product image
type ProductImage struct {
	URL       string `json:"url"`
	AltText   string `json:"alt_text,omitempty"`
	SortOrder int    `json:"sort_order"`
	IsPrimary bool   `json:"is_primary"`
}

// GetLocalizedName returns the name in the specified locale or fallback
func (p *Product) GetLocalizedName(locale string) string {
	if name, ok := p.Name[locale]; ok {
		return name
	}
	// Fallback to first available
	for _, name := range p.Name {
		return name
	}
	return ""
}

// GetLocalizedDescription returns the description in the specified locale or fallback
func (p *Product) GetLocalizedDescription(locale string) string {
	if desc, ok := p.Description[locale]; ok {
		return desc
	}
	// Fallback to first available
	for _, desc := range p.Description {
		return desc
	}
	return ""
}

// IsActive returns true if product is active
func (p *Product) IsActive() bool {
	return p.Status == ProductStatusActive && p.DeletedAt == nil
}

// NewProduct creates a new product with defaults
func NewProduct(tenantID uuid.UUID, sku string) *Product {
	now := time.Now()
	return &Product{
		ID:          uuid.New(),
		TenantID:    tenantID,
		ProductType: ProductTypeSimple,
		SKU:         sku,
		Name:        make(map[string]string),
		Description: make(map[string]string),
		CategoryIDs: []uuid.UUID{},
		Attributes:  []ProductAttribute{},
		Status:      ProductStatusDraft,
		Images:      []ProductImage{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// IsVariantParent returns true if product is a variant parent
func (p *Product) IsVariantParent() bool {
	return p.ProductType == ProductTypeVariantParent
}

// IsVariant returns true if product is a variant
func (p *Product) IsVariant() bool {
	return p.ProductType == ProductTypeVariant
}

// GetEffectiveName returns the effective name for a variant (with inheritance from parent)
func (p *Product) GetEffectiveName(parent *Product, locale string) string {
	// Variant has own name?
	if p.IsVariant() {
		if name, ok := p.Name[locale]; ok && name != "" {
			return name
		}
		// Fallback to parent name (axis values appended by service layer)
		if parent != nil {
			return parent.GetLocalizedName(locale)
		}
	}
	return p.GetLocalizedName(locale)
}

// CreateProductRequest represents a request to create a product
type CreateProductRequest struct {
	ProductType ProductType                  `json:"product_type,omitempty"` // Defaults to 'simple'
	SKU         string                       `json:"sku" binding:"required,min=1,max=100"`
	Name        map[string]string            `json:"name" binding:"required"`
	Description map[string]string            `json:"description,omitempty"`
	CategoryIDs []uuid.UUID                  `json:"category_ids,omitempty"`
	Attributes  []ProductAttribute           `json:"attributes,omitempty"`
	Status      ProductStatus                `json:"status,omitempty"`
	Images      []ProductImage               `json:"images,omitempty"`
	VariantAxes []CreateVariantAxis          `json:"variant_axes,omitempty"` // Only for variant_parent
}

// UpdateProductRequest represents a request to update a product
type UpdateProductRequest struct {
	Name        map[string]string    `json:"name,omitempty"`
	Description map[string]string    `json:"description,omitempty"`
	CategoryIDs []uuid.UUID          `json:"category_ids,omitempty"`
	Attributes  []ProductAttribute   `json:"attributes,omitempty"`
	Status      *ProductStatus       `json:"status,omitempty"`
	Images      []ProductImage       `json:"images,omitempty"`
}

// ProductFilter represents filter options for listing products
type ProductFilter struct {
	TenantID        uuid.UUID
	CategoryID      *uuid.UUID
	IncludeChildren bool   // Include products from child categories
	Status          *ProductStatus
	ProductType     *ProductType // Filter by product type
	ParentID        *uuid.UUID   // Filter variants by parent
	ExcludeVariants bool         // Exclude variant products from list (show only simple + variant_parent)
	Search          *string      // Searches in SKU, name
	SKUs            []string
	Limit           int
	Offset          int
}
