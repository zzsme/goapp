package user

import (
	"errors"
	"time"
)

// User 领域实体 - 富血模型，包含业务行为
type User struct {
	id        int64
	username  string
	email     Email
	password  Password
	firstName string
	lastName  string
	isActive  bool
	isAdmin   bool
	createdAt time.Time
	updatedAt time.Time
}

// NewUser 创建新用户（用于注册场景）
func NewUser(username, email, password, firstName, lastName string) (*User, error) {
	// 验证用户名
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	// 创建Email值对象（自动验证）
	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	// 创建Password值对象（自动验证和加密）
	passwordVO, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		username:  username,
		email:     emailVO,
		password:  passwordVO,
		firstName: firstName,
		lastName:  lastName,
		isActive:  true,
		isAdmin:   false,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Reconstruct 从持久化数据重建User对象（用于从数据库加载）
func Reconstruct(
	id int64,
	username string,
	email Email,
	password Password,
	firstName, lastName string,
	isActive, isAdmin bool,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:        id,
		username:  username,
		email:     email,
		password:  password,
		firstName: firstName,
		lastName:  lastName,
		isActive:  isActive,
		isAdmin:   isAdmin,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// 业务行为方法

// Activate 激活用户
func (u *User) Activate() {
	u.isActive = true
	u.updatedAt = time.Now()
}

// Deactivate 停用用户
func (u *User) Deactivate() {
	u.isActive = false
	u.updatedAt = time.Now()
}

// PromoteToAdmin 提升为管理员
func (u *User) PromoteToAdmin() {
	u.isAdmin = true
	u.updatedAt = time.Now()
}

// DemoteFromAdmin 取消管理员权限
func (u *User) DemoteFromAdmin() {
	u.isAdmin = false
	u.updatedAt = time.Now()
}

// ChangePassword 修改密码
func (u *User) ChangePassword(newPassword string) error {
	pwd, err := NewPassword(newPassword)
	if err != nil {
		return err
	}
	u.password = pwd
	u.updatedAt = time.Now()
	return nil
}

// Authenticate 验证密码
func (u *User) Authenticate(password string) bool {
	return u.password.Verify(password)
}

// UpdateProfile 更新用户资料
func (u *User) UpdateProfile(firstName, lastName string) {
	if firstName != "" {
		u.firstName = firstName
	}
	if lastName != "" {
		u.lastName = lastName
	}
	u.updatedAt = time.Now()
}

// UpdateUsername 更新用户名
func (u *User) UpdateUsername(username string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	u.username = username
	u.updatedAt = time.Now()
	return nil
}

// UpdateEmail 更新邮箱
func (u *User) UpdateEmail(email string) error {
	emailVO, err := NewEmail(email)
	if err != nil {
		return err
	}
	u.email = emailVO
	u.updatedAt = time.Now()
	return nil
}

// Getter 方法（只读访问）

// ID 返回用户ID
func (u *User) ID() int64 {
	return u.id
}

// Username 返回用户名
func (u *User) Username() string {
	return u.username
}

// Email 返回邮箱（返回值对象）
func (u *User) Email() Email {
	return u.email
}

// EmailValue 返回邮箱字符串
func (u *User) EmailValue() string {
	return u.email.Value()
}

// Password 返回密码（返回值对象）
func (u *User) Password() Password {
	return u.password
}

// FirstName 返回名字
func (u *User) FirstName() string {
	return u.firstName
}

// LastName 返回姓氏
func (u *User) LastName() string {
	return u.lastName
}

// IsActive 返回是否激活
func (u *User) IsActive() bool {
	return u.isActive
}

// IsAdmin 返回是否为管理员
func (u *User) IsAdmin() bool {
	return u.isAdmin
}

// CreatedAt 返回创建时间
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt 返回更新时间
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
