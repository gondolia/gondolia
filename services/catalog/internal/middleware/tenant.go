package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gondolia/gondolia/services/catalog/internal/repository"
)

// TenantMiddleware extracts tenant from request and validates it
func TenantMiddleware(tenantRepo repository.TenantRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant from header or subdomain
		tenantCode := c.GetHeader("X-Tenant-ID")
		if tenantCode == "" {
			// Try to extract from host (subdomain.domain.com)
			// For now, require header
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "TENANT_REQUIRED",
					"message": "X-Tenant-ID header is required",
				},
			})
			return
		}

		// Get tenant
		tenant, err := tenantRepo.GetByCode(c.Request.Context(), tenantCode)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_TENANT",
					"message": "invalid tenant",
				},
			})
			return
		}

		if !tenant.IsActive {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "TENANT_INACTIVE",
					"message": "tenant is not active",
				},
			})
			return
		}

		// Set tenant in context
		c.Set("tenant", tenant)
		c.Set(ContextKeyTenantID, tenant.ID)

		c.Next()
	}
}
