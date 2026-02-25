package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/middleware"
	"github.com/gondolia/gondolia/services/catalog/internal/service"
)

// CategoryHandler handles category endpoints
type CategoryHandler struct {
	categoryService *service.CategoryService
	productService  *service.ProductService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService *service.CategoryService, productService *service.ProductService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		productService:  productService,
	}
}

// GetTree handles GET /categories (returns tree)
func (h *CategoryHandler) GetTree(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	categories, err := h.categoryService.GetTree(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

// List handles GET /categories/list (with pagination)
func (h *CategoryHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	filter := domain.CategoryFilter{
		TenantID: tenantID,
		Limit:    100,
		Offset:   0,
	}

	// Parse query parameters
	if c.Query("parent_id") != "" {
		parentID, err := uuid.Parse(c.Query("parent_id"))
		if err == nil {
			filter.ParentID = &parentID
		}
	}

	if c.Query("active") != "" {
		active := c.Query("active") == "true"
		filter.Active = &active
	}

	if c.Query("search") != "" {
		search := c.Query("search")
		filter.Search = &search
	}

	// Pagination
	if limit := parseInt(c.Query("limit"), 100); limit > 0 {
		filter.Limit = limit
	}
	if offset := parseInt(c.Query("offset"), 0); offset >= 0 {
		filter.Offset = offset
	}

	categories, total, err := h.categoryService.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   categories,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// Get handles GET /categories/:id
func (h *CategoryHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid category ID",
			},
		})
		return
	}

	// Check if ancestors should be included (for breadcrumbs)
	includeAncestors := c.Query("include_ancestors") == "true"
	
	var category *domain.Category
	if includeAncestors {
		category, err = h.categoryService.GetByIDWithAncestors(c.Request.Context(), id)
	} else {
		category, err = h.categoryService.GetByID(c.Request.Context(), id)
	}

	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if domain.IsNotFoundError(err) {
			status = http.StatusNotFound
			code = "NOT_FOUND"
		}
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Create handles POST /categories
func (h *CategoryHandler) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req domain.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	category, err := h.categoryService.Create(c.Request.Context(), tenantID, req)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if domain.IsValidationError(err) {
			status = http.StatusBadRequest
			code = "VALIDATION_ERROR"
		}
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// Update handles PUT /categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid category ID",
			},
		})
		return
	}

	var req domain.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	category, err := h.categoryService.Update(c.Request.Context(), id, req)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if domain.IsNotFoundError(err) {
			status = http.StatusNotFound
			code = "NOT_FOUND"
		} else if domain.IsValidationError(err) {
			status = http.StatusBadRequest
			code = "VALIDATION_ERROR"
		}
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Delete handles DELETE /categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid category ID",
			},
		})
		return
	}

	if err := h.categoryService.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if domain.IsNotFoundError(err) {
			status = http.StatusNotFound
			code = "NOT_FOUND"
		} else if domain.IsValidationError(err) {
			status = http.StatusBadRequest
			code = "VALIDATION_ERROR"
		}
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetProducts handles GET /categories/:id/products
func (h *CategoryHandler) GetProducts(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid category ID",
			},
		})
		return
	}

	// Verify category exists
	_, err = h.categoryService.GetByID(c.Request.Context(), categoryID)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if domain.IsNotFoundError(err) {
			status = http.StatusNotFound
			code = "NOT_FOUND"
		}
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	filter := domain.ProductFilter{
		TenantID:        tenantID,
		CategoryID:      &categoryID,
		IncludeChildren: c.Query("include_children") == "true",
		Limit:           50,
		Offset:          0,
	}

	// Parse query parameters
	if c.Query("status") != "" {
		status := domain.ProductStatus(c.Query("status"))
		filter.Status = &status
	}

	if c.Query("search") != "" {
		search := c.Query("search")
		filter.Search = &search
	}

	// Pagination
	if limit := parseInt(c.Query("limit"), 50); limit > 0 {
		filter.Limit = limit
	}
	if offset := parseInt(c.Query("offset"), 0); offset >= 0 {
		filter.Offset = offset
	}

	products, total, err := h.productService.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   products,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// UpdateSortOrder handles PATCH /categories/:id/sort
func (h *CategoryHandler) UpdateSortOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid category ID",
			},
		})
		return
	}

	var req struct {
		SortOrder int `json:"sort_order" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	category, err := h.categoryService.UpdateSortOrder(c.Request.Context(), id, req.SortOrder)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if domain.IsNotFoundError(err) {
			status = http.StatusNotFound
			code = "NOT_FOUND"
		}
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

// AddProduct handles POST /categories/:id/products
func (h *CategoryHandler) AddProduct(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid category ID",
			},
		})
		return
	}

	var req struct {
		ProductID uuid.UUID `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	if err := h.categoryService.AddProduct(c.Request.Context(), categoryID, req.ProductID); err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if domain.IsNotFoundError(err) {
			status = http.StatusNotFound
			code = "NOT_FOUND"
		} else if domain.IsValidationError(err) {
			status = http.StatusBadRequest
			code = "VALIDATION_ERROR"
		}
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveProduct handles DELETE /categories/:id/products/:productId
func (h *CategoryHandler) RemoveProduct(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid category ID",
			},
		})
		return
	}

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid product ID",
			},
		})
		return
	}

	if err := h.categoryService.RemoveProduct(c.Request.Context(), categoryID, productID); err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if domain.IsNotFoundError(err) {
			status = http.StatusNotFound
			code = "NOT_FOUND"
		}
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}
