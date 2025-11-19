# DDD项目结构说明

## 目录树

```
goapp/
├── main.go                           # 应用入口
├── go.mod
├── go.sum
│
├── docs/                             # 文档
│   ├── ARCHITECTURE.md               # 整体架构文档
│   ├── DDD_ARCHITECTURE.md           # DDD架构设计
│   ├── DDD_MIGRATION_GUIDE.md        # DDD迁移指南
│   └── PROJECT_STRUCTURE.md          # 原项目结构
│
├── internal/                         # 内部代码
│   │
│   ├── domain/                       # 领域层 ⭐DDD核心
│   │   └── user/                     # User限界上下文
│   │       ├── user.go               # User实体（富血模型）
│   │       ├── email.go              # Email值对象
│   │       ├── password.go           # Password值对象
│   │       ├── repository.go         # 仓储接口定义
│   │       └── events.go             # 领域事件
│   │
│   ├── application/                  # 应用层
│   │   └── user_service.go           # 用户应用服务
│   │
│   ├── infrastructure/               # 基础设施层
│   │   └── persistence/              # 持久化
│   │       ├── user_po.go            # 用户持久化对象
│   │       ├── user_mapper.go        # 对象映射器
│   │       └── user_repository_impl.go # 仓储实现
│   │
│   ├── interfaces/                   # 接口层
│   │   └── http/                     # HTTP接口
│   │       └── user_controller.go    # 用户控制器
│   │
│   ├── app/                          # 应用配置（保留）
│   │   ├── config.go
│   │   ├── database.go
│   │   ├── logger.go
│   │   ├── redis.go
│   │   ├── services.go
│   │   ├── validator.go
│   │   └── errors/
│   │       ├── codes.go
│   │       └── error.go
│   │
│   ├── context/                      # 上下文（保留）
│   │   ├── api_context.go
│   │   └── request_id.go
│   │
│   ├── dto/                          # 数据传输对象（保留）
│   │   ├── product_dto.go
│   │   ├── response_dto.go
│   │   └── user_dto.go
│   │
│   ├── events/                       # 事件系统（保留）
│   │   ├── event_bus.go
│   │   └── events.go
│   │
│   ├── middleware/                   # 中间件（保留）
│   │   ├── auth.go
│   │   ├── events_middleware.go
│   │   ├── logger.go
│   │   └── response_formatter.go
│   │
│   ├── router/                       # 路由（保留）
│   │   └── router.go
│   │
│   ├── tasks/                        # 任务（保留）
│   │   └── tasks.go
│   │
│   ├── models/                       # 旧模型（待删除）
│   │   ├── product.go
│   │   └── user.go
│   │
│   ├── services/                     # 旧服务（待删除）
│   │   ├── monitor_service.go
│   │   ├── product_service.go
│   │   └── user_service.go
│   │
│   ├── repositories/                 # 旧仓储（待删除）
│   │   ├── product_repository.go
│   │   └── user_repository.go
│   │
│   └── controllers/                  # 旧控制器（待删除）
│       ├── monitor_controller.go
│       ├── product_controller.go
│       └── user_controller.go
│
└── utils/                            # 工具函数（保留）
    ├── array.go
    ├── convert.go
    ├── file.go
    ├── http.go
    ├── paginator.go
    ├── security.go
    ├── stringutil.go
    └── timeutil.go
```

## 分层说明

### 1. 领域层 (Domain Layer) 

**位置：** `internal/domain/`

**职责：**
- 核心业务逻辑
- 领域实体和值对象
- 领域服务
- 仓储接口定义
- 领域事件

**特点：**
- 不依赖任何外层
- 纯业务逻辑
- 高度可测试

**示例：**
```
domain/user/
├── user.go          # 实体：用户聚合根
├── email.go         # 值对象：邮箱
├── password.go      # 值对象：密码
├── repository.go    # 仓储接口
└── events.go        # 领域事件
```

### 2. 应用层 (Application Layer)

**位置：** `internal/application/`

**职责：**
- 业务流程编排
- 用例实现
- 事务管理
- 调用领域对象
- 发布领域事件

**特点：**
- 依赖领域层
- 不包含业务规则
- 协调领域对象完成任务

**示例：**
```
application/
├── user_service.go      # 用户应用服务
└── product_service.go   # 产品应用服务（待创建）
```

### 3. 基础设施层 (Infrastructure Layer)

**位置：** `internal/infrastructure/`

**职责：**
- 技术实现细节
- 数据库访问
- 外部服务调用
- 实现领域层定义的接口

**特点：**
- 实现仓储接口
- 处理技术复杂性
- 可替换实现

**示例：**
```
infrastructure/persistence/
├── user_po.go                # 持久化对象
├── user_mapper.go            # 领域对象 ↔ PO转换
├── user_repository_impl.go   # 仓储实现
└── product_repository_impl.go # （待创建）
```

### 4. 接口层 (Interfaces Layer)

**位置：** `internal/interfaces/`

**职责：**
- 对外API
- HTTP请求处理
- 数据格式转换
- 调用应用服务

**特点：**
- 依赖应用层
- 处理输入输出
- 格式化响应

**示例：**
```
interfaces/http/
├── user_controller.go     # 用户HTTP控制器
└── product_controller.go  # 产品HTTP控制器（待创建）
```

## 依赖关系

```
接口层 (Interfaces)
    ↓ 依赖
应用层 (Application)
    ↓ 依赖
领域层 (Domain) ⭐核心
    ↑ 实现
基础设施层 (Infrastructure)
```

**依赖倒置：**
- 领域层定义接口（Repository）
- 基础设施层实现接口
- 应用层依赖接口，不依赖具体实现

## User领域详解

### User实体 (domain/user/user.go)

```go
type User struct {
    id        int64
    username  string
    email     Email      // 值对象
    password  Password   // 值对象
    firstName string
    lastName  string
    isActive  bool
    isAdmin   bool
    createdAt time.Time
    updatedAt time.Time
}
```

**特点：**
- 字段私有（封装）
- 包含业务行为
- 通过方法访问和修改状态

**业务方法：**
- `Activate()` - 激活用户
- `Deactivate()` - 停用用户
- `Authenticate(password)` - 密码认证
- `ChangePassword(newPassword)` - 修改密码
- `UpdateProfile(firstName, lastName)` - 更新资料
- 等等...

### Email值对象 (domain/user/email.go)

```go
type Email struct {
    value string  // 私有不可变
}

func NewEmail(email string) (Email, error) {
    // 创建时验证
    if !isValidEmail(email) {
        return Email{}, errors.New("invalid email")
    }
    return Email{value: email}, nil
}
```

**优势：**
- 类型安全
- 验证集中
- 不可变性
- 表达力强

### Password值对象 (domain/user/password.go)

```go
type Password struct {
    hashedValue string  // 存储加密后的值
}

func NewPassword(plainPassword string) (Password, error) {
    // 验证 + 加密
    if len(plainPassword) < 8 {
        return Password{}, errors.New("password too short")
    }
    hash, _ := bcrypt.GenerateFromPassword(...)
    return Password{hashedValue: string(hash)}, nil
}

func (p Password) Verify(plainPassword string) bool {
    // 验证密码
}
```

### Repository接口 (domain/user/repository.go)

```go
type Repository interface {
    Save(user *User) error
    FindByID(id int64) (*User, error)
    FindByEmail(email Email) (*User, error)
    // ...
}
```

**实现：** `infrastructure/persistence/user_repository_impl.go`

### 持久化对象 (infrastructure/persistence/user_po.go)

```go
type UserPO struct {
    ID        int64          `gorm:"primaryKey"`
    Username  string         `gorm:"uniqueIndex"`
    Email     string         `gorm:"uniqueIndex"`
    Password  string         `gorm:"size:255"`
    // ... GORM标签
}
```

**转换：** `user_mapper.go` 负责 `User` ↔ `UserPO` 转换

## 数据流示例

### 创建用户流程

```
1. HTTP请求
   POST /api/users
   Body: {"username": "john", "email": "john@example.com", ...}
   
2. Controller (interfaces/http/user_controller.go)
   解析请求 → DTO
   
3. Application Service (application/user_service.go)
   CreateUser() → 创建领域对象 → 验证 → 保存
   
4. Domain (domain/user/)
   NewUser() → 自动验证Email、Password
   
5. Repository Interface (domain/user/repository.go)
   Save(user)
   
6. Repository Impl (infrastructure/persistence/user_repository_impl.go)
   User → UserPO → GORM → Database
   
7. 响应
   User → DTO → JSON → HTTP Response
```

## 待迁移模块

### Product领域（下一步）

建议结构：
```
domain/product/
├── product.go       # Product实体
├── sku.go           # SKU值对象
├── price.go         # Price值对象
├── stock.go         # Stock值对象
├── repository.go    # 仓储接口
└── events.go        # 领域事件

application/
└── product_service.go

infrastructure/persistence/
├── product_po.go
├── product_mapper.go
└── product_repository_impl.go

interfaces/http/
└── product_controller.go
```

## 旧文件清理计划

测试通过后删除：

```
✓ internal/models/user.go
✓ internal/services/user_service.go  
✓ internal/repositories/user_repository.go
✓ internal/controllers/user_controller.go
```

保留（Product迁移后再删除）：
```
⏳ internal/models/product.go
⏳ internal/services/product_service.go
⏳ internal/repositories/product_repository.go
⏳ internal/controllers/product_controller.go
```

## 关键原则

### 1. 单向依赖
外层可以依赖内层，内层不能依赖外层

### 2. 依赖倒置
领域层定义接口，基础设施层实现接口

### 3. 封装性
- 领域对象字段私有
- 通过方法访问和修改
- 值对象不可变

### 4. 职责分离
- 领域层：业务规则
- 应用层：流程编排
- 基础设施层：技术实现
- 接口层：API处理

## 参考文档

- [DDD_ARCHITECTURE.md](./DDD_ARCHITECTURE.md) - DDD架构详解
- [DDD_MIGRATION_GUIDE.md](./DDD_MIGRATION_GUIDE.md) - 迁移指南
- [ARCHITECTURE.md](./ARCHITECTURE.md) - 整体架构
