package middleware

import (
	"goapp/internal/context"
	"goapp/internal/dto"

	"github.com/gin-gonic/gin"
)

// ResponseFormatter is a middleware that ensures all responses follow the standard format
func ResponseFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get or create request ID before processing
		requestID := context.GetRequestID(c)
		c.Header("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Skip if response is already handled (e.g., by APIContext)
		if c.Writer.Written() {
			return
		}

		// Get response data from context
		data, exists := c.Get("response")
		if !exists {
			// If no response data is set, return empty success response
			c.JSON(200, dto.StandardResponse{
				Errno:     0,
				Errmsg:    "",
				Data:      nil,
				RequestID: requestID,
			})
			return
		}

		// Return standardized response
		c.JSON(c.Writer.Status(), dto.StandardResponse{
			Errno:     0,
			Errmsg:    "",
			Data:      data,
			RequestID: requestID,
		})
	}
}
