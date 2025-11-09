package errors

import (
	"errors"
	"fmt"
)

// AppError represents an application error with code, message, and optional cause and details
type AppError struct {
	Code    ErrorCode
	Message string
	Cause   error
	Details map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New creates a new AppError with the given error code and message
func New(code ErrorCode, message string) *AppError {
	if message == "" {
		message = GetStandardMessage(code)
	}

	return &AppError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// Wrap creates a new AppError that wraps an existing error
func Wrap(code ErrorCode, message string, err error) *AppError {
	if message == "" {
		message = GetStandardMessage(code)
	}

	return &AppError{
		Code:    code,
		Message: message,
		Cause:   err,
		Details: make(map[string]interface{}),
	}
}

// With adds a key-value pair to the error details
func (e *AppError) With(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithDetails adds multiple key-value pairs to the error details
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// IsCode checks if an error has a specific error code
func IsCode(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// Unwrap returns the wrapped error
func Unwrap(err error) error {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Cause
	}
	return nil
}

// GetCode extracts the error code from an error
func GetCode(err error) ErrorCode {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return InternalServer
}

// GetMessage extracts the error message from an error
func GetMessage(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}
	return err.Error()
}

// GetDetails extracts the error details from an error
func GetDetails(err error) map[string]interface{} {
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Details != nil {
		return appErr.Details
	}
	return nil
}

// AsAppError tries to convert an error to an AppError
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
