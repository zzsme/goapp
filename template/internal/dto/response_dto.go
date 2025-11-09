package dto

// Response represents a standard API response
type Response struct {
	Code      int         `json:"code"`      // Error code (0 for success)
	Message   string      `json:"message"`   // Error message or success message
	Data      interface{} `json:"data"`      // Response data (optional)
	RequestID string      `json:"requestId"` // Request ID for tracing
}

// PageResponse represents a paginated response
type PageResponse struct {
	Total    int64       `json:"total"`    // Total number of records
	Page     int         `json:"page"`     // Current page number
	PageSize int         `json:"pageSize"` // Number of records per page
	Data     interface{} `json:"data"`     // Page data
}

// ListResponse represents a list response
type ListResponse struct {
	Total int64       `json:"total"` // Total number of records
	Data  interface{} `json:"data"`  // List data
}

// TokenResponse represents an authentication token response
type TokenResponse struct {
	AccessToken  string `json:"accessToken"`  // JWT access token
	TokenType    string `json:"tokenType"`    // Token type (e.g., "Bearer")
	ExpiresIn    int64  `json:"expiresIn"`    // Token expiration time in seconds
	RefreshToken string `json:"refreshToken"` // Refresh token for obtaining new access tokens
}

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	Field   string `json:"field"`   // Field name for validation errors
	Message string `json:"message"` // Error message for this field
}
