package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gondolia/gondolia/services/catalog/internal/middleware"
	"github.com/gondolia/gondolia/services/catalog/internal/service"
)

// SearchHandler handles search endpoints
type SearchHandler struct {
	searchService *service.SearchService
	syncService   *service.SyncService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService *service.SearchService, syncService *service.SyncService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
		syncService:   syncService,
	}
}

// Search handles GET /search
func (h *SearchHandler) Search(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	query := c.Query("q")
	filters := make(map[string]any)

	// Parse filters
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if productType := c.Query("type"); productType != "" {
		filters["product_type"] = productType
	}
	if excludeType := c.Query("exclude_type"); excludeType != "" {
		filters["exclude_product_type"] = excludeType
	}
	if category := c.Query("category"); category != "" {
		// Collect category + all descendant IDs for hierarchical filtering
		categoryIDs := []string{category}
		if h.searchService != nil {
			if descendants, err := h.searchService.GetCategoryDescendantIDs(c.Request.Context(), tenantID, category); err == nil {
				categoryIDs = append(categoryIDs, descendants...)
			}
		}
		filters["category_ids"] = categoryIDs
		// Debug: log category IDs
	}

	// Pagination
	offset := parseInt(c.Query("offset"), 0)
	limit := parseInt(c.Query("limit"), 20)
	if limit > 100 {
		limit = 100
	}

	result, err := h.searchService.Search(c.Request.Context(), tenantID, query, filters, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "SEARCH_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// SyncPIM handles POST /sync/pim
func (h *SearchHandler) SyncPIM(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	fullSync := c.Query("full") == "true"

	result, err := h.syncService.SyncFromPIM(c.Request.Context(), tenantID, fullSync)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "SYNC_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
