package service

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
)

// MockProductRepository is a mock implementation for testing
type MockProductRepository struct {
	products map[uuid.UUID]*domain.Product
	bySKU    map[string]*domain.Product
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		products: make(map[uuid.UUID]*domain.Product),
		bySKU:    make(map[string]*domain.Product),
	}
}

func (m *MockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, ok := m.products[id]
	if !ok {
		return nil, domain.ErrProductNotFound
	}
	return product, nil
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*domain.Product, error) {
	key := tenantID.String() + ":" + sku
	product, ok := m.bySKU[key]
	if !ok {
		return nil, domain.ErrProductNotFound
	}
	return product, nil
}

func (m *MockProductRepository) List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error) {
	var results []domain.Product
	for _, p := range m.products {
		if p.TenantID == filter.TenantID && p.DeletedAt == nil {
			results = append(results, *p)
		}
	}
	return results, len(results), nil
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	m.products[product.ID] = product
	key := product.TenantID.String() + ":" + product.SKU
	m.bySKU[key] = product
	return nil
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	if _, ok := m.products[product.ID]; !ok {
		return domain.ErrProductNotFound
	}
	m.products[product.ID] = product
	return nil
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, ok := m.products[id]; !ok {
		return domain.ErrProductNotFound
	}
	delete(m.products, id)
	return nil
}

// Variant-specific methods (no-op for product tests)
func (m *MockProductRepository) GetProductWithVariants(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return m.GetByID(ctx, id)
}
func (m *MockProductRepository) ListVariants(ctx context.Context, parentID uuid.UUID, status ...domain.ProductStatus) ([]domain.Product, error) {
	return nil, nil
}
func (m *MockProductRepository) FindVariantByAxisValues(ctx context.Context, parentID uuid.UUID, axisValues map[string]string) (*domain.Product, error) {
	return nil, domain.ErrProductNotFound
}
func (m *MockProductRepository) GetAvailableAxisValues(ctx context.Context, parentID uuid.UUID, selected map[string]string) (map[string][]domain.AxisOption, error) {
	return nil, nil
}
func (m *MockProductRepository) SetVariantAxes(ctx context.Context, parentID uuid.UUID, axes []domain.VariantAxis) error {
	return nil
}
func (m *MockProductRepository) GetVariantAxes(ctx context.Context, parentID uuid.UUID) ([]domain.VariantAxis, error) {
	return nil, nil
}
func (m *MockProductRepository) SetAxisValues(ctx context.Context, variantID uuid.UUID, values []domain.AxisValueEntry) error {
	return nil
}
func (m *MockProductRepository) GetAxisValues(ctx context.Context, variantID uuid.UUID) ([]domain.AxisValueEntry, error) {
	return nil, nil
}

// MockPriceRepository for product service tests
type MockPriceRepository struct{}

func (m *MockPriceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	return nil, domain.ErrPriceNotFound
}
func (m *MockPriceRepository) ListByProduct(ctx context.Context, productID uuid.UUID) ([]domain.Price, error) {
	return nil, nil
}
func (m *MockPriceRepository) List(ctx context.Context, filter domain.PriceFilter) ([]domain.Price, int, error) {
	return nil, 0, nil
}
func (m *MockPriceRepository) Create(ctx context.Context, price *domain.Price) error { return nil }
func (m *MockPriceRepository) Update(ctx context.Context, price *domain.Price) error { return nil }
func (m *MockPriceRepository) Delete(ctx context.Context, id uuid.UUID) error        { return nil }
func (m *MockPriceRepository) CheckOverlap(ctx context.Context, price *domain.Price) (bool, error) {
	return false, nil
}

func TestProductService_Create(t *testing.T) {
	repo := NewMockProductRepository()
	service := NewProductService(repo, &MockPriceRepository{}, nil)
	ctx := context.Background()
	tenantID := uuid.New()

	req := domain.CreateProductRequest{
		SKU: "TEST-001",
		Name: map[string]string{
			"de": "Testprodukt",
			"en": "Test Product",
		},
		Description: map[string]string{
			"de": "Ein Testprodukt",
			"en": "A test product",
		},
		Status: domain.ProductStatusActive,
	}

	product, err := service.Create(ctx, tenantID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.SKU != req.SKU {
		t.Errorf("expected SKU %s, got %s", req.SKU, product.SKU)
	}

	if product.Status != domain.ProductStatusActive {
		t.Errorf("expected status active, got %s", product.Status)
	}

	// Try to create duplicate
	_, err = service.Create(ctx, tenantID, req)
	if err != domain.ErrProductAlreadyExists {
		t.Errorf("expected ErrProductAlreadyExists, got %v", err)
	}
}

func TestProductService_GetByID(t *testing.T) {
	repo := NewMockProductRepository()
	service := NewProductService(repo, &MockPriceRepository{}, nil)
	ctx := context.Background()

	// Test non-existent product
	_, err := service.GetByID(ctx, uuid.New())
	if err != domain.ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}

	// Create and retrieve product
	tenantID := uuid.New()
	product := domain.NewProduct(tenantID, "TEST-002")
	repo.Create(ctx, product)

	retrieved, err := service.GetByID(ctx, product.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrieved.ID != product.ID {
		t.Errorf("expected ID %s, got %s", product.ID, retrieved.ID)
	}
}

func TestProductService_Update(t *testing.T) {
	repo := NewMockProductRepository()
	service := NewProductService(repo, &MockPriceRepository{}, nil)
	ctx := context.Background()

	tenantID := uuid.New()
	product := domain.NewProduct(tenantID, "TEST-003")
	repo.Create(ctx, product)

	newStatus := domain.ProductStatusArchived
	req := domain.UpdateProductRequest{
		Status: &newStatus,
		Name: map[string]string{
			"de": "Aktualisiertes Produkt",
		},
	}

	updated, err := service.Update(ctx, product.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.Status != newStatus {
		t.Errorf("expected status %s, got %s", newStatus, updated.Status)
	}

	if updated.Name["de"] != "Aktualisiertes Produkt" {
		t.Errorf("expected updated name, got %s", updated.Name["de"])
	}
}

func TestProductService_Delete(t *testing.T) {
	repo := NewMockProductRepository()
	service := NewProductService(repo, &MockPriceRepository{}, nil)
	ctx := context.Background()

	tenantID := uuid.New()
	product := domain.NewProduct(tenantID, "TEST-004")
	repo.Create(ctx, product)

	err := service.Delete(ctx, product.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify product is deleted
	_, err = service.GetByID(ctx, product.ID)
	if err != domain.ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound after delete, got %v", err)
	}
}
