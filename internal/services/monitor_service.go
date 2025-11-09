package services

import (
	"goapp/internal/app"
	"goapp/internal/events"
	"goapp/internal/middleware"
	"runtime"
	"sync"
	"time"
)

// RouteStat represents statistics for a single route
type RouteStat struct {
	Route     string
	Count     int
	AvgTime   time.Duration
	AvgTimeMs float64
}

// MonitorService monitors application metrics and health
type MonitorService struct {
	requestStats      map[string]int
	errorCounts       map[int]int
	latencyHistogram  map[string][]time.Duration
	routeResponseTime map[string]time.Duration
	requestCount      int
	errorCount        int
	startTime         time.Time
	mutex             sync.RWMutex
}

// NewMonitorService creates a new monitoring service
func NewMonitorService() *MonitorService {
	service := &MonitorService{
		requestStats:      make(map[string]int),
		errorCounts:       make(map[int]int),
		latencyHistogram:  make(map[string][]time.Duration),
		routeResponseTime: make(map[string]time.Duration),
		startTime:         time.Now(),
	}

	// Register event handlers
	service.registerEventHandlers()

	// Start periodic metrics reporting
	go service.reportMetrics()

	return service
}

// registerEventHandlers subscribes to relevant events
func (s *MonitorService) registerEventHandlers() {
	// Handle request completion events
	events.Subscribe(middleware.RequestComplete, func(e events.Event) {
		if payload, ok := e.Payload.(middleware.EventPayload); ok {
			s.recordRequest(payload)
		}
	})

	// Handle request error events
	events.Subscribe(middleware.RequestError, func(e events.Event) {
		if payload, ok := e.Payload.(middleware.EventPayload); ok {
			s.recordError(payload)
		}
	})

	// Handle system events
	events.Subscribe(events.DatabaseError, func(e events.Event) {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.errorCount++
		app.Warn("Database error detected by monitor", "details", e.Payload)
	})

	// Handle authentication failures
	events.Subscribe(middleware.AuthFailed, func(e events.Event) {
		if payload, ok := e.Payload.(middleware.EventPayload); ok {
			app.Warn("Authentication failed",
				"ip", payload.IP,
				"path", payload.Path,
				"request_id", payload.RequestID,
			)
		}
	})
}

// recordRequest records request metrics
func (s *MonitorService) recordRequest(payload middleware.EventPayload) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Increment request count
	s.requestCount++

	// Track requests by path
	path := payload.Method + " " + payload.Path
	s.requestStats[path]++

	// Record latency
	if _, exists := s.latencyHistogram[path]; !exists {
		s.latencyHistogram[path] = make([]time.Duration, 0)
	}
	s.latencyHistogram[path] = append(s.latencyHistogram[path], payload.Latency)

	// Update average response time
	s.routeResponseTime[path] = s.calculateAvgLatency(path)
}

// recordError records error metrics
func (s *MonitorService) recordError(payload middleware.EventPayload) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Increment error count
	s.errorCount++

	// Track errors by status code
	s.errorCounts[payload.StatusCode]++
}

// calculateAvgLatency calculates the average latency for a path
func (s *MonitorService) calculateAvgLatency(path string) time.Duration {
	latencies := s.latencyHistogram[path]
	if len(latencies) == 0 {
		return 0
	}

	var total time.Duration
	for _, latency := range latencies {
		total += latency
	}
	return total / time.Duration(len(latencies))
}

// reportMetrics periodically reports system metrics
func (s *MonitorService) reportMetrics() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.RLock()
		uptime := time.Since(s.startTime).Round(time.Second)
		requestCount := s.requestCount
		errorCount := s.errorCount
		s.mutex.RUnlock()

		// Collect runtime metrics
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		// Log metrics
		app.Info("System metrics",
			"uptime", uptime,
			"goroutines", runtime.NumGoroutine(),
			"memory_used_mb", m.Alloc/1024/1024,
			"total_requests", requestCount,
			"error_count", errorCount,
		)

		// Publish metrics event
		events.Publish(events.EventType("system.metrics"), map[string]interface{}{
			"uptime":         uptime.String(),
			"goroutines":     runtime.NumGoroutine(),
			"memory_used_mb": m.Alloc / 1024 / 1024,
			"total_requests": requestCount,
			"error_rate":     float64(errorCount) / float64(requestCount+1) * 100,
			"timestamp":      time.Now(),
		})
	}
}

// GetStats returns current monitoring statistics
func (s *MonitorService) GetStats() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	routes := make([]RouteStat, 0, len(s.requestStats))
	for route, count := range s.requestStats {
		routes = append(routes, RouteStat{
			Route:     route,
			Count:     count,
			AvgTime:   s.routeResponseTime[route],
			AvgTimeMs: float64(s.routeResponseTime[route]) / float64(time.Millisecond),
		})
	}

	// Get error distribution
	errorDist := make(map[string]int)
	for code, count := range s.errorCounts {
		if code >= 500 {
			errorDist["server_errors"] += count
		} else if code >= 400 {
			errorDist["client_errors"] += count
		}
	}

	return map[string]interface{}{
		"uptime":             time.Since(s.startTime).Round(time.Second).String(),
		"total_requests":     s.requestCount,
		"error_count":        s.errorCount,
		"error_rate":         float64(s.errorCount) / float64(s.requestCount+1) * 100,
		"routes":             routes,
		"error_distribution": errorDist,
		"goroutines":         runtime.NumGoroutine(),
		"timestamp":          time.Now(),
	}
}
