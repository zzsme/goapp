package context

import (
	"net/http"

	"goapp/internal/app/errors"
	"goapp/internal/dto"

	"github.com/gin-gonic/gin"
)

// APIContext wraps the Gin context to provide standardized API responses
type APIContext struct {
	ctx *gin.Context
}

// GetAPIContext returns a new APIContext wrapping the Gin context
func GetAPIContext(c *gin.Context) *APIContext {
	return &APIContext{
		ctx: c,
	}
}

// Success sends a successful response with the provided data
func (ac *APIContext) Success(data interface{}) {
	ac.ctx.JSON(http.StatusOK, dto.StandardResponse{
		Errno:     0,
		Errmsg:    "",
		Data:      data,
		RequestID: GetRequestID(ac.ctx),
	})
}

// Error sends an error response with the provided error code and message
// For backward compatibility with existing code
func (ac *APIContext) Error(errno int, errmsg string) {
	// Map error codes to appropriate HTTP status codes
	httpStatus := errorToHTTPStatus(errors.ErrorCode(errno))

	ac.ctx.JSON(httpStatus, dto.StandardResponse{
		Errno:     errno,
		Errmsg:    errmsg,
		Data:      nil,
		RequestID: GetRequestID(ac.ctx),
	})

	// Abort the request to prevent further handlers from executing
	ac.ctx.Abort()
}

// ErrorWithAppError sends an error response with the provided AppError
func (ac *APIContext) ErrorWithAppError(err *errors.AppError) {
	// Map error codes to appropriate HTTP status codes
	httpStatus := errorToHTTPStatus(err.Code)

	ac.ctx.JSON(httpStatus, dto.StandardResponse{
		Errno:     int(err.Code),
		Errmsg:    err.Message,
		Data:      nil,
		RequestID: GetRequestID(ac.ctx),
	})

	// Abort the request to prevent further handlers from executing
	ac.ctx.Abort()
}

// ErrorWithCode sends an error response with the provided error code and message
func (ac *APIContext) ErrorWithCode(code errors.ErrorCode, message string) {
	ac.ErrorWithAppError(errors.NewError(code, message))
}

// errorToHTTPStatus maps internal error codes to HTTP status codes
func errorToHTTPStatus(code errors.ErrorCode) int {
	// Client errors (1000-1999)
	if errors.IsClientError(code) {
		switch code {
		case errors.BadRequest:
			return http.StatusBadRequest
		case errors.Unauthorized:
			return http.StatusUnauthorized
		case errors.Forbidden:
			return http.StatusForbidden
		case errors.NotFound:
			return http.StatusNotFound
		case errors.Conflict:
			return http.StatusConflict
		case errors.TooManyReq:
			return http.StatusTooManyRequests
		default:
			return http.StatusBadRequest
		}
	}

	// Server errors (5000-5999)
	if errors.IsServerError(code) {
		switch code {
		case errors.NotImplemented:
			return http.StatusNotImplemented
		default:
			return http.StatusInternalServerError
		}
	}

	// Default status code for any other error
	return http.StatusInternalServerError
}
