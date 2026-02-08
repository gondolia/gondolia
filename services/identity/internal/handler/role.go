package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
	"github.com/gondolia/gondolia/services/identity/internal/middleware"
	"github.com/gondolia/gondolia/services/identity/internal/service"
)

// RoleHandler handles role endpoints
type RoleHandler struct {
	roleService *service.RoleService
}

// NewRoleHandler creates a new role handler
func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// List returns a list of roles
func (h *RoleHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var query struct {
		CompanyID *uuid.UUID `form:"company_id"`
		IsSystem  *bool      `form:"is_system"`
		Limit     int        `form:"limit,default=50"`
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

	filter := domain.RoleFilter{
		TenantID:  tenantID,
		CompanyID: query.CompanyID,
		IsSystem:  query.IsSystem,
		Limit:     query.Limit,
		Offset:    query.Offset,
	}

	roles, total, err := h.roleService.List(c.Request.Context(), filter)
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
		"data":   roles,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// Get returns a single role
func (h *RoleHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid role ID",
			},
		})
		return
	}

	role, err := h.roleService.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrRoleNotFound {
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

	c.JSON(http.StatusOK, role)
}

// Create creates a new role
func (h *RoleHandler) Create(c *gin.Context) {
	var req domain.CreateRoleRequest
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

	role, err := h.roleService.Create(c.Request.Context(), tenantID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrRoleAlreadyExists {
			status = http.StatusConflict
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "CREATE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// Update updates a role
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid role ID",
			},
		})
		return
	}

	var req domain.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	role, err := h.roleService.Update(c.Request.Context(), id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrRoleNotFound {
			status = http.StatusNotFound
		} else if err == domain.ErrRoleIsSystem {
			status = http.StatusForbidden
		} else if err == domain.ErrRoleAlreadyExists {
			status = http.StatusConflict
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "UPDATE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, role)
}

// Delete deletes a role
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid role ID",
			},
		})
		return
	}

	if err := h.roleService.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrRoleNotFound {
			status = http.StatusNotFound
		} else if err == domain.ErrRoleIsSystem {
			status = http.StatusForbidden
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
