package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
	"github.com/gondolia/gondolia/services/identity/internal/repository"
)

// RoleService handles role operations
type RoleService struct {
	roleRepo repository.RoleRepository
}

// NewRoleService creates a new role service
func NewRoleService(roleRepo repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

// List returns a list of roles
func (s *RoleService) List(ctx context.Context, filter domain.RoleFilter) ([]domain.Role, int, error) {
	if filter.Limit == 0 {
		filter.Limit = 50
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return s.roleRepo.List(ctx, filter)
}

// GetByID returns a role by ID
func (s *RoleService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	return s.roleRepo.GetByID(ctx, id)
}

// Create creates a new role
func (s *RoleService) Create(ctx context.Context, tenantID uuid.UUID, req domain.CreateRoleRequest) (*domain.Role, error) {
	// Check if name already exists
	_, err := s.roleRepo.GetByName(ctx, tenantID, req.CompanyID, req.Name)
	if err == nil {
		return nil, domain.ErrRoleAlreadyExists
	}

	// Create role
	role := domain.NewRole(tenantID, req.Name)
	role.CompanyID = req.CompanyID
	role.Permissions = req.Permissions

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// Update updates a role
func (s *RoleService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateRoleRequest) (*domain.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cannot modify system roles
	if role.IsSystem {
		return nil, domain.ErrRoleIsSystem
	}

	if req.Name != nil {
		// Check if new name already exists
		existing, err := s.roleRepo.GetByName(ctx, role.TenantID, role.CompanyID, *req.Name)
		if err == nil && existing.ID != role.ID {
			return nil, domain.ErrRoleAlreadyExists
		}
		role.Name = *req.Name
	}

	if req.Permissions != nil {
		role.Permissions = req.Permissions
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// Delete deletes a role
func (s *RoleService) Delete(ctx context.Context, id uuid.UUID) error {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Cannot delete system roles
	if role.IsSystem {
		return domain.ErrRoleIsSystem
	}

	return s.roleRepo.Delete(ctx, id)
}

// CreateSystemRoles creates default system roles for a tenant
func (s *RoleService) CreateSystemRoles(ctx context.Context, tenantID uuid.UUID) error {
	// Admin role
	adminRole := domain.NewSystemRole(tenantID, "Administrator", domain.DefaultAdminPermissions())
	if err := s.roleRepo.Create(ctx, adminRole); err != nil {
		return err
	}

	// Standard user role
	userRole := domain.NewSystemRole(tenantID, "Benutzer", domain.DefaultUserPermissions())
	if err := s.roleRepo.Create(ctx, userRole); err != nil {
		return err
	}

	// Read-only role
	readOnlyPerms := map[domain.Permission]bool{
		domain.PermSeeOrders:   true,
		domain.PermSeeInvoices: true,
	}
	readOnlyRole := domain.NewSystemRole(tenantID, "Nur Lesen", readOnlyPerms)
	if err := s.roleRepo.Create(ctx, readOnlyRole); err != nil {
		return err
	}

	return nil
}
