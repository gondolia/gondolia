package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LivenessHandler returns OK if the service is running
func LivenessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// ReadinessHandler returns OK if the service is ready to receive traffic
func ReadinessHandler(c *gin.Context) {
	// TODO: Check database connection, redis, etc.
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// MetricsHandler returns Prometheus metrics
func MetricsHandler(c *gin.Context) {
	// TODO: Implement Prometheus metrics
	c.String(http.StatusOK, "# Metrics placeholder")
}
