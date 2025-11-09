package errors

import "fmt"

// AppError represents an application error with a code and message
type AppError struct {
	Code    ErrorCode
	Message string
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewError creates a new AppError with a custom message
func NewError(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// New creates a new AppError with the standard message for the error code
func New(code ErrorCode) *AppError {
	return &AppError{
		Code:    code,
		Message: GetStandardMessage(code),
	}
}

// IsClientError checks if the error code represents a client error
func IsClientError(code ErrorCode) bool {
	return code >= 1000 && code < 2000
}

// IsServerError checks if the error code represents a server error
func IsServerError(code ErrorCode) bool {
	return code >= 5000 && code < 6000
}
