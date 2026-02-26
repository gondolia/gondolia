package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/gondolia/gondolia/provider/search"
)

// SearchService handles product search operations
type SearchService struct {
	searchProvider search.SearchProvider
}

// NewSearchService creates a new search service
func NewSearchService(searchProvider search.SearchProvider) *SearchService {
	return &SearchService{
		searchProvider: searchProvider,
	}
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
