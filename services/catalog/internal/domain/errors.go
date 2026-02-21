package domain

import "errors"

// Domain errors
var (
	// Product errors
	ErrProductNotFound      = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product with this SKU already exists")
	ErrProductInvalidStatus = errors.New("invalid product status")

	// Category errors
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("category with this code already exists")
	ErrCategoryCircularRef   = errors.New("circular category reference detected")
	ErrCategoryHasProducts   = errors.New("category has products assigned")

	// Price errors
	ErrPriceNotFound      = errors.New("price not found")
	ErrPriceInvalidRange  = errors.New("invalid price date range")
	ErrPriceOverlap       = errors.New("price range overlaps with existing price")

	// Tenant errors
	ErrTenantNotFound  = errors.New("tenant not found")
	ErrTenantNotActive = errors.New("tenant is not active")

	// Variant errors
	ErrTooManyAxes                = errors.New("variant parent cannot have more than 4 axes")
	ErrDuplicateVariantCombination = errors.New("variant with these axis values already exists")

	// General errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrProductNotFound) ||
		errors.Is(err, ErrCategoryNotFound) ||
		errors.Is(err, ErrPriceNotFound) ||
		errors.Is(err, ErrTenantNotFound)
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrProductAlreadyExists) ||
		errors.Is(err, ErrCategoryAlreadyExists) ||
		errors.Is(err, ErrCategoryCircularRef) ||
		errors.Is(err, ErrPriceInvalidRange) ||
		errors.Is(err, ErrPriceOverlap) ||
		errors.Is(err, ErrProductInvalidStatus) ||
		errors.Is(err, ErrTooManyAxes) ||
		errors.Is(err, ErrDuplicateVariantCombination)
}
