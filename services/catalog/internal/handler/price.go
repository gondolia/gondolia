package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/middleware"
	"github.com/gondolia/gondolia/services/catalog/internal/service"
)

// PriceHandler handles price endpoints
type PriceHandler struct {
	priceService *service.PriceService
}

// NewPriceHandler creates a new price handler
func NewPriceHandler(priceService *service.PriceService) *PriceHandler {
	return &PriceHandler{
		priceService: priceService,
	}
}

// ListByProduct handles GET /products/:productId/prices
func (h *PriceHandler) ListByProduct(c *gin.Context) {
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

	prices, err := h.priceService.ListByProduct(c.Request.Context(), productID)
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

	c.JSON(http.StatusOK, gin.H{
		"data": prices,
	})
}

// Create handles POST /products/:productId/prices
func (h *PriceHandler) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	
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

	var req domain.CreatePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	price, err := h.priceService.Create(c.Request.Context(), tenantID, productID, req)
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

	c.JSON(http.StatusCreated, price)
}

// Update handles PUT /prices/:id
func (h *PriceHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid price ID",
			},
		})
		return
	}

	var req domain.UpdatePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	price, err := h.priceService.Update(c.Request.Context(), id, req)
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

	c.JSON(http.StatusOK, price)
}

// Delete handles DELETE /prices/:id
func (h *PriceHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid price ID",
			},
		})
		return
	}

	if err := h.priceService.Delete(c.Request.Context(), id); err != nil {
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
