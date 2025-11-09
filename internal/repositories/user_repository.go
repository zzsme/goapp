package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"goapp/internal/app"
	"goapp/internal/models"
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

// SQLUserRepository implements UserRepository interface using SQL database
type SQLUserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository() UserRepository {
	return &SQLUserRepository{
		db: app.DB,
	}
}

// Find retrieves a user by ID
func (r *SQLUserRepository) Find(id int64) (*models.User, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, is_active, is_admin, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var user models.User
	if err := user.ScanFromRow(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *SQLUserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, is_active, is_admin, created_at, updated_at
		FROM users
		WHERE email = ?
	`
	row := r.db.QueryRow(query, email)

	var user models.User
	if err := user.ScanFromRow(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &user, nil
}

// FindByUsername retrieves a user by username
func (r *SQLUserRepository) FindByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, is_active, is_admin, created_at, updated_at
		FROM users
		WHERE username = ?
	`
	row := r.db.QueryRow(query, username)

	var user models.User
	if err := user.ScanFromRow(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &user, nil
}

// FindAll retrieves users with pagination
func (r *SQLUserRepository) FindAll(limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, username, email, password, first_name, last_name, is_active, is_admin, created_at, updated_at
		FROM users
		ORDER BY id
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error finding users: %w", err)
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user := new(models.User)
		if err := user.ScanFromRows(rows); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// Create inserts a new user
func (r *SQLUserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password, first_name, last_name, is_active, is_admin, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.Exec(
		query,
		user.Username,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.IsAdmin,
	)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting inserted user ID: %w", err)
	}

	user.ID = id
	return nil
}

// Update updates an existing user
func (r *SQLUserRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, first_name = ?, last_name = ?, 
		    is_active = ?, is_admin = ?, updated_at = NOW()
		WHERE id = ?
	`
	result, err := r.db.Exec(
		query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.IsAdmin,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}

	return nil
}

// Delete removes a user by ID
func (r *SQLUserRepository) Delete(id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}
