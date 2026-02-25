package domain

import "errors"

// Common errors
var (
	ErrCartNotFound       = errors.New("cart not found")
	ErrCartItemNotFound   = errors.New("cart item not found")
	ErrCartNotActive      = errors.New("cart is not active")
	ErrInvalidQuantity    = errors.New("invalid quantity")
	ErrProductNotFound    = errors.New("product not found")
	ErrPriceNotAvailable  = errors.New("price not available")
	ErrInvalidCartOwner   = errors.New("invalid cart owner")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrValidationFailed   = errors.New("validation failed")
)

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrCartNotFound) ||
		errors.Is(err, ErrCartItemNotFound) ||
		errors.Is(err, ErrProductNotFound)
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrInvalidQuantity) ||
		errors.Is(err, ErrValidationFailed)
}
