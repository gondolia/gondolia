package domain

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	// MaxVariantAxes is the maximum number of variant axes per parent product.
	// 4 axes with 5 values each = 625 theoretical combinations — UX limit.
	MaxVariantAxes = 4

	// SoftVariantLimit is the recommended maximum number of variants per parent.
	// Above this, the product should probably be parametric instead.
	SoftVariantLimit = 200
)

// GenerateVariantName creates a display name from parent name + axis labels.
// Format: "Parent Name — Axis1 Label, Axis2 Label, ..."
// Returns parent name unchanged if no axis values are provided.
func GenerateVariantName(parentName map[string]string, axisValues []AxisValueEntry, locale string) string {
	name, ok := parentName[locale]
	if !ok {
		// Fallback to any available locale
		for _, v := range parentName {
			name = v
			break
		}
	}
	if len(axisValues) == 0 {
		return name
	}

	labels := make([]string, 0, len(axisValues))
	for _, av := range axisValues {
		label := av.OptionCode // fallback to code
		if len(av.OptionLabel) > 0 {
			if l, ok := av.OptionLabel[locale]; ok && l != "" {
				label = l
			} else {
				// fallback to any locale
				for _, l := range av.OptionLabel {
					if l != "" {
						label = l
						break
					}
				}
			}
		}
		labels = append(labels, label)
	}

	return fmt.Sprintf("%s — %s", name, strings.Join(labels, ", "))
}

// GetEffectiveImages returns variant images if present, otherwise parent images.
// Design decision: variant images REPLACE parent images (no merge).
func GetEffectiveImages(variantImages, parentImages []ProductImage) []ProductImage {
	if len(variantImages) > 0 {
		return variantImages
	}
	return parentImages
}

// VariantAxis defines a variant axis on a parent product or parametric product
type VariantAxis struct {
	ID            uuid.UUID         `json:"id"`
	ProductID     uuid.UUID         `json:"product_id"`
	AttributeCode string            `json:"attribute_code"`
	Position      int               `json:"position"`

	// Parametric fields (only used when InputType == "range")
	InputType string   `json:"input_type"`           // "select" or "range"
	MinValue  *float64 `json:"min_value,omitempty"`
	MaxValue  *float64 `json:"max_value,omitempty"`
	StepValue *float64 `json:"step_value,omitempty"`
	Unit      string   `json:"unit,omitempty"`

	// Resolved fields (populated by service layer)
	Label   map[string]string `json:"label,omitempty"`   // i18n label of the attribute
	Options []AxisOption      `json:"options,omitempty"` // Available values for this axis (select type)
}

// AxisOption represents a possible value for a variant axis
type AxisOption struct {
	Code     string            `json:"code"`
	Label    map[string]string `json:"label"`
	Position int               `json:"position"`

	// Dynamic: is this option available with the current selection?
	Available *bool `json:"available,omitempty"`
}

// AxisValueEntry stores the axis value for a specific variant
type AxisValueEntry struct {
	VariantID         uuid.UUID         `json:"variant_id,omitempty"`
	AxisID            uuid.UUID         `json:"axis_id,omitempty"`
	AxisAttributeCode string            `json:"axis_attribute_code"`
	OptionCode        string            `json:"option_code"`

	// Resolved fields (populated by service layer)
	AxisLabel   map[string]string `json:"axis_label,omitempty"`
	OptionLabel map[string]string `json:"option_label,omitempty"`
}

// ProductVariant is a compact representation of a variant within the parent response
type ProductVariant struct {
	ID         uuid.UUID              `json:"id"`
	SKU        string                 `json:"sku"`
	AxisValues map[string]string      `json:"axis_values"` // attribute_code -> option_code
	Status     ProductStatus          `json:"status"`
	Images     []ProductImage         `json:"images,omitempty"`
	Price      *VariantPrice          `json:"price,omitempty"`
	Stock      *VariantAvailability   `json:"availability,omitempty"`
}

// VariantPrice contains the resolved price of a variant
type VariantPrice struct {
	Net        float64     `json:"net"`
	Currency   string      `json:"currency"`
	TierPrices []TierPrice `json:"tier_prices,omitempty"`
}

// TierPrice represents a quantity-based price tier
type TierPrice struct {
	MinQuantity int     `json:"min_quantity"`
	Price       float64 `json:"price"`
}

// VariantAvailability contains the availability status of a variant
type VariantAvailability struct {
	InStock  bool `json:"in_stock"`
	Quantity *int `json:"quantity,omitempty"` // Optional: concrete amount
}

// CreateVariantParentRequest represents a request to create a variant parent product
type CreateVariantParentRequest struct {
	SKU          string                 `json:"sku" binding:"required,min=1,max=100"`
	Name         map[string]string      `json:"name" binding:"required"`
	Description  map[string]string      `json:"description,omitempty"`
	CategoryIDs  []uuid.UUID            `json:"category_ids,omitempty"`
	Attributes   []ProductAttribute     `json:"attributes,omitempty"`
	Images       []ProductImage         `json:"images,omitempty"`
	VariantAxes  []CreateVariantAxis    `json:"variant_axes" binding:"required,min=1"`
}

// CreateVariantAxis represents an axis definition in the create request
type CreateVariantAxis struct {
	AttributeCode string `json:"attribute_code" binding:"required"`
	Position      int    `json:"position"`
}

// CreateVariantRequest represents a request to create a variant under a parent
type CreateVariantRequest struct {
	SKU        string                 `json:"sku" binding:"required,min=1,max=100"`
	Name       map[string]string      `json:"name,omitempty"`
	AxisValues map[string]string      `json:"axis_values" binding:"required,min=1"` // attribute_code -> option_code
	Attributes []ProductAttribute     `json:"attributes,omitempty"`
	Images     []ProductImage         `json:"images,omitempty"`
}

// VariantSelectionRequest represents axis values for variant selection
type VariantSelectionRequest struct {
	AxisValues map[string]string `json:"axis_values"` // attribute_code -> option_code
}

// PriceRange represents min/max price for a variant parent
type PriceRange struct {
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Currency string  `json:"currency"`
}
