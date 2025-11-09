package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"goapp/internal/models"
	"goapp/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidPassword = errors.New("password must be at least 8 characters long")
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailTaken      = errors.New("email is already taken")
	ErrUsernameTaken   = errors.New("username is already taken")
)

// UserService handles business logic for user operations
type UserService struct {
	repo repositories.UserRepository
}

// NewUserService creates a new UserService
func NewUserService() *UserService {
	return &UserService{
		repo: repositories.NewUserRepository(),
	}
}

// CreateUser creates a new user with validation
func (s *UserService) CreateUser(user *models.User) error {
	// Validate email format
	if !isValidEmail(user.Email) {
		return ErrInvalidEmail
	}

	// Validate password length
	if len(user.Password) < 8 {
		return ErrInvalidPassword
	}

	// Check if email is already taken
	existingUser, err := s.repo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return ErrEmailTaken
	}

	// Check if username is already taken
	existingUser, err = s.repo.FindByUsername(user.Username)
	if err == nil && existingUser != nil {
		return ErrUsernameTaken
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	user.Password = string(hashedPassword)

	// Set default values
	user.IsActive = true
	user.IsAdmin = false

	// Create user
	return s.repo.Create(user)
}

// UpdateUser updates user information with validation
func (s *UserService) UpdateUser(user *models.User) error {
	// Validate email format if provided
	if user.Email != "" && !isValidEmail(user.Email) {
		return ErrInvalidEmail
	}

	// Check if user exists
	existingUser, err := s.repo.Find(user.ID)
	if err != nil {
		return ErrUserNotFound
	}

	// Check if new email is taken by another user
	if user.Email != "" && user.Email != existingUser.Email {
		otherUser, err := s.repo.FindByEmail(user.Email)
		if err == nil && otherUser != nil && otherUser.ID != user.ID {
			return ErrEmailTaken
		}
	}

	// Check if new username is taken by another user
	if user.Username != "" && user.Username != existingUser.Username {
		otherUser, err := s.repo.FindByUsername(user.Username)
		if err == nil && otherUser != nil && otherUser.ID != user.ID {
			return ErrUsernameTaken
		}
	}

	// Update user
	return s.repo.Update(user)
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id int64) (*models.User, error) {
	user, err := s.repo.Find(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// ListUsers retrieves a paginated list of users
func (s *UserService) ListUsers(page, pageSize int) ([]*models.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.repo.FindAll(pageSize, offset)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id int64) error {
	return s.repo.Delete(id)
}

// AuthenticateUser authenticates a user with email/username and password
func (s *UserService) AuthenticateUser(login, password string) (*models.User, error) {
	var user *models.User
	var err error

	// Try to find user by email or username
	if strings.Contains(login, "@") {
		user, err = s.repo.FindByEmail(login)
	} else {
		user, err = s.repo.FindByUsername(login)
	}

	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
