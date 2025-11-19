package user

import (
	"errors"
	"regexp"
)

// Email 值对象 - 封装邮箱验证逻辑
type Email struct {
	value string
}

// NewEmail 创建Email值对象，自动验证格式
func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, errors.New("email cannot be empty")
	}
	
	if !isValidEmail(email) {
		return Email{}, errors.New("invalid email format")
	}
	
	return Email{value: email}, nil
}

// Value 返回邮箱值
func (e Email) Value() string {
	return e.value
}

// String 实现Stringer接口
func (e Email) String() string {
	return e.value
}

// Equals 比较两个Email是否相等
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
