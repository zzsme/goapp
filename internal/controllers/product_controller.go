package controllers

import (
	"strconv"

	"goapp/internal/app"
	"goapp/internal/app/errors"
	"goapp/internal/context"
	"goapp/internal/dto"
	"goapp/internal/models"
	"goapp/internal/services"

	"github.com/gin-gonic/gin"
)

// ProductController handles HTTP requests for product operations
type ProductController struct {
	productService *services.ProductService
}

// NewProductController creates a new ProductController
func NewProductController() *ProductController {
	return &ProductController{
		productService: services.NewProductService(),
	}
}

// Register adds product routes to the router group
func (c *ProductController) Register(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		products.GET("", c.ListProducts)
		products.GET("/:id", c.GetProduct)
		products.POST("", c.CreateProduct)
		products.PUT("/:id", c.UpdateProduct)
		products.DELETE("/:id", c.DeleteProduct)
		products.PUT("/:id/stock", c.UpdateProductStock)
		products.GET("/category/:id", c.ListProductsByCategory)
	}
}

// ListProducts handles requests to list products with pagination
func (c *ProductController) ListProducts(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	var pagination dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, err.Error())
		return
	}

	products, err := c.productService.ListProducts(pagination.Page, pagination.PageSize)
	if err != nil {
		app.ErrorContext(ctx, "Failed to list products", "error", err)
		apiCtx.ErrorWithCode(errors.InternalServer, "Failed to retrieve products")
		return
	}

	// Convert products to response DTOs
	productResponses := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			SKU:         product.SKU,
			Stock:       product.Stock,
			CategoryID:  product.CategoryID,
			IsActive:    product.IsActive,
		}
	}

	response := dto.ProductsListResponse{
		Products:   productResponses,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalItems: len(products), // In a real app, you'd get this from the database
		TotalPages: (len(products) + pagination.PageSize - 1) / pagination.PageSize,
	}

	apiCtx.Success(response)
}

// GetProduct handles requests to get a specific product
func (c *ProductController) GetProduct(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, "Invalid product ID")
		return
	}

	product, err := c.productService.GetProduct(id)
	if err != nil {
		app.ErrorContext(ctx, "Failed to get product", "error", err, "id", id)
		apiCtx.ErrorWithCode(errors.NotFound, "Product not found")
		return
	}

	response := dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		IsActive:    product.IsActive,
	}

	apiCtx.Success(response)
}

// CreateProduct handles product creation requests
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	var req dto.ProductCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiCtx.ErrorWithCode(errors.Validation, err.Error())
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

	if err := c.productService.CreateProduct(product); err != nil {
		app.ErrorContext(ctx, "Failed to create product", "error", err)
		apiCtx.ErrorWithCode(errors.BadRequest, err.Error())
		return
	}

	response := dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		IsActive:    product.IsActive,
	}

	apiCtx.Success(response)
}

// UpdateProduct handles requests to update a product
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, "Invalid product ID")
		return
	}

	var req dto.ProductUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiCtx.ErrorWithCode(errors.Validation, err.Error())
		return
	}

	// Get existing product
	product, err := c.productService.GetProduct(id)
	if err != nil {
		apiCtx.ErrorWithCode(errors.NotFound, "Product not found")
		return
	}

	// Update fields if provided
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.SKU != nil {
		product.SKU = *req.SKU
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := c.productService.UpdateProduct(product); err != nil {
		app.ErrorContext(ctx, "Failed to update product", "error", err, "id", id)
		apiCtx.ErrorWithCode(errors.BadRequest, err.Error())
		return
	}

	response := dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		IsActive:    product.IsActive,
	}

	apiCtx.Success(response)
}

// DeleteProduct handles requests to delete a product
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, "Invalid product ID")
		return
	}

	if err := c.productService.DeleteProduct(id); err != nil {
		app.ErrorContext(ctx, "Failed to delete product", "error", err, "id", id)
		apiCtx.ErrorWithCode(errors.NotFound, "Product not found")
		return
	}

	apiCtx.Success(gin.H{"message": "Product deleted successfully"})
}

// UpdateProductStock handles requests to update a product's stock
func (c *ProductController) UpdateProductStock(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, "Invalid product ID")
		return
	}

	var req dto.ProductStockUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiCtx.ErrorWithCode(errors.Validation, err.Error())
		return
	}

	if err := c.productService.UpdateProductStock(id, req.Quantity); err != nil {
		app.ErrorContext(ctx, "Failed to update product stock", "error", err, "id", id)
		apiCtx.ErrorWithCode(errors.BadRequest, err.Error())
		return
	}

	apiCtx.Success(gin.H{"message": "Product stock updated successfully"})
}

// ListProductsByCategory handles requests to list products by category
func (c *ProductController) ListProductsByCategory(ctx *gin.Context) {
	apiCtx := context.GetAPIContext(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		apiCtx.ErrorWithCode(errors.BadRequest, "Invalid category ID")
		return
	}

	products, err := c.productService.ListProductsByCategory(id)
	if err != nil {
		app.ErrorContext(ctx, "Failed to list products by category", "error", err, "category_id", id)
		apiCtx.ErrorWithCode(errors.InternalServer, "Failed to retrieve products")
		return
	}

	// Convert products to response DTOs
	productResponses := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = dto.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			SKU:         product.SKU,
			Stock:       product.Stock,
			CategoryID:  product.CategoryID,
			IsActive:    product.IsActive,
		}
	}

	apiCtx.Success(gin.H{"products": productResponses})
}
