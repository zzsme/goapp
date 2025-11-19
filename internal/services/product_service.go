package services

import (
	"goapp/internal/app"
	"goapp/internal/models"
)

// ProductService 产品服务（简化版）
type ProductService struct{}

// NewProductService 创建产品服务
func NewProductService() *ProductService {
	return &ProductService{}
}

// GetAll 获取所有产品
func (s *ProductService) GetAll(page, pageSize int) ([]models.Product, error) {
	var products []models.Product
	db := app.GetDB()
	
	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(&products).Error
	return products, err
}

// GetByID 根据ID获取产品
func (s *ProductService) GetByID(id int64) (*models.Product, error) {
	var product models.Product
	db := app.GetDB()
	
	err := db.First(&product, id).Error
	return &product, err
}

// Create 创建产品
func (s *ProductService) Create(product *models.Product) error {
	db := app.GetDB()
	return db.Create(product).Error
}

// Update 更新产品
func (s *ProductService) Update(product *models.Product) error {
	db := app.GetDB()
	return db.Save(product).Error
}

// Delete 删除产品
func (s *ProductService) Delete(id int64) error {
	db := app.GetDB()
	return db.Delete(&models.Product{}, id).Error
}
