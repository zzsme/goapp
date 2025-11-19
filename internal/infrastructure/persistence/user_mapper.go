package persistence

import (
	"goapp/internal/domain/user"
)

// UserMapper 领域对象与持久化对象之间的转换器
type UserMapper struct{}

// ToEntity 将持久化对象转换为领域实体
func (m *UserMapper) ToEntity(po *UserPO) (*user.User, error) {
	// 重建Email值对象
	email, err := user.NewEmail(po.Email)
	if err != nil {
		return nil, err
	}

	// 从哈希值重建Password值对象
	password := user.NewPasswordFromHash(po.Password)

	// 使用Reconstruct重建领域对象
	return user.Reconstruct(
		po.ID,
		po.Username,
		email,
		password,
		po.FirstName,
		po.LastName,
		po.IsActive,
		po.IsAdmin,
		po.CreatedAt,
		po.UpdatedAt,
	), nil
}

// ToPO 将领域实体转换为持久化对象
func (m *UserMapper) ToPO(entity *user.User) *UserPO {
	return &UserPO{
		ID:        entity.ID(),
		Username:  entity.Username(),
		Email:     entity.EmailValue(),
		Password:  entity.Password().Hash(),
		FirstName: entity.FirstName(),
		LastName:  entity.LastName(),
		IsActive:  entity.IsActive(),
		IsAdmin:   entity.IsAdmin(),
		CreatedAt: entity.CreatedAt(),
		UpdatedAt: entity.UpdatedAt(),
	}
}

// ToEntityList 批量转换为领域实体列表
func (m *UserMapper) ToEntityList(poList []*UserPO) ([]*user.User, error) {
	entities := make([]*user.User, 0, len(poList))
	for _, po := range poList {
		entity, err := m.ToEntity(po)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}
	return entities, nil
}
