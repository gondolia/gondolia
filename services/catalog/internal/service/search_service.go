package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/provider/search"
	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/repository"
)

// SearchService handles product search operations
type SearchService struct {
	searchProvider search.SearchProvider
	categoryRepo   repository.CategoryRepository
}

// NewSearchService creates a new search service
func NewSearchService(searchProvider search.SearchProvider, categoryRepo repository.CategoryRepository) *SearchService {
	return &SearchService{
		searchProvider: searchProvider,
		categoryRepo:   categoryRepo,
	}
}

// GetCategoryDescendantIDs returns all descendant category IDs for a given category
func (s *SearchService) GetCategoryDescendantIDs(ctx context.Context, tenantID uuid.UUID, categoryID string) ([]string, error) {
	catUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, err
	}
	// Use List with parent filter to find direct children, then recurse
	// Simpler: query all categories and find descendants in memory
	filter := domain.CategoryFilter{
		TenantID: tenantID,
		Limit:    1000,
	}
	cats, _, err := s.categoryRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	// Build parent->children map
	childMap := make(map[uuid.UUID][]uuid.UUID)
	for _, cat := range cats {
		if cat.ParentID != nil {
			childMap[*cat.ParentID] = append(childMap[*cat.ParentID], cat.ID)
		}
	}
	// BFS to collect all descendants
	var ids []string
	queue := []uuid.UUID{catUUID}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, childID := range childMap[current] {
			ids = append(ids, childID.String())
			queue = append(queue, childID)
		}
	}
	return ids, nil
}

// Search searches for products
func (s *SearchService) Search(ctx context.Context, tenantID uuid.UUID, query string, filters map[string]any, offset, limit int) (*SearchResult, error) {
	searchQuery := search.SearchQuery{
		Query:  query,
		Offset: offset,
		Limit:  limit,
	}

	// Add tenant filter
	searchQuery.Filters = []search.Filter{
		{
			Field:    "tenant_id",
			Operator: "=",
			Value:    tenantID.String(),
		},
	}

	// Add additional filters
	for field, value := range filters {
		if field == "exclude_product_type" {
			searchQuery.Filters = append(searchQuery.Filters, search.Filter{
				Field:    "product_type",
				Operator: "!=",
				Value:    value,
			})
		} else {
			searchQuery.Filters = append(searchQuery.Filters, search.Filter{
				Field:    field,
				Operator: "=",
				Value:    value,
			})
		}
	}

	result, err := s.searchProvider.Search(ctx, "products", searchQuery)
	if err != nil {
		return nil, err
	}

	return &SearchResult{
		Hits:      result.Hits,
		TotalHits: result.TotalHits,
		Facets:    result.Facets,
		Offset:    offset,
		Limit:     limit,
	}, nil
}

// SearchResult represents search results
type SearchResult struct {
	Hits      []search.Document         `json:"hits"`
	TotalHits int                       `json:"total_hits"`
	Facets    map[string]map[string]int `json:"facets,omitempty"`
	Offset    int                       `json:"offset"`
	Limit     int                       `json:"limit"`
}
