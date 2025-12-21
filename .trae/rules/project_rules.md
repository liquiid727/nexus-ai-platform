# 项目结构
ai-gateway/
├── cmd/
│   └── server/
│       └── main.go          # 应用入口
├── configs/                 # 配置文件
│   ├── config.yaml          # 主配置文件
│   └── config.develop.yaml  # 开发环境覆盖配置
├── internal/                # 私有应用代码
│   ├── config/              # 配置加载
│   │   └── config.go
│   ├── handler/             # HTTP处理器 (Controller层)
│   │   ├── chat.go
│   │   ├── auth.go
│   │   └── quota.go
│   ├── grpc/                # gRPC处理器
│   │   ├── proxy.go
│   │   └── quota.go
│   ├── middleware/          # HTTP中间件
│   │   ├── auth.go
│   │   ├── rate_limiter.go
│   │   └── logger.go
│   ├── router/              # 路由注册
│   │   └── router.go
│   ├── service/             # 业务逻辑 (Service层)
│   │   ├── proxy.go
│   │   ├── quota.go
│   │   └── cost.go
│   ├── repository/          # 数据仓库 (Repository层，解耦DB)
│   │   └── entity/          # 数据库实体模型
│   ├── dao/                 # 数据访问对象 (具体DB实现)
│   └── pkg/                 # 内部共享库 (logger, error等)
│       ├── logger/
│       └── util/
├── pkg/                     # 公共库（可被外部项目引用）
├── scripts/                 # 构建与部署脚本
├── deployments/             # 部署配置 (Docker, K8s)
├── docs/                    # 文档 (API Swagger, SQL等)
├── .env.example             # 环境变量示例
├── .gitignore
├── go.mod                   # Go模块定义
├── Makefile                 # 构建管理命令
└── README.md

# *** 项目开发规范

## 1. 技术栈与选型 (Tech Stack)
| 类别 | 技术/库 | 说明 |
| :--- | :--- | :--- |
| **语言** | Go 1.22+ | 强制使用最新稳定版特性 |
| **Web框架** | Echo v4 | 高性能 Web 框架 |
| **ORM** | GORM v2 | 必须配合 Entity 定义使用 |
| **配置管理** | Viper | 支持 YAML/Env 自动加载 |
| **日志** | Zap | 必须使用 `internal/pkg/logger` 封装 |
| **数据库** | MySQL 8.0+ | 核心业务数据存储 |
| **缓存** | Redis 6.0+ | 缓存与分布式锁 |
| **消息队列** | Kafka 2.8+ | 异步解耦与削峰填谷 |
| **认证授权** | JWT + RBAC | 标准化身份验证 |
| **监控** | Prometheus + Grafana | 业务指标与系统监控 |
| **链路追踪** | OpenTelemetry | 分布式链路追踪 (Jaeger/Zipkin) |

## 2. 架构设计原则 (Architecture)
遵循 **分层架构** 与 **依赖倒置** 原则：
- **Handler 层**：仅处理 HTTP 请求/响应、参数校验、调用 Service。**禁止包含业务逻辑**。
- **Service 层**：包含核心业务逻辑。依赖 Repository 接口，而非具体实现。
- **Repository/DAO 层**：负责数据的持久化操作。**禁止包含业务逻辑**。
- **Entity 层**：定义数据库模型，不应包含 JSON tag（DTO 与 Entity 分离）。

**禁止行为**：
- 禁止在 Handler 层直接访问数据库。
- 禁止在 Service 层直接引入 `gorm.DB`，必须通过 Repository/DAO 接口操作。
- 禁止在循环中进行 SQL 查询（N+1 问题），应使用批量查询。

## 3. 编码规范 (Coding Style)
- **格式化**：保存时必须自动运行 `goimports` 和 `gofumpt`。
- **命名规范**：
  - 文件名：使用 `snake_case` (如 `user_service.go`)。
  - 接口名：包含方法名的 `er` 结尾 (如 `Reader`, `Worker`) 或 `I` 前缀 (如 `IUserService`)。
  - 变量名：使用 `camelCase`，首字母缩写词保持一致 (如 `userID`, `xmlHTTPRequest`)。
- **函数长度**：原则上单函数不超过 **50行**，复杂逻辑需拆分。
- **注释**：所有 Exported (首字母大写) 的函数、结构体、接口必须包含注释。

## 4. 错误处理与日志 (Error & Logging)
- **错误处理**：
  - 必须处理所有返回的 `error`，禁止使用 `_` 忽略。
  - 使用 `internal/pkg/error` 统一封装业务错误码。
  - Service 层返回的 error 应当包含堆栈信息 (使用 `pkg/errors` 或 Go 1.13+ wrapping)。
- **日志记录**：
  - 使用结构化日志 (Structured Logging)。
  - **Error 级别**：必须记录，且包含 Stack Trace。
  - **Info 级别**：关键业务节点（如订单创建、支付成功）。
  - **Debug 级别**：开发调试信息，生产环境通常关闭。
  - 禁止使用 `fmt.Println`，必须使用 `logger`。

## 5. 数据库规范 (Database)
- **迁移管理**：所有 Schema 变更必须通过 SQL 迁移脚本或 GORM Migration 工具进行，禁止手动修改线上库。
- **字段规范**：
  - 表名：复数形式 `snake_case` (如 `users`, `order_items`)。
  - 必须包含 `created_at`, `updated_at` 时间戳。
  - 逻辑删除使用 `deleted_at` (GORM Soft Delete)。
- **索引**：外键字段、查询频繁字段必须建立索引。

## 6. 测试规范 (Testing)
- **单元测试**：
  - 核心业务逻辑 (Service 层) 覆盖率需达到 **80%+**。
  - 使用 `Table-Driven Tests` (表格驱动测试) 模式。
  - 使用 `gomock` 或 `testify/mock` 隔离外部依赖 (DB, Redis, API)。
- **集成测试**：关键 API 接口需编写集成测试。

## 7. 版本控制与协作 (Git Flow)
- **分支模型**：
  - `main`: 生产环境稳定分支。
  - `develop`: 开发主分支。
  - `feature/xxx`: 功能开发分支。
  - `fix/xxx`: Bug 修复分支。
- **提交信息 (Conventional Commits)**：
  - `feat`: 新功能
  - `fix`: 修复 Bug
  - `docs`: 文档变更
  - `style`: 代码格式调整 (不影响逻辑)
  - `refactor`: 代码重构
  - `test`: 测试相关
  - `chore`: 构建过程或辅助工具变更
  - **示例**: `feat(auth): add jwt token refresh mechanism`

## 8. API 设计规范
- 遵循 RESTful 风格。
- 响应结构统一：
  ```json
  {
    "code": 0,          // 业务状态码，0表示成功
    "message": "success", // 提示信息
    "data": { ... },    // 业务数据
    "request_id": "..." // 追踪ID
  }
  ```
