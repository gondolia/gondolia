package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LivenessHandler handles liveness probe
func LivenessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// ReadinessHandler handles readiness probe
func ReadinessHandler(c *gin.Context) {
	// TODO: Check database connection, cache, etc.
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// MetricsHandler handles metrics endpoint
func MetricsHandler(c *gin.Context) {
	// TODO: Return Prometheus metrics
	c.String(http.StatusOK, "# Metrics\n")
}
