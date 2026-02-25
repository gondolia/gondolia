package domain

import "errors"

// Domain errors
var (
	// Order errors
	ErrOrderNotFound           = errors.New("order not found")
	ErrOrderAlreadyExists      = errors.New("order with this number already exists")
	ErrOrderInvalidStatus      = errors.New("invalid order status")
	ErrOrderCannotBeCancelled  = errors.New("order cannot be cancelled in current status")
	ErrOrderInvalidTransition  = errors.New("invalid order status transition")

	// Cart errors
	ErrCartNotFound            = errors.New("cart not found")
	ErrCartEmpty               = errors.New("cart is empty")
	ErrCartValidationFailed    = errors.New("cart validation failed")

	// Tenant errors
	ErrTenantNotFound  = errors.New("tenant not found")
	ErrTenantNotActive = errors.New("tenant is not active")

	// General errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrOrderNotFound) ||
		errors.Is(err, ErrCartNotFound) ||
		errors.Is(err, ErrTenantNotFound)
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrOrderAlreadyExists) ||
		errors.Is(err, ErrOrderInvalidStatus) ||
		errors.Is(err, ErrOrderCannotBeCancelled) ||
		errors.Is(err, ErrOrderInvalidTransition) ||
		errors.Is(err, ErrCartEmpty) ||
		errors.Is(err, ErrCartValidationFailed)
}
