package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/cart/internal/domain"
	"github.com/gondolia/gondolia/services/cart/internal/middleware"
	"github.com/gondolia/gondolia/services/cart/internal/service"
)

// CartHandler handles cart endpoints
type CartHandler struct {
	cartService *service.CartService
}

// NewCartHandler creates a new cart handler
func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// GetCart handles GET /cart
func (h *CartHandler) GetCart(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	sessionID := middleware.GetSessionID(c)

	cart, err := h.cartService.GetOrCreateCart(c.Request.Context(), tenantID, userID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddItem handles POST /cart/items
func (h *CartHandler) AddItem(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	sessionID := middleware.GetSessionID(c)

	var req domain.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	cart, err := h.cartService.AddItem(c.Request.Context(), tenantID, userID, sessionID, req)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if domain.IsValidationError(err) {
			status = http.StatusBadRequest
			code = "VALIDATION_ERROR"
		} else if domain.IsNotFoundError(err) {
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

	c.JSON(http.StatusOK, cart)
}

// UpdateItem handles PATCH /cart/items/:id
func (h *CartHandler) UpdateItem(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid item ID",
			},
		})
		return
	}

	var req domain.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	cart, err := h.cartService.UpdateItemQuantity(c.Request.Context(), itemID, req.Quantity)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if domain.IsValidationError(err) {
			status = http.StatusBadRequest
			code = "VALIDATION_ERROR"
		} else if domain.IsNotFoundError(err) {
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

	c.JSON(http.StatusOK, cart)
}

// RemoveItem handles DELETE /cart/items/:id
func (h *CartHandler) RemoveItem(c *gin.Context) {
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid item ID",
			},
		})
		return
	}

	cart, err := h.cartService.RemoveItem(c.Request.Context(), itemID)
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

	c.JSON(http.StatusOK, cart)
}

// ClearCart handles DELETE /cart
func (h *CartHandler) ClearCart(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	sessionID := middleware.GetSessionID(c)

	if err := h.cartService.ClearCart(c.Request.Context(), tenantID, userID, sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ValidateCart handles POST /cart/validate
func (h *CartHandler) ValidateCart(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	sessionID := middleware.GetSessionID(c)

	cart, err := h.cartService.ValidateCart(c.Request.Context(), tenantID, userID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// CompleteCart handles POST /cart/complete - marks cart as completed after order creation
func (h *CartHandler) CompleteCart(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	sessionID := middleware.GetSessionID(c)

	if err := h.cartService.CompleteCart(c.Request.Context(), tenantID, userID, sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart marked as completed",
	})
}
