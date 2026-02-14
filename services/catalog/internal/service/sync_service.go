package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/provider/pim"
	"github.com/gondolia/gondolia/provider/search"
	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/repository"
)

// SyncService handles PIM synchronization and search indexing
type SyncService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	pimProvider  pim.PIMProvider
	searchProvider search.SearchProvider
}

// NewSyncService creates a new sync service
func NewSyncService(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	pimProvider pim.PIMProvider,
	searchProvider search.SearchProvider,
) *SyncService {
	return &SyncService{
		productRepo:    productRepo,
		categoryRepo:   categoryRepo,
		pimProvider:    pimProvider,
		searchProvider: searchProvider,
	}
}

// SyncFromPIM synchronizes products and categories from PIM
func (s *SyncService) SyncFromPIM(ctx context.Context, tenantID uuid.UUID, fullSync bool) (*SyncResult, error) {
	result := &SyncResult{
		StartedAt: time.Now(),
	}

	// Sync categories first
	if err := s.syncCategories(ctx, tenantID, result); err != nil {
		result.Error = err.Error()
		result.CompletedAt = time.Now()
		return result, err
	}

	// Sync products
	if err := s.syncProducts(ctx, tenantID, fullSync, result); err != nil {
		result.Error = err.Error()
		result.CompletedAt = time.Now()
		return result, err
	}

	result.CompletedAt = time.Now()
	return result, nil
}

// syncCategories syncs categories from PIM
func (s *SyncService) syncCategories(ctx context.Context, tenantID uuid.UUID, result *SyncResult) error {
	categories, err := s.pimProvider.FetchCategories(ctx)
	if err != nil {
		return err
	}

	for _, pimCat := range categories {
		category, err := s.categoryRepo.GetByCode(ctx, tenantID, pimCat.Code)
		
		if err == domain.ErrCategoryNotFound {
			// Create new category
			category = domain.NewCategory(tenantID, pimCat.Code)
			category.Name = pimCat.Labels
			category.PIMCode = &pimCat.Code
			now := time.Now()
			category.LastSyncedAt = &now
			
			if err := s.categoryRepo.Create(ctx, category); err != nil {
				result.CategoriesFailed++
				continue
			}
			result.CategoriesCreated++
		} else if err == nil {
			// Update existing category
			category.Name = pimCat.Labels
			now := time.Now()
			category.LastSyncedAt = &now
			category.UpdatedAt = now
			
			if err := s.categoryRepo.Update(ctx, category); err != nil {
				result.CategoriesFailed++
				continue
			}
			result.CategoriesUpdated++
		} else {
			result.CategoriesFailed++
		}
	}

	return nil
}

// syncProducts syncs products from PIM
func (s *SyncService) syncProducts(ctx context.Context, tenantID uuid.UUID, fullSync bool, result *SyncResult) error {
	filter := pim.ProductFilter{
		Limit: 100,
	}

	if !fullSync {
		// Incremental sync - only products updated in last 24h
		since := time.Now().Add(-24 * time.Hour)
		filter.UpdatedSince = &since
	}

	for {
		page, err := s.pimProvider.FetchProducts(ctx, filter)
		if err != nil {
			return err
		}

		for _, pimProduct := range page.Products {
			if err := s.syncProduct(ctx, tenantID, pimProduct, result); err != nil {
				result.ProductsFailed++
				continue
			}
		}

		// Check if there are more pages
		if page.NextCursor == "" {
			break
		}
		filter.Cursor = page.NextCursor
	}

	return nil
}

// syncProduct syncs a single product
func (s *SyncService) syncProduct(ctx context.Context, tenantID uuid.UUID, pimProduct pim.Product, result *SyncResult) error {
	product, err := s.productRepo.GetBySKU(ctx, tenantID, pimProduct.Identifier)
	
	// Convert PIM product to domain product
	name := make(map[string]string)
	description := make(map[string]string)
	
	for _, val := range pimProduct.Values["name"] {
		if val.Locale != "" {
			if str, ok := val.Data.(string); ok {
				name[val.Locale] = str
			}
		}
	}
	
	for _, val := range pimProduct.Values["description"] {
		if val.Locale != "" {
			if str, ok := val.Data.(string); ok {
				description[val.Locale] = str
			}
		}
	}

	if err == domain.ErrProductNotFound {
		// Create new product
		product = domain.NewProduct(tenantID, pimProduct.Identifier)
		product.Name = name
		product.Description = description
		product.PIMIdentifier = &pimProduct.Identifier
		now := time.Now()
		product.LastSyncedAt = &now
		
		if pimProduct.Enabled {
			product.Status = domain.ProductStatusActive
		}
		
		if err := s.productRepo.Create(ctx, product); err != nil {
			return err
		}
		result.ProductsCreated++
	} else if err == nil {
		// Update existing product
		product.Name = name
		product.Description = description
		now := time.Now()
		product.LastSyncedAt = &now
		product.UpdatedAt = now
		
		if pimProduct.Enabled {
			product.Status = domain.ProductStatusActive
		} else {
			product.Status = domain.ProductStatusArchived
		}
		
		if err := s.productRepo.Update(ctx, product); err != nil {
			return err
		}
		result.ProductsUpdated++
	} else {
		return err
	}

	return nil
}

// IndexProduct indexes a product in the search engine
func (s *SyncService) IndexProduct(ctx context.Context, product *domain.Product) error {
	doc := search.Document{
		"id":          product.ID.String(),
		"tenant_id":   product.TenantID.String(),
		"sku":         product.SKU,
		"name":        product.Name,
		"description": product.Description,
		"status":      string(product.Status),
		"created_at":  product.CreatedAt.Unix(),
		"updated_at":  product.UpdatedAt.Unix(),
	}

	_, err := s.searchProvider.IndexDocuments(ctx, "products", []search.Document{doc})
	return err
}

// RemoveProductFromIndex removes a product from the search index
func (s *SyncService) RemoveProductFromIndex(ctx context.Context, productID uuid.UUID) error {
	_, err := s.searchProvider.DeleteDocuments(ctx, "products", []string{productID.String()})
	return err
}

// SyncResult represents the result of a PIM sync operation
type SyncResult struct {
	StartedAt         time.Time `json:"started_at"`
	CompletedAt       time.Time `json:"completed_at"`
	ProductsCreated   int       `json:"products_created"`
	ProductsUpdated   int       `json:"products_updated"`
	ProductsFailed    int       `json:"products_failed"`
	CategoriesCreated int       `json:"categories_created"`
	CategoriesUpdated int       `json:"categories_updated"`
	CategoriesFailed  int       `json:"categories_failed"`
	Error             string    `json:"error,omitempty"`
}
