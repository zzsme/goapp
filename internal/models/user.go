package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:100;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:255;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"` // Never expose password in JSON responses
	FirstName string         `json:"first_name" gorm:"size:100"`
	LastName  string         `json:"last_name" gorm:"size:100"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	IsAdmin   bool           `json:"is_admin" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete support
}

// TableName returns the database table name for the User model
func (User) TableName() string {
	return "users"
}
