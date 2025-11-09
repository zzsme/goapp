package dto

// UserCreateRequest represents the data needed to create a new user
type UserCreateRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// UserUpdateRequest represents the data needed to update a user
type UserUpdateRequest struct {
	Username  *string `json:"username" binding:"omitempty,min=3,max=50"`
	Email     *string `json:"email" binding:"omitempty,email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	IsActive  *bool   `json:"is_active"`
}

// UserUpdatePasswordRequest represents the data needed to update a user's password
type UserUpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// UserLoginRequest represents the data needed for user login
type UserLoginRequest struct {
	Login    string `json:"login" binding:"required"` // Can be email or username
	Password string `json:"password" binding:"required"`
}

// UserResponse represents the user data to be returned in API responses
type UserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  bool   `json:"is_active"`
	IsAdmin   bool   `json:"is_admin"`
}

// UsersListResponse represents a paginated list of users
type UsersListResponse struct {
	Users      []UserResponse `json:"users"`
	TotalItems int            `json:"total_items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// UserLoginResponse represents the response after successful login
type UserLoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=10" binding:"min=1,max=100"`
}
