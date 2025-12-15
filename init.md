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
│   ├── handler/             # HTTP处理器
│   │   ├── chat.go
│   │   ├── auth.go
│   │   └── quota.go
│   ├── middleware/          # 中间件
│   │   ├── auth.go
│   │   ├── rate_limiter.go
│   │   └── logger.go
│   ├── router/              # 路由注册
│   │   └── router.go
│   ├── service/             # 业务逻辑
│   │   ├── proxy.go
│   │   ├── quota.go
│   │   └── cost.go
│   ├── model/               # 数据模型（对应数据库）
│   │   └── entity/
│   ├── dao/                 # 数据访问
│   └── pkg/                 # 可公开的内部库
│       ├── logger/
│       └── util/
├── pkg/                     # 公共库（可被外部引用）
├── scripts/                 # 部署脚本
├── deployments/             # 部署配置
├── .env.example             # 环境变量示例
├── .gitignore
├── go.mod                   # Go模块
├── Makefile                 # 构建命令
└── README.md