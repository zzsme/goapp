package repositories

import (
	"errors"
	"fmt"

	"goapp/internal/app"
	"goapp/internal/models"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Find(id int64) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindAll(limit, offset int) ([]*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int64) error
}

// GormUserRepository implements UserRepository interface using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository() UserRepository {
	return &GormUserRepository{
		db: app.GetDB(),
	}
}

// Find retrieves a user by ID
func (r *GormUserRepository) Find(id int64) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("error finding user: %w", result.Error)
	}
	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *GormUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("error finding user: %w", result.Error)
	}
	return &user, nil
}

// FindByUsername retrieves a user by username
func (r *GormUserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("error finding user: %w", result.Error)
	}
	return &user, nil
}

// FindAll retrieves users with pagination
func (r *GormUserRepository) FindAll(limit, offset int) ([]*models.User, error) {
	var users []*models.User
	result := r.db.Offset(offset).Limit(limit).Order("id").Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding users: %w", result.Error)
	}
	return users, nil
}

// Create inserts a new user
func (r *GormUserRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return fmt.Errorf("error creating user: %w", result.Error)
	}
	return nil
}

// Update updates an existing user
func (r *GormUserRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return fmt.Errorf("error updating user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}
	return nil
}

// Delete removes a user by ID
func (r *GormUserRepository) Delete(id int64) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("error deleting user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}
	return nil
}
