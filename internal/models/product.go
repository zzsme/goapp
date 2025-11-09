package models

import (
	"time"

	"gorm.io/gorm"
)

// Product represents a product in the system
type Product struct {
	ID          int64          `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:255;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	SKU         string         `json:"sku" gorm:"uniqueIndex;size:100;not null"`
	Stock       int            `json:"stock" gorm:"not null;default:0"`
	CategoryID  int64          `json:"category_id" gorm:"index"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete support
}

// TableName returns the database table name for the Product model
func (Product) TableName() string {
	return "products"
}
