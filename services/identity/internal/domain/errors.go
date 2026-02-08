package domain

import "errors"

// Domain errors
var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserNotActive     = errors.New("user is not active")
	ErrUserDeleted       = errors.New("user has been deleted")
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrSSORequired       = errors.New("SSO login required for this user")
	ErrInvitationPending = errors.New("invitation has not been accepted yet")
	ErrInvitationExpired = errors.New("invitation has expired")
	ErrInvitationInvalid = errors.New("invalid invitation token")

	// Company errors
	ErrCompanyNotFound      = errors.New("company not found")
	ErrCompanyNotActive     = errors.New("company is not active")
	ErrCompanyAlreadyExists = errors.New("company with this SAP number already exists")

	// Role errors
	ErrRoleNotFound       = errors.New("role not found")
	ErrRoleIsSystem       = errors.New("system roles cannot be modified")
	ErrRoleAlreadyExists  = errors.New("role with this name already exists")

	// UserCompany errors
	ErrUserNotInCompany     = errors.New("user is not assigned to this company")
	ErrUserAlreadyInCompany = errors.New("user is already assigned to this company")
	ErrNoCompanyAssigned    = errors.New("user has no company assigned")

	// Auth errors
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrTokenExpired         = errors.New("token has expired")
	ErrTokenInvalid         = errors.New("invalid token")
	ErrTokenRevoked         = errors.New("token has been revoked")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")

	// Password reset errors
	ErrPasswordResetExpired  = errors.New("password reset link has expired")
	ErrPasswordResetUsed     = errors.New("password reset link has already been used")
	ErrPasswordResetNotFound = errors.New("password reset token not found")
	ErrPasswordTooWeak       = errors.New("password does not meet requirements")

	// Tenant errors
	ErrTenantNotFound = errors.New("tenant not found")
	ErrTenantNotActive = errors.New("tenant is not active")

	// General errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrCompanyNotFound) ||
		errors.Is(err, ErrRoleNotFound) ||
		errors.Is(err, ErrTenantNotFound) ||
		errors.Is(err, ErrRefreshTokenNotFound) ||
		errors.Is(err, ErrPasswordResetNotFound)
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrUserAlreadyExists) ||
		errors.Is(err, ErrCompanyAlreadyExists) ||
		errors.Is(err, ErrRoleAlreadyExists) ||
		errors.Is(err, ErrUserAlreadyInCompany) ||
		errors.Is(err, ErrPasswordTooWeak)
}

// IsAuthError checks if error is an authentication error
func IsAuthError(err error) bool {
	return errors.Is(err, ErrInvalidCredentials) ||
		errors.Is(err, ErrTokenExpired) ||
		errors.Is(err, ErrTokenInvalid) ||
		errors.Is(err, ErrTokenRevoked) ||
		errors.Is(err, ErrUserNotActive) ||
		errors.Is(err, ErrSSORequired) ||
		errors.Is(err, ErrInvitationPending)
}
