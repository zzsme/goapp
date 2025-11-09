package models

import (
	"database/sql"
	"time"
)

// Product represents a product in the system
type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	SKU         string    `json:"sku"`
	Stock       int       `json:"stock"`
	CategoryID  int64     `json:"category_id"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName returns the database table name for the Product model
func (Product) TableName() string {
	return "products"
}

// ScanFromRow scans a database row into a Product model
func (p *Product) ScanFromRow(row *sql.Row) error {
	return row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.SKU,
		&p.Stock,
		&p.CategoryID,
		&p.IsActive,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
}

// ScanFromRows scans a database row from a rows result set into a Product model
func (p *Product) ScanFromRows(rows *sql.Rows) error {
	return rows.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.SKU,
		&p.Stock,
		&p.CategoryID,
		&p.IsActive,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
}
