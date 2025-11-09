package middleware

import (
	"goapp/internal/app"
	"goapp/internal/context"
	"goapp/internal/events"
	"time"

	"github.com/gin-gonic/gin"
)

// EventPayload represents the data structure for HTTP request events
type EventPayload struct {
	Method     string
	Path       string
	StatusCode int
	RequestID  string
	UserID     interface{}
	IP         string
	Latency    time.Duration
	Error      interface{}
	Time       time.Time
}

// HTTP request event types
const (
	RequestStarted  events.EventType = "http.request_started"
	RequestComplete events.EventType = "http.request_complete"
	RequestError    events.EventType = "http.request_error"
	RateLimited     events.EventType = "http.rate_limited"
	AuthFailed      events.EventType = "http.auth_failed"
)

// EventsMiddleware emits events for HTTP requests
func EventsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip certain paths that don't need events
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		// Get start time
		startTime := time.Now()
		requestID := context.GetRequestID(c)

		// Get user ID if available
		var userID interface{}
		if id, exists := c.Get("user_id"); exists {
			userID = id
		}

		// Emit request started event
		events.Publish(RequestStarted, EventPayload{
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			RequestID: requestID,
			UserID:    userID,
			IP:        c.ClientIP(),
			Time:      startTime,
		})

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get response status
		statusCode := c.Writer.Status()

		// Create event payload
		payload := EventPayload{
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			StatusCode: statusCode,
			RequestID:  requestID,
			UserID:     userID,
			IP:         c.ClientIP(),
			Latency:    latency,
			Time:       startTime,
		}

		// Get any errors
		if len(c.Errors) > 0 {
			payload.Error = c.Errors.Last().Err

			// Emit request error event
			events.Publish(RequestError, payload)

			app.ErrorContext(c, "Request error",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", statusCode,
				"latency", latency,
				"error", payload.Error,
			)
		} else if statusCode >= 400 && statusCode < 500 {
			if statusCode == 401 {
				events.Publish(AuthFailed, payload)
			} else if statusCode == 429 {
				events.Publish(RateLimited, payload)
			} else {
				// Client errors
				app.WarnContext(c, "Client error",
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"status", statusCode,
					"latency", latency,
				)
			}
		}

		// Always emit request complete event
		events.Publish(RequestComplete, payload)

		// For excessive latency, log a warning
		if latency > 500*time.Millisecond {
			app.WarnContext(c, "Slow request",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", statusCode,
				"latency", latency,
			)
		}
	}
}
