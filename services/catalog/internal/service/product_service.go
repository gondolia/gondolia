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
	priceRepo     repository.PriceRepository
	attrTransRepo repository.AttributeTranslationRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo repository.ProductRepository, priceRepo repository.PriceRepository, attrTransRepo repository.AttributeTranslationRepository) *ProductService {
	return &ProductService{
		productRepo:   productRepo,
		priceRepo:     priceRepo,
		attrTransRepo: attrTransRepo,
	}
}

// GetByID retrieves a product by ID
func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product.ProductType == domain.ProductTypeBundle {
		s.enrichBundleListData(ctx, product)
	}
	return product, nil
}

// GetBySKU retrieves a product by SKU
func (s *ProductService) GetBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*domain.Product, error) {
	return s.productRepo.GetBySKU(ctx, tenantID, sku)
}

// List retrieves products with filtering and pagination.
// For variant_parent products, enriches with variant_count, price_range and variant_summary.
func (s *ProductService) List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error) {
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

	// Enrich variant_parent products with variant metadata
	// Enrich bundle products with base price from prices table
	for i := range products {
		if products[i].ProductType == domain.ProductTypeVariantParent {
			s.enrichVariantParentListData(ctx, &products[i])
		} else if products[i].ProductType == domain.ProductTypeBundle {
			s.enrichBundleListData(ctx, &products[i])
		}
	}

	return products, total, nil
}

// enrichVariantParentListData adds variant_count, price_range and variant_summary
// to a variant_parent product for list/card display.
func (s *ProductService) enrichVariantParentListData(ctx context.Context, product *domain.Product) {
	variants, err := s.productRepo.ListVariants(ctx, product.ID, domain.ProductStatusActive)
	if err != nil || len(variants) == 0 {
		return
	}

	// Variant count
	count := len(variants)
	product.VariantCount = &count

	// Price range — collect from prices table
	var minPrice, maxPrice float64
	var currency string
	priceFound := false
	for _, v := range variants {
		prices, err := s.priceRepo.ListByProduct(ctx, v.ID)
		if err != nil || len(prices) == 0 {
			continue
		}
		basePrice := prices[0] // First price (lowest min_quantity)
		if !priceFound || basePrice.Price < minPrice {
			minPrice = basePrice.Price
		}
		if !priceFound || basePrice.Price > maxPrice {
			maxPrice = basePrice.Price
		}
		currency = basePrice.Currency
		priceFound = true
	}
	if priceFound {
		product.VariantPriceRange = &domain.PriceRange{
			Min:      minPrice,
			Max:      maxPrice,
			Currency: currency,
		}
	}

	// Variant summary — collect unique axis value labels per axis
	summary := make(map[string][]string)
	seen := make(map[string]map[string]bool) // axis -> set of codes
	for _, v := range variants {
		axisValues, err := s.productRepo.GetAxisValues(ctx, v.ID)
		if err != nil {
			continue
		}
		for _, av := range axisValues {
			if seen[av.AxisAttributeCode] == nil {
				seen[av.AxisAttributeCode] = make(map[string]bool)
			}
			if !seen[av.AxisAttributeCode][av.OptionCode] {
				seen[av.AxisAttributeCode][av.OptionCode] = true
				label := formatOptionLabelSimple(av.OptionCode)
				summary[av.AxisAttributeCode] = append(summary[av.AxisAttributeCode], label)
			}
		}
	}
	if len(summary) > 0 {
		product.VariantSummary = summary
	}
}

// enrichBundleListData adds base_price to a bundle product for list/card display.
// For fixed mode: uses the stored price directly.
// For configurable mode: calculates minimum price using min_quantity for each component.
func (s *ProductService) enrichBundleListData(ctx context.Context, product *domain.Product) {
	prices, err := s.priceRepo.ListByProduct(ctx, product.ID)
	if err != nil || len(prices) == 0 {
		return
	}
	// Use the first price entry (lowest min_quantity)
	bp := prices[0].Price
	bc := prices[0].Currency
	product.BasePrice = &bp
	product.BaseCurrency = &bc
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

// UpdateStatus updates the status of a product
func (s *ProductService) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProductStatus) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	product.Status = status
	product.UpdatedAt = time.Now()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// AddAttribute adds or updates an attribute on a product
func (s *ProductService) AddAttribute(ctx context.Context, id uuid.UUID, attr domain.ProductAttribute) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if attribute already exists - if so, update it
	found := false
	for i, existing := range product.Attributes {
		if existing.Key == attr.Key {
			product.Attributes[i] = attr
			found = true
			break
		}
	}

	// If not found, append it
	if !found {
		product.Attributes = append(product.Attributes, attr)
	}

	product.UpdatedAt = time.Now()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateAttribute updates an existing attribute on a product
func (s *ProductService) UpdateAttribute(ctx context.Context, id uuid.UUID, key string, attr domain.ProductAttribute) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Find and update the attribute
	found := false
	for i, existing := range product.Attributes {
		if existing.Key == key {
			product.Attributes[i] = attr
			found = true
			break
		}
	}

	if !found {
		return nil, domain.ErrAttributeNotFound
	}

	product.UpdatedAt = time.Now()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteAttribute removes an attribute from a product
func (s *ProductService) DeleteAttribute(ctx context.Context, id uuid.UUID, key string) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Filter out the attribute to delete
	newAttributes := make([]domain.ProductAttribute, 0, len(product.Attributes))
	found := false
	for _, attr := range product.Attributes {
		if attr.Key != key {
			newAttributes = append(newAttributes, attr)
		} else {
			found = true
		}
	}

	if !found {
		return nil, domain.ErrAttributeNotFound
	}

	product.Attributes = newAttributes
	product.UpdatedAt = time.Now()

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
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
