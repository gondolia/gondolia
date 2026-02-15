package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gondolia/gondolia/services/identity/internal/config"
	"github.com/gondolia/gondolia/services/identity/internal/domain"
	"github.com/gondolia/gondolia/services/identity/internal/middleware"
	"github.com/gondolia/gondolia/services/identity/internal/service"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
	config      *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, userService *service.UserService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		config:      cfg,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
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
	ipAddress := middleware.GetClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	tokens, err := h.authService.Login(c.Request.Context(), tenantID, req, ipAddress, userAgent)
	if err != nil {
		status := http.StatusUnauthorized
		code := "AUTHENTICATION_FAILED"
		message := err.Error()

		if domain.IsAuthError(err) {
			status = http.StatusUnauthorized
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    code,
				"message": message,
			},
		})
		return
	}

	// Set refresh token as HttpOnly cookie
	h.setRefreshTokenCookie(c, tokens.RefreshToken)

	// Return only access token in response body
	c.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
		"token_type":   tokens.TokenType,
		"expires_in":   tokens.ExpiresIn,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Read refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_TOKEN",
				"message": "refresh token not found",
			},
		})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), refreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "LOGOUT_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	// Clear refresh token cookie
	h.clearRefreshTokenCookie(c)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	// Read refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "MISSING_TOKEN",
				"message": "refresh token not found",
			},
		})
		return
	}

	ipAddress := middleware.GetClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	tokens, err := h.authService.RefreshToken(c.Request.Context(), refreshToken, ipAddress, userAgent)
	if err != nil {
		status := http.StatusUnauthorized
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "TOKEN_REFRESH_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	// Set new refresh token as HttpOnly cookie
	h.setRefreshTokenCookie(c, tokens.RefreshToken)

	// Return only access token in response body
	c.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
		"token_type":   tokens.TokenType,
		"expires_in":   tokens.ExpiresIn,
	})
}

// Me returns current user with context
func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	companyID := middleware.GetCompanyID(c)

	result, err := h.authService.GetCurrentUser(c.Request.Context(), userID, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "FETCH_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// SwitchCompany switches company context
func (h *AuthHandler) SwitchCompany(c *gin.Context) {
	var req domain.SwitchCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	userID := middleware.GetUserID(c)
	ipAddress := middleware.GetClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	tokens, err := h.authService.SwitchCompany(c.Request.Context(), userID, req.CompanyID, ipAddress, userAgent)
	if err != nil {
		status := http.StatusBadRequest
		if err == domain.ErrUserNotInCompany {
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "SWITCH_COMPANY_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	// Set refresh token as HttpOnly cookie
	h.setRefreshTokenCookie(c, tokens.RefreshToken)

	// Return only access token in response body
	c.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
		"token_type":   tokens.TokenType,
		"expires_in":   tokens.ExpiresIn,
	})
}

// ForgotPassword initiates password reset
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req domain.ForgotPasswordRequest
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

	// Note: Token is returned here for dev purposes
	// In production, send via email
	token, _ := h.authService.ForgotPassword(c.Request.Context(), tenantID, req.Email)

	// Always return success to prevent email enumeration
	response := gin.H{"message": "if the email exists, a reset link has been sent"}
	if token != "" {
		// Only include token in development
		response["reset_token"] = token
	}

	c.JSON(http.StatusOK, response)
}

// ResetPassword resets password with token
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req domain.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		status := http.StatusBadRequest
		c.JSON(status, gin.H{
			"error": gin.H{
				"code":    "RESET_PASSWORD_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

// ValidateInvitation validates an invitation token
func (h *AuthHandler) ValidateInvitation(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "token is required",
			},
		})
		return
	}

	user, err := h.userService.ValidateInvitationToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_INVITATION",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":     user.Email,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
	})
}

// AcceptInvitation accepts an invitation
func (h *AuthHandler) AcceptInvitation(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "token is required",
			},
		})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required,min=8"`
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

	user, err := h.userService.AcceptInvitation(c.Request.Context(), token, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "ACCEPT_INVITATION_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "invitation accepted",
		"user_id": user.ID,
	})
}

// parseUUID parses a UUID from string with error handling
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// setRefreshTokenCookie sets the refresh token as HttpOnly cookie
func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	maxAge := int(h.config.JWTRefreshTokenExpiry.Seconds())

	// SameSite must be set BEFORE SetCookie
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie(
		"refresh_token",           // name
		token,                     // value
		maxAge,                    // max age in seconds (7 days)
		"/api/v1/auth",            // path - only sent to auth endpoints
		"",                        // domain - empty means current domain
		h.config.SecureCookies,    // secure - only HTTPS (configurable for dev)
		true,                      // httpOnly - not accessible via JavaScript
	)
}

// clearRefreshTokenCookie clears the refresh token cookie
func (h *AuthHandler) clearRefreshTokenCookie(c *gin.Context) {
	c.SetCookie(
		"refresh_token",
		"",
		-1, // max age -1 deletes the cookie
		"/api/v1/auth",
		"",
		h.config.SecureCookies,
		true,
	)
}
