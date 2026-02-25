package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT token and extracts user_id
// If no token is provided, treats as guest and uses session_id
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get JWT token from Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// Check for internal service-to-service X-User-ID header
			if userIDHeader := c.GetHeader("X-User-ID"); userIDHeader != "" {
				userID, err := uuid.Parse(userIDHeader)
				if err == nil {
					c.Set(ContextKeyUserID, userID)
					// Also pass through session ID if present
					if sessionID := c.GetHeader("X-Session-ID"); sessionID != "" {
						c.Set(ContextKeySessionID, sessionID)
					}
					c.Next()
					return
				}
			}
			// No token and no X-User-ID - treat as guest, get or create session ID
			sessionID := c.GetHeader("X-Session-ID")
			if sessionID == "" {
				// Create new session ID
				sessionID = uuid.New().String()
			}
			c.Set(ContextKeySessionID, sessionID)
			c.Next()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			// Invalid token - treat as guest
			sessionID := c.GetHeader("X-Session-ID")
			if sessionID == "" {
				sessionID = uuid.New().String()
			}
			c.Set(ContextKeySessionID, sessionID)
			c.Next()
			return
		}

		// Extract user_id from claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userIDStr, ok := claims["user_id"].(string); ok {
				userID, err := uuid.Parse(userIDStr)
				if err == nil {
					c.Set(ContextKeyUserID, userID)
				}
			}
		}

		c.Next()
	}
}
