package domain

import (
	"time"

	"github.com/google/uuid"
)

// AttributeTranslation represents a translated attribute definition
type AttributeTranslation struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	AttributeKey string   `json:"attribute_key"` // e.g., "thickness_mm", "voltage"
	Locale      string    `json:"locale"`        // e.g., "de", "en", "fr"
	DisplayName string    `json:"display_name"`  // e.g., "Dicke", "Spannung"
	Unit        *string   `json:"unit,omitempty"` // e.g., "mm", "V"
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductAttributeWithTranslation extends ProductAttribute with translation info
type ProductAttributeWithTranslation struct {
	Key         string        `json:"key"`
	Type        AttributeType `json:"type"`
	Value       any           `json:"value"`
	DisplayName string        `json:"display_name,omitempty"` // Translated name
	Unit        *string       `json:"unit,omitempty"`         // Translated unit
}

// NewAttributeTranslation creates a new attribute translation
func NewAttributeTranslation(tenantID uuid.UUID, attributeKey, locale, displayName string) *AttributeTranslation {
	now := time.Now()
	return &AttributeTranslation{
		ID:          uuid.New(),
		TenantID:    tenantID,
		AttributeKey: attributeKey,
		Locale:      locale,
		DisplayName: displayName,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// CreateAttributeTranslationRequest represents a request to create an attribute translation
type CreateAttributeTranslationRequest struct {
	AttributeKey string  `json:"attribute_key" binding:"required,min=1,max=100"`
	Locale      string  `json:"locale" binding:"required,len=2"`
	DisplayName string  `json:"display_name" binding:"required,min=1,max=200"`
	Unit        *string `json:"unit,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateAttributeTranslationRequest represents a request to update an attribute translation
type UpdateAttributeTranslationRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Unit        *string `json:"unit,omitempty"`
	Description *string `json:"description,omitempty"`
}

// AttributeTranslationFilter represents filter options for listing attribute translations
type AttributeTranslationFilter struct {
	TenantID     uuid.UUID
	AttributeKey *string
	Locale       *string
	Limit        int
	Offset       int
}
