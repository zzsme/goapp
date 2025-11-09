# GoWK 项目结构详解

本文档详细说明GoWK框架的项目结构和各组件的职责。

## 顶层目录

| 目录/文件 | 描述 |
|---------|------|
| `main.go` | 应用程序入口点，负责初始化核心组件、事件系统、数据库连接等 |
| `go.mod` | Go模块定义文件，包含项目依赖管理 |
| `go.sum` | 依赖版本锁定文件 |
| `internal/` | 内部包目录，包含应用的主要代码 |
| `logs/` | 日志文件目录 |
| `utils/` | 通用工具函数库 |

## internal/ 目录

internal目录包含应用的所有核心代码，以下是各子目录的详细说明：

### app/ - 核心应用组件

核心应用组件包含应用程序的基础功能：

| 文件/目录 | 描述 |
|---------|------|
| `errors/` | 错误处理系统，定义错误码和错误类型 |
| `config.go` | 配置管理，处理应用的配置加载和访问 |
| `database.go` | 数据库初始化和连接管理 |
| `redis.go` | Redis初始化和连接管理 |
| `logger.go` | 日志管理，提供结构化日志记录 |
| `services.go` | 服务注册和管理 |
| `validator.go` | 请求验证器，处理输入验证 |

### context/ - 请求上下文

处理HTTP请求的上下文信息：

| 文件 | 描述 |
|-----|------|
| `api_context.go` | API上下文封装，提供标准化的响应格式 |
| `request_id.go` | 请求ID管理，用于请求跟踪 |

### controllers/ - API控制器

控制器层负责处理HTTP请求和返回响应：

| 文件 | 描述 |
|-----|------|
| `monitor_controller.go` | 监控相关API接口 |
| `product_controller.go` | 产品管理API接口 |
| `user_controller.go` | 用户管理API接口 |

### dto/ - 数据传输对象

定义请求和响应的数据结构：

| 文件 | 描述 |
|-----|------|
| `product_dto.go` | 产品相关的DTO定义 |
| `response_dto.go` | 通用响应DTO定义 |
| `user_dto.go` | 用户相关的DTO定义 |

### events/ - 事件系统

实现事件驱动架构：

| 文件 | 描述 |
|-----|------|
| `event_bus.go` | 事件总线实现，处理事件的发布和订阅 |
| `events.go` | 事件类型定义 |

### middleware/ - HTTP中间件

提供请求处理的中间件：

| 文件 | 描述 |
|-----|------|
| `auth.go` | 认证中间件，处理用户认证 |
| `events_middleware.go` | 事件中间件，记录请求事件 |
| `logger.go` | 日志中间件，记录请求日志 |
| `response_formatter.go` | 响应格式化中间件，统一响应格式 |

### models/ - 数据模型

定义数据库实体模型：

| 文件 | 描述 |
|-----|------|
| `product.go` | 产品数据模型 |
| `user.go` | 用户数据模型 |

### repositories/ - 数据访问层

处理数据库操作：

| 文件 | 描述 |
|-----|------|
| `product_repository.go` | 产品数据访问 |
| `user_repository.go` | 用户数据访问 |

### router/ - 路由管理

处理API路由注册：

| 文件 | 描述 |
|-----|------|
| `router.go` | 路由配置和注册 |

### services/ - 业务逻辑层

实现业务逻辑：

| 文件 | 描述 |
|-----|------|
| `monitor_service.go` | 监控服务实现 |
| `product_service.go` | 产品服务实现 |
| `user_service.go` | 用户服务实现 |

### tasks/ - 后台任务

处理后台定时任务和一次性任务：

| 文件 | 描述 |
|-----|------|
| `tasks.go` | 任务定义和执行 |

## utils/ 目录

通用工具函数库，包含各种辅助功能：

| 文件 | 描述 |
|-----|------|
| `array.go` | 数组和切片操作函数 |
| `convert.go` | 类型转换工具 |
| `file.go` | 文件操作工具 |
| `http.go` | HTTP相关工具 |
| `paginator.go` | 分页工具 |
| `security.go` | 安全相关功能（加密、哈希等） |
| `stringutil.go` | 字符串处理工具 |
| `timeutil.go` | 时间处理工具 |

## 依赖关系

GoWK框架采用了清晰的分层架构：

1. 控制器依赖服务层
2. 服务层依赖仓库层
3. 仓库层依赖数据库

这种分层架构使代码更易于维护和测试，同时也符合依赖倒置原则，使各层之间解耦。
