package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/middleware"
	"github.com/gondolia/gondolia/services/catalog/internal/service"
)

// ProductHandler handles product endpoints
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// List handles GET /products
func (h *ProductHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	filter := domain.ProductFilter{
		TenantID:        tenantID,
		Limit:           50,
		Offset:          0,
		ExcludeVariants: true, // Default: show only simple + variant_parent products
	}

	// Parse query parameters
	if c.Query("category_id") != "" {
		categoryID, err := uuid.Parse(c.Query("category_id"))
		if err == nil {
			filter.CategoryID = &categoryID
		}
	}

	if c.Query("status") != "" {
		status := domain.ProductStatus(c.Query("status"))
		filter.Status = &status
	}

	if c.Query("search") != "" {
		search := c.Query("search")
		filter.Search = &search
	}

	if c.Query("product_type") != "" {
		pt := domain.ProductType(c.Query("product_type"))
		filter.ProductType = &pt
	}

	// Pagination
	if limit := parseInt(c.Query("limit"), 50); limit > 0 {
		filter.Limit = limit
	}
	if offset := parseInt(c.Query("offset"), 0); offset >= 0 {
		filter.Offset = offset
	}

	// Check if locale is requested for translated attributes
	locale := c.Query("locale")
	if locale != "" {
		products, total, err := h.productService.ListWithTranslations(c.Request.Context(), filter, locale)
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
		return
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

// Get handles GET /products/:id
func (h *ProductHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid product ID",
			},
		})
		return
	}

	// Check if locale is requested for translated attributes
	locale := c.Query("locale")
	if locale != "" {
		product, err := h.productService.GetByIDWithTranslations(c.Request.Context(), id, locale)
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
		c.JSON(http.StatusOK, product)
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), id)
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

	c.JSON(http.StatusOK, product)
}

// Create handles POST /products
func (h *ProductHandler) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req domain.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	product, err := h.productService.Create(c.Request.Context(), tenantID, req)
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

	c.JSON(http.StatusCreated, product)
}

// Update handles PUT /products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid product ID",
			},
		})
		return
	}

	var req domain.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	product, err := h.productService.Update(c.Request.Context(), id, req)
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

	c.JSON(http.StatusOK, product)
}

// Delete handles DELETE /products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid product ID",
			},
		})
		return
	}

	if err := h.productService.Delete(c.Request.Context(), id); err != nil {
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
