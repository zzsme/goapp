# GoWK 框架标准化指南

本文档描述了 GoWK 框架的标准化规范，包括命名约定、响应格式、错误处理等方面的最佳实践。

## 响应格式标准化

### 1. 统一的响应字段

将原有的 `Errno/Errmsg` 格式规范化为更标准的命名：

```go
// 旧格式
{
    "errno": 0,
    "errmsg": "success",
    "data": null,
    "requestId": "xxx"
}

// 新格式
{
    "code": 0,
    "message": "success",
    "data": null,
    "requestId": "xxx"
}
```

理由：
- `code/message` 是更通用的命名约定
- 符合 REST API 设计最佳实践
- 提供更清晰的语义

### 2. 错误码规范

错误码采用分层设计：

- 0: 成功
- 1xxx: 客户端错误
- 5xxx: 服务器错误

示例：
```go
const (
    Success        = 0    // 成功
    BadRequest     = 1000 // 请求参数错误
    Unauthorized   = 1001 // 未授权
    InternalServer = 5000 // 服务器内部错误
)
```

### 3. HTTP 状态码映射

将内部错误码映射到合适的 HTTP 状态码：

```go
func errorToHTTPStatus(code ErrorCode) int {
    switch code {
    case BadRequest:
        return http.StatusBadRequest
    case Unauthorized:
        return http.StatusUnauthorized
    case InternalServer:
        return http.StatusInternalServerError
    default:
        return http.StatusInternalServerError
    }
}
```

## 命名规范

### 1. 包名

- 使用小写单词
- 避免下划线
- 简短且有意义

```go
// 好的例子
package config
package middleware
package handlers

// 避免的例子
package configUtils
package middleware_handlers
```

### 2. 接口名

- 通常以 "er" 结尾
- 单一职责原则

```go
// 好的例子
type Reader interface { ... }
type Writer interface { ... }
type Handler interface { ... }

// 避免的例子
type DBInterface interface { ... }
type UtilityInterface interface { ... }
```

### 3. 变量命名

- 驼峰命名法
- 简短但有描述性
- 缩写保持一致性

```go
// 好的例子
userID   // 而不是 userId 或 user_id
apiClient // 而不是 APIClient 或 api_client
db       // 数据库连接可以用简短名称

// 避免的例子
u        // 除非在非常小的作用域内
data     // 太过宽泛
```

### 4. 常量命名

- 使用驼峰命名法
- 包级常量可以使用全大写下划线

```go
// 包级常量
const (
    MaxConnections = 100
    DefaultTimeout = time.Second * 30
)

// 特定类型的常量
const (
    StatusActive   = "active"
    StatusInactive = "inactive"
)
```

## 代码组织

### 1. 文件结构

每个包应该有一个清晰的职责：

```
internal/
  ├── app/        # 应用核心组件
  ├── models/     # 数据模型定义
  ├── services/   # 业务逻辑
  └── handlers/   # 请求处理器
```

### 2. 依赖注入

推荐使用构造函数注入依赖：

```go
// 好的例子
type UserService struct {
    repo  Repository
    cache Cache
}

func NewUserService(repo Repository, cache Cache) *UserService {
    return &UserService{
        repo:  repo,
        cache: cache,
    }
}

// 避免的例子
type UserService struct {
    repo  *UserRepository // 直接依赖具体类型
    cache *RedisCache    // 直接依赖具体类型
}
```

### 3. 接口定义

- 接口应该在使用者一方定义
- 保持接口小而精确

```go
// 好的例子
type UserRepository interface {
    FindByID(id string) (*User, error)
    Save(user *User) error
}

// 避免的例子
type Repository interface {
    // 过于宽泛的接口
    Find(query string) (interface{}, error)
    Save(data interface{}) error
    Delete(id interface{}) error
}
```

## 错误处理最佳实践

### 1. 错误包装

在跨层传递错误时，添加上下文信息：

```go
// 好的例子
if err := db.Query(sql); err != nil {
    return errors.Wrap(errors.Database, "failed to query user data", err)
}

// 避免的例子
if err := db.Query(sql); err != nil {
    return err // 丢失上下文
}
```

### 2. 错误检查

使用辅助函数进行错误类型检查：

```go
// 好的例子
if errors.IsCode(err, errors.NotFound) {
    // 处理未找到的情况
}

// 避免的例子
if err.Error() == "not found" { // 脆弱的字符串比较
    // ...
}
```

## 日志规范

### 1. 日志级别使用

- DEBUG: 调试信息
- INFO: 正常操作信息
- WARN: 需要注意但不是错误的情况
- ERROR: 错误信息

### 2. 结构化日志

使用字段而不是格式化字符串：

```go
// 好的例子
logger.Info("User logged in", 
    "userID", user.ID,
    "ip", request.IP,
)

// 避免的例子
logger.Infof("User %s logged in from %s", 
    user.ID, 
    request.IP,
)
```

## 注释规范

### 1. 包注释

每个包都应该有包注释：

```go
// Package config 提供应用配置管理功能。
// 它支持从环境变量和配置文件加载配置，
// 并提供实时配置更新能力。
package config
```

### 2. 导出项注释

所有导出的类型、变量、常量、函数都应该有注释：

```go
// UserService 提供用户相关的业务逻辑实现。
type UserService struct { ... }

// ErrNotFound 表示请求的资源不存在。
var ErrNotFound = errors.New("resource not found")

// HandleRequest 处理传入的 HTTP 请求并返回响应。
func HandleRequest(req *Request) (*Response, error) { ... }
```

## 测试规范

### 1. 测试文件命名

- 测试文件与被测试文件对应
- 以 _test.go 结尾

```
user_service.go
user_service_test.go
```

### 2. 测试函数命名

使用 "Test" 前缀加被测试函数名：

```go
func TestUserService_Create(t *testing.T) { ... }
func TestUserService_Update(t *testing.T) { ... }
```

### 3. 表格驱动测试

使用表格驱动的方式组织测试用例：

```go
func TestCalculate(t *testing.T) {
    tests := []struct {
        name     string
        input    int
        expected int
        wantErr  bool
    }{
        {
            name:     "positive number",
            input:    1,
            expected: 2,
            wantErr:  false,
        },
        // 更多测试用例...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试逻辑
        })
    }
}
```

## 版本控制

### 1. 语义化版本

遵循语义化版本规范 (SemVer)：

- MAJOR.MINOR.PATCH
- MAJOR: 不兼容的 API 变更
- MINOR: 向后兼容的功能性新增
- PATCH: 向后兼容的问题修正

### 2. Git 提交信息

使用规范的提交信息格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型包括：
- feat: 新功能
- fix: 修复
- docs: 文档更改
- style: 格式调整
- refactor: 重构
- test: 测试相关
- chore: 构建过程或辅助工具的变动

通过遵循这些标准化规范，我们可以：

1. 提高代码可维护性
2. 降低团队协作成本
3. 提供一致的用户体验
4. 简化问题诊断和修复流程
5. 加快新成员融入团队的速度

这些规范不是一成不变的，应该根据项目实际情况和团队反馈持续改进和调整。
