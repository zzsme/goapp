package application

import (
	"errors"
	"fmt"
	"strings"

	"goapp/internal/domain/user"
	"goapp/internal/events"
)

// UserService 用户应用服务 - 编排业务流程
type UserService struct {
	userRepo user.Repository
}

// NewUserService 创建用户应用服务
func NewUserService(userRepo user.Repository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(username, email, password, firstName, lastName string) (*user.User, error) {
	// 1. 创建领域对象（自动验证）
	userEntity, err := user.NewUser(username, email, password, firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	// 2. 检查邮箱唯一性
	emailVO, _ := user.NewEmail(email)
	exists, err := s.userRepo.ExistsByEmail(emailVO)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, errors.New("email is already taken")
	}

	// 3. 检查用户名唯一性
	exists, err = s.userRepo.ExistsByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, errors.New("username is already taken")
	}

	// 4. 保存用户
	savedUser, err := s.userRepo.Save(userEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// 5. 发布领域事件
	events.Publish(events.UserCreated, user.UserCreated{
		UserID:    savedUser.ID(),
		Username:  savedUser.Username(),
		Email:     savedUser.EmailValue(),
		CreatedAt: savedUser.CreatedAt(),
	})

	return savedUser, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id int64, username, email, firstName, lastName *string, isActive *bool) (*user.User, error) {
	// 1. 查找用户
	userEntity, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 2. 更新用户名
	if username != nil && *username != "" && *username != userEntity.Username() {
		// 检查新用户名是否已被使用
		exists, err := s.userRepo.ExistsByUsername(*username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username existence: %w", err)
		}
		if exists {
			return nil, errors.New("username is already taken")
		}
		if err := userEntity.UpdateUsername(*username); err != nil {
			return nil, err
		}
	}

	// 3. 更新邮箱
	if email != nil && *email != "" && *email != userEntity.EmailValue() {
		// 检查新邮箱是否已被使用
		emailVO, err := user.NewEmail(*email)
		if err != nil {
			return nil, err
		}
		exists, err := s.userRepo.ExistsByEmail(emailVO)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return nil, errors.New("email is already taken")
		}
		if err := userEntity.UpdateEmail(*email); err != nil {
			return nil, err
		}
	}

	// 4. 更新资料
	if firstName != nil || lastName != nil {
		fn := userEntity.FirstName()
		ln := userEntity.LastName()
		if firstName != nil {
			fn = *firstName
		}
		if lastName != nil {
			ln = *lastName
		}
		userEntity.UpdateProfile(fn, ln)
	}

	// 5. 更新激活状态
	if isActive != nil {
		if *isActive {
			userEntity.Activate()
		} else {
			userEntity.Deactivate()
		}
	}

	// 6. 保存更新
	updatedUser, err := s.userRepo.Save(userEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// 7. 发布事件
	events.Publish(events.UserUpdated, user.UserUpdated{
		UserID:    updatedUser.ID(),
		Username:  updatedUser.Username(),
		UpdatedAt: updatedUser.UpdatedAt(),
	})

	return updatedUser, nil
}

// GetUser 获取用户
func (s *UserService) GetUser(id int64) (*user.User, error) {
	userEntity, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return userEntity, nil
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(page, pageSize int) ([]*user.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.userRepo.FindAll(pageSize, offset)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int64) error {
	// 检查用户是否存在
	userEntity, err := s.userRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// 删除用户
	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// 发布事件
	events.Publish(events.UserDeleted, user.UserDeleted{
		UserID:    id,
		DeletedAt: userEntity.UpdatedAt(),
	})

	return nil
}

// AuthenticateUser 用户认证
func (s *UserService) AuthenticateUser(login, password string) (*user.User, error) {
	var userEntity *user.User
	var err error

	// 判断是邮箱还是用户名登录
	if strings.Contains(login, "@") {
		emailVO, err := user.NewEmail(login)
		if err != nil {
			return nil, errors.New("invalid email format")
		}
		userEntity, err = s.userRepo.FindByEmail(emailVO)
	} else {
		userEntity, err = s.userRepo.FindByUsername(login)
	}

	if err != nil {
		return nil, errors.New("user not found")
	}

	// 验证密码
	if !userEntity.Authenticate(password) {
		return nil, errors.New("invalid password")
	}

	// 检查用户是否激活
	if !userEntity.IsActive() {
		return nil, errors.New("user is not active")
	}

	return userEntity, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id int64, oldPassword, newPassword string) error {
	// 查找用户
	userEntity, err := s.userRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// 验证旧密码
	if !userEntity.Authenticate(oldPassword) {
		return errors.New("invalid old password")
	}

	// 修改密码（领域对象内部验证）
	if err := userEntity.ChangePassword(newPassword); err != nil {
		return err
	}

	// 保存
	updatedUser, err := s.userRepo.Save(userEntity)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	// 发布事件
	events.Publish(events.PasswordChanged, user.PasswordChanged{
		UserID:    updatedUser.ID(),
		UpdatedAt: updatedUser.UpdatedAt(),
	})

	return nil
}
