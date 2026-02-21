package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/repository"
)

// ParametricService handles parametric product business logic
type ParametricService struct {
	productRepo    repository.ProductRepository
	pricingRepo    repository.ParametricPricingRepository
	axisOptionRepo repository.AxisOptionRepository
	skuMappingRepo repository.SKUMappingRepository
}

// NewParametricService creates a new parametric service
func NewParametricService(
	productRepo repository.ProductRepository,
	pricingRepo repository.ParametricPricingRepository,
	axisOptionRepo repository.AxisOptionRepository,
	skuMappingRepo repository.SKUMappingRepository,
) *ParametricService {
	return &ParametricService{
		productRepo:    productRepo,
		pricingRepo:    pricingRepo,
		axisOptionRepo: axisOptionRepo,
		skuMappingRepo: skuMappingRepo,
	}
}

// CalculatePrice computes the price for a parametric product given user parameters
func (s *ParametricService) CalculatePrice(ctx context.Context, productID uuid.UUID, req domain.ParametricPriceRequest) (*domain.ParametricPriceResponse, error) {
	// Load product
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product.ProductType != domain.ProductTypeParametric {
		return nil, domain.ErrProductNotFound
	}

	// Load axes
	axes, err := s.productRepo.GetVariantAxes(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Validate selections (all select axes must have a value)
	if err := domain.ValidateSelections(axes, req.Selections); err != nil {
		return nil, err
	}

	// Validate range parameters
	if err := domain.ValidateParameters(axes, req.Parameters); err != nil {
		return nil, err
	}

	// Load pricing formula
	pricing, err := s.pricingRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Try to resolve SKU mapping from selections
	var skuMapping *domain.SKUMapping
	if len(req.Selections) > 0 {
		m, err := s.skuMappingRepo.FindBySelections(ctx, productID, req.Selections)
		if err != nil && err != domain.ErrSKUMappingNotFound {
			return nil, err
		}
		skuMapping = m
	}

	// Calculate
	qty := req.Quantity
	if qty <= 0 {
		qty = 1
	}

	return pricing.CalculatePrice(req.Parameters, qty, skuMapping)
}

// GetPricing returns the parametric pricing config for a product
func (s *ParametricService) GetPricing(ctx context.Context, productID uuid.UUID) (*domain.ParametricPricing, error) {
	return s.pricingRepo.GetByProductID(ctx, productID)
}

// GetAxisOptions returns the select options for all axes of a parametric product
func (s *ParametricService) GetAxisOptions(ctx context.Context, productID uuid.UUID) (map[uuid.UUID][]domain.AxisOption, error) {
	return s.axisOptionRepo.ListByProductID(ctx, productID)
}

// GetSKUMappings returns all SKU mappings for a parametric product
func (s *ParametricService) GetSKUMappings(ctx context.Context, productID uuid.UUID) ([]domain.SKUMapping, error) {
	return s.skuMappingRepo.ListByProductID(ctx, productID)
}
