package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/identity/internal/domain"
	"github.com/gondolia/gondolia/services/identity/internal/middleware"
	"github.com/gondolia/gondolia/services/identity/internal/service"
)

// CompanyHandler handles company endpoints
type CompanyHandler struct {
	companyService *service.CompanyService
}

// NewCompanyHandler creates a new company handler
func NewCompanyHandler(companyService *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: companyService}
}

// List returns a list of companies
func (h *CompanyHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var query struct {
		SAPCompanyNumber *string `form:"sap_company_number"`
		IsActive         *bool   `form:"is_active"`
		Search           *string `form:"search"`
		Limit            int     `form:"limit,default=20"`
		Offset           int     `form:"offset,default=0"`
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

	filter := domain.CompanyFilter{
		TenantID:         tenantID,
		SAPCompanyNumber: query.SAPCompanyNumber,
		IsActive:         query.IsActive,
		Search:           query.Search,
		Limit:            query.Limit,
		Offset:           query.Offset,
	}

	companies, total, err := h.companyService.List(c.Request.Context(), filter)
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
		"data":   companies,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// Get returns a single company
func (h *CompanyHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid company ID",
			},
		})
		return
	}

	company, err := h.companyService.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrCompanyNotFound {
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

	c.JSON(http.StatusOK, company)
}

// Create creates a new company
func (h *CompanyHandler) Create(c *gin.Context) {
	var req domain.CreateCompanyRequest
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

	company, err := h.companyService.Create(c.Request.Context(), tenantID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrCompanyAlreadyExists {
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

	c.JSON(http.StatusCreated, company)
}

// Update updates a company
func (h *CompanyHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid company ID",
			},
		})
		return
	}

	var req domain.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	company, err := h.companyService.Update(c.Request.Context(), id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrCompanyNotFound {
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

	c.JSON(http.StatusOK, company)
}

// Delete soft deletes a company
func (h *CompanyHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid company ID",
			},
		})
		return
	}

	if err := h.companyService.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrCompanyNotFound {
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

// ListUsers returns users assigned to a company
func (h *CompanyHandler) ListUsers(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid company ID",
			},
		})
		return
	}

	users, err := h.companyService.ListUsers(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "LIST_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// AddUser adds a user to a company
func (h *CompanyHandler) AddUser(c *gin.Context) {
	companyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid company ID",
			},
		})
		return
	}

	var req domain.AddUserToCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	if err := h.companyService.AddUser(c.Request.Context(), companyID, req.UserID, req.RoleID, req.UserType); err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrUserAlreadyInCompany {
			status = http.StatusConflict
		} else if err == domain.ErrUserNotFound || err == domain.ErrCompanyNotFound || err == domain.ErrRoleNotFound {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "ADD_USER_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user added to company"})
}

// UpdateUserRole updates a user's role in a company
func (h *CompanyHandler) UpdateUserRole(c *gin.Context) {
	companyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid company ID",
			},
		})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	var req struct {
		RoleID uuid.UUID `json:"role_id" binding:"required"`
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

	if err := h.companyService.UpdateUserRole(c.Request.Context(), companyID, userID, req.RoleID); err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrUserNotInCompany {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "UPDATE_ROLE_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user role updated"})
}

// RemoveUser removes a user from a company
func (h *CompanyHandler) RemoveUser(c *gin.Context) {
	companyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid company ID",
			},
		})
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid user ID",
			},
		})
		return
	}

	if err := h.companyService.RemoveUser(c.Request.Context(), companyID, userID); err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrUserNotInCompany {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "REMOVE_USER_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}
