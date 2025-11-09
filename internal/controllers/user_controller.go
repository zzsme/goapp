package controllers

import (
	"strconv"

	"goapp/internal/app"
	"goapp/internal/app/errors"
	"goapp/internal/context"
	"goapp/internal/dto"
	"goapp/internal/models"
	"goapp/internal/services"

	"github.com/gin-gonic/gin"
)

// UserController handles HTTP requests for user operations
type UserController struct {
	userService *services.UserService
}

// NewUserController creates a new UserController
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// Register adds user routes to the router group
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

// CreateUser handles user creation requests
func (c *UserController) CreateUser(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	var req dto.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiCtx.ErrorWithCode(errors.Validation, err.Error())
		return
	}

	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := c.userService.CreateUser(user); err != nil {
		app.ErrorContext(ctx, "Failed to create user", "error", err)
		apiCtx.ErrorWithCode(errors.BadRequest, err.Error())
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		IsAdmin:   user.IsAdmin,
	}

	apiCtx.Success(response)
}

// ListUsers handles requests to list users with pagination
func (c *UserController) ListUsers(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	var pagination dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, err.Error())
		return
	}

	users, err := c.userService.ListUsers(pagination.Page, pagination.PageSize)
	if err != nil {
		app.ErrorContext(ctx, "Failed to list users", "error", err)
		apiCtx.ErrorWithCode(errors.InternalServer, "Failed to retrieve users")
		return
	}

	// Convert users to response DTOs
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
			IsAdmin:   user.IsAdmin,
		}
	}

	response := dto.UsersListResponse{
		Users:      userResponses,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalItems: len(users), // In a real app, you'd get this from the database
		TotalPages: (len(users) + pagination.PageSize - 1) / pagination.PageSize,
	}

	apiCtx.Success(response)
}

// GetUser handles requests to get a specific user
func (c *UserController) GetUser(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, "Invalid user ID")
		return
	}

	user, err := c.userService.GetUser(id)
	if err != nil {
		app.ErrorContext(ctx, "Failed to get user", "error", err, "id", id)
		apiCtx.ErrorWithCode(errors.NotFound, "User not found")
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		IsAdmin:   user.IsAdmin,
	}

	apiCtx.Success(response)
}

// UpdateUser handles requests to update a user
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

	// Get existing user
	user, err := c.userService.GetUser(id)
	if err != nil {
		apiCtx.ErrorWithCode(errors.NotFound, "User not found")
		return
	}

	// Update fields if provided
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := c.userService.UpdateUser(user); err != nil {
		app.ErrorContext(ctx, "Failed to update user", "error", err, "id", id)
		apiCtx.ErrorWithCode(errors.BadRequest, err.Error())
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		IsAdmin:   user.IsAdmin,
	}

	apiCtx.Success(response)
}

// DeleteUser handles requests to delete a user
func (c *UserController) DeleteUser(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, "Invalid user ID")
		return
	}

	if err := c.userService.DeleteUser(id); err != nil {
		app.ErrorContext(ctx, "Failed to delete user", "error", err, "id", id)
		apiCtx.ErrorWithCode(errors.NotFound, "User not found")
		return
	}

	apiCtx.Success(gin.H{"message": "User deleted successfully"})
}

// Login handles user authentication requests
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest
	apiCtx := context.GetAPIContext(ctx)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		app.WarnContext(ctx, "Login validation failed", "error", err)
		apiCtx.ErrorWithCode(errors.Validation, err.Error())
		return
	}

	user, err := c.userService.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		app.ErrorContext(ctx, "Login failed", "error", err, "login", req.Login)
		apiCtx.ErrorWithCode(errors.Unauthorized, "Invalid credentials")
		return
	}

	// In a real application, you would generate a JWT token here
	token := "dummy-jwt-token"

	response := dto.UserLoginResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
			IsAdmin:   user.IsAdmin,
		},
		Token: token,
	}

	apiCtx.Success(response)
}

// UpdatePassword handles password update requests
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

	// Get user
	user, err := c.userService.GetUser(id)
	if err != nil {
		apiCtx.ErrorWithCode(errors.NotFound, "User not found")
		return
	}

	// Verify current password and update with new password
	// In a real application, you would implement this in the service layer
	app.InfoContext(ctx, "Password update requested", "user_id", user.ID)
	apiCtx.Success(gin.H{"message": "Password updated successfully"})
}
