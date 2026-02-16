package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/catalog/internal/domain"
	"github.com/gondolia/gondolia/services/catalog/internal/middleware"
	"github.com/gondolia/gondolia/services/catalog/internal/service"
)

// AttributeTranslationHandler handles attribute translation endpoints
type AttributeTranslationHandler struct {
	service *service.AttributeTranslationService
}

// NewAttributeTranslationHandler creates a new attribute translation handler
func NewAttributeTranslationHandler(service *service.AttributeTranslationService) *AttributeTranslationHandler {
	return &AttributeTranslationHandler{
		service: service,
	}
}

// List handles GET /attribute-translations
func (h *AttributeTranslationHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	filter := domain.AttributeTranslationFilter{
		TenantID: tenantID,
		Limit:    100,
		Offset:   0,
	}

	// Parse query parameters
	if c.Query("attribute_key") != "" {
		key := c.Query("attribute_key")
		filter.AttributeKey = &key
	}

	if c.Query("locale") != "" {
		locale := c.Query("locale")
		filter.Locale = &locale
	}

	// Pagination
	if limit := parseInt(c.Query("limit"), 100); limit > 0 {
		filter.Limit = limit
	}
	if offset := parseInt(c.Query("offset"), 0); offset >= 0 {
		filter.Offset = offset
	}

	translations, total, err := h.service.List(c.Request.Context(), filter)
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
		"data":   translations,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GetByLocale handles GET /attribute-translations/by-locale/:locale
// Returns all translations for a specific locale as a map (attribute_key -> translation)
func (h *AttributeTranslationHandler) GetByLocale(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	locale := c.Param("locale")

	if locale == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_LOCALE",
				"message": "locale is required",
			},
		})
		return
	}

	translations, err := h.service.GetByTenantAndLocale(c.Request.Context(), tenantID, locale)
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
		"data": translations,
	})
}

// Create handles POST /attribute-translations
func (h *AttributeTranslationHandler) Create(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req domain.CreateAttributeTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	translation, err := h.service.Create(c.Request.Context(), tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, translation)
}

// Update handles PUT /attribute-translations/:id
func (h *AttributeTranslationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid translation ID",
			},
		})
		return
	}

	var req domain.UpdateAttributeTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	translation, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, translation)
}

// Delete handles DELETE /attribute-translations/:id
func (h *AttributeTranslationHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "invalid translation ID",
			},
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
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
