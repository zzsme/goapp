package persistence

import (
	"errors"
	"fmt"

	"goapp/internal/domain/user"

	"gorm.io/gorm"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	db     *gorm.DB
	mapper *UserMapper
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepositoryImpl{
		db:     db,
		mapper: &UserMapper{},
	}
}

// Save 保存用户（新建或更新），返回保存后的实体
func (r *UserRepositoryImpl) Save(entity *user.User) (*user.User, error) {
	po := r.mapper.ToPO(entity)

	// 如果ID为0，表示新建；否则为更新
	if po.ID == 0 {
		result := r.db.Create(po)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to create user: %w", result.Error)
		}
		// 重新从数据库加载以获取生成的ID和其他字段
		return r.mapper.ToEntity(po)
	}

	result := r.db.Save(po)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user with ID %d not found", po.ID)
	}
	
	// 返回更新后的实体
	return r.mapper.ToEntity(po)
}

// FindByID 根据ID查找用户
func (r *UserRepositoryImpl) FindByID(id int64) (*user.User, error) {
	var po UserPO
	result := r.db.First(&po, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return r.mapper.ToEntity(&po)
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepositoryImpl) FindByEmail(email user.Email) (*user.User, error) {
	var po UserPO
	result := r.db.Where("email = ?", email.Value()).First(&po)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found", email.Value())
		}
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return r.mapper.ToEntity(&po)
}

// FindByUsername 根据用户名查找用户
func (r *UserRepositoryImpl) FindByUsername(username string) (*user.User, error) {
	var po UserPO
	result := r.db.Where("username = ?", username).First(&po)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return r.mapper.ToEntity(&po)
}

// FindAll 查询用户列表（分页）
func (r *UserRepositoryImpl) FindAll(limit, offset int) ([]*user.User, error) {
	var poList []*UserPO
	result := r.db.Offset(offset).Limit(limit).Order("id").Find(&poList)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find users: %w", result.Error)
	}

	return r.mapper.ToEntityList(poList)
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(id int64) error {
	result := r.db.Delete(&UserPO{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}
	return nil
}

// ExistsByEmail 检查邮箱是否已存在
func (r *UserRepositoryImpl) ExistsByEmail(email user.Email) (bool, error) {
	var count int64
	result := r.db.Model(&UserPO{}).Where("email = ?", email.Value()).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check email existence: %w", result.Error)
	}
	return count > 0, nil
}

// ExistsByUsername 检查用户名是否已存在
func (r *UserRepositoryImpl) ExistsByUsername(username string) (bool, error) {
	var count int64
	result := r.db.Model(&UserPO{}).Where("username = ?", username).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check username existence: %w", result.Error)
	}
	return count > 0, nil
}
