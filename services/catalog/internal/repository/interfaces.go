package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
)

// TenantRepository defines the interface for tenant data access
type TenantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	GetByCode(ctx context.Context, code string) (*domain.Tenant, error)
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*domain.Product, error)
	List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error)
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error // Soft delete

	// Variant-specific methods
	GetProductWithVariants(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ListVariants(ctx context.Context, parentID uuid.UUID, status ...domain.ProductStatus) ([]domain.Product, error)
	FindVariantByAxisValues(ctx context.Context, parentID uuid.UUID, axisValues map[string]string) (*domain.Product, error)
	GetAvailableAxisValues(ctx context.Context, parentID uuid.UUID, selected map[string]string) (map[string][]domain.AxisOption, error)

	// Variant axes management
	SetVariantAxes(ctx context.Context, parentID uuid.UUID, axes []domain.VariantAxis) error
	GetVariantAxes(ctx context.Context, parentID uuid.UUID) ([]domain.VariantAxis, error)

	// Axis values for variants
	SetAxisValues(ctx context.Context, variantID uuid.UUID, values []domain.AxisValueEntry) error
	GetAxisValues(ctx context.Context, variantID uuid.UUID) ([]domain.AxisValueEntry, error)
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetByIDWithAncestors(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Category, error)
	GetTree(ctx context.Context, tenantID uuid.UUID) ([]domain.Category, error)
	List(ctx context.Context, filter domain.CategoryFilter) ([]domain.Category, int, error)
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error // Soft delete
	HasProducts(ctx context.Context, id uuid.UUID) (bool, error)
	GetAncestors(ctx context.Context, id uuid.UUID) ([]domain.Category, error)
}

// PriceRepository defines the interface for price data access
type PriceRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Price, error)
	ListByProduct(ctx context.Context, productID uuid.UUID) ([]domain.Price, error)
	List(ctx context.Context, filter domain.PriceFilter) ([]domain.Price, int, error)
	Create(ctx context.Context, price *domain.Price) error
	Update(ctx context.Context, price *domain.Price) error
	Delete(ctx context.Context, id uuid.UUID) error // Soft delete
	CheckOverlap(ctx context.Context, price *domain.Price) (bool, error)
}

// ParametricPricingRepository defines the interface for parametric pricing data access
type ParametricPricingRepository interface {
	GetByProductID(ctx context.Context, productID uuid.UUID) (*domain.ParametricPricing, error)
	Create(ctx context.Context, pricing *domain.ParametricPricing) error
	Update(ctx context.Context, pricing *domain.ParametricPricing) error
	Delete(ctx context.Context, productID uuid.UUID) error
}

// AxisOptionRepository defines the interface for parametric axis option data access
type AxisOptionRepository interface {
	ListByAxisID(ctx context.Context, axisID uuid.UUID) ([]domain.AxisOption, error)
	ListByProductID(ctx context.Context, productID uuid.UUID) (map[uuid.UUID][]domain.AxisOption, error)
}

// SKUMappingRepository defines the interface for parametric SKU mapping data access
type SKUMappingRepository interface {
	FindBySelections(ctx context.Context, productID uuid.UUID, selections map[string]string) (*domain.SKUMapping, error)
	ListByProductID(ctx context.Context, productID uuid.UUID) ([]domain.SKUMapping, error)
}

// AttributeTranslationRepository defines the interface for attribute translation data access
type AttributeTranslationRepository interface {
	GetByKey(ctx context.Context, tenantID uuid.UUID, attributeKey, locale string) (*domain.AttributeTranslation, error)
	GetByTenantAndLocale(ctx context.Context, tenantID uuid.UUID, locale string) (map[string]*domain.AttributeTranslation, error)
	List(ctx context.Context, filter domain.AttributeTranslationFilter) ([]domain.AttributeTranslation, int, error)
	Create(ctx context.Context, translation *domain.AttributeTranslation) error
	Update(ctx context.Context, translation *domain.AttributeTranslation) error
	Delete(ctx context.Context, id uuid.UUID) error
}
