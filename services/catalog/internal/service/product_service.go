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
	productRepo   repository.ProductRepository
	attrTransRepo repository.AttributeTranslationRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo repository.ProductRepository, attrTransRepo repository.AttributeTranslationRepository) *ProductService {
	return &ProductService{
		productRepo:   productRepo,
		attrTransRepo: attrTransRepo,
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
	
	// Set product type (default to simple if not specified)
	if req.ProductType != "" {
		product.ProductType = req.ProductType
	}
	
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

	// If variant_parent with axes, create axes
	if product.ProductType == domain.ProductTypeVariantParent && len(req.VariantAxes) > 0 {
		axes := make([]domain.VariantAxis, len(req.VariantAxes))
		for i, axisReq := range req.VariantAxes {
			axes[i] = domain.VariantAxis{
				ID:            uuid.New(),
				ProductID:     product.ID,
				AttributeCode: axisReq.AttributeCode,
				Position:      axisReq.Position,
			}
		}
		if err := s.productRepo.SetVariantAxes(ctx, product.ID, axes); err != nil {
			return nil, err
		}
		product.VariantAxes = axes
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

// GetByIDWithTranslations retrieves a product by ID with translated attributes
func (s *ProductService) GetByIDWithTranslations(ctx context.Context, id uuid.UUID, locale string) (*domain.ProductWithTranslatedAttributes, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.enrichWithTranslations(ctx, product, locale), nil
}

// ListWithTranslations retrieves products with translated attributes
func (s *ProductService) ListWithTranslations(ctx context.Context, filter domain.ProductFilter, locale string) ([]domain.ProductWithTranslatedAttributes, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 200 {
		filter.Limit = 200
	}

	products, total, err := s.productRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Load translations once for all products
	var translations map[string]*domain.AttributeTranslation
	if locale != "" && len(products) > 0 {
		translations, _ = s.attrTransRepo.GetByTenantAndLocale(ctx, filter.TenantID, locale)
	}

	// Enrich products with translations
	result := make([]domain.ProductWithTranslatedAttributes, len(products))
	for i, p := range products {
		product := p // Create a copy
		result[i] = domain.ProductWithTranslatedAttributes{
			Product: product,
		}
		if translations != nil {
			result[i].TranslatedAttributes = s.translateAttributes(product.Attributes, translations)
		}
	}

	return result, total, nil
}

// enrichWithTranslations enriches a single product with translations
func (s *ProductService) enrichWithTranslations(ctx context.Context, product *domain.Product, locale string) *domain.ProductWithTranslatedAttributes {
	result := &domain.ProductWithTranslatedAttributes{
		Product: *product,
	}

	if locale != "" && len(product.Attributes) > 0 {
		translations, err := s.attrTransRepo.GetByTenantAndLocale(ctx, product.TenantID, locale)
		if err == nil {
			result.TranslatedAttributes = s.translateAttributes(product.Attributes, translations)
		}
	}

	return result
}

// translateAttributes translates product attributes using the translation map
func (s *ProductService) translateAttributes(attributes []domain.ProductAttribute, translations map[string]*domain.AttributeTranslation) []domain.ProductAttributeWithTranslation {
	result := make([]domain.ProductAttributeWithTranslation, len(attributes))

	for i, attr := range attributes {
		result[i] = domain.ProductAttributeWithTranslation{
			Key:   attr.Key,
			Type:  attr.Type,
			Value: attr.Value,
		}

		// Apply translation if available
		if trans, ok := translations[attr.Key]; ok {
			result[i].DisplayName = trans.DisplayName
			result[i].Unit = trans.Unit
		}
	}

	return result
}
