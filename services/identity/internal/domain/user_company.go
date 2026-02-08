package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserCompany represents the relationship between a user and a company
type UserCompany struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CompanyID uuid.UUID `json:"company_id"`
	RoleID    *uuid.UUID `json:"role_id,omitempty"`
	UserType  UserType  `json:"user_type"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Loaded relations (not always populated)
	Company *Company `json:"company,omitempty"`
	Role    *Role    `json:"role,omitempty"`
	User    *User    `json:"user,omitempty"`
}

// NewUserCompany creates a new user-company relationship
func NewUserCompany(userID, companyID uuid.UUID, roleID *uuid.UUID, userType UserType) *UserCompany {
	now := time.Now()
	return &UserCompany{
		ID:        uuid.New(),
		UserID:    userID,
		CompanyID: companyID,
		RoleID:    roleID,
		UserType:  userType,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddUserToCompanyRequest represents a request to add a user to a company
type AddUserToCompanyRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	RoleID   uuid.UUID `json:"role_id" binding:"required"`
	UserType UserType  `json:"user_type"`
}

// InviteUserRequest represents a request to invite a new user to a company
type InviteUserRequest struct {
	Email     string    `json:"email" binding:"required,email"`
	FirstName string    `json:"firstname" binding:"required,min=1,max=100"`
	LastName  string    `json:"lastname" binding:"required,min=1,max=100"`
	Phone     *string   `json:"phone,omitempty"`
	Mobile    *string   `json:"mobile,omitempty"`
	CompanyID uuid.UUID `json:"company_id" binding:"required"`
	RoleID    uuid.UUID `json:"role_id" binding:"required"`
}

// UpdateUserCompanyRequest represents a request to update a user's role in a company
type UpdateUserCompanyRequest struct {
	RoleID   *uuid.UUID `json:"role_id,omitempty"`
	UserType *UserType  `json:"user_type,omitempty"`
}
