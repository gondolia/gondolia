package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/repository"
)

// PriceService handles price business logic
type PriceService struct {
	priceRepo   repository.PriceRepository
	productRepo repository.ProductRepository
}

// NewPriceService creates a new price service
func NewPriceService(priceRepo repository.PriceRepository, productRepo repository.ProductRepository) *PriceService {
	return &PriceService{
		priceRepo:   priceRepo,
		productRepo: productRepo,
	}
}

// GetByID retrieves a price by ID
func (s *PriceService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Price, error) {
	return s.priceRepo.GetByID(ctx, id)
}

// ListByProduct retrieves all prices for a product
func (s *PriceService) ListByProduct(ctx context.Context, productID uuid.UUID) ([]domain.Price, error) {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return s.priceRepo.ListByProduct(ctx, productID)
}

// List retrieves prices with filtering and pagination
func (s *PriceService) List(ctx context.Context, filter domain.PriceFilter) ([]domain.Price, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 200 {
		filter.Limit = 200
	}
	return s.priceRepo.List(ctx, filter)
}

// Create creates a new price
func (s *PriceService) Create(ctx context.Context, tenantID, productID uuid.UUID, req domain.CreatePriceRequest) (*domain.Price, error) {
	// Verify product exists and belongs to tenant
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product.TenantID != tenantID {
		return nil, domain.ErrProductNotFound
	}

	// Validate date range
	if req.ValidFrom != nil && req.ValidTo != nil {
		if req.ValidFrom.After(*req.ValidTo) {
			return nil, domain.ErrPriceInvalidRange
		}
	}

	// Create price
	price := domain.NewPrice(tenantID, productID, req.Price, req.Currency)
	price.CustomerGroupID = req.CustomerGroupID
	
	if req.MinQuantity != nil {
		price.MinQuantity = *req.MinQuantity
	}
	price.ValidFrom = req.ValidFrom
	price.ValidTo = req.ValidTo

	// Check for overlapping prices
	hasOverlap, err := s.priceRepo.CheckOverlap(ctx, price)
	if err != nil {
		return nil, err
	}
	if hasOverlap {
		return nil, domain.ErrPriceOverlap
	}

	if err := s.priceRepo.Create(ctx, price); err != nil {
		return nil, err
	}

	return price, nil
}

// Update updates a price
func (s *PriceService) Update(ctx context.Context, id uuid.UUID, req domain.UpdatePriceRequest) (*domain.Price, error) {
	// Get existing price
	price, err := s.priceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.MinQuantity != nil {
		price.MinQuantity = *req.MinQuantity
	}
	if req.Price != nil {
		price.Price = *req.Price
	}
	if req.Currency != nil {
		price.Currency = *req.Currency
	}
	if req.ValidFrom != nil {
		price.ValidFrom = req.ValidFrom
	}
	if req.ValidTo != nil {
		price.ValidTo = req.ValidTo
	}

	// Validate date range
	if price.ValidFrom != nil && price.ValidTo != nil {
		if price.ValidFrom.After(*price.ValidTo) {
			return nil, domain.ErrPriceInvalidRange
		}
	}

	price.UpdatedAt = time.Now()

	// Check for overlapping prices
	hasOverlap, err := s.priceRepo.CheckOverlap(ctx, price)
	if err != nil {
		return nil, err
	}
	if hasOverlap {
		return nil, domain.ErrPriceOverlap
	}

	if err := s.priceRepo.Update(ctx, price); err != nil {
		return nil, err
	}

	return price, nil
}

// Delete soft-deletes a price
func (s *PriceService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.priceRepo.Delete(ctx, id)
}
