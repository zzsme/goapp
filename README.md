# GoApp - DDD架构示例项目

> 采用领域驱动设计（DDD）的Go Web应用

## 🏗️ 项目结构

```
goapp/
├── main.go                      # 应用入口
├── go.mod                       # Go模块定义
│
├── docs/                        # 📚 文档
│   └── DDD_GUIDE.md            # DDD架构完整指南
│
├── internal/                    # 内部代码
│   │
│   ├── domain/                 # 🎯 领域层（核心业务逻辑）
│   │   └── user/
│   │       ├── user.go         # User实体（富血模型）
│   │       ├── email.go        # Email值对象
│   │       ├── password.go     # Password值对象
│   │       ├── repository.go   # 仓储接口
│   │       └── events.go       # 领域事件
│   │
│   ├── application/            # 📋 应用层（业务流程编排）
│   │   └── user_service.go
│   │
│   ├── infrastructure/         # 🔧 基础设施层（技术实现）
│   │   └── persistence/
│   │       ├── user_po.go
│   │       ├── user_mapper.go
│   │       └── user_repository_impl.go
│   │
│   ├── interfaces/             # 🌐 接口层（对外API）
│   │   └── http/
│   │       └── user_controller.go
│   │
│   ├── app/                    # ⚙️ 应用配置
│   │   ├── config.go
│   │   ├── database.go
│   │   ├── logger.go
│   │   └── errors/
│   │
│   ├── middleware/             # 🔀 中间件
│   ├── router/                 # 🛣️ 路由
│   ├── events/                 # 📡 事件系统
│   ├── dto/                    # 📦 数据传输对象
│   ├── context/                # 🔗 上下文
│   ├── tasks/                  # ⏰ 任务
│   │
│   └── _legacy/                # 🗂️ 旧代码（待清理）
│       ├── models/
│       ├── services/
│       ├── repositories/
│       └── controllers/
│
└── utils/                      # 🛠️ 工具函数
```

## 🎯 DDD四层架构

```
┌─────────────────────────────────────┐
│   接口层 (Interfaces)                │
│   处理HTTP请求、数据格式转换           │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   应用层 (Application)               │
│   编排业务流程、调用领域对象           │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   领域层 (Domain) ⭐                 │
│   核心业务逻辑、业务规则              │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   基础设施层 (Infrastructure)         │
│   数据库访问、外部服务调用             │
└───────────────────────────────────────┘
```

## ✨ 核心特性

### 富血模型（Rich Domain Model）
实体包含数据和业务行为，不是简单的数据容器：

```go
type User struct {
    id       int64
    email    Email      // 值对象
    password Password   // 值对象
}

// 业务行为
func (u *User) Activate()
func (u *User) Authenticate(password string) bool
func (u *User) ChangePassword(newPassword string) error
```

### 值对象（Value Objects）
封装验证逻辑，类型安全：

```go
type Email struct {
    value string
}

func NewEmail(email string) (Email, error) {
    // 创建时自动验证
    if !isValidEmail(email) {
        return Email{}, errors.New("invalid email")
    }
    return Email{value: email}, nil
}
```

### 依赖倒置（Dependency Inversion）
领域层定义接口，基础设施层实现：

```go
// 领域层定义
type Repository interface {
    Save(user *User) (*User, error)
    FindByID(id int64) (*User, error)
}

// 基础设施层实现
type UserRepositoryImpl struct {
    db *gorm.DB
}
```

## 🚀 快速开始

### 安装依赖
```bash
go mod download
```

### 运行项目
```bash
go run main.go
```

### 运行测试
```bash
go test ./...
```

## 📖 文档

- [DDD架构完整指南](docs/DDD_GUIDE.md) - 详细的DDD架构说明

## 🔄 迁移状态

### ✅ 已完成
- User领域 - 完全采用DDD架构

### ⏳ 进行中
- Product领域 - 计划迁移
- Order领域 - 待规划

## 🎓 学习资源

### 核心概念
- **实体（Entity）** - 有唯一标识的业务对象
- **值对象（Value Object）** - 没有唯一标识的不可变对象
- **聚合（Aggregate）** - 一组相关对象的集合
- **仓储（Repository）** - 持久化接口
- **领域服务（Domain Service）** - 跨实体的业务逻辑
- **应用服务（Application Service）** - 业务流程编排

### 推荐阅读
- 《领域驱动设计》- Eric Evans
- 《实现领域驱动设计》- Vaughn Vernon

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 License

MIT License
