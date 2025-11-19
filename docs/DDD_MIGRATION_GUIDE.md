# DDD架构迁移指南

## 概述

本文档说明如何将项目从传统分层架构迁移到轻量级DDD架构。

## 已完成的重构

### ✅ User领域重构完成

**新增文件：**

```
internal/
├── domain/user/                      # 领域层
│   ├── user.go                       # User实体（富血模型）
│   ├── email.go                      # Email值对象
│   ├── password.go                   # Password值对象
│   ├── repository.go                 # 仓储接口
│   └── events.go                     # 领域事件定义
│
├── application/                      # 应用层
│   └── user_service.go               # 用户应用服务
│
├── infrastructure/persistence/       # 基础设施层
│   ├── user_po.go                    # 持久化对象
│   ├── user_mapper.go                # 对象转换器
│   └── user_repository_impl.go       # 仓储实现
│
└── interfaces/http/                  # 接口层
    └── user_controller.go            # HTTP控制器
```

## 新旧架构对比

### 旧架构（传统分层）

```go
// internal/models/user.go - 贫血模型
type User struct {
    ID       int64
    Email    string
    Password string
}

// internal/services/user_service.go - 业务逻辑分散
func (s *UserService) CreateUser(user *models.User) error {
    // 验证逻辑在这里
    if !isValidEmail(user.Email) {
        return errors.New("invalid email")
    }
    
    // 密码加密在这里
    hash, _ := bcrypt.GenerateFromPassword(...)
    user.Password = string(hash)
    
    // 保存
    return s.repo.Create(user)
}
```

### 新架构（轻量DDD）

```go
// internal/domain/user/user.go - 富血模型
type User struct {
    id       int64
    email    Email      // 值对象（自动验证）
    password Password   // 值对象（自动加密）
}

// 业务行为在实体内部
func (u *User) Activate() {
    u.isActive = true
}

func (u *User) Authenticate(password string) bool {
    return u.password.Verify(password)
}

// internal/application/user_service.go - 流程编排
func (s *UserService) CreateUser(username, email, password string) (*user.User, error) {
    // 创建领域对象（自动验证）
    userEntity, err := user.NewUser(username, email, password)
    if err != nil {
        return nil, err
    }
    
    // 业务规则检查
    exists, _ := s.userRepo.ExistsByEmail(userEntity.Email())
    if exists {
        return nil, errors.New("email already taken")
    }
    
    // 保存
    s.userRepo.Save(userEntity)
    
    return userEntity, nil
}
```

## 核心改进点

### 1. 富血模型 vs 贫血模型

**旧方式（贫血）：**
```go
// 数据和行为分离
user := &models.User{Email: "test@example.com"}
userService.ValidateEmail(user.Email)  // 验证在外部
userService.HashPassword(user)         // 加密在外部
```

**新方式（富血）：**
```go
// 数据和行为封装
email, _ := user.NewEmail("test@example.com") // 创建时自动验证
password, _ := user.NewPassword("secret123")  // 创建时自动加密
user := user.NewUser("john", email, password)
```

### 2. 值对象封装验证

**旧方式：**
```go
// 验证逻辑散落各处
func CreateUser(email string) error {
    if !regexp.Match(..., email) {
        return errors.New("invalid email")
    }
    // ...
}

func UpdateUser(email string) error {
    if !regexp.Match(..., email) {  // 重复的验证
        return errors.New("invalid email")
    }
    // ...
}
```

**新方式：**
```go
// 验证逻辑集中在值对象
type Email struct {
    value string
}

func NewEmail(email string) (Email, error) {
    if !isValidEmail(email) {
        return Email{}, errors.New("invalid email")
    }
    return Email{value: email}, nil
}

// 使用时自动验证
email, err := user.NewEmail("test@example.com")
```

### 3. 依赖倒置

**旧方式：**
```go
// Service直接依赖具体实现
type UserService struct {
    repo *repositories.GormUserRepository  // 具体实现
}
```

**新方式：**
```go
// 领域层定义接口
// domain/user/repository.go
type Repository interface {
    Save(user *User) error
    FindByID(id int64) (*User, error)
}

// 应用层依赖接口
type UserService struct {
    userRepo user.Repository  // 接口
}

// 基础设施层实现接口
type UserRepositoryImpl struct {
    db *gorm.DB
}
```

## 使用示例

### 创建用户

```go
// 初始化仓储
db := app.GetDB()
userRepo := persistence.NewUserRepository(db)
eventBus := events.GetEventBus()

// 初始化应用服务
userService := application.NewUserService(userRepo, eventBus)

// 创建用户
userEntity, err := userService.CreateUser(
    "johndoe",
    "john@example.com",
    "password123",
    "John",
    "Doe",
)
if err != nil {
    // 处理错误
}

// 使用领域对象
fmt.Println(userEntity.Username())
fmt.Println(userEntity.EmailValue())
```

### 用户认证

```go
// 认证
userEntity, err := userService.AuthenticateUser(
    "john@example.com",
    "password123",
)
if err != nil {
    // 认证失败
}

// 检查权限
if userEntity.IsAdmin() {
    // 管理员操作
}
```

### 更新用户

```go
// 更新用户信息
username := "newalias"
email := "newemail@example.com"

updatedUser, err := userService.UpdateUser(
    userId,
    &username,  // 可选参数
    &email,     // 可选参数
    nil,        // firstName不更新
    nil,        // lastName不更新
    nil,        // isActive不更新
)
```

### 修改密码

```go
err := userService.ChangePassword(
    userId,
    "old password",
    "new password",
)
```

## 待迁移内容

### ⏳ Product领域（待重构）

当前Product仍使用旧架构，建议按以下步骤重构：

1. 创建 `internal/domain/product/` 目录
2. 实现Product实体和值对象（SKU、Price、Stock等）
3. 定义Product仓储接口
4. 实现Product持久化层
5. 重构Product应用服务

### ⏳ 旧文件清理

重构完成后可以删除：
- `internal/models/user.go` （已被领域对象替代）
- `internal/services/user_service.go` （已被应用服务替代）
- `internal/repositories/user_repository.go` （已被新仓储替代）
- `internal/controllers/user_controller.go` （已被新控制器替代）

**注意：** 暂时保留这些文件以确保兼容性，待测试通过后再删除。

## 测试策略

### 单元测试

**领域对象测试（不依赖基础设施）：**
```go
func TestUser_Authenticate(t *testing.T) {
    user, _ := user.NewUser("john", "john@example.com", "password123")
    
    // 测试正确密码
    assert.True(t, user.Authenticate("password123"))
    
    // 测试错误密码
    assert.False(t, user.Authenticate("wrongpassword"))
}

func TestEmail_Validation(t *testing.T) {
    // 测试有效邮箱
    _, err := user.NewEmail("valid@example.com")
    assert.NoError(t, err)
    
    // 测试无效邮箱
    _, err = user.NewEmail("invalid-email")
    assert.Error(t, err)
}
```

**应用服务测试（Mock仓储）：**
```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Save(u *user.User) error {
    args := m.Called(u)
    return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockRepo.On("ExistsByEmail", mock.Anything).Return(false, nil)
    mockRepo.On("Save", mock.Anything).Return(nil)
    
    service := application.NewUserService(mockRepo, eventBus)
    
    userEntity, err := service.CreateUser("john", "john@example.com", "password123", "", "")
    assert.NoError(t, err)
    assert.NotNil(t, userEntity)
}
```

## 常见问题

### Q: 为什么需要PO和领域对象分离？
A: 领域对象关注业务逻辑，PO关注数据持久化。分离避免GORM标签污染领域模型，使领域层更纯粹。

### Q: 值对象增加了代码复杂度，值得吗？
A: 值得。虽然代码量略增，但收益明显：
- 类型安全（编译时发现错误）
- 验证集中（不会遗漏）
- 代码可读性更高
- 易于重构和测试

### Q: 旧代码如何兼容？
A: 当前新旧代码可以并存。建议：
1. 新功能使用新架构
2. 旧功能逐步迁移
3. 测试通过后删除旧代码

### Q: 一定要用DDD吗？
A: 不一定。DDD适合：
- 业务逻辑复杂的项目
- 需要长期维护的项目
- 团队规模较大的项目

简单CRUD可能不需要完整DDD。

## 下一步计划

1. [ ] **测试新User领域**
   - 编写单元测试
   - 编写集成测试
   - 性能测试

2. [ ] **重构Product领域**
   - 按User模式重构
   - 实现Product值对象

3. [ ] **清理旧代码**
   - 删除旧User相关文件
   - 更新文档

4. [ ] **扩展功能**
   - 实现Order领域
   - 跨领域协作

## 参考文档

- [DDD_ARCHITECTURE.md](./DDD_ARCHITECTURE.md) - DDD架构设计详解
- [ARCHITECTURE.md](./ARCHITECTURE.md) - 整体架构文档
