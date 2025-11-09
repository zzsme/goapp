# GoWK 架构设计文档

## 整体架构

GoWK采用分层架构设计，主要分为以下几层：

```mermaid
graph TB
    subgraph HTTP层
        A[Router] --> B[Middleware]
        B --> C[Controller]
    end
    
    subgraph 业务层
        C --> D[Service]
        D --> E[Repository]
    end
    
    subgraph 数据层
        E --> F[Database]
        E --> G[Redis]
    end
    
    subgraph 事件系统
        H[EventBus] --> I[EventHandlers]
        D --> H
    end
    
    subgraph 监控系统
        J[Metrics Collector] --> K[Monitor Service]
        B --> J
    end
    
    subgraph 核心组件
        L[Logger]
        M[Config]
        N[Validator]
        O[ErrorHandler]
    end
```

## 核心组件

### 1. HTTP处理流程

```mermaid
sequenceDiagram
    participant Client
    participant Router
    participant Middleware
    participant Controller
    participant Service
    participant Repository
    participant Database

    Client->>Router: HTTP Request
    Router->>Middleware: 路由匹配
    Middleware->>Middleware: 认证/日志/监控
    Middleware->>Controller: 处理请求
    Controller->>Service: 调用业务逻辑
    Service->>Repository: 数据操作
    Repository->>Database: SQL查询
    Database-->>Repository: 返回数据
    Repository-->>Service: 封装数据
    Service-->>Controller: 业务处理
    Controller-->>Client: HTTP Response
```

### 2. 事件系统

```mermaid
graph LR
    A[Publisher] -->|发布事件| B[EventBus]
    B -->|分发| C[Handler1]
    B -->|分发| D[Handler2]
    B -->|分发| E[Handler3]
```

### 3. 错误处理

```mermaid
graph TD
    A[错误发生] -->|创建| B[AppError]
    B -->|包含| C[ErrorCode]
    B -->|包含| D[Message]
    B -->|转换| E[HTTP Response]
    
    subgraph 错误类型
        F[Client Errors]
        G[Server Errors]
    end
```

## 数据流

### 1. 请求处理流程

```mermaid
sequenceDiagram
    participant C as Client
    participant M as Middleware
    participant H as Handler
    participant S as Service
    participant D as Database

    C->>M: HTTP Request
    M->>M: 请求预处理
    M->>H: 处理请求
    H->>S: 业务逻辑
    S->>D: 数据操作
    D-->>S: 返回结果
    S-->>H: 处理结果
    H-->>M: 格式化响应
    M-->>C: HTTP Response
```

### 2. 事件处理流程

```mermaid
sequenceDiagram
    participant S as Service
    participant E as EventBus
    participant H as Handlers
    participant L as Logger

    S->>E: 发布事件
    E->>H: 分发事件
    H->>L: 记录日志
    H->>S: 处理回调
```

## 组件交互

### 1. 控制器与服务层

```mermaid
classDiagram
    class Controller {
        +Service service
        +HandleRequest()
        +ValidateInput()
        +FormatResponse()
    }
    class Service {
        +Repository repo
        +ProcessBusiness()
        +ValidateLogic()
    }
    class Repository {
        +Database db
        +CRUD()
    }
    Controller --> Service
    Service --> Repository
```

### 2. 中间件链

```mermaid
graph LR
    A[Request] -->|1| B[Auth]
    B -->|2| C[Logger]
    C -->|3| D[Monitor]
    D -->|4| E[Handler]
    E -->|5| F[Response]
```

## 配置管理

```mermaid
graph TD
    A[配置文件] -->|加载| B[Config组件]
    B -->|注入| C[Services]
    B -->|注入| D[Database]
    B -->|注入| E[Redis]
    B -->|注入| F[Logger]
```

## 监控系统

```mermaid
graph TD
    A[请求] -->|记录| B[Metrics收集器]
    B -->|统计| C[监控服务]
    C -->|展示| D[监控接口]
    C -->|告警| E[告警系统]
```

## 安全架构

```mermaid
graph TD
    A[请求] -->|1| B[TLS终止]
    B -->|2| C[认证中间件]
    C -->|3| D[授权检查]
    D -->|4| E[输入验证]
    E -->|5| F[业务处理]
    F -->|6| G[响应加密]
