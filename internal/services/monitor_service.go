package services

import (
	"sync"
	"time"
)

// MonitorService 监控服务（简化版）
type MonitorService struct {
	startTime      time.Time
	totalRequests  int64
	errorCount     int64
	mu             sync.RWMutex
}

// NewMonitorService 创建监控服务
func NewMonitorService() *MonitorService {
	return &MonitorService{
		startTime: time.Now(),
	}
}

// GetStats 获取统计信息
func (m *MonitorService) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	uptime := time.Since(m.startTime)
	errorRate := 0.0
	if m.totalRequests > 0 {
		errorRate = float64(m.errorCount) / float64(m.totalRequests) * 100
	}

	return map[string]interface{}{
		"uptime":         uptime,
		"total_requests": m.totalRequests,
		"error_count":    m.errorCount,
		"error_rate":     errorRate,
	}
}
