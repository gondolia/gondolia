package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/middleware"
	"github.com/gondolia/gondolia/services/catalog/internal/service"
)

// VariantHandler handles variant product endpoints
type VariantHandler struct {
	variantService *service.VariantService
}

// NewVariantHandler creates a new variant handler
func NewVariantHandler(variantService *service.VariantService) *VariantHandler {
	return &VariantHandler{
		variantService: variantService,
	}
}

// CreateVariantParent handles POST /products (with product_type=variant_parent)
func (h *VariantHandler) CreateVariantParent(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req domain.CreateVariantParentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	product, err := h.variantService.CreateVariantParent(c.Request.Context(), tenantID, req)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if err == domain.ErrProductAlreadyExists {
			status = http.StatusConflict
			code = "PRODUCT_EXISTS"
		}
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": product})
}

// CreateVariant handles POST /products/:parent_id/variants
func (h *VariantHandler) CreateVariant(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	parentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid parent product ID",
			},
		})
		return
	}

	var req domain.CreateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	variant, err := h.variantService.CreateVariant(c.Request.Context(), tenantID, parentID, req)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if err == domain.ErrProductAlreadyExists {
			status = http.StatusConflict
			code = "PRODUCT_EXISTS"
		} else if err == domain.ErrProductNotFound {
			status = http.StatusNotFound
			code = "PARENT_NOT_FOUND"
		}
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": variant})
}

// GetProductWithVariants handles GET /products/:id (detects variant_parent and includes variants)
func (h *VariantHandler) GetProductWithVariants(c *gin.Context) {
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

	product, err := h.variantService.GetProductWithVariants(c.Request.Context(), id)
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

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// ListVariants handles GET /products/:id/variants
func (h *VariantHandler) ListVariants(c *gin.Context) {
	parentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid parent product ID",
			},
		})
		return
	}

	// Optional status filter
	var statusFilters []domain.ProductStatus
	if statusStr := c.Query("status"); statusStr != "" {
		statusFilters = append(statusFilters, domain.ProductStatus(statusStr))
	}

	variants, err := h.variantService.ListVariants(c.Request.Context(), parentID, statusFilters...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": variants})
}

// SelectVariant handles GET /products/:id/variants/select?axis=value&axis=value
func (h *VariantHandler) SelectVariant(c *gin.Context) {
	parentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid parent product ID",
			},
		})
		return
	}

	// Parse axis values from query params
	axisValues := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 && key != "status" {
			axisValues[key] = values[0]
		}
	}

	if len(axisValues) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_AXIS_VALUES",
				"message": "at least one axis value must be provided",
			},
		})
		return
	}

	variant, err := h.variantService.SelectVariant(c.Request.Context(), parentID, axisValues)
	if err != nil {
		status := http.StatusNotFound
		code := "VARIANT_NOT_FOUND"
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    code,
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": variant})
}

// GetAvailableAxisValues handles GET /products/:id/variants/available?axis=value
func (h *VariantHandler) GetAvailableAxisValues(c *gin.Context) {
	parentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid parent product ID",
			},
		})
		return
	}

	// Parse selected axis values from query params
	selected := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			selected[key] = values[0]
		}
	}

	available, err := h.variantService.GetAvailableAxisValues(c.Request.Context(), parentID, selected)
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
		"selected":  selected,
		"available": available,
	})
}
