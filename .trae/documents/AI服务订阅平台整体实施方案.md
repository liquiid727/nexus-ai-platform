## 总览与仓库现状
- 当前目录`d:\code\AI_GATEWAY`为空，视为全新工程；方案将从零规划并适配后续代码骨架搭建。
- 满足功能清单与技术规范，分三阶段演进：单体MVP→核心服务拆分→全面微服务化。

## 技术栈选择
- 后端：
  - MVP：`Node.js + NestJS`（快速交付、丰富生态，内置OpenAPI、Swagger、JWT、RBAC实践成熟）。
  - 高性能路径（阶段2/3）：`Go + gRPC`（核心用量计量、计费引擎、API网关扩展）。
- 通信：`gRPC + Protobuf`（服务间契约）；对外 `REST + OpenAPI`（开发者门户与SDK生成）。
- 数据：
  - 事务库：`PostgreSQL`
  - 时序指标：`TimescaleDB`（Postgres扩展）
  - 日志搜索：`Elasticsearch`
  - 缓存：`Redis`
  - 事件总线：`Kafka`
- 网关与服务网格：`Envoy/Kong`（对外网关）+ `Istio`（服务网格）
- 配置中心与密钥：`Consul` + `Vault`（或云KMS）
- 监控与追踪：`Prometheus + Grafana`、`Jaeger`、`Sentry`
- 前端与后台：`React + Ant Design`（管理后台与仪表板）

## 领域划分与微服务边界
- `iam`（用户/认证/授权/RBAC/2FA/会话）
- `catalog`（AI服务目录与规格）
- `subscription`（套餐、生命周期、续费、升降级）
- `usage`（用量计量、实时聚合、阈值预警）
- `billing`（计费引擎、周期账单、明细、发票）
- `payment`（多支付网关、退款、对账）
- `api-portal`（密钥管理、访问策略、速率限制、SDK/Swagger）
- `notification`（模板化消息、邮件/短信/WebSocket/推送）
- `admin`（审计、配置中心、财务看板、运维监控）
- `security`（WAF、防护规则、审计日志、合规）

## 核心能力设计要点
- 用户与认证（iam）：
  - 邮箱/密码注册+邮件验证、手机号+短信验证码、第三方OAuth2（微信/GitHub/Google）
  - 2FA（TOTP + 备用恢复码）、设备与会话管理、登录历史
  - RBAC：角色（user、org_admin、sys_admin）+细粒度权限；策略缓存与审计
- 订阅与套餐（subscription/catalog）：
  - 多级套餐（免费/基础/专业/企业定制）+弹性计费（次/量/时长）
  - 生命周期：开通、升级/降级、自动续费、宽限期（Grace）策略
- 用量与消费（usage）：
  - 维度：API调用、计算时长、存储用量；实时聚合与迟到事件处理
  - 阈值预警：80/90/100%多渠道通知；项目/团队维度成本分析
- 账单与支付（billing/payment）：
  - 周期账单（自然月/自定义）与明细关联；税率与地区规则
  - 网关：支付宝/微信支付/信用卡（Stripe）/PayPal；发票自动开具；退款与对账
- API管理门户（api-portal）：
  - JWT/OAuth2密钥发放；IP白名单、速率限制、配额
  - Swagger UI、交互式请求测试台；SDK生成（TS/Go/Python/Java）
- 通知中心（notification）：
  - 模板化事务邮件、短信、站内信；WebSocket实时；移动推送（FCM/APNs）
  - 收件箱：分类、已读/未读、查询与归档
- 管理后台与运维（admin）：
  - 审计日志、服务配置中心、财务看板
  - 分布式追踪（Jaeger）、异常报警（Sentry）、自动扩缩容策略
- 安全体系（security）：
  - AES-256敏感数据加密、密钥轮换；PCI DSS路径；WAF与漏洞扫描

## 数据与存储建模
- 事务表：`users, user_profiles, sessions, devices, roles, permissions, orgs, plans, subscriptions, invoices, payments, refunds, transactions, api_keys, policies, notifications, audit_logs`
- 时序：`usage_metrics(metric, value, ts, dims{service,team,project})`；`cost_metrics`
- 索引与检索：日志与审计写入`Elasticsearch`（支持异常交易检测）
- 多租户：`tenant_id`列级隔离+策略；阶段3支持schema/实例级隔离

## 通信与接口契约（Protobuf）
- 包结构：`iam.v1`, `catalog.v1`, `subscription.v1`, `usage.v1`, `billing.v1`, `payment.v1`, `api.v1`, `notification.v1`, `admin.v1`
- 示例（缩略）：
  - `iam.v1.AuthService`: `RegisterEmail`, `VerifyEmail`, `RegisterPhone`, `VerifySMS`, `OAuthLogin`, `IssueToken`, `Enable2FA`, `ListSessions`
  - `subscription.v1.SubscriptionService`: `Create`, `Upgrade`, `Downgrade`, `SetAutoRenew`, `Cancel`, `GetState`
  - `usage.v1.UsageService`: `IngestEvent`, `GetRealtimeUsage`, `SetThreshold`, `ListAlerts`
  - `billing.v1.BillingService`: `GenerateInvoice`, `ListInvoices`, `AttachUsageDetail`
  - `payment.v1.PaymentService`: `CreatePayment`, `Refund`, `WebhookNotify`, `Reconcile`
- 契约治理：使用`Buf`进行Protobuflint与版本管理；向外暴露`REST`由网关与`OpenAPI`桥接

## 部署与基础设施
- 容器化：每服务`Dockerfile`；`docker-compose`用于MVP本地联调
- K8s：阶段2起部署至`Kubernetes`，使用`Helm`与`Kustomize`；多区域拓扑
- 服务网格：`Istio`流量治理（熔断、重试、金丝雀/蓝绿）
- 配置中心：`Consul`；密钥管理`Vault/KMS`；统一`env`注入
- 流水线：`GitHub Actions/GitLab CI`（单测/安全扫描/镜像构建/部署）
- 混沌工程：`Chaos Mesh/Litmus`按阶段引入

## 监控与告警
- 指标：应用/业务/基础设施三层；`Prometheus`采集、`Grafana`看板
- 追踪：`Jaeger`；错误：`Sentry`
- 业务指标：注册转化率、订阅留存率、ARPU
- 技术指标：依赖拓扑、数据库查询性能、消息队列积压
- 预警：80/90/100%用量阈值通知（邮件/短信/站内信）

## 质量保障与性能
- 目标：`P99 < 200ms`，`10,000+ TPS`（核心链路Go/gRPC，缓存与批量聚合）
- 性能测试：`k6/Gatling`压测；契约与兼容性测试（`grpcurl`、`Dredd`）
- 可靠性：`99.95% SLA`、跨AZ容灾（读写分离、RPO/RTO目标）
- 安全：年度渗透测试、SOC2流程与证据库、依赖扫描（SCA）、容器镜像扫描

## 演进路线与里程碑
- 阶段1（3个月）：单体MVP
  - NestJS整合：认证/RBAC、服务目录、订阅、用量采集、账单生成（简化）、支付沙盒（Stripe/PayPal/微信/支付宝）、通知与后台最小版
  - 交付：设计文档、Proto契约、基础看板、性能测试报告（MVP规模）、故障恢复预案
- 阶段2（6个月）：核心服务拆分
  - 拆分`usage/billing/payment`为Go微服务；引入Kafka、TimescaleDB；上Istio与Prometheus；蓝绿部署
  - 交付：扩展契约与SDK、性能报告（万级TPS）、对账与异常检测完善
- 阶段3（12个月）：全面微服务化
  - 全域微服务、跨区域部署、多租户隔离增强、混沌工程与合规体系完善
  - 交付：完整SLA、PCI/SOC2路径、全量看板与审计

## 交付物清单（每模块）
- 详细设计文档（目标、时序、数据模型、异常与恢复）
- 接口契约（Protobuf定义与版本治理，OpenAPI映射）
- 性能测试报告（场景、数据、指标、瓶颈与优化建议）
- 故障恢复预案（回滚、降级、数据修复、RPO/RTO）
- 监控仪表板配置（Grafana面板JSON、告警规则）

## 目录结构建议（MVP→微服务）
- `apps/gateway`（NestJS对外API与开发者门户）
- `services/iam`、`services/subscription`（TS）
- `services/usage`、`services/billing`、`services/payment`（Go）
- `proto/`（Buf管理契约）
- `deploy/`（Docker/K8s/Helm）
- `infra/`（Istio、Consul、Vault、Grafana、Prometheus、Jaeger、Sentry）
- `docs/`（设计、测试、预案、看板）

## 风险与合规要点
- 支付合规（PCI DSS）与数据保护（AES-256、字段脱敏、密钥轮换）
- WAF与速率限制防护：突发流量与机器人；异常交易检测（规则+模型）
- 多区域一致性与成本控制：数据复制、延迟与费用权衡

## 下一步
- 经确认后：初始化仓库骨架（MVP）、创建Proto契约与基础CI、搭建开发环境与最小可用看板；随后迭代各模块功能。