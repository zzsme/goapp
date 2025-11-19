package product

import "fmt"

// SKU 库存单位值对象
type SKU struct {
	value string
}

// NewSKU 创建SKU值对象
func NewSKU(id int64) SKU {
	return SKU{value: fmt.Sprintf("SKU-%d", id)}
}

// NewSKUFromString 从字符串创建SKU
func NewSKUFromString(sku string) SKU {
	return SKU{value: sku}
}

// Value 返回SKU字符串
func (s SKU) Value() string {
	return s.value
}
