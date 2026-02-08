package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/identity/internal/auth"
	"github.com/gondolia/gondolia/services/identity/internal/domain"
)

// Context keys
const (
	ContextKeyClaims   = "claims"
	ContextKeyUserID   = "user_id"
	ContextKeyTenantID = "tenant_id"
	ContextKeyCompanyID = "company_id"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "missing authorization header",
				},
			})
			return
		}

		// Extract Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "invalid authorization header format",
				},
			})
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "invalid or expired token",
				},
			})
			return
		}

		// Set claims in context
		c.Set(ContextKeyClaims, claims)
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyTenantID, claims.TenantID)
		c.Set(ContextKeyCompanyID, claims.CompanyID)

		c.Next()
	}
}

// RequirePermission middleware checks if user has required permission
func RequirePermission(permission domain.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get(ContextKeyClaims)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "not authenticated",
				},
			})
			return
		}

		accessClaims := claims.(*auth.AccessTokenClaims)

		// Check permission
		hasPermission := false
		for _, p := range accessClaims.Permissions {
			if p == string(permission) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "insufficient permissions",
				},
			})
			return
		}

		c.Next()
	}
}

// GetClaims returns the JWT claims from context
func GetClaims(c *gin.Context) *auth.AccessTokenClaims {
	claims, exists := c.Get(ContextKeyClaims)
	if !exists {
		return nil
	}
	return claims.(*auth.AccessTokenClaims)
}

// GetUserID returns the user ID from context
func GetUserID(c *gin.Context) uuid.UUID {
	id, _ := c.Get(ContextKeyUserID)
	return id.(uuid.UUID)
}

// GetTenantID returns the tenant ID from context
func GetTenantID(c *gin.Context) uuid.UUID {
	id, _ := c.Get(ContextKeyTenantID)
	return id.(uuid.UUID)
}

// GetCompanyID returns the company ID from context
func GetCompanyID(c *gin.Context) uuid.UUID {
	id, _ := c.Get(ContextKeyCompanyID)
	return id.(uuid.UUID)
}

// GetClientIP returns the client IP address
func GetClientIP(c *gin.Context) string {
	// Check X-Forwarded-For first (for proxied requests)
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	return c.ClientIP()
}
