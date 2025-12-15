-- ----------------------------
-- 1. 组织架构表
-- ----------------------------
CREATE TABLE departments (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()), -- 部门ID，如 'dept_platform'
    name VARCHAR(100) NOT NULL,                  -- 部门名称
    parent_id VARCHAR(36),                       -- 上级部门ID，支持树形结构
    cost_center_code VARCHAR(50),                -- 财务成本中心代码
    monthly_token_budget BIGINT,                 -- 月度Token总预算
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id),
    INDEX idx_cost_center (cost_center_code)
) COMMENT = '部门信息表，用于成本归集和配额分配';

CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()), -- 用户ID，与内部SSO系统一致
    employee_id VARCHAR(50) UNIQUE NOT NULL,     -- 工号
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    department_id VARCHAR(36) NOT NULL,          -- 所属部门
    role ENUM('admin', 'user', 'billing') DEFAULT 'user',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE RESTRICT,
    INDEX idx_department (department_id),
    INDEX idx_email (email)
) COMMENT = '用户表，从内部SSO同步，关联部门';

-- ----------------------------
-- 2. 模型与路由配置表 (动态配置)
-- ----------------------------
CREATE TABLE ai_models (
    id VARCHAR(50) PRIMARY KEY,                  -- 网关内模型标识，如 'gpt-4-turbo'
    description VARCHAR(255),
    default_provider_id VARCHAR(36),             -- 默认供应商KEY
    routing_strategy ENUM('load_balance', 'priority', 'lowest_cost') DEFAULT 'load_balance',
    is_active BOOLEAN DEFAULT TRUE,
    config_json JSON,                            -- 扩展配置：温度、最大token等默认值
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) COMMENT = 'AI模型路由配置表';

CREATE TABLE api_provider_keys (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    model_id VARCHAR(50) NOT NULL,               -- 对应的网关模型标识
    provider_name VARCHAR(50) NOT NULL,          -- 供应商，如 'openai', 'azure', 'claude'
    api_key_encrypted TEXT NOT NULL,             -- 加密存储的API Key
    endpoint_base_url VARCHAR(500),              -- 供应商API端点（如有）
    priority INT DEFAULT 1,                      -- 优先级，数字越小优先级越高
    max_rpm INT,                                 -- 供应商限速：每分钟最大请求数
    max_tpm INT,                                 -- 供应商限速：每分钟最大Token数
    cost_per_million_input DECIMAL(12,6),        -- 输入Token成本（美元/百万）
    cost_per_million_output DECIMAL(12,6),       -- 输出Token成本（美元/百万）
    is_enabled BOOLEAN DEFAULT TRUE,             -- 是否启用
    health_status ENUM('healthy', 'slow', 'down') DEFAULT 'healthy',
    last_used_at TIMESTAMP NULL,
    failure_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (model_id) REFERENCES ai_models(id) ON DELETE CASCADE,
    INDEX idx_model_enabled (model_id, is_enabled, priority),
    INDEX idx_health_status (health_status)
) COMMENT = '供应商API密钥池，支持同一模型多KEY负载均衡与故障转移';

-- ----------------------------
-- 3. 配额管理表 (核心控制)
-- ----------------------------
CREATE TABLE token_quotas (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    department_id VARCHAR(36) NOT NULL,          -- 配额所属部门
    model_id VARCHAR(50),                        -- NULL表示该部门对所有模型的总配额
    quota_type ENUM('total', 'daily', 'monthly') DEFAULT 'monthly',
    max_tokens BIGINT NOT NULL,                  -- 周期内最大Token数
    used_tokens BIGINT DEFAULT 0,                -- 当前周期已用（需要定期重置或累加）
    reset_at TIMESTAMP,                          -- 下次重置时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE,
    FOREIGN KEY (model_id) REFERENCES ai_models(id) ON DELETE CASCADE,
    UNIQUE KEY uk_department_model_type (department_id, model_id, quota_type),
    INDEX idx_reset_at (reset_at)
) COMMENT = '多维度Token配额表，支持部门、模型、周期的组合限额';

-- ----------------------------
-- 4. 审计与日志表 (核心追溯)
-- ----------------------------
CREATE TABLE request_logs (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()), -- 唯一请求ID，用于全链路追踪
    request_id VARCHAR(100),                     -- 网关生成的唯一标识
    user_id VARCHAR(36) NOT NULL,                -- 请求用户
    department_id VARCHAR(36) NOT NULL,          -- 请求部门（冗余存储，避免关联查询）
    model_id VARCHAR(50) NOT NULL,               -- 请求的模型
    provider_key_id VARCHAR(36),                 -- 实际使用的供应商KEY
    provider_name VARCHAR(50),                   -- 供应商名（冗余）
    
    -- 请求详情
    prompt TEXT,                                 -- 原始请求提示词（考虑安全可脱敏存储）
    request_body JSON,                           -- 完整请求体
    request_model VARCHAR(100),                  -- 用户实际请求的模型名
    
    -- 响应详情
    response_body LONGTEXT,                      -- 完整响应体（非流式）
    response_status INT,                         -- HTTP状态码
    error_message TEXT,
    
    -- Token与成本核算
    prompt_tokens INT DEFAULT 0,
    completion_tokens INT DEFAULT 0,
    total_tokens INT DEFAULT 0,
    estimated_cost_usd DECIMAL(12,6) DEFAULT 0,
    
    -- 性能指标
    gateway_latency_ms INT,                      -- 网关内部处理耗时
    total_latency_ms INT,                        -- 总耗时（用户感受到的）
    upstream_latency_ms INT,                     -- 上游供应商处理耗时
    
    -- 时间戳
    requested_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6), -- 高精度时间戳
    responded_at TIMESTAMP(6) NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (model_id) REFERENCES ai_models(id),
    FOREIGN KEY (provider_key_id) REFERENCES api_provider_keys(id),
    INDEX idx_time_range (requested_at),         -- 按时间范围查询（最重要）
    INDEX idx_user_department (user_id, department_id),
    INDEX idx_model_provider (model_id, provider_name),
    INDEX idx_request_id (request_id)
) COMMENT = '请求审计日志表，所有AI调用的黄金数据源，数据量巨大，考虑分区';

-- ----------------------------
-- 5. 成本核算与分析表 (聚合数据)
-- ----------------------------
CREATE TABLE monthly_costs (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    department_id VARCHAR(36) NOT NULL,
    `year_month` CHAR(7) NOT NULL,                 -- 年月，格式 '2024-05'
    model_id VARCHAR(50) NOT NULL,
    total_requests INT DEFAULT 0,
    total_prompt_tokens BIGINT DEFAULT 0,
    total_completion_tokens BIGINT DEFAULT 0,
    total_tokens BIGINT DEFAULT 0,
    estimated_cost_usd DECIMAL(14,4) DEFAULT 0,  -- 月度总成本
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (model_id) REFERENCES ai_models(id),
    UNIQUE KEY uk_department_month_model (department_id, `year_month`, model_id)
) COMMENT = '月度成本聚合表，用于快速生成部门账单和报告';