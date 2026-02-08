package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// UserType represents the type of user in a company
type UserType int

const (
	UserTypeAdmin    UserType = 0 // Company Admin
	UserTypeSalesRep UserType = 1 // Sales Representative
	UserTypeInvited  UserType = 2 // Invited User
)

func (t UserType) String() string {
	switch t {
	case UserTypeAdmin:
		return "admin"
	case UserTypeSalesRep:
		return "salesrep"
	case UserTypeInvited:
		return "invited"
	default:
		return "unknown"
	}
}

// User represents a user in the system (formerly customer)
type User struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`

	// Status
	IsActive      bool `json:"is_active"`
	IsImported    bool `json:"is_imported"`
	IsSalesMaster bool `json:"is_salesmaster"`
	SSOOnly       bool `json:"sso_only"`

	// SAP Mapping
	SAPUserID         *string `json:"sap_user_id,omitempty"`
	SAPCustomerNumber *string `json:"sap_customer_number,omitempty"`

	// Profile
	Email           string  `json:"email"`
	PasswordHash    *string `json:"-"` // Never expose
	FirstName       string  `json:"firstname"`
	LastName        string  `json:"lastname"`
	Phone           *string `json:"phone,omitempty"`
	Mobile          *string `json:"mobile,omitempty"`
	DefaultLanguage string  `json:"default_language"`

	// Company Context
	DefaultCompanyID *uuid.UUID `json:"default_company_id,omitempty"`

	// Invitation
	InvitationToken *string    `json:"-"` // Never expose
	InvitedAt       *time.Time `json:"invited_at,omitempty"`

	// Tracking
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	// Loaded relations (not always populated)
	Companies []UserCompany `json:"companies,omitempty"`
}

// Name returns the full name
func (u *User) Name() string {
	return u.FirstName + " " + u.LastName
}

// HasOpenInvitation returns true if user has pending invitation
func (u *User) HasOpenInvitation() bool {
	return u.PasswordHash == nil && u.InvitationToken != nil
}

// HasExpiredInvitation returns true if invitation is older than 30 days
func (u *User) HasExpiredInvitation() bool {
	if !u.HasOpenInvitation() || u.InvitedAt == nil {
		return false
	}
	return time.Since(*u.InvitedAt) > 30*24*time.Hour
}

// CleanSAPCustomerNumber removes leading zeros from SAP customer number
func (u *User) CleanSAPCustomerNumber() string {
	if u.SAPCustomerNumber == nil {
		return ""
	}
	return strings.TrimLeft(*u.SAPCustomerNumber, "0")
}

// CanLogin checks if user is allowed to login with password
func (u *User) CanLogin() error {
	if !u.IsActive {
		return ErrUserNotActive
	}
	if u.SSOOnly {
		return ErrSSORequired
	}
	if u.HasOpenInvitation() {
		return ErrInvitationPending
	}
	if u.DeletedAt != nil {
		return ErrUserDeleted
	}
	return nil
}

// NewUser creates a new user with defaults
func NewUser(tenantID uuid.UUID, email, firstName, lastName string) *User {
	now := time.Now()
	return &User{
		ID:              uuid.New(),
		TenantID:        tenantID,
		Email:           strings.ToLower(strings.TrimSpace(email)),
		FirstName:       strings.TrimSpace(firstName),
		LastName:        strings.TrimSpace(lastName),
		DefaultLanguage: "de",
		IsActive:        false,
		IsImported:      false,
		IsSalesMaster:   false,
		SSOOnly:         false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	FirstName string  `json:"firstname" binding:"required,min=1,max=100"`
	LastName  string  `json:"lastname" binding:"required,min=1,max=100"`
	Phone     *string `json:"phone,omitempty"`
	Mobile    *string `json:"mobile,omitempty"`
	Password  *string `json:"password,omitempty"` // If not provided, invitation flow
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	FirstName       *string `json:"firstname,omitempty"`
	LastName        *string `json:"lastname,omitempty"`
	Phone           *string `json:"phone,omitempty"`
	Mobile          *string `json:"mobile,omitempty"`
	DefaultLanguage *string `json:"default_language,omitempty"`
	IsActive        *bool   `json:"is_active,omitempty"`
}

// UserFilter represents filter options for listing users
type UserFilter struct {
	TenantID  uuid.UUID
	CompanyID *uuid.UUID
	Email     *string
	IsActive  *bool
	Search    *string // Searches in name, email
	Limit     int
	Offset    int
}
