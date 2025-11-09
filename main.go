package main

import (
	"fmt"
	"goapp/internal/app"
	"goapp/internal/events"
	"goapp/internal/router"
	"goapp/internal/services"
	"goapp/internal/tasks"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultPort = 8080
)

func main() {
	// Initialize core components
	app.InitLogger()
	fmt.Println("Logger initialized successfully")

	// Initialize event system
	events.InitEventBus()
	fmt.Println("Event system initialized successfully")

	// Try to initialize database (but continue if it fails)
	tryInitialize("Database", func() {
		app.InitDB()
	})

	// Try to initialize Redis (but continue if it fails)
	tryInitialize("Redis", func() {
		app.InitRedis()
	})

	// Initialize validator
	app.InitValidator()
	fmt.Println("Validator initialized successfully")

	// Handle command-line tasks
	args := os.Args
	if len(args) > 1 && args[1] == "task" {
		if len(args) < 3 {
			fmt.Println("Please specify a task name, e.g.: go run main.go task cleanup")
			return
		}
		taskName := args[2]
		tasks.RunTask(taskName)
		return
	}

	// Initialize monitoring service
	monitor := services.NewMonitorService()

	// Emit system start event
	events.Publish(events.SystemStarted, map[string]interface{}{
		"port": defaultPort,
		"mode": app.ConfigData.Server.Mode,
	})

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the web server
	r := router.SetupRouter()
	port := defaultPort

	// Start server in a goroutine
	go func() {
		fmt.Printf("🚀 Server starting on port %d\n", port)
		if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
			app.Error("Server failed to start", "error", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	fmt.Println("\n🛑 Shutting down server...")

	// Emit system shutdown event
	events.Publish(events.SystemShutdown, map[string]interface{}{
		"time": fmt.Sprintf("%v", app.ConfigData.Server.Mode),
	})

	// Print final stats
	stats := monitor.GetStats()
	fmt.Printf("📊 Final Stats:\n")
	fmt.Printf("- Uptime: %v\n", stats["uptime"])
	fmt.Printf("- Total Requests: %v\n", stats["total_requests"])
	fmt.Printf("- Error Rate: %.2f%%\n", stats["error_rate"])
}

// tryInitialize attempts to initialize a component but continues if it fails
func tryInitialize(componentName string, initFunc func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("WARNING: Failed to initialize %s: %v\n", componentName, r)
			fmt.Printf("The application will continue without %s functionality\n", componentName)
		} else {
			fmt.Printf("%s initialized successfully\n", componentName)
		}
	}()

	initFunc()
}
