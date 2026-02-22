package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrBundleNotFound        = errors.New("bundle not found")
	ErrInvalidBundleMode     = errors.New("invalid bundle mode")
	ErrInvalidPriceMode      = errors.New("invalid price mode")
	ErrInvalidComponentType  = errors.New("invalid component type - must be simple, variant, or parametric")
	ErrComponentNotFound     = errors.New("component not found")
	ErrInvalidQuantity       = errors.New("invalid quantity")
	ErrQuantityOutOfRange    = errors.New("quantity out of allowed range")
	ErrBundleNesting         = errors.New("bundle nesting not allowed")
	ErrVariantParentInBundle = errors.New("variant_parent not allowed in bundle - use specific variant")
)

// BundleMode defines how customers can configure component quantities
type BundleMode string

const (
	BundleModeFixed        BundleMode = "fixed"        // Admin sets fixed quantities
	BundleModeConfigurable BundleMode = "configurable" // Customer can adjust quantities
)

// BundlePriceMode defines how bundle price is calculated
type BundlePriceMode string

const (
	BundlePriceModeComputed BundlePriceMode = "computed" // Sum of component prices
	BundlePriceModeFixed    BundlePriceMode = "fixed"    // Fixed bundle price
)

// BundleComponent represents a component in a bundle
type BundleComponent struct {
	ID                 uuid.UUID              `json:"id"`
	TenantID           uuid.UUID              `json:"tenant_id"`
	BundleProductID    uuid.UUID              `json:"bundle_product_id"`
	ComponentProductID uuid.UUID              `json:"component_product_id"`
	Quantity           int                    `json:"quantity"`
	MinQuantity        *int                   `json:"min_quantity,omitempty"`
	MaxQuantity        *int                   `json:"max_quantity,omitempty"`
	SortOrder          int                    `json:"sort_order"`
	DefaultParameters  map[string]interface{} `json:"default_parameters,omitempty"` // For parametric components
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`

	// Populated when loading components with product details
	Product *Product `json:"product,omitempty"`
}

// BundleComponentRequest represents a component in a bundle creation/update request
type BundleComponentRequest struct {
	ComponentProductID uuid.UUID              `json:"component_product_id" binding:"required"`
	Quantity           int                    `json:"quantity" binding:"required,min=1"`
	MinQuantity        *int                   `json:"min_quantity,omitempty"`
	MaxQuantity        *int                   `json:"max_quantity,omitempty"`
	SortOrder          int                    `json:"sort_order"`
	DefaultParameters  map[string]interface{} `json:"default_parameters,omitempty"`
}

// SetBundleComponentsRequest represents the request to set bundle components
type SetBundleComponentsRequest struct {
	Components []BundleComponentRequest `json:"components" binding:"required,dive"`
}

// BundlePriceRequest represents a request to calculate bundle price
type BundlePriceRequest struct {
	Components []BundleComponentPriceRequest `json:"components" binding:"required,dive"`
}

// BundleComponentPriceRequest represents a component in a price calculation request
type BundleComponentPriceRequest struct {
	ComponentID uuid.UUID              `json:"component_id" binding:"required"` // component record ID, not product ID
	Quantity    int                    `json:"quantity" binding:"required,min=1"`
	Parameters  map[string]float64     `json:"parameters,omitempty"` // For parametric components (range axes)
	Selections  map[string]string      `json:"selections,omitempty"` // For parametric components (select axes)
}

// BundlePriceResponse represents the response to a price calculation request
type BundlePriceResponse struct {
	PriceMode  BundlePriceMode                  `json:"price_mode"`
	Components []BundleComponentPriceResponse   `json:"components,omitempty"`
	Total      float64                          `json:"total"`
	Currency   string                           `json:"currency"`
}

// BundleComponentPriceResponse represents a component's calculated price
type BundleComponentPriceResponse struct {
	ComponentID uuid.UUID `json:"component_id"`
	ProductID   uuid.UUID `json:"product_id"`
	SKU         string    `json:"sku,omitempty"`
	UnitPrice   float64   `json:"unit_price"`
	Quantity    int       `json:"quantity"`
	LineTotal   float64   `json:"line_total"`
}

// ValidateComponent checks if a product can be used as a bundle component
func ValidateComponent(product *Product) error {
	if product == nil {
		return ErrComponentNotFound
	}

	// Check product type
	switch product.ProductType {
	case ProductTypeSimple, ProductTypeVariant, ProductTypeParametric:
		// Valid component types
		return nil
	case ProductTypeVariantParent:
		return ErrVariantParentInBundle
	case ProductTypeBundle:
		return ErrBundleNesting
	default:
		return fmt.Errorf("%w: %s", ErrInvalidComponentType, product.ProductType)
	}
}

// ValidateComponentQuantity validates quantity against min/max constraints
func ValidateComponentQuantity(component *BundleComponent, requestedQty int) error {
	if requestedQty <= 0 {
		return ErrInvalidQuantity
	}

	if component.MinQuantity != nil && requestedQty < *component.MinQuantity {
		return fmt.Errorf("%w: minimum %d", ErrQuantityOutOfRange, *component.MinQuantity)
	}

	if component.MaxQuantity != nil && requestedQty > *component.MaxQuantity {
		return fmt.Errorf("%w: maximum %d", ErrQuantityOutOfRange, *component.MaxQuantity)
	}

	return nil
}

// ValidateBundleMode validates a bundle mode string
func ValidateBundleMode(mode string) error {
	switch BundleMode(mode) {
	case BundleModeFixed, BundleModeConfigurable:
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrInvalidBundleMode, mode)
	}
}

// ValidateBundlePriceMode validates a bundle price mode string
func ValidateBundlePriceMode(mode string) error {
	switch BundlePriceMode(mode) {
	case BundlePriceModeComputed, BundlePriceModeFixed:
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrInvalidPriceMode, mode)
	}
}
