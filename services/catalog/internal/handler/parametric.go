package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/service"
)

// ParametricHandler handles parametric product endpoints
type ParametricHandler struct {
	parametricService *service.ParametricService
}

// NewParametricHandler creates a new parametric handler
func NewParametricHandler(parametricService *service.ParametricService) *ParametricHandler {
	return &ParametricHandler{
		parametricService: parametricService,
	}
}

// CalculatePrice handles POST /products/:id/calculate-price
func (h *ParametricHandler) CalculatePrice(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_ID", "message": "Invalid product ID"},
		})
		return
	}

	var req domain.ParametricPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{"code": "INVALID_REQUEST", "message": err.Error()},
		})
		return
	}

	result, err := h.parametricService.CalculatePrice(c.Request.Context(), productID, req)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) || errors.Is(err, domain.ErrParametricPricingNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "Product or pricing not found"},
			})
			return
		}
		if errors.Is(err, domain.ErrParameterOutOfRange) || errors.Is(err, domain.ErrParameterInvalidStep) || errors.Is(err, domain.ErrMissingParameter) {
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
