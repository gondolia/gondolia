package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/repository"
)

// VariantService handles variant product business logic
type VariantService struct {
	productRepo repository.ProductRepository
	priceRepo   repository.PriceRepository
}

// NewVariantService creates a new variant service
func NewVariantService(productRepo repository.ProductRepository, priceRepo repository.PriceRepository) *VariantService {
	return &VariantService{
		productRepo: productRepo,
		priceRepo:   priceRepo,
	}
}

// CreateVariantParent creates a new variant parent product with axes
func (s *VariantService) CreateVariantParent(ctx context.Context, tenantID uuid.UUID, req domain.CreateVariantParentRequest) (*domain.Product, error) {
	// Check if product with SKU already exists
	existing, err := s.productRepo.GetBySKU(ctx, tenantID, req.SKU)
	if err == nil && existing != nil {
		return nil, domain.ErrProductAlreadyExists
	}

	// Validate: at least one axis required
	if len(req.VariantAxes) == 0 {
		return nil, fmt.Errorf("variant parent must have at least one variant axis")
	}

	// Validate: max 4 axes (design decision 11.1)
	if len(req.VariantAxes) > domain.MaxVariantAxes {
		return nil, domain.ErrTooManyAxes
	}

	// Create parent product
	product := domain.NewProduct(tenantID, req.SKU)
	product.ProductType = domain.ProductTypeVariantParent
	product.Name = req.Name
	product.Description = req.Description
	product.CategoryIDs = req.CategoryIDs
	product.Attributes = req.Attributes
	product.Images = req.Images
	product.Status = domain.ProductStatusDraft // Parents start as draft

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	// Create variant axes
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
		return nil, fmt.Errorf("failed to set variant axes: %w", err)
	}

	product.VariantAxes = axes

	return product, nil
}

// CreateVariant creates a new variant under a parent product
func (s *VariantService) CreateVariant(ctx context.Context, tenantID uuid.UUID, parentID uuid.UUID, req domain.CreateVariantRequest) (*domain.Product, error) {
	// Get parent product
	parent, err := s.productRepo.GetByID(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("parent product not found: %w", err)
	}

	// Validate parent is variant_parent type
	if parent.ProductType != domain.ProductTypeVariantParent {
		return nil, fmt.Errorf("parent product is not a variant parent")
	}

	// Check if variant with SKU already exists
	existing, err := s.productRepo.GetBySKU(ctx, tenantID, req.SKU)
	if err == nil && existing != nil {
		return nil, domain.ErrProductAlreadyExists
	}

	// Get parent's axes
	axes, err := s.productRepo.GetVariantAxes(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent axes: %w", err)
	}

	// Validate: all axes must have values
	if err := s.validateAxisValues(axes, req.AxisValues); err != nil {
		return nil, err
	}

	// Check uniqueness: no other variant with same axis values
	existingVariant, err := s.productRepo.FindVariantByAxisValues(ctx, parentID, req.AxisValues)
	if err == nil && existingVariant != nil {
		return nil, fmt.Errorf("variant with these axis values already exists: %s", existingVariant.SKU)
	}

	// Create variant product
	variant := domain.NewProduct(tenantID, req.SKU)
	variant.ProductType = domain.ProductTypeVariant
	variant.ParentID = &parentID
	variant.Status = domain.ProductStatusDraft

	// Inherit from parent (if not explicitly set)
	variant.Name = s.inheritName(req.Name, parent.Name)
	variant.Description = parent.Description // Variants inherit parent description
	variant.CategoryIDs = parent.CategoryIDs
	variant.Attributes = s.mergeAttributes(parent.Attributes, req.Attributes)
	variant.Images = req.Images

	if err := s.productRepo.Create(ctx, variant); err != nil {
		return nil, err
	}

	// Set axis values
	axisValueEntries := make([]domain.AxisValueEntry, 0, len(req.AxisValues))
	for _, axis := range axes {
		if optCode, ok := req.AxisValues[axis.AttributeCode]; ok {
			axisValueEntries = append(axisValueEntries, domain.AxisValueEntry{
				VariantID:         variant.ID,
				AxisID:            axis.ID,
				AxisAttributeCode: axis.AttributeCode,
				OptionCode:        optCode,
			})
		}
	}

	if err := s.productRepo.SetAxisValues(ctx, variant.ID, axisValueEntries); err != nil {
		return nil, fmt.Errorf("failed to set axis values: %w", err)
	}

	variant.AxisValues = axisValueEntries

	return variant, nil
}

// GetProductWithVariants retrieves a product with all its variants (if parent)
func (s *VariantService) GetProductWithVariants(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := s.productRepo.GetProductWithVariants(ctx, id)
	if err != nil {
		return nil, err
	}

	// Enrich bundle products with base price
	if product.ProductType == domain.ProductTypeBundle && s.priceRepo != nil {
		prices, err := s.priceRepo.ListByProduct(ctx, product.ID)
		if err == nil && len(prices) > 0 {
			bp := prices[0].Price
			bc := prices[0].Currency
			product.BasePrice = &bp
			product.BaseCurrency = &bc
		}
	}

	// Enrich each variant with price and availability data
	if product.ProductType == domain.ProductTypeVariantParent && s.priceRepo != nil {
		for i := range product.Variants {
			prices, err := s.priceRepo.ListByProduct(ctx, product.Variants[i].ID)
			if err != nil || len(prices) == 0 {
				continue
			}

			// Build tier prices; pick base price = lowest min_quantity without customer group
			tierPrices := make([]domain.TierPrice, 0, len(prices))
			var basePrice float64
			var currency string
			for _, p := range prices {
				tierPrices = append(tierPrices, domain.TierPrice{
					MinQuantity: p.MinQuantity,
					Price:       p.Price,
				})
				// Prefer public price (no customer group) with min_quantity == 1
				if p.CustomerGroupID == nil && (basePrice == 0 || p.MinQuantity < 2) {
					basePrice = p.Price
					currency = strings.TrimSpace(p.Currency)
				}
			}
			// Fallback: use first price if no public price found
			if basePrice == 0 {
				basePrice = prices[0].Price
				currency = strings.TrimSpace(prices[0].Currency)
			}

			product.Variants[i].Price = &domain.VariantPrice{
				Net:        basePrice,
				Currency:   currency,
				TierPrices: tierPrices,
			}

			inStock := true
			product.Variants[i].Stock = &domain.VariantAvailability{
				InStock: inStock,
			}
		}
	}

	return product, nil
}

// ListVariants returns all variants for a parent product
func (s *VariantService) ListVariants(ctx context.Context, parentID uuid.UUID, status ...domain.ProductStatus) ([]domain.Product, error) {
	return s.productRepo.ListVariants(ctx, parentID, status...)
}

// SelectVariant finds a variant by axis values, resolves images and generates display name
func (s *VariantService) SelectVariant(ctx context.Context, parentID uuid.UUID, axisValues map[string]string) (*domain.Product, error) {
	variant, err := s.productRepo.FindVariantByAxisValues(ctx, parentID, axisValues)
	if err != nil {
		return nil, fmt.Errorf("variant not found for selected values: %w", err)
	}

	// Load axis values for the variant (with labels for name generation)
	values, err := s.productRepo.GetAxisValues(ctx, variant.ID)
	if err == nil {
		// Enrich axis values with labels (using shared helpers from label_helpers.go)
		for i := range values {
			values[i].OptionLabel = formatOptionLabel(values[i].OptionCode)
			values[i].AxisLabel = formatAxisLabel(values[i].AxisAttributeCode)
		}
		variant.AxisValues = values
	}

	// Load prices for this variant (including tier/staffel prices)
	if s.priceRepo != nil {
		prices, err := s.priceRepo.ListByProduct(ctx, variant.ID)
		if err == nil && len(prices) > 0 {
			tierPrices := make([]domain.TierPrice, 0, len(prices))
			var basePrice float64
			var currency string
			for _, p := range prices {
				tierPrices = append(tierPrices, domain.TierPrice{
					MinQuantity: p.MinQuantity,
					Price:       p.Price,
				})
				if p.MinQuantity <= 1 || basePrice == 0 {
					basePrice = p.Price
					currency = p.Currency
				}
			}
			variant.Price = &domain.VariantPrice{
				Net:        basePrice,
				Currency:   currency,
				TierPrices: tierPrices,
			}
		}
	}

	// Load parent for image inheritance and name generation
	parent, err := s.productRepo.GetByID(ctx, parentID)
	if err == nil {
		// Design decision 11.4: variant images replace parent images
		variant.Images = domain.GetEffectiveImages(variant.Images, parent.Images)

		// Design decision 11.5: auto-generate name if variant has no own name
		for locale := range parent.Name {
			if existing, ok := variant.Name[locale]; !ok || existing == "" || existing == parent.Name[locale] {
				generated := domain.GenerateVariantName(parent.Name, values, locale)
				if variant.Name == nil {
					variant.Name = make(map[string]string)
				}
				variant.Name[locale] = generated
			}
		}
	}

	return variant, nil
}

// GetAvailableAxisValues returns available axis options based on current selection
func (s *VariantService) GetAvailableAxisValues(ctx context.Context, parentID uuid.UUID, selected map[string]string) (map[string][]domain.AxisOption, error) {
	return s.productRepo.GetAvailableAxisValues(ctx, parentID, selected)
}

// validateAxisValues checks that all required axes have values
func (s *VariantService) validateAxisValues(axes []domain.VariantAxis, values map[string]string) error {
	for _, axis := range axes {
		if _, ok := values[axis.AttributeCode]; !ok {
			return fmt.Errorf("missing value for required axis: %s", axis.AttributeCode)
		}
	}

	// Check for extra values not in axes
	for attrCode := range values {
		found := false
		for _, axis := range axes {
			if axis.AttributeCode == attrCode {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown axis attribute: %s", attrCode)
		}
	}

	return nil
}

// inheritName inherits name from parent if variant name is empty
func (s *VariantService) inheritName(variantName, parentName map[string]string) map[string]string {
	if len(variantName) > 0 {
		return variantName
	}
	// Empty map means: inherit from parent (service layer will append axis values)
	result := make(map[string]string)
	for locale, name := range parentName {
		result[locale] = name
	}
	return result
}

// inheritMap inherits map values from parent if empty
func (s *VariantService) inheritMap(variant, parent map[string]string) map[string]string {
	if len(variant) > 0 {
		return variant
	}
	return parent
}

// mergeAttributes merges parent and variant attributes (variant overrides)
func (s *VariantService) mergeAttributes(parentAttrs, variantAttrs []domain.ProductAttribute) []domain.ProductAttribute {
	// Start with parent attributes
	merged := make(map[string]domain.ProductAttribute)
	for _, attr := range parentAttrs {
		merged[attr.Key] = attr
	}

	// Override with variant attributes
	for _, attr := range variantAttrs {
		merged[attr.Key] = attr
	}

	// Convert back to slice
	result := make([]domain.ProductAttribute, 0, len(merged))
	for _, attr := range merged {
		result = append(result, attr)
	}

	return result
}
