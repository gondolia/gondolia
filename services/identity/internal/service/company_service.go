package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
	"github.com/gondolia/gondolia/services/identity/internal/repository"
)

// CompanyService handles company operations
type CompanyService struct {
	companyRepo     repository.CompanyRepository
	userCompanyRepo repository.UserCompanyRepository
	roleRepo        repository.RoleRepository
	userRepo        repository.UserRepository
}

// NewCompanyService creates a new company service
func NewCompanyService(
	companyRepo repository.CompanyRepository,
	userCompanyRepo repository.UserCompanyRepository,
	roleRepo repository.RoleRepository,
	userRepo repository.UserRepository,
) *CompanyService {
	return &CompanyService{
		companyRepo:     companyRepo,
		userCompanyRepo: userCompanyRepo,
		roleRepo:        roleRepo,
		userRepo:        userRepo,
	}
}

// List returns a list of companies
func (s *CompanyService) List(ctx context.Context, filter domain.CompanyFilter) ([]domain.Company, int, error) {
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return s.companyRepo.List(ctx, filter)
}

// GetByID returns a company by ID
func (s *CompanyService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	return s.companyRepo.GetByID(ctx, id)
}

// GetBySAPNumber returns a company by SAP number
func (s *CompanyService) GetBySAPNumber(ctx context.Context, tenantID uuid.UUID, sapNumber string) (*domain.Company, error) {
	return s.companyRepo.GetBySAPNumber(ctx, tenantID, sapNumber)
}

// Create creates a new company
func (s *CompanyService) Create(ctx context.Context, tenantID uuid.UUID, req domain.CreateCompanyRequest) (*domain.Company, error) {
	// Check if SAP number already exists
	_, err := s.companyRepo.GetBySAPNumber(ctx, tenantID, req.SAPCompanyNumber)
	if err == nil {
		return nil, domain.ErrCompanyAlreadyExists
	}

	// Create company
	company := domain.NewCompany(tenantID, req.SAPCompanyNumber, req.Name)
	company.Email = req.Email
	company.Street = req.Street
	company.HouseNumber = req.HouseNumber
	company.ZIP = req.ZIP
	company.City = req.City
	company.SAPCustomerGroup = req.SAPCustomerGroup
	company.SAPShippingPlant = req.SAPShippingPlant
	company.SAPOffice = req.SAPOffice

	if req.Currency != nil {
		company.Currency = *req.Currency
	}
	if req.Country != nil {
		company.Country = *req.Country
	}

	if err := s.companyRepo.Create(ctx, company); err != nil {
		return nil, err
	}

	return company, nil
}

// Update updates a company
func (s *CompanyService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateCompanyRequest) (*domain.Company, error) {
	company, err := s.companyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.Description != nil {
		company.Description = req.Description
	}
	if req.Email != nil {
		company.Email = req.Email
	}
	if req.Street != nil {
		company.Street = req.Street
	}
	if req.HouseNumber != nil {
		company.HouseNumber = req.HouseNumber
	}
	if req.ZIP != nil {
		company.ZIP = req.ZIP
	}
	if req.City != nil {
		company.City = req.City
	}
	if req.Country != nil {
		company.Country = *req.Country
	}
	if req.Phone != nil {
		company.Phone = req.Phone
	}
	if req.Fax != nil {
		company.Fax = req.Fax
	}
	if req.URL != nil {
		company.URL = req.URL
	}
	if req.DesiredDeliveryDays != nil {
		company.DesiredDeliveryDays = req.DesiredDeliveryDays
	}
	if req.DefaultShippingNote != nil {
		company.DefaultShippingNote = req.DefaultShippingNote
	}
	if req.DisableOrderFeature != nil {
		company.DisableOrderFeature = *req.DisableOrderFeature
	}
	if req.CustomPrimaryColor != nil {
		company.CustomPrimaryColor = req.CustomPrimaryColor
	}
	if req.CustomSecondaryColor != nil {
		company.CustomSecondaryColor = req.CustomSecondaryColor
	}
	if req.IsActive != nil {
		company.IsActive = *req.IsActive
	}

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return nil, err
	}

	return company, nil
}

// Delete soft deletes a company
func (s *CompanyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.companyRepo.Delete(ctx, id)
}

// ListUsers returns users assigned to a company
func (s *CompanyService) ListUsers(ctx context.Context, companyID uuid.UUID) ([]domain.UserCompany, error) {
	return s.userCompanyRepo.ListByCompany(ctx, companyID)
}

// AddUser adds a user to a company
func (s *CompanyService) AddUser(ctx context.Context, companyID, userID, roleID uuid.UUID, userType domain.UserType) error {
	// Check user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check company exists
	_, err = s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return err
	}

	// Check role exists
	_, err = s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	// Check not already in company
	_, err = s.userCompanyRepo.GetByUserAndCompany(ctx, userID, companyID)
	if err == nil {
		return domain.ErrUserAlreadyInCompany
	}

	// Create relationship
	uc := domain.NewUserCompany(userID, companyID, &roleID, userType)
	return s.userCompanyRepo.Create(ctx, uc)
}

// UpdateUserRole updates a user's role in a company
func (s *CompanyService) UpdateUserRole(ctx context.Context, companyID, userID, roleID uuid.UUID) error {
	uc, err := s.userCompanyRepo.GetByUserAndCompany(ctx, userID, companyID)
	if err != nil {
		return err
	}

	// Check role exists
	_, err = s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	uc.RoleID = &roleID
	return s.userCompanyRepo.Update(ctx, uc)
}

// RemoveUser removes a user from a company
func (s *CompanyService) RemoveUser(ctx context.Context, companyID, userID uuid.UUID) error {
	// Remove from company
	if err := s.userCompanyRepo.Delete(ctx, userID, companyID); err != nil {
		return err
	}

	// Check if user has other companies
	count, err := s.userCompanyRepo.CountByUser(ctx, userID)
	if err != nil {
		return err
	}

	// If no other companies, deactivate user (unless salesmaster)
	if count == 0 {
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return nil // User might not exist, that's ok
		}

		if !user.IsSalesMaster {
			user.IsActive = false
			_ = s.userRepo.Update(ctx, user)
		}
	}

	return nil
}
