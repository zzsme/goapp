# GoWK 错误处理系统

GoWK 框架提供了一套完整的错误处理机制，旨在统一错误管理、简化错误处理流程、提高代码可维护性并提供清晰的错误反馈。

## 错误码系统

GoWK 使用一个集中化的错误码系统，将错误分为客户端错误（1xxx）和服务器错误（5xxx）：

### 客户端错误（1000-1999）

| 错误码 | 常量名称 | 描述 |
|--------|----------|------|
| 1000 | BadRequest | 请求参数错误或格式不正确 |
| 1001 | Unauthorized | 未提供身份验证或身份验证失败 |
| 1003 | Forbidden | 权限不足，禁止访问 |
| 1004 | NotFound | 请求的资源不存在 |
| 1005 | Validation | 数据验证失败 |
| 1009 | Conflict | 资源冲突（如重复创建） |
| 1029 | TooManyReq | 请求过于频繁，被限流 |
| 1030 | InvalidToken | 提供的令牌无效 |
| 1031 | ExpiredToken | 提供的令牌已过期 |
| 1032 | InvalidFormat | 数据格式错误 |

### 服务器错误（5000-5999）

| 错误码 | 常量名称 | 描述 |
|--------|----------|------|
| 5000 | InternalServer | 服务器内部错误 |
| 5001 | Database | 数据库操作失败 |
| 5002 | Cache | 缓存操作失败 |
| 5003 | NotImplemented | 请求的功能未实现 |
| 5004 | ThirdParty | 第三方服务调用失败 |
| 5005 | Config | 配置错误 |

## 错误类型

GoWK 框架定义了 `AppError` 类型，用于表示应用程序错误：

```go
// AppError 表示应用程序错误
type AppError struct {
    Code    ErrorCode              // 错误码
    Message string                 // 错误消息
    Cause   error                  // 原始错误（可选）
    Details map[string]interface{} // 详细信息（可选）
}
```

## 创建错误

GoWK 提供了多种创建错误的方法：

### 创建新错误

```go
// 创建简单错误
err := errors.New(errors.BadRequest, "Invalid user ID format")

// 包装现有错误
if result, err := db.Exec(query, args...); err != nil {
    return errors.Wrap(errors.Database, "Failed to insert user", err)
}

// 添加错误详情
err := errors.New(errors.Validation, "Validation failed").
    With("field", "email").
    With("reason", "invalid format")

// 或批量添加详情
err := errors.New(errors.Validation, "Validation failed").
    WithDetails(map[string]interface{}{
        "field": "email",
        "reason": "invalid format",
    })
```

## 在 API 控制器中使用错误处理

GoWK 的 `APIContext` 提供了多种方法来发送错误响应：

### 基本错误响应

```go
// 使用错误码常量
func (c *UserController) GetUser(ctx *context.APIContext) {
    id := ctx.Param("id")
    if id == "" {
        ctx.ErrorWithCode(errors.BadRequest, "User ID is required")
        return
    }
    
    // ...处理请求...
}
```

### 使用 AppError

```go
func (c *ProductController) Create(ctx *context.APIContext) {
    var dto productDTO
    if err := ctx.BindJSON(&dto); err != nil {
        ctx.ErrorWithAppError(errors.Wrap(errors.BadRequest, "Invalid request body", err))
        return
    }
    
    product, err := c.service.CreateProduct(dto)
    if err != nil {
        // 直接传递服务层返回的 AppError
        if appErr, ok := err.(*errors.AppError); ok {
            ctx.ErrorWithAppError(appErr)
        } else {
            ctx.ErrorWithCode(errors.InternalServer, "Failed to create product")
        }
        return
    }
    
    ctx.Success(product)
}
```

## 错误检查

GoWK 提供了方便的错误检查方法：

```go
// 检查错误码
if errors.IsCode(err, errors.NotFound) {
    // 处理资源不存在的情况
}

// 获取原始错误（Unwrap）
if origErr := errors.Unwrap(err); origErr != nil {
    // 处理原始错误
}
```

## 标准错误响应格式

所有 API 错误响应都会采用统一格式：

```json
{
  "code": 1000,
  "message": "Invalid parameters",
  "data": {
    "field": "email",
    "reason": "invalid format"
  },
  "requestId": "bce68d7d-0ace-49dd-b3ed-2977xxxx"
}
```

其中：

- `code`：错误码，0 表示成功，其他值表示错误
- `message`：人类可读的错误消息
- `data`：可选的详细错误信息或额外数据
- `requestId`：请求唯一标识，用于跟踪和调试

## 最佳实践

1. **使用预定义错误码**：始终使用 `errors` 包中定义的错误码常量，而不是硬编码数字。
2. **提供有意义的错误消息**：错误消息应清晰描述问题，但不要泄露敏感信息。
3. **在适当的层次处理错误**：通常在服务层创建错误，在控制器层决定如何响应错误。
4. **包装底层错误**：总是使用 `Wrap` 方法包装第三方库或底层错误，以保留错误链。
5. **添加上下文信息**：使用 `With` 或 `WithDetails` 方法添加有助于调试的上下文信息。
6. **日志记录**：在记录错误日志时，确保包含原始错误信息和堆栈跟踪。

## 错误处理流程图

```
请求 → 中间件 → 控制器 → 服务层 → 存储库
                 ↑          ↑         ↑
                 |          |         |
              验证错误    业务错误   数据错误
                 |          |         |
                 └──────────┴─────────┘
                          ↓
                    创建 AppError
                          ↓
                  使用 APIContext 
                   发送错误响应
                          ↓
                 中间件记录错误日志
```

通过这种统一的错误处理系统，GoWK 框架确保了错误管理的一致性，提高了代码的可维护性，并为客户端提供了清晰的错误反馈。
