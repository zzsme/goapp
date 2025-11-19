package controllers

import (
	"goapp/internal/app"
	"goapp/internal/application"
	"goapp/internal/infrastructure/persistence"
	httpcontroller "goapp/internal/interfaces/http"
)

// NewUserController 创建用户控制器（适配器模式）
func NewUserController() *httpcontroller.UserController {
	// 初始化仓储
	userRepo := persistence.NewUserRepository(app.GetDB())
	
	// 初始化应用服务
	userService := application.NewUserService(userRepo)
	
	// 创建HTTP控制器
	return httpcontroller.NewUserController(userService)
}
