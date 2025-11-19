package product

import (
	"errors"
	"time"
)

// Product 产品实体 - 富血模型
type Product struct {
	id          int64
	name        string
	description string
	price       Money      // 值对象
	sku         SKU        // 值对象
	stock       int
	categoryID  int64
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
}

// NewProduct 创建新产品
func NewProduct(name, description string, price float64, sku, categoryID int64, stock int) (*Product, error) {
	// 验证产品名称
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}

	// 创建Money值对象
	moneyVO, err := NewMoney(price)
	if err != nil {
		return nil, err
	}

	// 创建SKU值对象
	skuVO := NewSKU(sku)

	// 验证库存
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	now := time.Now()
	return &Product{
		name:        name,
		description: description,
		price:       moneyVO,
		sku:         skuVO,
		stock:       stock,
		categoryID:  categoryID,
		isActive:    true,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Reconstruct 从持久化数据重建Product对象
func Reconstruct(
	id int64,
	name, description string,
	price Money,
	sku SKU,
	stock int,
	categoryID int64,
	isActive bool,
	createdAt, updatedAt time.Time,
) *Product {
	return &Product{
		id:          id,
		name:        name,
		description: description,
		price:       price,
		sku:         sku,
		stock:       stock,
		categoryID:  categoryID,
		isActive:    isActive,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// 业务行为方法

// Activate 激活产品
func (p *Product) Activate() {
	p.isActive = true
	p.updatedAt = time.Now()
}

// Deactivate 停用产品
func (p *Product) Deactivate() {
	p.isActive = false
	p.updatedAt = time.Now()
}

// UpdatePrice 更新价格
func (p *Product) UpdatePrice(newPrice float64) error {
	money, err := NewMoney(newPrice)
	if err != nil {
		return err
	}
	p.price = money
	p.updatedAt = time.Now()
	return nil
}

// AddStock 增加库存
func (p *Product) AddStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	p.stock += quantity
	p.updatedAt = time.Now()
	return nil
}

// ReduceStock 减少库存
func (p *Product) ReduceStock(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if p.stock < quantity {
		return errors.New("insufficient stock")
	}
	p.stock -= quantity
	p.updatedAt = time.Now()
	return nil
}

// UpdateInfo 更新产品信息
func (p *Product) UpdateInfo(name, description string) error {
	if name != "" {
		p.name = name
	}
	if description != "" {
		p.description = description
	}
	p.updatedAt = time.Now()
	return nil
}

// IsAvailable 检查产品是否可用
func (p *Product) IsAvailable() bool {
	return p.isActive && p.stock > 0
}

// Getter方法

// ID 返回产品ID
func (p *Product) ID() int64 {
	return p.id
}

// Name 返回产品名称
func (p *Product) Name() string {
	return p.name
}

// Description 返回产品描述
func (p *Product) Description() string {
	return p.description
}

// Price 返回价格值对象
func (p *Product) Price() Money {
	return p.price
}

// PriceValue 返回价格数值
func (p *Product) PriceValue() float64 {
	return p.price.Value()
}

// SKU 返回SKU值对象
func (p *Product) SKU() SKU {
	return p.sku
}

// SKUValue 返回SKU字符串
func (p *Product) SKUValue() string {
	return p.sku.Value()
}

// Stock 返回库存
func (p *Product) Stock() int {
	return p.stock
}

// CategoryID 返回分类ID
func (p *Product) CategoryID() int64 {
	return p.categoryID
}

// IsActive 返回是否激活
func (p *Product) IsActive() bool {
	return p.isActive
}

// CreatedAt 返回创建时间
func (p *Product) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt 返回更新时间
func (p *Product) UpdatedAt() time.Time {
	return p.updatedAt
}
