package user

import "time"

// UserCreated 用户创建事件
type UserCreated struct {
	UserID    int64
	Username  string
	Email     string
	CreatedAt time.Time
}

// UserUpdated 用户更新事件
type UserUpdated struct {
	UserID    int64
	Username  string
	UpdatedAt time.Time
}

// UserDeleted 用户删除事件
type UserDeleted struct {
	UserID    int64
	DeletedAt time.Time
}

// UserActivated 用户激活事件
type UserActivated struct {
	UserID    int64
	UpdatedAt time.Time
}

// UserDeactivated 用户停用事件
type UserDeactivated struct {
	UserID    int64
	UpdatedAt time.Time
}

// PasswordChanged 密码修改事件
type PasswordChanged struct {
	UserID    int64
	UpdatedAt time.Time
}
