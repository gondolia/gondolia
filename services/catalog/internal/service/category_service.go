package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/repository"
)

// CategoryService handles category business logic
type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// GetByID retrieves a category by ID
func (s *CategoryService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

// GetByIDWithAncestors retrieves a category by ID including its ancestors (for breadcrumbs)
func (s *CategoryService) GetByIDWithAncestors(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	return s.categoryRepo.GetByIDWithAncestors(ctx, id)
}

// GetByCode retrieves a category by code
func (s *CategoryService) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Category, error) {
	return s.categoryRepo.GetByCode(ctx, tenantID, code)
}

// GetTree retrieves the complete category tree
func (s *CategoryService) GetTree(ctx context.Context, tenantID uuid.UUID) ([]domain.Category, error) {
	return s.categoryRepo.GetTree(ctx, tenantID)
}

// List retrieves categories with filtering and pagination
func (s *CategoryService) List(ctx context.Context, filter domain.CategoryFilter) ([]domain.Category, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 100
	}
	if filter.Limit > 500 {
		filter.Limit = 500
	}
	return s.categoryRepo.List(ctx, filter)
}

// Create creates a new category
func (s *CategoryService) Create(ctx context.Context, tenantID uuid.UUID, req domain.CreateCategoryRequest) (*domain.Category, error) {
	// Check if category with code already exists
	existing, err := s.categoryRepo.GetByCode(ctx, tenantID, req.Code)
	if err == nil && existing != nil {
		return nil, domain.ErrCategoryAlreadyExists
	}

	// Validate parent if specified
	if req.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, err
		}
		// Ensure parent is in same tenant
		if parent.TenantID != tenantID {
			return nil, domain.ErrCategoryNotFound
		}
	}

	// Create category
	category := domain.NewCategory(tenantID, req.Code)
	category.ParentID = req.ParentID
	category.Name = req.Name
	
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.Image != nil {
		category.Image = req.Image
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	if req.Active != nil {
		category.Active = *req.Active
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// Update updates a category
func (s *CategoryService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateCategoryRequest) (*domain.Category, error) {
	// Get existing category
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate parent if being changed
	if req.ParentID != nil {
		// Prevent circular reference
		if *req.ParentID == id {
			return nil, domain.ErrCategoryCircularRef
		}
		
		parent, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, err
		}
		
		// Ensure parent is in same tenant
		if parent.TenantID != category.TenantID {
			return nil, domain.ErrCategoryNotFound
		}
		
		category.ParentID = req.ParentID
	}

	// Update fields
	if req.Name != nil {
		category.Name = req.Name
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.Image != nil {
		category.Image = req.Image
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	if req.Active != nil {
		category.Active = *req.Active
	}

	category.UpdatedAt = time.Now()

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete soft-deletes a category
func (s *CategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if category has products assigned
	hasProducts, err := s.categoryRepo.HasProducts(ctx, id)
	if err != nil {
		return err
	}
	if hasProducts {
		return domain.ErrCategoryHasProducts
	}

	return s.categoryRepo.Delete(ctx, id)
}
