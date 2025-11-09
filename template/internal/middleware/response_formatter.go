package middleware

import (
	"net/http"
	"time"

	"gowk/internal/context"
	"gowk/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ResponseFormatter is a middleware that formats API responses
func ResponseFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Create an API context
		apiCtx := &context.APIContext{
			GinContext: c,
			RequestID:  requestID,
		}

		// Store it in the Gin context
		c.Set("apiContext", apiCtx)

		// Set the request ID header
		c.Writer.Header().Set("X-Request-ID", requestID)

		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Log the request if an error occurred
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Log the error
				// app.Error("Request error", "path", c.Request.URL.Path, "error", e.Error())

				// If no response was sent yet, send an error response
				if !c.Writer.Written() {
					c.JSON(http.StatusInternalServerError, dto.Response{
						Code:      5000, // Internal server error
						Message:   "An unexpected error occurred",
						RequestID: requestID,
					})
				}
			}
			return
		}

		// If no response was sent by the handler, format and send one
		if !c.Writer.Written() {
			// Get the response from the context if it exists
			var responseData interface{}
			if data, exists := c.Get("response"); exists {
				responseData = data
			}

			c.JSON(http.StatusOK, dto.Response{
				Code:      0,
				Message:   "success",
				Data:      responseData,
				RequestID: requestID,
			})
		}

		// Log request duration for monitoring
		duration := time.Since(startTime)
		if duration > time.Second*1 {
			// Log slow requests
			// app.Warn("Slow request", "path", c.Request.URL.Path, "method", c.Request.Method, "duration", duration.String())
		}
	}
}
