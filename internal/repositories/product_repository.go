package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"goapp/internal/app"
	"goapp/internal/models"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Find(id int64) (*models.Product, error)
	FindByCategory(categoryID int64) ([]*models.Product, error)
	FindAll(limit, offset int) ([]*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id int64) error
}

// SQLProductRepository implements ProductRepository interface using SQL database
type SQLProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository() ProductRepository {
	return &SQLProductRepository{
		db: app.DB,
	}
}

// Find retrieves a product by ID
func (r *SQLProductRepository) Find(id int64) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, sku, stock, category_id, is_active, created_at, updated_at
		FROM products
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var product models.Product
	if err := product.ScanFromRow(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("error finding product: %w", err)
	}

	return &product, nil
}

// FindByCategory retrieves products by category ID
func (r *SQLProductRepository) FindByCategory(categoryID int64) ([]*models.Product, error) {
	query := `
		SELECT id, name, description, price, sku, stock, category_id, is_active, created_at, updated_at
		FROM products
		WHERE category_id = ?
		ORDER BY name
	`
	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error finding products by category: %w", err)
	}
	defer rows.Close()

	products := make([]*models.Product, 0)
	for rows.Next() {
		product := new(models.Product)
		if err := product.ScanFromRows(rows); err != nil {
			return nil, fmt.Errorf("error scanning product: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// FindAll retrieves products with pagination
func (r *SQLProductRepository) FindAll(limit, offset int) ([]*models.Product, error) {
	query := `
		SELECT id, name, description, price, sku, stock, category_id, is_active, created_at, updated_at
		FROM products
		ORDER BY name
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error finding products: %w", err)
	}
	defer rows.Close()

	products := make([]*models.Product, 0)
	for rows.Next() {
		product := new(models.Product)
		if err := product.ScanFromRows(rows); err != nil {
			return nil, fmt.Errorf("error scanning product: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// Create inserts a new product
func (r *SQLProductRepository) Create(product *models.Product) error {
	query := `
		INSERT INTO products (name, description, price, sku, stock, category_id, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.Exec(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.SKU,
		product.Stock,
		product.CategoryID,
		product.IsActive,
	)
	if err != nil {
		return fmt.Errorf("error creating product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting inserted product ID: %w", err)
	}

	product.ID = id
	return nil
}

// Update updates an existing product
func (r *SQLProductRepository) Update(product *models.Product) error {
	query := `
		UPDATE products 
		SET name = ?, description = ?, price = ?, sku = ?, stock = ?,
		    category_id = ?, is_active = ?, updated_at = NOW()
		WHERE id = ?
	`
	result, err := r.db.Exec(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.SKU,
		product.Stock,
		product.CategoryID,
		product.IsActive,
		product.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", product.ID)
	}

	return nil
}

// Delete removes a product by ID
func (r *SQLProductRepository) Delete(id int64) error {
	query := "DELETE FROM products WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}

	return nil
}
