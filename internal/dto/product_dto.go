package dto

// ProductCreateRequest represents the data needed to create a new product
type ProductCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,min=0"`
	SKU         string  `json:"sku" binding:"required"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	CategoryID  int64   `json:"category_id" binding:"required"`
}

// ProductUpdateRequest represents the data needed to update a product
type ProductUpdateRequest struct {
	Name        *string  `json:"name" binding:"omitempty"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty,min=0"`
	SKU         *string  `json:"sku"`
	Stock       *int     `json:"stock" binding:"omitempty,min=0"`
	CategoryID  *int64   `json:"category_id"`
	IsActive    *bool    `json:"is_active"`
}

// ProductResponse represents the product data to be returned in API responses
type ProductResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
	Stock       int     `json:"stock"`
	CategoryID  int64   `json:"category_id"`
	IsActive    bool    `json:"is_active"`
}

// ProductsListResponse represents a paginated list of products
type ProductsListResponse struct {
	Products   []ProductResponse `json:"products"`
	TotalItems int               `json:"total_items"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// ProductStockUpdateRequest represents the data needed to update product stock
type ProductStockUpdateRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}
