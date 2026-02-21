package domain

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

var (
	ErrParametricPricingNotFound = errors.New("parametric pricing not found")
	ErrParameterOutOfRange       = errors.New("parameter value out of range")
	ErrParameterInvalidStep      = errors.New("parameter value does not match step")
	ErrMissingParameter          = errors.New("missing required parameter")
	ErrInvalidFormulaType        = errors.New("invalid formula type")
	ErrSKUMappingNotFound        = errors.New("no SKU mapping found for selected options")
	ErrMissingSelection          = errors.New("missing required selection")
)

// FormulaType defines how parametric pricing is calculated
type FormulaType string

const (
	FormulaTypeFixed          FormulaType = "fixed"
	FormulaTypePerUnit        FormulaType = "per_unit"
	FormulaTypePerM2          FormulaType = "per_m2"
	FormulaTypePerRunningMeter FormulaType = "per_running_meter"
)

// ParametricPricing defines pricing rules for a parametric product
type ParametricPricing struct {
	ID            uuid.UUID  `json:"id"`
	ProductID     uuid.UUID  `json:"product_id"`
	FormulaType   string     `json:"formula_type"`
	BasePrice     float64    `json:"base_price"`
	UnitPrice     *float64   `json:"unit_price,omitempty"`
	Currency      string     `json:"currency"`
	MinOrderValue *float64   `json:"min_order_value,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ParametricPriceRequest is the input for price calculation
type ParametricPriceRequest struct {
	Parameters map[string]float64 `json:"parameters"` // axis_code -> numeric value
	Selections map[string]string  `json:"selections"` // axis_code -> option code (for select axes)
	Quantity   int                `json:"quantity"`
}

// ParametricPriceResponse is the result of a price calculation
type ParametricPriceResponse struct {
	SKU        string             `json:"sku"`
	UnitPrice  float64            `json:"unit_price"`
	TotalPrice float64            `json:"total_price"`
	Currency   string             `json:"currency"`
	Quantity   int                `json:"quantity"`
	Breakdown  map[string]float64 `json:"breakdown,omitempty"`
}

// SKUMapping maps a combination of select axis values to a specific SKU with pricing
type SKUMapping struct {
	ID         uuid.UUID          `json:"id"`
	ProductID  uuid.UUID          `json:"product_id"`
	Selections map[string]string  `json:"selections"`
	SKU        string             `json:"sku"`
	UnitPrice  float64            `json:"unit_price"`
	BasePrice  float64            `json:"base_price"`
	Stock      *int               `json:"stock,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

// ValidateParameters checks that all range axes have valid values
func ValidateParameters(axes []VariantAxis, params map[string]float64) error {
	for _, axis := range axes {
		if axis.InputType != "range" {
			continue
		}

		val, ok := params[axis.AttributeCode]
		if !ok {
			return fmt.Errorf("%w: %s", ErrMissingParameter, axis.AttributeCode)
		}

		if axis.MinValue != nil && val < *axis.MinValue {
			return fmt.Errorf("%w: %s must be >= %.2f", ErrParameterOutOfRange, axis.AttributeCode, *axis.MinValue)
		}
		if axis.MaxValue != nil && val > *axis.MaxValue {
			return fmt.Errorf("%w: %s must be <= %.2f", ErrParameterOutOfRange, axis.AttributeCode, *axis.MaxValue)
		}

		if axis.StepValue != nil && *axis.StepValue > 0 {
			minVal := 0.0
			if axis.MinValue != nil {
				minVal = *axis.MinValue
			}
			remainder := math.Mod(val-minVal, *axis.StepValue)
			if remainder > 0.001 && math.Abs(remainder-*axis.StepValue) > 0.001 {
				return fmt.Errorf("%w: %s must be in steps of %.2f", ErrParameterInvalidStep, axis.AttributeCode, *axis.StepValue)
			}
		}
	}
	return nil
}

// CalculatePrice computes the price based on formula type, parameters, and SKU mapping
// If skuMapping is provided, its unit_price and base_price override the pricing defaults.
func (pp *ParametricPricing) CalculatePrice(params map[string]float64, quantity int, skuMapping *SKUMapping) (*ParametricPriceResponse, error) {
	if quantity <= 0 {
		quantity = 1
	}

	breakdown := make(map[string]float64)
	var unitPrice float64

	// Use SKU mapping prices if available, otherwise fall back to pricing table
	basePrice := pp.BasePrice
	unitPriceVal := 0.0
	if pp.UnitPrice != nil {
		unitPriceVal = *pp.UnitPrice
	}
	sku := ""

	if skuMapping != nil {
		basePrice = skuMapping.BasePrice
		unitPriceVal = skuMapping.UnitPrice
		sku = skuMapping.SKU
	}

	switch FormulaType(pp.FormulaType) {
	case FormulaTypeFixed:
		unitPrice = basePrice
		breakdown["base_price"] = basePrice

	case FormulaTypePerUnit:
		unitPrice = basePrice + unitPriceVal
		breakdown["base_price"] = basePrice
		breakdown["unit_price"] = unitPriceVal

	case FormulaTypePerM2:
		lengthMM := params["length_mm"]
		widthMM := params["width_mm"]
		areaM2 := (lengthMM / 1000.0) * (widthMM / 1000.0)

		unitPrice = basePrice + (areaM2 * unitPriceVal)

		breakdown["base_price"] = basePrice
		breakdown["area_m2"] = math.Round(areaM2*10000) / 10000
		breakdown["price_per_m2"] = unitPriceVal

	case FormulaTypePerRunningMeter:
		lengthM := params["length_m"]
		unitPrice = basePrice + (lengthM * unitPriceVal)
		breakdown["base_price"] = basePrice
		breakdown["length_m"] = lengthM
		breakdown["price_per_m"] = unitPriceVal

	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidFormulaType, pp.FormulaType)
	}

	// Round to 2 decimal places
	unitPrice = math.Round(unitPrice*100) / 100
	totalPrice := math.Round(unitPrice*float64(quantity)*100) / 100

	// Apply minimum order value
	if pp.MinOrderValue != nil && totalPrice < *pp.MinOrderValue {
		breakdown["min_order_surcharge"] = *pp.MinOrderValue - totalPrice
		totalPrice = *pp.MinOrderValue
	}

	return &ParametricPriceResponse{
		SKU:        sku,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
		Currency:   pp.Currency,
		Quantity:   quantity,
		Breakdown:  breakdown,
	}, nil
}

// ValidateSelections checks that all select axes have a value provided
func ValidateSelections(axes []VariantAxis, selections map[string]string) error {
	for _, axis := range axes {
		if axis.InputType != "select" {
			continue
		}
		if _, ok := selections[axis.AttributeCode]; !ok {
			return fmt.Errorf("%w: %s", ErrMissingSelection, axis.AttributeCode)
		}
	}
	return nil
}
