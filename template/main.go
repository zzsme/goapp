package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gowk/internal/app"
	"gowk/internal/events"
	"gowk/internal/router"
	"gowk/internal/tasks"
)

const (
	defaultPort = 8080
)

func main() {
	// 初始化核心组件
	app.InitLogger()
	fmt.Println("Logger initialized successfully")

	// 初始化事件系统
	events.InitEventBus()
	fmt.Println("Event system initialized successfully")

	// 尝试初始化数据库（失败不会中断启动）
	tryInitialize("Database", func() {
		app.InitDB()
	})

	// 尝试初始化Redis（失败不会中断启动）
	tryInitialize("Redis", func() {
		app.InitRedis()
	})

	// 初始化验证器
	app.InitValidator()
	fmt.Println("Validator initialized successfully")

	// 处理命令行任务
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

	// 发布系统启动事件
	events.Publish(events.SystemStarted, map[string]interface{}{
		"port": defaultPort,
		"mode": app.ConfigData.Server.Mode,
	})

	// 优雅关闭设置
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动Web服务器
	r := router.SetupRouter()
	port := defaultPort

	// 在goroutine中启动服务器
	go func() {
		fmt.Printf("🚀 Server starting on port %d\n", port)
		if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
			app.Error("Server failed to start", "error", err)
		}
	}()

	// 等待关闭信号
	<-quit
	fmt.Println("\n🛑 Shutting down server...")

	// 发布系统关闭事件
	events.Publish(events.SystemShutdown, map[string]interface{}{
		"time": fmt.Sprintf("%v", app.Now()),
	})

	// 打印最终统计信息
	fmt.Printf("📊 Server shutdown complete\n")
}

// tryInitialize尝试初始化组件，但如果失败则继续执行
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
