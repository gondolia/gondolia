package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/middleware"
	"github.com/gondolia/gondolia/services/catalog/internal/service"
)

// BundleHandler handles bundle product endpoints
type BundleHandler struct {
	bundleService *service.BundleService
}

// NewBundleHandler creates a new bundle handler
func NewBundleHandler(bundleService *service.BundleService) *BundleHandler {
	return &BundleHandler{
		bundleService: bundleService,
	}
}

// GetComponents handles GET /products/:id/bundle-components
func (h *BundleHandler) GetComponents(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_ID", "message": "Invalid product ID"},
		})
		return
	}

	components, err := h.bundleService.GetComponents(c.Request.Context(), productID)
	if err != nil {
		if errors.Is(err, domain.ErrBundleNotFound) || errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "Bundle not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to get components"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"components": components})
}

// SetComponents handles PUT /products/:id/bundle-components
func (h *BundleHandler) SetComponents(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_ID", "message": "Invalid product ID"},
		})
		return
	}

	// Get tenant from context
	tenantID := middleware.GetTenantID(c)

	var req domain.SetBundleComponentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()},
		})
		return
	}

	err = h.bundleService.SetComponents(c.Request.Context(), productID, tenantID, req.Components)
	if err != nil {
		if errors.Is(err, domain.ErrBundleNotFound) || errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "Bundle not found"},
			})
			return
		}
		if errors.Is(err, domain.ErrInvalidComponentType) ||
			errors.Is(err, domain.ErrBundleNesting) ||
			errors.Is(err, domain.ErrVariantParentInBundle) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{"code": "INVALID_COMPONENT", "message": err.Error()},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to set components"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Components updated successfully"})
}

// AddComponent handles POST /products/:id/bundle-components
func (h *BundleHandler) AddComponent(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_ID", "message": "Invalid product ID"},
		})
		return
	}

	tenantID := middleware.GetTenantID(c)

	var req domain.BundleComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()},
		})
		return
	}

	component, err := h.bundleService.AddComponent(c.Request.Context(), productID, tenantID, req)
	if err != nil {
		if errors.Is(err, domain.ErrBundleNotFound) || errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": err.Error()},
			})
			return
		}
		if errors.Is(err, domain.ErrInvalidComponentType) ||
			errors.Is(err, domain.ErrBundleNesting) ||
			errors.Is(err, domain.ErrVariantParentInBundle) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{"code": "INVALID_COMPONENT", "message": err.Error()},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to add component"},
		})
		return
	}

	c.JSON(http.StatusCreated, component)
}

// UpdateComponent handles PUT /products/:id/bundle-components/:compId
func (h *BundleHandler) UpdateComponent(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_ID", "message": "Invalid product ID"},
		})
		return
	}

	componentID, err := uuid.Parse(c.Param("compId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_ID", "message": "Invalid component ID"},
		})
		return
	}

	var req domain.BundleComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()},
		})
		return
	}

	component, err := h.bundleService.UpdateComponent(c.Request.Context(), productID, componentID, req)
	if err != nil {
		if errors.Is(err, domain.ErrComponentNotFound) || errors.Is(err, domain.ErrBundleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": err.Error()},
			})
			return
		}
		if errors.Is(err, domain.ErrInvalidComponentType) ||
			errors.Is(err, domain.ErrBundleNesting) ||
			errors.Is(err, domain.ErrVariantParentInBundle) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{"code": "INVALID_COMPONENT", "message": err.Error()},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to update component"},
		})
		return
	}

	c.JSON(http.StatusOK, component)
}

// DeleteComponent handles DELETE /products/:id/bundle-components/:compId
func (h *BundleHandler) DeleteComponent(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_ID", "message": "Invalid product ID"},
		})
		return
	}

	componentID, err := uuid.Parse(c.Param("compId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_ID", "message": "Invalid component ID"},
		})
		return
	}

	if err := h.bundleService.DeleteComponent(c.Request.Context(), productID, componentID); err != nil {
		if errors.Is(err, domain.ErrComponentNotFound) || errors.Is(err, domain.ErrBundleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": err.Error()},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to delete component"},
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// CalculatePrice handles POST /bundles/:id/calculate-price
func (h *BundleHandler) CalculatePrice(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_ID", "message": "Invalid product ID"},
		})
		return
	}

	var req domain.BundlePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()},
		})
		return
	}

	result, err := h.bundleService.CalculatePrice(c.Request.Context(), productID, req)
	if err != nil {
		if errors.Is(err, domain.ErrBundleNotFound) || errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "Bundle not found"},
			})
			return
		}
		if errors.Is(err, domain.ErrComponentNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{"code": "INVALID_COMPONENT", "message": err.Error()},
			})
			return
		}
		if errors.Is(err, domain.ErrQuantityOutOfRange) || errors.Is(err, domain.ErrInvalidQuantity) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{"code": "INVALID_QUANTITY", "message": err.Error()},
			})
			return
		}
		// Parametric pricing errors
		if errors.Is(err, domain.ErrParameterOutOfRange) ||
			errors.Is(err, domain.ErrParameterInvalidStep) ||
			errors.Is(err, domain.ErrMissingParameter) ||
			errors.Is(err, domain.ErrMissingSelection) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{"code": "INVALID_PARAMETERS", "message": err.Error()},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to calculate price"},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
