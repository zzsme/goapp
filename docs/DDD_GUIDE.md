# DDD架构完整指南

> 本项目已采用轻量级DDD（领域驱动设计）架构

---

## 目录

- [架构概述](#架构概述)
- [目录结构](#目录结构)
- [核心概念](#核心概念)
- [迁移对比](#迁移对比)
- [使用示例](#使用示例)
- [问题修复](#问题修复)
- [下一步计划](#下一步计划)

---

## 架构概述

### 四层架构

```
┌─────────────────────────────────────┐
│   接口层 (Interfaces Layer)          │
│   - HTTP Controllers                │
│   - DTO (数据传输对象)                │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   应用层 (Application Layer)         │
│   - Application Services            │
│   - 业务流程编排                      │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   领域层 (Domain Layer) ⭐核心        │
│   - 实体 (Entities)                  │
│   - 值对象 (Value Objects)           │
│   - 仓储接口 (Repository Interfaces) │
│   - 领域事件 (Domain Events)         │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   基础设施层 (Infrastructure Layer)   │
│   - 仓储实现 (Repository Impl)        │
│   - 持久化对象 (PO)                   │
│   - 数据库访问                        │
└───────────────────────────────────────┘
```

### 依赖关系

```
接口层 → 应用层 → 领域层 ← 基础设施层
```

- **外层依赖内层**：接口层、应用层、基础设施层都依赖领域层
- **依赖倒置**：领域层定义接口，基础设施层实现接口

---

## 目录结构

### DDD新增目录

```
internal/
├── domain/                    # 领域层（核心业务逻辑）
│   └── user/
│       ├── user.go           # User实体
│       ├── email.go          # Email值对象
│       ├── password.go       # Password值对象
│       ├── repository.go     # 仓储接口
│       └── events.go         # 领域事件
│
├── application/               # 应用层（业务流程编排）
│   └── user_service.go
│
├── infrastructure/            # 基础设施层（技术实现）
│   └── persistence/
│       ├── user_po.go
│       ├── user_mapper.go
│       └── user_repository_impl.go
│
└── interfaces/                # 接口层（对外API）
    └── http/
        └── user_controller.go
```

### 完整项目结构

```
goapp/
├── main.go
├── docs/
│   └── DDD_GUIDE.md          # 本文档
│
├── internal/
│   ├── domain/               # 🆕 DDD领域层
│   ├── application/          # 🆕 DDD应用层
│   ├── infrastructure/       # 🆕 DDD基础设施层
│   ├── interfaces/           # 🆕 DDD接口层
│   │
│   ├── app/                  # 应用配置
│   ├── context/              # 上下文
│   ├── dto/                  # 数据传输对象
│   ├── events/               # 事件系统
│   ├── middleware/           # 中间件
│   ├── router/               # 路由
│   ├── tasks/                # 任务
│   │
│   ├── models/               # ⏳ 旧模型（待清理）
│   ├── services/             # ⏳ 旧服务（待清理）
│   ├── repositories/         # ⏳ 旧仓储（待清理）
│   └── controllers/          # ⏳ 旧控制器（待清理）
│
└── utils/                    # 工具函数
```

---

## 核心概念

### 1. 实体（Entity）- 富血模型

**特点：**
- 包含数据和业务行为
- 字段私有，通过方法访问
- 封装业务规则

**示例：**
```go
type User struct {
    id       int64
    email    Email      // 值对象
    password Password   // 值对象
    isActive bool
}

// 业务行为
func (u *User) Activate() {
    u.isActive = true
    u.updatedAt = time.Now()
}

func (u *User) Authenticate(password string) bool {
    return u.password.Verify(password)
}
```

### 2. 值对象（Value Object）

**特点：**
- 不可变
- 封装验证逻辑
- 类型安全

**Email值对象：**
```go
type Email struct {
    value string  // 私有不可变
}

func NewEmail(email string) (Email, error) {
    if !isValidEmail(email) {
        return Email{}, errors.New("invalid email")
    }
    return Email{value: email}, nil
}

func (e Email) Value() string {
    return e.value
}
```

**优势：**
- ✅ 验证逻辑集中
- ✅ 类型安全（不会误传string）
- ✅ 代码更具表达力

### 3. 仓储（Repository）

**接口在领域层定义：**
```go
// internal/domain/user/repository.go
type Repository interface {
    Save(user *User) (*User, error)
    FindByID(id int64) (*User, error)
    FindByEmail(email Email) (*User, error)
}
```

**实现在基础设施层：**
```go
// internal/infrastructure/persistence/user_repository_impl.go
type UserRepositoryImpl struct {
    db     *gorm.DB
    mapper *UserMapper
}

func (r *UserRepositoryImpl) Save(entity *user.User) (*user.User, error) {
    // 领域对象 → PO → 数据库
    // ...
}
```

### 4. 应用服务（Application Service）

**职责：**
- 编排业务流程
- 调用领域对象
- 事务管理
- 发布领域事件

**示例：**
```go
func (s *UserService) CreateUser(username, email, password string) (*user.User, error) {
    // 1. 创建领域对象（自动验证）
    userEntity, err := user.NewUser(username, email, password)
    
    // 2. 业务规则检查
    exists, _ := s.userRepo.ExistsByEmail(userEntity.Email())
    if exists {
        return nil, errors.New("email already taken")
    }
    
    // 3. 保存
    savedUser, _ := s.userRepo.Save(userEntity)
    
    // 4. 发布事件
    s.eventBus.Publish(events.UserCreated, ...)
    
    return savedUser, nil
}
```

### 5. 持久化对象（PO）vs 领域对象

**为什么分离？**
- 领域对象：纯业务逻辑
- 持久化对象：数据存储（GORM标签）
- 避免技术细节污染领域模型

```go
// PO（带GORM标签）
type UserPO struct {
    ID       int64  `gorm:"primaryKey"`
    Email    string `gorm:"uniqueIndex"`
    Password string `gorm:"size:255"`
}

// 领域对象（纯业务）
type User struct {
    id       int64
    email    Email
    password Password
}

// Mapper转换
func (m *UserMapper) ToEntity(po *UserPO) (*User, error)
func (m *UserMapper) ToPO(entity *User) *UserPO
```

---

## 迁移对比

### 旧架构（贫血模型）

```go
// models/user.go - 只有数据
type User struct {
    ID       int64
    Email    string
    Password string
}

// services/user_service.go - 业务逻辑分散
func (s *UserService) CreateUser(user *models.User) error {
    // 验证在这里
    if !isValidEmail(user.Email) {
        return errors.New("invalid email")
    }
    
    // 加密在这里
    hash, _ := bcrypt.GenerateFromPassword(...)
    user.Password = string(hash)
    
    return s.repo.Create(user)
}
```

### 新架构（DDD）

```go
// domain/user/user.go - 富血模型
type User struct {
    id       int64
    email    Email      // 值对象（自动验证）
    password Password   // 值对象（自动加密）
}

func (u *User) Authenticate(password string) bool {
    return u.password.Verify(password)
}

// application/user_service.go - 流程编排
func (s *UserService) CreateUser(username, email, password string) (*user.User, error) {
    // 创建领域对象（自动验证和加密）
    userEntity, _ := user.NewUser(username, email, password)
    
    // 业务规则检查
    exists, _ := s.userRepo.ExistsByEmail(userEntity.Email())
    if exists {
        return nil, errors.New("email already taken")
    }
    
    // 保存
    return s.userRepo.Save(userEntity)
}
```

---

## 使用示例

### 创建用户

```go
// 初始化
db := app.GetDB()
userRepo := persistence.NewUserRepository(db)
eventBus := events.GetEventBus()
userService := application.NewUserService(userRepo, eventBus)

// 创建用户
user, err := userService.CreateUser(
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
fmt.Println(user.ID())       // 数据库生成的ID
fmt.Println(user.Username()) // johndoe
fmt.Println(user.EmailValue()) // john@example.com
```

### 用户认证

```go
user, err := userService.AuthenticateUser(
    "john@example.com",
    "password123",
)
if err != nil {
    // 认证失败
}

if user.IsAdmin() {
    // 管理员操作
}
```

### 更新用户

```go
username := "newalias"
email := "new@example.com"

user, err := userService.UpdateUser(
    userId,
    &username,  // 可选
    &email,     // 可选
    nil,        // firstName不更新
    nil,        // lastName不更新
    nil,        // isActive不更新
)
```

---

## 问题修复

### 关键问题：Save方法返回值

**问题：**
原实现中`Save()`返回`error`，导致新创建的用户无法获取数据库生成的ID。

**影响：**
- 创建用户后ID为0
- 领域事件中UserID不正确

**解决方案：**

```go
// 修改前
Save(user *User) error

// 修改后
Save(user *User) (*User, error)
```

**实现：**

```go
func (r *UserRepositoryImpl) Save(entity *user.User) (*user.User, error) {
    po := r.mapper.ToPO(entity)

    if po.ID == 0 {
        // 新建
        r.db.Create(po)
        // 返回包含生成ID的实体
        return r.mapper.ToEntity(po)
    }

    // 更新
    r.db.Save(po)
    return r.mapper.ToEntity(po)
}
```

**应用服务调整：**

```go
// CreateUser
savedUser, err := s.userRepo.Save(userEntity)
if err != nil {
    return nil, err
}

// 发布事件时使用正确的ID
s.eventBus.Publish(events.UserCreated, user.UserCreated{
    UserID: savedUser.ID(),  // 现在有正确的ID
    // ...
})
```

**修改的文件：**
- `internal/domain/user/repository.go`
- `internal/infrastructure/persistence/user_repository_impl.go`
- `internal/application/user_service.go`

---

## 下一步计划

### 1. 测试

```go
// 领域对象单元测试
func TestUser_Authenticate(t *testing.T) {
    user, _ := user.NewUser("john", "john@example.com", "password123", "", "")
    
    assert.True(t, user.Authenticate("password123"))
    assert.False(t, user.Authenticate("wrongpassword"))
}

// 应用服务测试（Mock仓储）
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockRepo.On("Save", mock.Anything).Return(&user.User{...}, nil)
    
    service := NewUserService(mockRepo, eventBus)
    user, err := service.CreateUser(...)
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### 2. 重构Product领域

按相同模式重构：
```
domain/product/
├── product.go       # Product实体
├── sku.go           # SKU值对象
├── price.go         # Price值对象
├── repository.go    # 仓储接口
└── events.go        # 领域事件
```

### 3. 清理旧代码

测试通过后删除：
- `internal/models/user.go`
- `internal/services/user_service.go`
- `internal/repositories/user_repository.go`
- `internal/controllers/user_controller.go`

### 4. 扩展功能

- 实现Order领域
- 处理领域间协作
- 使用领域事件解耦

---

## 核心优势

### ✅ 业务逻辑清晰
业务规则在
