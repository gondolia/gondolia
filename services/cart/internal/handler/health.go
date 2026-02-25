package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LivenessHandler handles liveness probe
func LivenessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

// ReadinessHandler handles readiness probe
func ReadinessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// MetricsHandler handles metrics endpoint
func MetricsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"metrics": "placeholder",
	})
}
