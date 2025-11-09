package errors

// ErrorCode represents a unique error code in the system
type ErrorCode int

// Error codes for API responses
const (
	// Success represents a successful operation
	Success ErrorCode = 0

	// Client errors (1000-1999)
	BadRequest    ErrorCode = 1000 // Invalid parameters or request
	Unauthorized  ErrorCode = 1001 // Authentication required
	Forbidden     ErrorCode = 1003 // Permission denied
	NotFound      ErrorCode = 1004 // Resource not found
	Validation    ErrorCode = 1005 // Validation error
	Conflict      ErrorCode = 1009 // Resource conflict
	TooManyReq    ErrorCode = 1029 // Too many requests
	InvalidToken  ErrorCode = 1030 // Invalid token
	ExpiredToken  ErrorCode = 1031 // Token expired
	InvalidFormat ErrorCode = 1032 // Invalid data format

	// Server errors (5000-5999)
	InternalServer ErrorCode = 5000 // Internal server error
	Database       ErrorCode = 5001 // Database error
	Cache          ErrorCode = 5002 // Cache error
	NotImplemented ErrorCode = 5003 // Not implemented
	ThirdParty     ErrorCode = 5004 // Third-party service error
	Config         ErrorCode = 5005 // Configuration error
)

// Standard error messages
var standardMessages = map[ErrorCode]string{
	BadRequest:     "Invalid parameters",
	Unauthorized:   "Authentication required",
	Forbidden:      "Permission denied",
	NotFound:       "Resource not found",
	InternalServer: "Internal server error",
	Validation:     "Validation error",
	Database:       "Database error",
	InvalidToken:   "Invalid token",
	ExpiredToken:   "Token expired",
	TooManyReq:     "Too many requests",
	InvalidFormat:  "Invalid data format",
	NotImplemented: "Feature not implemented",
	ThirdParty:     "Third-party service error",
	Config:         "Configuration error",
}

// GetStandardMessage returns the standard message for an error code
func GetStandardMessage(code ErrorCode) string {
	if msg, ok := standardMessages[code]; ok {
		return msg
	}
	return "Unknown error"
}
