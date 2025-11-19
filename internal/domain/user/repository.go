package user

// Repository 用户仓储接口 - 在领域层定义，在基础设施层实现
type Repository interface {
	// Save 保存用户（新建或更新），返回保存后的用户实体（包含生成的ID）
	Save(user *User) (*User, error)

	// FindByID 根据ID查找用户
	FindByID(id int64) (*User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(email Email) (*User, error)

	// FindByUsername 根据用户名查找用户
	FindByUsername(username string) (*User, error)

	// FindAll 查询用户列表（分页）
	FindAll(limit, offset int) ([]*User, error)

	// Delete 删除用户
	Delete(id int64) error

	// ExistsByEmail 检查邮箱是否已存在
	ExistsByEmail(email Email) (bool, error)

	// ExistsByUsername 检查用户名是否已存在
	ExistsByUsername(username string) (bool, error)
}
