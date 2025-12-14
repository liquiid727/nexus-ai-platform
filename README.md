# AI 服务订阅平台（AI_GATEWAY） 
这是一版MVP ，后续将会更改为go开发使用


- 目标：构建覆盖认证、订阅、用量、计费、支付、API门户、通知、后台与安全体系的完整平台。
- 架构演进：阶段1 单体MVP → 阶段2 核心服务拆分 → 阶段3 全面微服务化。
- 目录概览：见下方结构与各模块README。

## 目录结构（初始骨架）
- apps/gateway
- services/iam
- services/subscription
- services/usage
- services/billing
- services/payment
- notification
- admin
- security
- proto
- deploy
- infra
- docs

## 快速开始（后续补充）
- 安装与开发环境说明将在`docs/`中提供。
- 接口契约在`proto/`，CI将进行Proto Lint。

