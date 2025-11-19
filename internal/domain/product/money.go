package product

import "errors"

// Money 金额值对象
type Money struct {
	value float64
}

// NewMoney 创建Money值对象
func NewMoney(amount float64) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("price cannot be negative")
	}
	return Money{value: amount}, nil
}

// Value 返回金额值
func (m Money) Value() float64 {
	return m.value
}
