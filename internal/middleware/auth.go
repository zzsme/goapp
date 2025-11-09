package middleware

import (
	"strings"

	"goapp/internal/app"
	"goapp/internal/app/errors"
	"goapp/internal/context"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles authentication for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiCtx := context.GetAPIContext(c)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			app.WarnContext(c, "Missing authorization header")
			apiCtx.ErrorWithCode(errors.Unauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		// Extract token from Bearer header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.WarnContext(c, "Invalid authorization format", "header", authHeader)
			apiCtx.ErrorWithCode(errors.Unauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]

		// In a real application, you would validate the JWT token here
		// For now, we'll just check if it's our dummy token
		if token != "dummy-jwt-token" {
			app.ErrorContext(c, "Invalid token", "token", token)
			apiCtx.ErrorWithCode(errors.Unauthorized, "Invalid token")
			c.Abort()
			return
		}

		// You could set user information in the context here
		// c.Set("user_id", claims.UserID)

		c.Next()
	}
}

// AdminMiddleware ensures the user is an admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiCtx := context.GetAPIContext(c)

		// In a real application, you would get the user from the JWT claims
		// and check their admin status

		// For demonstration purposes, we'll just check a header
		isAdmin := c.GetHeader("X-Is-Admin")
		if isAdmin != "true" {
			app.ErrorContext(c, "Unauthorized admin access attempt")
			apiCtx.ErrorWithCode(errors.Forbidden, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
