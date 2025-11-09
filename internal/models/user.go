package models

import (
	"database/sql"
	"time"
)

// User represents a user in the system
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose password in JSON responses
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName returns the database table name for the User model
func (User) TableName() string {
	return "users"
}

// ScanFromRow scans a database row into a User model
func (u *User) ScanFromRow(row *sql.Row) error {
	return row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.FirstName,
		&u.LastName,
		&u.IsActive,
		&u.IsAdmin,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
}

// ScanFromRows scans a database row from a rows result set into a User model
func (u *User) ScanFromRows(rows *sql.Rows) error {
	return rows.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.FirstName,
		&u.LastName,
		&u.IsActive,
		&u.IsAdmin,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
}
