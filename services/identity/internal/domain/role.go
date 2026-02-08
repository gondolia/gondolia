package domain

import (
	"time"

	"github.com/google/uuid"
)

// Permission represents a permission key
type Permission string

// Permission constants
const (
	// Company Management
	PermManageCompany       Permission = "company.manage"
	PermManageUsersAndRoles Permission = "company.manage-users-and-roles-and-permissions"
	PermManageSettings      Permission = "company.manage-settings"
	PermManageAddresses     Permission = "company.manage-addresses"
	PermManageCustomSKUs    Permission = "company.manage-custom-skus"
	PermManageWatchlists    Permission = "company.manage-watchlists"

	// Order Data
	PermSeeOrders      Permission = "company.order-data.see-orders"
	PermSeeInvoices    Permission = "company.order-data.see-invoices"
	PermSeeShipments   Permission = "company.order-data.see-shipments"
	PermSeeReshipments Permission = "company.order-data.see-reshipments"
	PermSeeCredits     Permission = "company.order-data.see-credits"

	// Sales
	PermCreateOrder Permission = "sales.create-order"
)

// AllPermissions returns all available permissions
func AllPermissions() []Permission {
	return []Permission{
		PermManageCompany,
		PermManageUsersAndRoles,
		PermManageSettings,
		PermManageAddresses,
		PermManageCustomSKUs,
		PermManageWatchlists,
		PermSeeOrders,
		PermSeeInvoices,
		PermSeeShipments,
		PermSeeReshipments,
		PermSeeCredits,
		PermCreateOrder,
	}
}

// Role represents a role with permissions
type Role struct {
	ID        uuid.UUID            `json:"id"`
	TenantID  uuid.UUID            `json:"tenant_id"`
	CompanyID *uuid.UUID           `json:"company_id,omitempty"` // NULL = System Role

	Name        string              `json:"name"`
	Permissions map[Permission]bool `json:"permissions"`
	IsSystem    bool                `json:"is_system"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HasPermission checks if role has a specific permission
func (r *Role) HasPermission(perm Permission) bool {
	if r.Permissions == nil {
		return false
	}
	return r.Permissions[perm]
}

// PermissionStrings returns permissions as string slice (for JWT)
func (r *Role) PermissionStrings() []string {
	var perms []string
	for perm, granted := range r.Permissions {
		if granted {
			perms = append(perms, string(perm))
		}
	}
	return perms
}

// NewRole creates a new role with defaults
func NewRole(tenantID uuid.UUID, name string) *Role {
	now := time.Now()
	return &Role{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        name,
		Permissions: make(map[Permission]bool),
		IsSystem:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewSystemRole creates a new system role
func NewSystemRole(tenantID uuid.UUID, name string, permissions map[Permission]bool) *Role {
	role := NewRole(tenantID, name)
	role.IsSystem = true
	role.Permissions = permissions
	return role
}

// DefaultAdminPermissions returns permissions for admin role
func DefaultAdminPermissions() map[Permission]bool {
	return map[Permission]bool{
		PermManageCompany:       true,
		PermManageUsersAndRoles: true,
		PermManageSettings:      true,
		PermManageAddresses:     true,
		PermManageCustomSKUs:    true,
		PermManageWatchlists:    true,
		PermSeeOrders:           true,
		PermSeeInvoices:         true,
		PermSeeShipments:        true,
		PermSeeReshipments:      true,
		PermSeeCredits:          true,
		PermCreateOrder:         true,
	}
}

// DefaultUserPermissions returns permissions for standard user role
func DefaultUserPermissions() map[Permission]bool {
	return map[Permission]bool{
		PermManageWatchlists: true,
		PermSeeOrders:        true,
		PermSeeInvoices:      true,
		PermSeeShipments:     true,
		PermCreateOrder:      true,
	}
}

// CreateRoleRequest represents a request to create a role
type CreateRoleRequest struct {
	Name        string              `json:"name" binding:"required,min=1,max=100"`
	CompanyID   *uuid.UUID          `json:"company_id,omitempty"`
	Permissions map[Permission]bool `json:"permissions" binding:"required"`
}

// UpdateRoleRequest represents a request to update a role
type UpdateRoleRequest struct {
	Name        *string             `json:"name,omitempty"`
	Permissions map[Permission]bool `json:"permissions,omitempty"`
}

// RoleFilter represents filter options for listing roles
type RoleFilter struct {
	TenantID  uuid.UUID
	CompanyID *uuid.UUID // NULL = only system roles
	IsSystem  *bool
	Limit     int
	Offset    int
}
