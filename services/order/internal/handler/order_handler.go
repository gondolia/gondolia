package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/order/internal/domain"
	"github.com/gondolia/gondolia/services/order/internal/middleware"
	"github.com/gondolia/gondolia/services/order/internal/service"
)

// OrderHandler handles order endpoints
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// Checkout handles POST /checkout
func (h *OrderHandler) Checkout(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	sessionID := c.GetHeader("X-Session-ID")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "user authentication required",
			},
		})
		return
	}

	var req domain.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	order, err := h.orderService.Checkout(c.Request.Context(), tenantID, *userID, sessionID, &req)
	if err != nil {
		if domain.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": err.Error(),
				},
			})
			return
		}
		if domain.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": order,
	})
}

// List handles GET /orders
func (h *OrderHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "user authentication required",
			},
		})
		return
	}

	filter := domain.OrderFilter{
		TenantID: tenantID,
		UserID:   userID,
		Limit:    50,
		Offset:   0,
	}

	// Parse query parameters
	if c.Query("status") != "" {
		status := domain.OrderStatus(c.Query("status"))
		filter.Status = &status
	}

	// Pagination
	if limit := parseInt(c.Query("limit"), 50); limit > 0 {
		filter.Limit = limit
	}
	if offset := parseInt(c.Query("offset"), 0); offset >= 0 {
		filter.Offset = offset
	}

	orders, total, err := h.orderService.ListOrders(c.Request.Context(), filter)
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
		"data":   orders,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// Get handles GET /orders/:id
func (h *OrderHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "user authentication required",
			},
		})
		return
	}

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid order ID",
			},
		})
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		if domain.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "order not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Check if order belongs to user
	if order.UserID != *userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "access denied",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": order,
	})
}

// Cancel handles PATCH /orders/:id/cancel
func (h *OrderHandler) Cancel(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "user authentication required",
			},
		})
		return
	}

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid order ID",
			},
		})
		return
	}

	order, err := h.orderService.CancelOrder(c.Request.Context(), orderID, *userID)
	if err != nil {
		if domain.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "order not found",
				},
			})
			return
		}
		if domain.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": err.Error(),
				},
			})
			return
		}
		if err == domain.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "access denied",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": order,
	})
}
