package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Password 值对象 - 封装密码验证和加密逻辑
type Password struct {
	hashedValue string
}

// NewPassword 从明文密码创建Password值对象，自动加密
func NewPassword(plainPassword string) (Password, error) {
	if len(plainPassword) < 8 {
		return Password{}, errors.New("password must be at least 8 characters long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hashedValue: string(hash)}, nil
}

// NewPasswordFromHash 从已加密的密码创建Password值对象（用于从数据库加载）
func NewPasswordFromHash(hashedPassword string) Password {
	return Password{hashedValue: hashedPassword}
}

// Verify 验证明文密码是否匹配
func (p Password) Verify(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashedValue), []byte(plainPassword))
	return err == nil
}

// Hash 返回加密后的密码（用于持久化）
func (p Password) Hash() string {
	return p.hashedValue
}

// IsEmpty 检查密码是否为空
func (p Password) IsEmpty() bool {
	return p.hashedValue == ""
}
