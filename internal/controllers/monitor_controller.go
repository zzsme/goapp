package controllers

import (
	"goapp/internal/context"
	"goapp/internal/services"

	"github.com/gin-gonic/gin"
)

// MonitorController 监控控制器
type MonitorController struct {
	monitorService *services.MonitorService
}

// NewMonitorController 创建监控控制器
func NewMonitorController() *MonitorController {
	return &MonitorController{
		monitorService: services.NewMonitorService(),
	}
}

// Register 注册路由
func (mc *MonitorController) Register(router *gin.RouterGroup) {
	router.GET("/stats", mc.GetStats)
}

// GetStats 获取系统统计信息
func (mc *MonitorController) GetStats(c *gin.Context) {
	apiCtx := context.GetAPIContext(c)
	stats := mc.monitorService.GetStats()
	apiCtx.Success(stats)
}
