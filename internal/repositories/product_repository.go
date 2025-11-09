package repositories

import (
	"errors"
	"fmt"

	"goapp/internal/app"
	"goapp/internal/models"

	"gorm.io/gorm"
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

// GormProductRepository implements ProductRepository interface using GORM
type GormProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository() ProductRepository {
	return &GormProductRepository{
		db: app.GetDB(),
	}
}

// Find retrieves a product by ID
func (r *GormProductRepository) Find(id int64) (*models.Product, error) {
	var product models.Product
	result := r.db.First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("error finding product: %w", result.Error)
	}
	return &product, nil
}

// FindByCategory retrieves products by category ID
func (r *GormProductRepository) FindByCategory(categoryID int64) ([]*models.Product, error) {
	var products []*models.Product
	result := r.db.Where("category_id = ?", categoryID).Order("name").Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding products by category: %w", result.Error)
	}
	return products, nil
}

// FindAll retrieves products with pagination
func (r *GormProductRepository) FindAll(limit, offset int) ([]*models.Product, error) {
	var products []*models.Product
	result := r.db.Offset(offset).Limit(limit).Order("name").Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("error finding products: %w", result.Error)
	}
	return products, nil
}

// Create inserts a new product
func (r *GormProductRepository) Create(product *models.Product) error {
	result := r.db.Create(product)
	if result.Error != nil {
		return fmt.Errorf("error creating product: %w", result.Error)
	}
	return nil
}

// Update updates an existing product
func (r *GormProductRepository) Update(product *models.Product) error {
	result := r.db.Save(product)
	if result.Error != nil {
		return fmt.Errorf("error updating product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", product.ID)
	}
	return nil
}

// Delete removes a product by ID
func (r *GormProductRepository) Delete(id int64) error {
	result := r.db.Delete(&models.Product{}, id)
	if result.Error != nil {
		return fmt.Errorf("error deleting product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}
	return nil
}
