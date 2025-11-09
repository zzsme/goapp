# GoWK Framework

GoWK 是一个用 Go 语言编写的轻量级 Web 开发框架，专为构建高效、可扩展的 API 服务和 Web 应用而设计。

## 特性

- **模块化架构**：遵循清晰的目录结构和代码组织方式，便于维护和扩展
- **标准化错误处理**：统一的错误码系统和错误响应格式
- **中间件系统**：灵活的请求处理流程和功能扩展
- **事件驱动**：内置的事件系统用于系统内部的解耦通信
- **数据库集成**：集成 GORM，轻松处理数据库操作
- **缓存支持**：集成 Redis 客户端，方便实现数据缓存
- **配置管理**：灵活的配置加载系统，支持环境变量和配置文件
- **日志系统**：结构化的日志记录和管理
- **任务系统**：便于执行定时任务和一次性命令行任务
- **请求上下文**：丰富的请求上下文，简化接口编写

## 目录结构

```
├── config.yaml         # 配置文件
├── main.go             # 程序入口
├── go.mod              # Go模块定义
├── go.sum              # 依赖版本锁定
├── internal/           # 内部代码（不对外暴露）
│   ├── app/            # 应用核心组件
│   │   ├── config.go   # 配置管理
│   │   ├── database.go # 数据库连接
│   │   ├── logger.go   # 日志系统
│   │   ├── redis.go    # Redis连接
│   │   └── validator.go # 数据验证
│   ├── context/        # 请求上下文
│   │   └── api_context.go # API上下文
│   ├── controllers/    # 控制器
│   ├── dto/            # 数据传输对象
│   ├── events/         # 事件系统
│   ├── middleware/     # 中间件
│   ├── models/         # 数据模型
│   ├── repositories/   # 数据仓库
│   ├── router/         # 路由定义
│   ├── services/       # 业务服务
│   └── tasks/          # 任务系统
├── logs/               # 日志文件
└── utils/              # 工具函数
```

## 快速开始

### 安装

```bash
# 克隆模板仓库
git clone https://github.com/username/gowk-template.git myproject

# 进入项目目录
cd myproject

# 初始化 Go 模块
go mod edit -module github.com/yourusername/myproject

# 下载依赖
go mod tidy
```

### 运行

```bash
# 直接运行
go run main.go

# 或构建后运行
go build -o app
./app
```

### 创建 API 接口

1. 在 `internal/models` 中定义数据模型
2. 在 `internal/dto` 中定义请求和响应 DTO
3. 在 `internal/repositories` 中实现数据访问逻辑
4. 在 `internal/services` 中实现业务逻辑
5. 在 `internal/controllers` 中实现控制器
6. 在 `internal/router/router.go` 中注册路由

## 错误处理

GoWK 使用统一的错误码系统，由 `internal/app/errors` 包管理。错误码分为客户端错误（1xxx）和服务器错误（5xxx）两大类：

- **客户端错误**：如参数错误(1000)、未授权(1001)、访问禁止(1003)等
- **服务器错误**：如内部服务器错误(5000)、数据库错误(5001)等

### 标准错误响应格式

```json
{
  "code": 1000,
  "message": "Invalid parameters",
  "requestId": "bce68d7d-0ace-49dd-b3ed-2977xxxx"
}
```

### 在控制器中使用错误处理

```go
func (c *UserController) GetUser(ctx *context.APIContext) {
    userID := ctx.Param("id")
    if userID == "" {
        ctx.ErrorWithCode(errors.BadRequest, "User ID is required")
        return
    }
    
    user, err := c.userService.GetByID(userID)
    if err != nil {
        ctx.ErrorWithAppError(errors.Wrap(errors.Database, "Failed to get user", err))
        return
    }
    
    ctx.Success(user)
}
```

## 配置管理

GoWK 使用 `viper` 处理配置，支持从环境变量或配置文件读取配置。项目会自动加载项目根目录下的 `config.yaml` 文件。

## 使用中间件

```go
// 在 router.go 中
r.Use(middleware.ResponseFormatter())
r.Use(middleware.Logger())
r.Use(middleware.Recovery())

// 对特定路由组使用中间件
api := r.Group("/api")
api.Use(middleware.Auth())
{
    api.GET("/users", GetAPIContext(userController.List))
    // ...
}
```

## 事件系统

GoWK 提供了简单的事件总线系统用于系统内部组件间的解耦通信。

```go
// 注册事件处理器
events.Subscribe(events.UserCreated, func(data map[string]interface{}) {
    // 处理用户创建事件
})

// 发布事件
events.Publish(events.UserCreated, map[string]interface{}{
    "userID": user.ID,
    "email": user.Email,
})
```

## 命令行任务

GoWK 支持通过命令行执行一次性任务：

```bash
# 运行名为 "cleanup" 的任务
go run main.go task cleanup
```

## 许可证

MIT
