package context

import (
	"net/http"

	"gowk/internal/app/errors"
	"gowk/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// APIContext represents the context for API requests
type APIContext struct {
	GinContext *gin.Context
	RequestID  string
	UserID     uint
	Role       string
}

// GetAPIContext retrieves the API context from the gin context
func GetAPIContext(c *gin.Context) *APIContext {
	if apiCtx, exists := c.Get("apiContext"); exists {
		return apiCtx.(*APIContext)
	}

	// Create a new API context if not exists
	apiCtx := &APIContext{
		GinContext: c,
		RequestID:  uuid.New().String(),
	}
	c.Set("apiContext", apiCtx)
	return apiCtx
}

// Success sends a successful response with data
func (ctx *APIContext) Success(data interface{}) {
	ctx.GinContext.JSON(http.StatusOK, dto.Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		RequestID: ctx.RequestID,
	})
}

// Error sends an error response with code and message
func (ctx *APIContext) Error(code int, message string) {
	status := http.StatusInternalServerError
	if code >= 1000 && code < 2000 {
		status = http.StatusBadRequest
	}

	ctx.GinContext.AbortWithStatusJSON(status, dto.Response{
		Code:      code,
		Message:   message,
		RequestID: ctx.RequestID,
	})
}

// ErrorWithCode sends an error response using a predefined error code
func (ctx *APIContext) ErrorWithCode(code errors.ErrorCode, message string) {
	status := http.StatusInternalServerError
	if code >= 1000 && code < 2000 {
		status = http.StatusBadRequest
	}

	// Use standard message if none provided
	if message == "" {
		message = errors.GetStandardMessage(code)
	}

	ctx.GinContext.AbortWithStatusJSON(status, dto.Response{
		Code:      int(code),
		Message:   message,
		RequestID: ctx.RequestID,
	})
}

// ErrorWithAppError sends an error response using an AppError
func (ctx *APIContext) ErrorWithAppError(err *errors.AppError) {
	status := http.StatusInternalServerError
	if err.Code >= 1000 && err.Code < 2000 {
		status = http.StatusBadRequest
	}

	ctx.GinContext.AbortWithStatusJSON(status, dto.Response{
		Code:      int(err.Code),
		Message:   err.Message,
		Data:      err.Details,
		RequestID: ctx.RequestID,
	})
}

// Param gets a URL parameter
func (ctx *APIContext) Param(key string) string {
	return ctx.GinContext.Param(key)
}

// Query gets a query parameter
func (ctx *APIContext) Query(key string) string {
	return ctx.GinContext.Query(key)
}

// BindJSON binds request body to struct
func (ctx *APIContext) BindJSON(obj interface{}) error {
	return ctx.GinContext.ShouldBindJSON(obj)
}

// GetHeader gets a request header
func (ctx *APIContext) GetHeader(key string) string {
	return ctx.GinContext.GetHeader(key)
}

// SetUserID sets the user ID in the context
func (ctx *APIContext) SetUserID(userID uint) {
	ctx.UserID = userID
}

// SetRole sets the user role in the context
func (ctx *APIContext) SetRole(role string) {
	ctx.Role = role
}
