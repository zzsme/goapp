package controllers

import (
	"goapp/internal/context"
	"goapp/internal/dto"
	"goapp/internal/models"
	"goapp/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProductController 产品控制器
type ProductController struct {
	productService *services.ProductService
}

// NewProductController 创建产品控制器
func NewProductController() *ProductController {
	return &ProductController{
		productService: services.NewProductService(),
	}
}

// Register 注册路由
func (pc *ProductController) Register(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		products.GET("", pc.List)
		products.GET("/:id", pc.Get)
		products.POST("", pc.Create)
		products.PUT("/:id", pc.Update)
		products.DELETE("/:id", pc.Delete)
	}
}

// List 获取产品列表
func (pc *ProductController) List(c *gin.Context) {
	apiCtx := context.GetAPIContext(c)
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	products, err := pc.productService.GetAll(page, pageSize)
	if err != nil {
		apiCtx.Error(500, "Failed to get products")
		return
	}

	apiCtx.Success(products)
}

// Get 获取单个产品
func (pc *ProductController) Get(c *gin.Context) {
	apiCtx := context.GetAPIContext(c)
	
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	product, err := pc.productService.GetByID(id)
	if err != nil {
		apiCtx.Error(404, "Product not found")
		return
	}

	apiCtx.Success(product)
}

// Create 创建产品
func (pc *ProductController) Create(c *gin.Context) {
	apiCtx := context.GetAPIContext(c)
	
	var req dto.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiCtx.Error(400, "Invalid request")
		return
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
	}

	if err := pc.productService.Create(product); err != nil {
		apiCtx.Error(500, "Failed to create product")
		return
	}

	apiCtx.Success(product)
}

// Update 更新产品
func (pc *ProductController) Update(c *gin.Context) {
	apiCtx := context.GetAPIContext(c)
	
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	
	var req dto.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiCtx.Error(400, "Invalid request")
		return
	}

	product, err := pc.productService.GetByID(id)
	if err != nil {
		apiCtx.Error(404, "Product not found")
		return
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if err := pc.productService.Update(product); err != nil {
		apiCtx.Error(500, "Failed to update product")
		return
	}

	apiCtx.Success(product)
}

// Delete 删除产品
func (pc *ProductController) Delete(c *gin.Context) {
	apiCtx := context.GetAPIContext(c)
	
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	
	if err := pc.productService.Delete(id); err != nil {
		apiCtx.Error(500, "Failed to delete product")
		return
	}

	apiCtx.Success(gin.H{"message": "Product deleted successfully"})
}
