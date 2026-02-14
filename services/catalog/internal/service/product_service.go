package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/repository"
)

// ProductService handles product business logic
type ProductService struct {
	productRepo repository.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// GetByID retrieves a product by ID
func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

// GetBySKU retrieves a product by SKU
func (s *ProductService) GetBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*domain.Product, error) {
	return s.productRepo.GetBySKU(ctx, tenantID, sku)
}

// List retrieves products with filtering and pagination
func (s *ProductService) List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 200 {
		filter.Limit = 200
	}
	return s.productRepo.List(ctx, filter)
}

// Create creates a new product
func (s *ProductService) Create(ctx context.Context, tenantID uuid.UUID, req domain.CreateProductRequest) (*domain.Product, error) {
	// Check if product with SKU already exists
	existing, err := s.productRepo.GetBySKU(ctx, tenantID, req.SKU)
	if err == nil && existing != nil {
		return nil, domain.ErrProductAlreadyExists
	}

	// Create product
	product := domain.NewProduct(tenantID, req.SKU)
	product.Name = req.Name
	product.Description = req.Description
	product.CategoryIDs = req.CategoryIDs
	product.Attributes = req.Attributes
	product.Images = req.Images

	if req.Status != "" {
		product.Status = req.Status
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// Update updates a product
func (s *ProductService) Update(ctx context.Context, id uuid.UUID, req domain.UpdateProductRequest) (*domain.Product, error) {
	// Get existing product
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		product.Name = req.Name
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.CategoryIDs != nil {
		product.CategoryIDs = req.CategoryIDs
	}
	if req.Attributes != nil {
		product.Attributes = req.Attributes
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	if req.Images != nil {
		product.Images = req.Images
	}

	product.UpdatedAt = time.Now()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// Delete soft-deletes a product
func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.productRepo.Delete(ctx, id)
}
