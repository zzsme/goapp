### DDD轻量级架构设计文档

## 概述

本项目采用**轻量级DDD（领域驱动设计）架构**，在保持简单性的同时引入DDD核心思想。

## 架构分层

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
│   - 领域服务 (Domain Services)       │
│   - 仓储接口 (Repository Interfaces) │
│   - 领域事件 (Domain Events)         │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   基础设施层 (Infrastructure Layer)   │
│   - 仓储实现 (Repository Impl)        │
│   - 持久化对象 (PO)                   │
│   - 数据库访问                        │
│   - 外部服务                          │
└───────────────────────────────────────┘
```

## 目录结构

```
internal/
├── domain/                           # 领域层（核心业务逻辑）
│   └── user/                         # User限界上下文
│       ├── user.go                   # User实体（富血模型）
│       ├── email.go                  # Email值对象
│       ├── password.go               # Password值对象
│       ├── repository.go             # 仓储接口
│       └── events.go                 # 领域事件
│
├── application/                      # 应用层（业务流程编排）
│   └── user_service.go               # 用户应用服务
│
├── infrastructure/                   # 基础设施层（技术实现）
│   └── persistence/
│       ├── user_po.go                # 用户持久化对象
│       ├── user_mapper.go            # 领域对象↔持久化对象转换
│       └── user_repository_impl.go   # 仓储接口实现
│
└── interfaces/                       # 接口层（对外API）
    └── http/
        └── user_controller.go        # HTTP控制器
```

## 核心概念

### 1. 领域实体 (Entity)

**特点：**
- 富血模型（包含数据和业务行为）
- 封装业务规则
- 通过方法修改状态

**示例：**
```go
type User struct {
    id        int64
    username  string
    email     Email      // 值对象
    password  Password   // 值对象
    isActive  bool
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

### 2. 值对象 (Value Object)

**特点：**
- 不可变对象
- 无标识符
- 封装验证逻辑
- 通过值比较相等性

**示例：**
```go
type Email struct {
    value string
}

func NewEmail(email string) (Email, error) {
    if !isValidEmail(email) {
        return Email{}, errors.New("invalid email format")
    }
    return Email{value: email}, nil
}
```

**优势：**
- ✅ 验证逻辑集中在一处
- ✅ 类型安全（不会误传string）
- ✅ 代码更具表达力

### 3. 仓储接口 (Repository Interface)

**特点：**
- 在领域层定义接口
- 在基础设施层实现
- 依赖倒置原则

```go
// domain/user/repository.go
type Repository interface {
    Save(user *User) error
    FindByID(id int64) (*User, error)
    FindByEmail(email Email) (*User, error)
}

// infrastructure/persistence/user_repository_impl.go
type UserRepositoryImpl struct {
    db *gorm.DB
}

func (r *UserRepositoryImpl) Save(user *User) error {
    // 实现细节
}
```

### 4. 应用服务 (Application Service)

**职责：**
- 编排业务流程
- 调用领域对象
- 事务管理
- 发布领域事件

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
    s.userRepo.Save(userEntity)
    
    // 4. 发布事件
    s.eventBus.Publish(events.UserCreated, ...)
    
    return userEntity, nil
}
```

### 5. 持久化对象 (PO) vs 领域对象

**分离原因：**
- 领域对象关注业务逻辑
- 持久化对象关注数据存储
- 避免GORM标签污染领域模型

```go
// 持久化对象（带GORM标签）
type UserPO struct {
    ID       int64  `gorm:"primaryKey"`
    Email    string `gorm:"uniqueIndex"`
    Password string `gorm:"size:255"`
}

// 领域对象（纯业务逻辑）
type User struct {
    id       int64
    email    Email
    password Password
}

// Mapper负责转换
func (m *UserMapper) ToEntity(po *UserPO) (*User, error)
func (m *UserMapper) ToPO(entity *User) *UserPO
```

## 数据流

### 创建用户流程

```
HTTP Request
    ↓
Controller (interfaces/http/user_controller.go)
    ├─ 解析请求 (DTO)
    ↓
Application Service (application/user_service.go)
    ├─ 创建领域对象
    ├─ 业务规则验证
    ├─ 调用仓储保存
    ├─ 发布领域事件
    ↓
Repository Interface (domain/user/repository.go)
    ↓
Repository Impl (infrastructure/persistence/user_repository_impl.go)
    ├─ 领域对象 → PO
    ├─ 保存到数据库
    ↓
Database
```

## DDD核心优势

### ✅ 业务逻辑清晰
- 业务规则在领域对象内部
- 代码即文档

### ✅ 类型安全
- Email、Password等值对象
- 编译时发现错误

### ✅ 易于测试
- 领域对象不依赖基础设施
- 可以mock仓储接口

### ✅ 易于维护
- 职责分离
- 修改影响范围小

### ✅ 扩展性好
- 新增领域概念容易
- 支持复杂业务场景

## 与传统分层对比

| 维度 | 传统分层 | 轻量DDD |
|-----|---------|---------|
| Model | 贫血模型 | 富血模型 |
| 验证 | Service层散乱 | 值对象集中 |
| 职责 | 技术驱动 | 业务驱动 |
| 依赖 | 向下依赖 | 依赖倒置 |
| 测试 | 难以单测 | 易于单测 |
| 复杂度 | ⭐⭐ | ⭐⭐⭐ |
| 可维护性 | ⭐⭐ | ⭐⭐⭐⭐ |

## 最佳实践

### 1. 领域对象私有字段
```go
type User struct {
    id    int64  // 私有
    email Email  // 私有
}

// 通过方法访问
func (u *User) ID() int64 { return u.id }
```

### 2. 值对象不可变
```go
type Email struct {
    value string // 私有且不可修改
}

// 只提供读取方法
func (e Email) Value() string { return e.value }
```

### 3. 构造函数验证
```go
func NewUser(username, email, password string) (*User, error) {
    // 在创建时就验证
    emailVO, err := NewEmail(email)
    if err != nil {
        return nil, err
    }
    // ...
}
```

### 4. 领域事件解耦
```go
// 创建用户后发布事件
s.eventBus.Publish(events.UserCreated, user.UserCreated{
    UserID: userEntity.ID(),
    Email:  userEntity.EmailValue(),
})
```

## 扩展建议

### 阶段1：巩固User领域
- [ ] 完善用户相关业务规则
- [ ] 添加更多值对象（如UserProfile）
- [ ] 增强领域事件

### 阶段2：重构Product领域
- [ ] 按相同模式重构Product
- [ ] 添加Product值对象（SKU、Price等）
- [ ] 实现Product仓储

### 阶段3：跨领域协作
- [ ] 实现Order领域
- [ ] 处理领域间的依赖关系
- [ ] 使用领域事件解耦

### 阶段4：高级模式（可选）
- [ ] CQRS（命令查询分离）
- [ ] Event Sourcing（事件溯源）
- [ ] 聚合根设计

## 参考资料

- 《领域驱动设计》- Eric Evans
- 《实现领域驱动设计》- Vaughn Vernon
- [DDD Community](https://www.domainlanguage.com/)
