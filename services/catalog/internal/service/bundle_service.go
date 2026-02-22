package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/repository"
)

// BundleService handles bundle product business logic
type BundleService struct {
	bundleRepo     repository.BundleRepository
	productRepo    repository.ProductRepository
	priceRepo      repository.PriceRepository
	parametricSvc  *ParametricService
}

// NewBundleService creates a new bundle service
func NewBundleService(
	bundleRepo repository.BundleRepository,
	productRepo repository.ProductRepository,
	priceRepo repository.PriceRepository,
	parametricSvc *ParametricService,
) *BundleService {
	return &BundleService{
		bundleRepo:    bundleRepo,
		productRepo:   productRepo,
		priceRepo:     priceRepo,
		parametricSvc: parametricSvc,
	}
}

// GetComponents retrieves all components for a bundle
func (s *BundleService) GetComponents(ctx context.Context, bundleProductID uuid.UUID) ([]domain.BundleComponent, error) {
	// Verify bundle exists
	bundle, err := s.productRepo.GetByID(ctx, bundleProductID)
	if err != nil {
		return nil, err
	}
	if bundle.ProductType != domain.ProductTypeBundle {
		return nil, domain.ErrBundleNotFound
	}

	// Get components
	components, err := s.bundleRepo.GetComponents(ctx, bundleProductID)
	if err != nil {
		return nil, err
	}

	// Load product details for each component
	for i := range components {
		product, err := s.productRepo.GetByID(ctx, components[i].ComponentProductID)
		if err != nil {
			return nil, err
		}

		// For parametric products, load variant axes and pricing config
		if product.ProductType == domain.ProductTypeParametric {
			axes, err := s.productRepo.GetVariantAxes(ctx, product.ID)
			if err == nil {
				product.VariantAxes = axes
			}
			if s.parametricSvc != nil {
				// Load axis options for select-type axes
				axisOptions, err := s.parametricSvc.GetAxisOptions(ctx, product.ID)
				if err == nil && axisOptions != nil {
					for i := range product.VariantAxes {
						if opts, ok := axisOptions[product.VariantAxes[i].ID]; ok {
							product.VariantAxes[i].Options = opts
						}
					}
				}
				pricing, err := s.parametricSvc.GetPricing(ctx, product.ID)
				if err == nil {
					product.ParametricPricing = pricing
				}
			}
		}

		components[i].Product = product
	}

	return components, nil
}

// SetComponents sets/updates the components for a bundle
func (s *BundleService) SetComponents(ctx context.Context, bundleProductID uuid.UUID, tenantID uuid.UUID, requests []domain.BundleComponentRequest) error {
	// Verify bundle exists
	bundle, err := s.productRepo.GetByID(ctx, bundleProductID)
	if err != nil {
		return err
	}
	if bundle.ProductType != domain.ProductTypeBundle {
		return domain.ErrBundleNotFound
	}

	// Validate and convert requests to components
	components := make([]domain.BundleComponent, len(requests))
	for i, req := range requests {
		// Load component product
		product, err := s.productRepo.GetByID(ctx, req.ComponentProductID)
		if err != nil {
			return fmt.Errorf("component product %s: %w", req.ComponentProductID, err)
		}

		// Validate component type
		if err := domain.ValidateComponent(product); err != nil {
			return fmt.Errorf("component %s (%s): %w", product.SKU, product.ProductType, err)
		}

		// Build component
		component := domain.BundleComponent{
			TenantID:           tenantID,
			BundleProductID:    bundleProductID,
			ComponentProductID: req.ComponentProductID,
			Quantity:           req.Quantity,
			MinQuantity:        req.MinQuantity,
			MaxQuantity:        req.MaxQuantity,
			SortOrder:          req.SortOrder,
			DefaultParameters:  req.DefaultParameters,
		}

		// Validate quantity constraints
		if component.MinQuantity != nil && component.Quantity < *component.MinQuantity {
			return fmt.Errorf("component %s: default quantity %d is less than min %d", product.SKU, component.Quantity, *component.MinQuantity)
		}
		if component.MaxQuantity != nil && component.Quantity > *component.MaxQuantity {
			return fmt.Errorf("component %s: default quantity %d exceeds max %d", product.SKU, component.Quantity, *component.MaxQuantity)
		}
		if component.MinQuantity != nil && component.MaxQuantity != nil && *component.MinQuantity > *component.MaxQuantity {
			return fmt.Errorf("component %s: min_quantity %d exceeds max_quantity %d", product.SKU, *component.MinQuantity, *component.MaxQuantity)
		}

		components[i] = component
	}

	// Save components
	return s.bundleRepo.SetComponents(ctx, bundleProductID, components)
}

// CalculatePrice calculates the total price for a bundle based on customer selections
func (s *BundleService) CalculatePrice(ctx context.Context, bundleProductID uuid.UUID, req domain.BundlePriceRequest) (*domain.BundlePriceResponse, error) {
	// Load bundle product
	bundle, err := s.productRepo.GetByID(ctx, bundleProductID)
	if err != nil {
		return nil, err
	}
	if bundle.ProductType != domain.ProductTypeBundle {
		return nil, domain.ErrBundleNotFound
	}
	if bundle.BundlePriceMode == nil {
		return nil, fmt.Errorf("bundle %s has no price mode set", bundle.SKU)
	}

	priceMode := domain.BundlePriceMode(*bundle.BundlePriceMode)

	// Fixed price mode - just return the bundle's base price
	if priceMode == domain.BundlePriceModeFixed {
		prices, err := s.priceRepo.ListByProduct(ctx, bundleProductID)
		if err != nil {
			return nil, err
		}
		if len(prices) == 0 {
			return nil, fmt.Errorf("bundle %s has fixed price mode but no price defined", bundle.SKU)
		}

		// Use first active price (TODO: add tier pricing support)
		basePrice := prices[0].Price
		currency := prices[0].Currency

		return &domain.BundlePriceResponse{
			PriceMode:  priceMode,
			Total:      basePrice,
			Currency:   currency,
			Components: []domain.BundleComponentPriceResponse{}, // Empty for fixed mode
		}, nil
	}

	// Computed price mode - calculate sum of component prices
	components, err := s.bundleRepo.GetComponents(ctx, bundleProductID)
	if err != nil {
		return nil, err
	}

	// Build component ID map for quick lookup
	componentMap := make(map[uuid.UUID]*domain.BundleComponent)
	for i := range components {
		componentMap[components[i].ID] = &components[i]
	}

	var componentResponses []domain.BundleComponentPriceResponse
	var total float64
	var currency string

	for _, reqComp := range req.Components {
		// Find component
		component, ok := componentMap[reqComp.ComponentID]
		if !ok {
			return nil, fmt.Errorf("%w: %s", domain.ErrComponentNotFound, reqComp.ComponentID)
		}

		// Validate quantity (for configurable bundles)
		if bundle.BundleMode != nil && domain.BundleMode(*bundle.BundleMode) == domain.BundleModeConfigurable {
			if err := domain.ValidateComponentQuantity(component, reqComp.Quantity); err != nil {
				return nil, fmt.Errorf("component %s: %w", reqComp.ComponentID, err)
			}
		}

		// Load component product
		product, err := s.productRepo.GetByID(ctx, component.ComponentProductID)
		if err != nil {
			return nil, err
		}

		// Calculate component price based on product type
		var unitPrice float64
		var sku string

		switch product.ProductType {
		case domain.ProductTypeSimple, domain.ProductTypeVariant:
			// Get price from prices table
			prices, err := s.priceRepo.ListByProduct(ctx, product.ID)
			if err != nil {
				return nil, err
			}
			if len(prices) == 0 {
				return nil, fmt.Errorf("component %s has no price defined", product.SKU)
			}
			unitPrice = prices[0].Price
			currency = prices[0].Currency
			sku = product.SKU

		case domain.ProductTypeParametric:
			// Calculate parametric price
			priceReq := domain.ParametricPriceRequest{
				Parameters: reqComp.Parameters,
				Selections: reqComp.Selections,
				Quantity:   reqComp.Quantity,
			}
			priceResp, err := s.parametricSvc.CalculatePrice(ctx, product.ID, priceReq)
			if err != nil {
				return nil, fmt.Errorf("component %s parametric pricing: %w", product.SKU, err)
			}
			unitPrice = priceResp.UnitPrice
			currency = priceResp.Currency
			sku = priceResp.SKU
		}

		lineTotal := unitPrice * float64(reqComp.Quantity)
		total += lineTotal

		componentResponses = append(componentResponses, domain.BundleComponentPriceResponse{
			ComponentID: reqComp.ComponentID,
			ProductID:   product.ID,
			SKU:         sku,
			UnitPrice:   unitPrice,
			Quantity:    reqComp.Quantity,
			LineTotal:   lineTotal,
		})
	}

	return &domain.BundlePriceResponse{
		PriceMode:  priceMode,
		Components: componentResponses,
		Total:      total,
		Currency:   currency,
	}, nil
}
