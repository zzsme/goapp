package context

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "X-Request-ID"

// GetRequestID returns the request ID from the context
// If no request ID exists, it generates a new one
func GetRequestID(c *gin.Context) string {
	// Check if we already have a request ID in the context
	if requestID, exists := c.Get(requestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}

	// Check if the request ID is in the headers
	if requestID := c.GetHeader(requestIDKey); requestID != "" {
		// Store for future use
		c.Set(requestIDKey, requestID)
		return requestID
	}

	// Generate a new request ID
	requestID := uuid.New().String()
	c.Set(requestIDKey, requestID)
	c.Header(requestIDKey, requestID)
	return requestID
}
