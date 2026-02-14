package domain

import (
	"time"

	"github.com/google/uuid"
)

// Price represents a B2B contract price for a product
type Price struct {
	ID              uuid.UUID  `json:"id"`
	TenantID        uuid.UUID  `json:"tenant_id"`
	ProductID       uuid.UUID  `json:"product_id"`
	CustomerGroupID *uuid.UUID `json:"customer_group_id,omitempty"` // nil = base price
	MinQuantity     int        `json:"min_quantity"`
	Price           float64    `json:"price"`
	Currency        string     `json:"currency"`
	ValidFrom       *time.Time `json:"valid_from,omitempty"`
	ValidTo         *time.Time `json:"valid_to,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// IsValid returns true if the price is valid at the given time
func (p *Price) IsValid(at time.Time) bool {
	if p.DeletedAt != nil {
		return false
	}
	if p.ValidFrom != nil && at.Before(*p.ValidFrom) {
		return false
	}
	if p.ValidTo != nil && at.After(*p.ValidTo) {
		return false
	}
	return true
}

// IsCurrentlyValid returns true if the price is valid now
func (p *Price) IsCurrentlyValid() bool {
	return p.IsValid(time.Now())
}

// NewPrice creates a new price with defaults
func NewPrice(tenantID, productID uuid.UUID, price float64, currency string) *Price {
	now := time.Now()
	return &Price{
		ID:          uuid.New(),
		TenantID:    tenantID,
		ProductID:   productID,
		MinQuantity: 1,
		Price:       price,
		Currency:    currency,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// CreatePriceRequest represents a request to create a price
type CreatePriceRequest struct {
	CustomerGroupID *uuid.UUID `json:"customer_group_id,omitempty"`
	MinQuantity     *int       `json:"min_quantity,omitempty"`
	Price           float64    `json:"price" binding:"required,gt=0"`
	Currency        string     `json:"currency" binding:"required,len=3"`
	ValidFrom       *time.Time `json:"valid_from,omitempty"`
	ValidTo         *time.Time `json:"valid_to,omitempty"`
}

// UpdatePriceRequest represents a request to update a price
type UpdatePriceRequest struct {
	MinQuantity     *int       `json:"min_quantity,omitempty"`
	Price           *float64   `json:"price,omitempty"`
	Currency        *string    `json:"currency,omitempty"`
	ValidFrom       *time.Time `json:"valid_from,omitempty"`
	ValidTo         *time.Time `json:"valid_to,omitempty"`
}

// PriceFilter represents filter options for listing prices
type PriceFilter struct {
	TenantID        uuid.UUID
	ProductID       *uuid.UUID
	CustomerGroupID *uuid.UUID
	ValidAt         *time.Time
	Currency        *string
	Limit           int
	Offset          int
}
