package http

import (
	"strconv"

	"goapp/internal/app"
	"goapp/internal/app/errors"
	"goapp/internal/application"
	"goapp/internal/context"
	"goapp/internal/dto"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器 - 处理HTTP请求
type UserController struct {
	userService *application.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService *application.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Register 注册路由
func (c *UserController) Register(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("", c.CreateUser)
		users.GET("", c.ListUsers)
		users.GET("/:id", c.GetUser)
		users.PUT("/:id", c.UpdateUser)
		users.DELETE("/:id", c.DeleteUser)
		users.POST("/login", c.Login)
		users.PUT("/:id/password", c.UpdatePassword)
	}
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	var req dto.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiCtx.ErrorWithCode(errors.Validation, err.Error())
		return
	}

	// 调用应用服务
	userEntity, err := c.userService.CreateUser(
		req.Username,
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		app.ErrorContext(ctx, "Failed to create user", "error", err)
		apiCtx.ErrorWithCode(errors.BadRequest, err.Error())
		return
	}

	// 转换为响应DTO
	response := dto.UserResponse{
		ID:        userEntity.ID(),
		Username:  userEntity.Username(),
		Email:     userEntity.EmailValue(),
		FirstName: userEntity.FirstName(),
		LastName:  userEntity.LastName(),
		IsActive:  userEntity.IsActive(),
		IsAdmin:   userEntity.IsAdmin(),
	}

	apiCtx.Success(response)
}

// ListUsers 获取用户列表
func (c *UserController) ListUsers(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	var pagination dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, err.Error())
		return
	}

	// 调用应用服务
	users, err := c.userService.ListUsers(pagination.Page, pagination.PageSize)
	if err != nil {
		app.ErrorContext(ctx, "Failed to list users", "error", err)
		apiCtx.ErrorWithCode(errors.InternalServer, "Failed to retrieve users")
		return
	}

	// 转换为响应DTO列表
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID(),
			Username:  user.Username(),
			Email:     user.EmailValue(),
			FirstName: user.FirstName(),
			LastName:  user.LastName(),
			IsActive:  user.IsActive(),
			IsAdmin:   user.IsAdmin(),
		}
	}

	response := dto.UsersListResponse{
		Users:      userResponses,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalItems: len(users),
		TotalPages: (len(users) + pagination.PageSize - 1) / pagination.PageSize,
	}

	apiCtx.Success(response)
}

// GetUser 获取单个用户
func (c *UserController) GetUser(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, "Invalid user ID")
		return
	}

	// 调用应用服务
	userEntity, err := c.userService.GetUser(id)
	if err != nil {
		app.ErrorContext(ctx, "Failed to get user", "error", err, "id", id)
		apiCtx.ErrorWithCode(errors.NotFound, "User not found")
		return
	}

	// 转换为响应DTO
	response := dto.UserResponse{
		ID:        userEntity.ID(),
		Username:  userEntity.Username(),
		Email:     userEntity.EmailValue(),
		FirstName: userEntity.FirstName(),
		LastName:  userEntity.LastName(),
		IsActive:  userEntity.IsActive(),
		IsAdmin:   userEntity.IsAdmin(),
	}

	apiCtx.Success(response)
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, "Invalid user ID")
		return
	}

	var req dto.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiCtx.ErrorWithCode(errors.Validation, err.Error())
		return
	}

	// 调用应用服务
	userEntity, err := c.userService.UpdateUser(
		id,
		req.Username,
		req.Email,
		req.FirstName,
		req.LastName,
		req.IsActive,
	)
	if err != nil {
		app.ErrorContext(ctx, "Failed to update user", "error", err, "id", id)
		apiCtx.ErrorWithCode(errors.BadRequest, err.Error())
		return
	}

	// 转换为响应DTO
	response := dto.UserResponse{
		ID:        userEntity.ID(),
		Username:  userEntity.Username(),
		Email:     userEntity.EmailValue(),
		FirstName: userEntity.FirstName(),
		LastName:  userEntity.LastName(),
		IsActive:  userEntity.IsActive(),
		IsAdmin:   userEntity.IsAdmin(),
	}

	apiCtx.Success(response)
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, "Invalid user ID")
		return
	}

	// 调用应用服务
	if err := c.userService.DeleteUser(id); err != nil {
		app.ErrorContext(ctx, "Failed to delete user", "error", err, "id", id)
		apiCtx.ErrorWithCode(errors.NotFound, err.Error())
		return
	}

	apiCtx.Success(gin.H{"message": "User deleted successfully"})
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest
	apiCtx := context.GetAPIContext(ctx)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		app.WarnContext(ctx, "Login validation failed", "error", err)
		apiCtx.ErrorWithCode(errors.Validation, err.Error())
		return
	}

	// 调用应用服务进行认证
	userEntity, err := c.userService.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		app.ErrorContext(ctx, "Login failed", "error", err, "login", req.Login)
		apiCtx.ErrorWithCode(errors.Unauthorized, "Invalid credentials")
		return
	}

	// 在实际应用中，这里应该生成JWT token
	token := "dummy-jwt-token"

	response := dto.UserLoginResponse{
		User: dto.UserResponse{
			ID:        userEntity.ID(),
			Username:  userEntity.Username(),
			Email:     userEntity.EmailValue(),
			FirstName: userEntity.FirstName(),
			LastName:  userEntity.LastName(),
			IsActive:  userEntity.IsActive(),
			IsAdmin:   userEntity.IsAdmin(),
		},
		Token: token,
	}

	apiCtx.Success(response)
}

// UpdatePassword 修改密码
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, "Invalid user ID")
		return
	}

	var req dto.UserUpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiCtx.ErrorWithCode(errors.Validation, err.Error())
		return
	}

	// 调用应用服务修改密码
	if err := c.userService.ChangePassword(id, req.CurrentPassword, req.NewPassword); err != nil {
		app.ErrorContext(ctx, "Failed to update password", "error", err, "user_id", id)
		apiCtx.ErrorWithCode(errors.BadRequest, err.Error())
		return
	}

	apiCtx.Success(gin.H{"message": "Password updated successfully"})
}
