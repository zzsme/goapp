package services

import (
	"fmt"
	"goapp/internal/app"
	"goapp/internal/models"
	"goapp/internal/repositories"
)

// ProductService handles business logic for product operations
type ProductService struct {
	productRepo repositories.ProductRepository
}

// NewProductService creates a new ProductService
func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repositories.NewProductRepository(),
	}
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(id int64) (*models.Product, error) {
	app.Debug("Getting product", "id", id)
	return s.productRepo.Find(id)
}

// ListProducts retrieves products with pagination
func (s *ProductService) ListProducts(page, pageSize int) ([]*models.Product, error) {
	app.Debug("Listing products", "page", page, "page_size", pageSize)

	// Ensure page is positive
	if page < 1 {
		page = 1
	}

	// Set reasonable default and bounds for page size
	if pageSize < 1 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	return s.productRepo.FindAll(pageSize, offset)
}

// ListProductsByCategory retrieves products by category
func (s *ProductService) ListProductsByCategory(categoryID int64) ([]*models.Product, error) {
	app.Debug("Listing products by category", "category_id", categoryID)
	return s.productRepo.FindByCategory(categoryID)
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(product *models.Product) error {
	app.Debug("Creating product", "name", product.Name)

	// Validate SKU uniqueness (would typically check against DB)
	if product.SKU == "" {
		return fmt.Errorf("SKU cannot be empty")
	}

	// Set defaults for new product
	product.IsActive = true

	return s.productRepo.Create(product)
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(product *models.Product) error {
	app.Debug("Updating product", "id", product.ID)

	// Ensure product exists
	_, err := s.productRepo.Find(product.ID)
	if err != nil {
		return err
	}

	// Update the product
	return s.productRepo.Update(product)
}

// DeleteProduct removes a product by ID
func (s *ProductService) DeleteProduct(id int64) error {
	app.Debug("Deleting product", "id", id)
	return s.productRepo.Delete(id)
}

// UpdateProductStock updates only the stock quantity of a product
func (s *ProductService) UpdateProductStock(id int64, quantity int) error {
	app.Debug("Updating product stock", "id", id, "quantity", quantity)

	// Ensure product exists
	product, err := s.productRepo.Find(id)
	if err != nil {
		return err
	}

	// Update stock quantity
	product.Stock = quantity

	return s.productRepo.Update(product)
}
