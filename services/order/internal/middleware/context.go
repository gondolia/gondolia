package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Context keys
const (
	ContextKeyTenantID = "tenant_id"
	ContextKeyUserID   = "user_id"
)

// GetTenantID returns the tenant ID from context
func GetTenantID(c *gin.Context) uuid.UUID {
	id, _ := c.Get(ContextKeyTenantID)
	return id.(uuid.UUID)
}

// GetUserID returns the user ID from context
func GetUserID(c *gin.Context) *uuid.UUID {
	id, exists := c.Get(ContextKeyUserID)
	if !exists {
		return nil
	}
	userID := id.(uuid.UUID)
	return &userID
}

// GetClientIP returns the client IP address
func GetClientIP(c *gin.Context) string {
	// Check X-Forwarded-For first (for proxied requests)
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		return xff
	}
	return c.ClientIP()
}
