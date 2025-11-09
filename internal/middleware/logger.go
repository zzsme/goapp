package middleware

import (
	"fmt"
	"time"

	"goapp/internal/app"
	"goapp/internal/context"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware is a middleware that logs HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get or generate request ID using our utility
		requestID := context.GetRequestID(c)

		// Add request ID to response headers
		c.Header("X-Request-ID", requestID)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		if raw != "" {
			path = path + "?" + raw
		}

		// Process request
		c.Next()

		// Stop timer
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		// Get error if any
		var errorMsg string
		if len(c.Errors) > 0 {
			errorMsg = c.Errors.String()
		}

		if statusCode >= 500 {
			app.ErrorContext(c, "HTTP Request",
				"status", statusCode,
				"latency", latency,
				"client_ip", clientIP,
				"method", method,
				"path", path,
				"error", errorMsg,
			)
		} else if statusCode >= 400 {
			app.WarnContext(c, "HTTP Request",
				"status", statusCode,
				"latency", latency,
				"client_ip", clientIP,
				"method", method,
				"path", path,
				"error", errorMsg,
			)
		} else {
			app.InfoContext(c, "HTTP Request",
				"status", statusCode,
				"latency", latency,
				"client_ip", clientIP,
				"method", method,
				"path", path,
			)
		}

		// Also log to console for development environment
		if app.ConfigData.Server.Mode == "debug" {
			fmt.Printf("%s | %3d | %13v | %15s | %-7s %s\n",
				time.Now().Format("2006/01/02 15:04:05"),
				statusCode,
				latency,
				clientIP,
				method,
				path,
			)
		}
	}
}

// RecoveryMiddleware recovers from any panics and logs the panic
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				app.ErrorContext(c, "Panic recovered",
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)

				// Get API context for standardized response
				apiCtx := context.GetAPIContext(c)
				apiCtx.Error(5000, "Internal server error")
				c.Abort()
			}
		}()

		c.Next()
	}
}
