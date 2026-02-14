package domain

import (
	"time"

	"github.com/google/uuid"
)

// Category represents a product category
type Category struct {
	ID        uuid.UUID         `json:"id"`
	TenantID  uuid.UUID         `json:"tenant_id"`
	Code      string            `json:"code"`
	ParentID  *uuid.UUID        `json:"parent_id,omitempty"`
	Name      map[string]string `json:"name"` // locale -> name
	SortOrder int               `json:"sort_order"`
	Active    bool              `json:"active"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	DeletedAt *time.Time        `json:"deleted_at,omitempty"`

	// PIM Integration
	PIMCode      *string    `json:"pim_code,omitempty"`
	LastSyncedAt *time.Time `json:"last_synced_at,omitempty"`

	// Loaded relations (not always populated)
	Children []Category `json:"children,omitempty"`
	
	// Computed fields
	ProductCount int `json:"product_count,omitempty"`
}

// GetLocalizedName returns the name in the specified locale or fallback
func (c *Category) GetLocalizedName(locale string) string {
	if name, ok := c.Name[locale]; ok {
		return name
	}
	// Fallback to first available
	for _, name := range c.Name {
		return name
	}
	return ""
}

// IsRoot returns true if category has no parent
func (c *Category) IsRoot() bool {
	return c.ParentID == nil
}

// NewCategory creates a new category with defaults
func NewCategory(tenantID uuid.UUID, code string) *Category {
	now := time.Now()
	return &Category{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Code:      code,
		Name:      make(map[string]string),
		SortOrder: 0,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
		Children:  []Category{},
	}
}

// CreateCategoryRequest represents a request to create a category
type CreateCategoryRequest struct {
	Code      string            `json:"code" binding:"required,min=1,max=100"`
	ParentID  *uuid.UUID        `json:"parent_id,omitempty"`
	Name      map[string]string `json:"name" binding:"required"`
	SortOrder *int              `json:"sort_order,omitempty"`
	Active    *bool             `json:"active,omitempty"`
}

// UpdateCategoryRequest represents a request to update a category
type UpdateCategoryRequest struct {
	ParentID  *uuid.UUID        `json:"parent_id,omitempty"`
	Name      map[string]string `json:"name,omitempty"`
	SortOrder *int              `json:"sort_order,omitempty"`
	Active    *bool             `json:"active,omitempty"`
}

// CategoryFilter represents filter options for listing categories
type CategoryFilter struct {
	TenantID uuid.UUID
	ParentID *uuid.UUID
	Active   *bool
	Search   *string // Searches in code, name
	Limit    int
	Offset   int
}
