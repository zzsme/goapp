package persistence

import (
	"time"

	"gorm.io/gorm"
)

// UserPO 用户持久化对象 - 对应数据库表结构
type UserPO struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:100;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:255;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"` // 加密后的密码
	FirstName string         `json:"first_name" gorm:"size:100"`
	LastName  string         `json:"last_name" gorm:"size:100"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	IsAdmin   bool           `json:"is_admin" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // 软删除支持
}

// TableName 指定表名
func (UserPO) TableName() string {
	return "users"
}
