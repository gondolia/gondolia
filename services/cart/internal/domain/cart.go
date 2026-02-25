package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CartStatus represents the status of a cart
type CartStatus string

const (
	CartStatusActive    CartStatus = "active"
	CartStatusMerged    CartStatus = "merged"
	CartStatusCompleted CartStatus = "completed"
)

// ProductType represents the type of product in the cart
type ProductType string

const (
	ProductTypeSimple     ProductType = "simple"
	ProductTypeVariant    ProductType = "variant"
	ProductTypeBundle     ProductType = "bundle"
	ProductTypeParametric ProductType = "parametric"
)

// Cart represents a shopping cart
type Cart struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`    // NULL for guest carts
	SessionID *string    `json:"session_id,omitempty"` // For guest carts
	Status    CartStatus `json:"status"`
	Subtotal  float64    `json:"subtotal"`             // Computed from items
	Currency  string     `json:"currency"`             // Currency of cart items
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Items     []CartItem `json:"items,omitempty"` // Populated when fetching cart with items
}

// CartItem represents an item in a cart
type CartItem struct {
	ID            uuid.UUID       `json:"id"`
	CartID        uuid.UUID       `json:"cart_id"`
	ProductID     uuid.UUID       `json:"product_id"`
	VariantID     *uuid.UUID      `json:"variant_id,omitempty"`
	ProductType   ProductType     `json:"product_type"`
	ProductName   string          `json:"product_name"`
	SKU           string          `json:"sku"`
	ImageURL      string          `json:"image_url,omitempty"`
	Quantity      int             `json:"quantity"`
	UnitPrice     float64         `json:"unit_price"`
	TotalPrice    float64         `json:"total_price"`
	Currency      string          `json:"currency"`
	Configuration *Configuration  `json:"configuration,omitempty"` // JSONB for Bundle/Parametric configs
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// Configuration stores product configuration for bundles and parametric items
type Configuration struct {
	// For bundles: component selections
	BundleComponents []BundleComponentSelection `json:"bundle_components,omitempty"`

	// For parametric: parameter values
	ParametricParams map[string]interface{} `json:"parametric_params,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling to handle both
// frontend (camelCase) and backend (snake_case) formats
func (c *Configuration) UnmarshalJSON(data []byte) error {
	// Try to unmarshal into a flexible structure
	var raw struct {
		// snake_case (backend format)
		BundleComponents []json.RawMessage      `json:"bundle_components"`
		ParametricParams map[string]interface{} `json:"parametric_params"`
		// camelCase (frontend format)
		BundleComponentsCamel []json.RawMessage `json:"bundleComponents"`
		// Frontend parametric format: { parameters: {...}, selections: {...} }
		Parameters map[string]interface{} `json:"parameters"`
		Selections map[string]string      `json:"selections"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Handle parametric configuration
	// Priority: parametric_params > parameters/selections
	if raw.ParametricParams != nil {
		c.ParametricParams = raw.ParametricParams
	} else if raw.Parameters != nil || raw.Selections != nil {
		// Frontend sends { parameters: {...}, selections: {...} }
		// Store them in ParametricParams for the service to use
		c.ParametricParams = make(map[string]interface{})
		if raw.Parameters != nil {
			c.ParametricParams["parameters"] = raw.Parameters
		}
		if raw.Selections != nil {
			c.ParametricParams["selections"] = raw.Selections
		}
	}

	// Use whichever format is present
	components := raw.BundleComponents
	if len(components) == 0 && len(raw.BundleComponentsCamel) > 0 {
		components = raw.BundleComponentsCamel
	}

	// Parse each component with flexible field names
	for _, compData := range components {
		var comp BundleComponentSelection
		if err := json.Unmarshal(compData, &comp); err != nil {
			return err
		}
		c.BundleComponents = append(c.BundleComponents, comp)
	}

	return nil
}

// BundleComponentSelection represents a selected component in a bundle
type BundleComponentSelection struct {
	ComponentID uuid.UUID  `json:"component_id"`
	ProductID   uuid.UUID  `json:"product_id"`
	VariantID   *uuid.UUID `json:"variant_id,omitempty"`
	Quantity    int        `json:"quantity"`
	// Additional fields for parametric components in bundles
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Selections map[string]string      `json:"selections,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling to handle both
// frontend (camelCase) and backend (snake_case) formats
func (b *BundleComponentSelection) UnmarshalJSON(data []byte) error {
	var raw struct {
		// snake_case
		ComponentID string                 `json:"component_id"`
		ProductID   string                 `json:"product_id"`
		VariantID   *string                `json:"variant_id"`
		Quantity    int                    `json:"quantity"`
		Parameters  map[string]interface{} `json:"parameters"`
		Selections  map[string]string      `json:"selections"`
		// camelCase
		ComponentIDCamel string  `json:"componentId"`
		ProductIDCamel   string  `json:"productId"`
		VariantIDCamel   *string `json:"variantId"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Use camelCase if snake_case is empty
	componentIDStr := raw.ComponentID
	if componentIDStr == "" {
		componentIDStr = raw.ComponentIDCamel
	}
	productIDStr := raw.ProductID
	if productIDStr == "" {
		productIDStr = raw.ProductIDCamel
	}
	variantIDStr := raw.VariantID
	if variantIDStr == nil {
		variantIDStr = raw.VariantIDCamel
	}

	// Parse UUIDs
	if componentIDStr != "" {
		id, err := uuid.Parse(componentIDStr)
		if err != nil {
			return err
		}
		b.ComponentID = id
	}
	if productIDStr != "" {
		id, err := uuid.Parse(productIDStr)
		if err != nil {
			return err
		}
		b.ProductID = id
	}
	if variantIDStr != nil && *variantIDStr != "" {
		id, err := uuid.Parse(*variantIDStr)
		if err != nil {
			return err
		}
		b.VariantID = &id
	}

	b.Quantity = raw.Quantity
	b.Parameters = raw.Parameters
	b.Selections = raw.Selections

	return nil
}

// AddItemRequest represents a request to add an item to cart
type AddItemRequest struct {
	ProductID     uuid.UUID      `json:"product_id" binding:"required"`
	VariantID     *uuid.UUID     `json:"variant_id,omitempty"`
	Quantity      int            `json:"quantity" binding:"required,min=1"`
	Configuration *Configuration `json:"configuration,omitempty"`
}

// UpdateItemRequest represents a request to update cart item quantity
type UpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// CartFilter represents filter options for cart queries
type CartFilter struct {
	TenantID  uuid.UUID
	UserID    *uuid.UUID
	SessionID *string
	Status    *CartStatus
}

// NewCart creates a new cart
func NewCart(tenantID uuid.UUID, userID *uuid.UUID, sessionID *string) *Cart {
	now := time.Now()
	return &Cart{
		ID:        uuid.New(),
		TenantID:  tenantID,
		UserID:    userID,
		SessionID: sessionID,
		Status:    CartStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
		Items:     []CartItem{},
	}
}

// IsActive returns true if cart is active
func (c *Cart) IsActive() bool {
	return c.Status == CartStatusActive
}

// TotalItems returns the total number of items in cart
func (c *Cart) TotalItems() int {
	total := 0
	for _, item := range c.Items {
		total += item.Quantity
	}
	return total
}

// TotalPrice returns the total price of all items
func (c *Cart) TotalPrice() float64 {
	total := 0.0
	for _, item := range c.Items {
		total += item.TotalPrice
	}
	return total
}

// GetCurrency returns the currency of the first item, or CHF as default
func (c *Cart) GetCurrency() string {
	if len(c.Items) > 0 {
		return c.Items[0].Currency
	}
	return "CHF"
}
