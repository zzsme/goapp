package router

import (
	"goapp/internal/app"
	"goapp/internal/context"
	"goapp/internal/controllers"
	"goapp/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the router and registers all routes
func SetupRouter() *gin.Engine {
	// Set Gin mode based on config
	gin.SetMode(app.ConfigData.Server.Mode)

	// Create router with default middleware
	router := gin.New()

	// Add custom middleware in correct order
	router.Use(middleware.LoggerMiddleware())   // First to set request ID
	router.Use(middleware.RecoveryMiddleware()) // Then recovery
	router.Use(middleware.CORSMiddleware())     // Then CORS
	router.Use(middleware.EventsMiddleware())   // Then event tracking
	router.Use(middleware.ResponseFormatter())  // Finally response formatting

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		apiCtx := context.GetAPIContext(c)
		apiCtx.Success(gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.GET("/ping", func(c *gin.Context) {
			apiCtx := context.GetAPIContext(c)
			apiCtx.Success(gin.H{"message": "pong"})
		})

		// User routes
		userController := controllers.NewUserController()
		users := v1.Group("/users")
		{
			// Public user endpoints
			users.POST("/register", userController.CreateUser)
			users.POST("/login", userController.Login)

			// Protected user endpoints
			protected := users.Use(middleware.AuthMiddleware())
			{
				protected.GET("", userController.ListUsers)
				protected.GET("/:id", userController.GetUser)
				protected.PUT("/:id", userController.UpdateUser)
				protected.PUT("/:id/password", userController.UpdatePassword)

				// Admin-only endpoints
				admin := protected.Use(middleware.AdminMiddleware())
				{
					admin.DELETE("/:id", userController.DeleteUser)
				}
			}
		}

		// Product routes
		productController := controllers.NewProductController()
		productController.Register(v1)

		// Monitoring routes (admin only)
		monitorController := controllers.NewMonitorController()
		admin := v1.Group("/admin")
		adminProtected := admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		monitorController.Register(adminProtected.(*gin.RouterGroup))
	}

	return router
}
