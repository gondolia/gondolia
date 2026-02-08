package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
	"github.com/gondolia/gondolia/services/identity/internal/middleware"
	"github.com/gondolia/gondolia/services/identity/internal/service"
)

// UserHandler handles user endpoints
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// List returns a list of users
func (h *UserHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var query struct {
		CompanyID *uuid.UUID `form:"company_id"`
		Email     *string    `form:"email"`
		IsActive  *bool      `form:"is_active"`
		Search    *string    `form:"search"`
		Limit     int        `form:"limit,default=20"`
		Offset    int        `form:"offset,default=0"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	filter := domain.UserFilter{
		TenantID:  tenantID,
		CompanyID: query.CompanyID,
		Email:     query.Email,
		IsActive:  query.IsActive,
		Search:    query.Search,
		Limit:     query.Limit,
		Offset:    query.Offset,
	}

	users, total, err := h.userService.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "LIST_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   users,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// Get returns a single user
func (h *UserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrUserNotFound {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "FETCH_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Create creates a new user
func (h *UserHandler) Create(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	tenantID := middleware.GetTenantID(c)

	user, err := h.userService.Create(c.Request.Context(), tenantID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrUserAlreadyExists {
			status = http.StatusConflict
		} else if err == domain.ErrPasswordTooWeak {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "CREATE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Update updates a user
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	user, err := h.userService.Update(c.Request.Context(), id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrUserNotFound {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "UPDATE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete soft deletes a user
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrUserNotFound {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "DELETE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// Activate activates a user
func (h *UserHandler) Activate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	if err := h.userService.Activate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "ACTIVATE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user activated"})
}

// Deactivate deactivates a user
func (h *UserHandler) Deactivate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	if err := h.userService.Deactivate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "DEACTIVATE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deactivated"})
}

// Invite invites a new user to a company
func (h *UserHandler) Invite(c *gin.Context) {
	var req domain.InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	tenantID := middleware.GetTenantID(c)

	user, err := h.userService.InviteUserToCompany(c.Request.Context(), tenantID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrUserAlreadyInCompany {
			status = http.StatusConflict
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "INVITE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": user.ID,
		"email":   user.Email,
		"message": "invitation sent",
	})
}
