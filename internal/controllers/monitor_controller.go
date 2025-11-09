package controllers

import (
	"goapp/internal/app/errors"
	"goapp/internal/context"
	"goapp/internal/services"

	"github.com/gin-gonic/gin"
)

// MonitorController handles monitoring and metrics endpoints
type MonitorController struct {
	monitorService *services.MonitorService
}

// NewMonitorController creates a new monitor controller
func NewMonitorController() *MonitorController {
	return &MonitorController{
		monitorService: services.NewMonitorService(),
	}
}

// Register registers routes for the controller
func (mc *MonitorController) Register(router *gin.RouterGroup) {
	metrics := router.Group("/metrics")
	{
		metrics.GET("", mc.GetMetrics)
		metrics.GET("/routes", mc.GetRouteMetrics)
		metrics.GET("/errors", mc.GetErrorMetrics)
	}
}

// GetMetrics returns overall system metrics
func (mc *MonitorController) GetMetrics(c *gin.Context) {
	apiCtx := context.GetAPIContext(c)
	stats := mc.monitorService.GetStats()

	// Format the response with only the top-level metrics
	apiCtx.Success(gin.H{
		"uptime":         stats["uptime"],
		"total_requests": stats["total_requests"],
		"error_count":    stats["error_count"],
		"error_rate":     stats["error_rate"],
		"goroutines":     stats["goroutines"],
	})
}

// GetRouteMetrics returns route-specific metrics
func (mc *MonitorController) GetRouteMetrics(c *gin.Context) {
	apiCtx := context.GetAPIContext(c)
	stats := mc.monitorService.GetStats()

	if routes, ok := stats["routes"].([]services.RouteStat); ok {
		apiCtx.Success(gin.H{
			"routes": routes,
		})
	} else {
		apiCtx.ErrorWithCode(errors.InternalServer, "Failed to retrieve route metrics")
	}
}

// GetErrorMetrics returns error-specific metrics
func (mc *MonitorController) GetErrorMetrics(c *gin.Context) {
	apiCtx := context.GetAPIContext(c)
	stats := mc.monitorService.GetStats()

	if errorDist, ok := stats["error_distribution"].(map[string]int); ok {
		apiCtx.Success(gin.H{
			"error_distribution": errorDist,
			"error_count":        stats["error_count"],
			"error_rate":         stats["error_rate"],
		})
	} else {
		apiCtx.ErrorWithCode(errors.InternalServer, "Failed to retrieve error metrics")
	}
}
