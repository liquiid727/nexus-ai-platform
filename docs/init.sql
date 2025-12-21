-- -----------------------------------------------------------------------------
-- Next AI Gateway Database Initialization Script
-- -----------------------------------------------------------------------------

-- =============================================================================
-- 1. 组织架构与用户核心 (Organization & Users)
-- =============================================================================

-- 1.1 部门表 (departments)
CREATE TABLE departments (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),                -- 部门ID，如 'dept_platform'
    name VARCHAR(100) NOT NULL,                                 -- 部门名称
    parent_id VARCHAR(36),                                      -- 上级部门ID，支持树形结构
    cost_center_code VARCHAR(50),                               -- 财务成本中心代码
    monthly_token_budget BIGINT,                                -- 月度Token总预算
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id),
    INDEX idx_cost_center (cost_center_code)
) COMMENT = '部门信息表，用于成本归集和配额分配';

-- 1.2 用户表 (users)
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),                -- 业务ID，如usr_xxxx
    internal_id BIGINT UNSIGNED AUTO_INCREMENT UNIQUE,          -- 内部ID（不暴露）
    
    -- 唯一标识（支持多登录方式）
    username VARCHAR(50) UNIQUE,                                -- 用户名，可空（允许其他方式注册）
    email VARCHAR(255) UNIQUE,                                  -- 邮箱，可空
    mobile VARCHAR(20) UNIQUE,                                  -- 手机号，可空

    -- 基础状态
    status TINYINT DEFAULT 1 COMMENT '1:正常 2:禁用 3:注销',
    type TINYINT DEFAULT 1 COMMENT '1:普通用户 2:VIP 3:管理员',
   
    employee_id VARCHAR(50) UNIQUE NOT NULL,                    -- 工号
    name VARCHAR(100) NOT NULL,
    -- email VARCHAR(255) NOT NULL,                             -- [Fixed] Duplicate column removed
    department_id VARCHAR(36) NOT NULL,                         -- 所属部门
    role ENUM('admin', 'user', 'billing') DEFAULT 'user',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL COMMENT '软删除',
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE RESTRICT,
    
    -- 索引
    INDEX idx_department (department_id),
    INDEX idx_email (email),
    INDEX idx_mobile (mobile),
    INDEX idx_status_created (status, created_at),
    INDEX idx_updated (updated_at)
) COMMENT = '用户表，从内部SSO同步，关联部门';

-- =============================================================================
-- 2. 用户扩展信息 (User Extensions)
-- =============================================================================

-- 2.1 用户密码表 (user_passwords)
CREATE TABLE user_passwords (
    user_id VARCHAR(32) PRIMARY KEY,
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    password_salt VARCHAR(64) NOT NULL COMMENT '密码盐值',
    algorithm VARCHAR(20) DEFAULT 'bcrypt' COMMENT '哈希算法',
    
    -- 安全信息
    last_changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    change_required BOOLEAN DEFAULT FALSE COMMENT '需要修改密码',
    failed_attempts INT DEFAULT 0,
    locked_until TIMESTAMP NULL,
    
    -- 历史密码（可选）
    previous_hash VARCHAR(255),
    
    -- 外键
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_locked (locked_until)
) ENGINE=InnoDB COMMENT='用户密码表 - 安全隔离';

-- 2.2 用户资料表 (user_profiles)
CREATE TABLE user_profiles (
    user_id VARCHAR(32) PRIMARY KEY,
    
    -- 基础资料
    nickname VARCHAR(100) COMMENT '昵称',
    avatar_url VARCHAR(500) COMMENT '头像URL',
    gender TINYINT DEFAULT 0 COMMENT '0:未知 1:男 2:女',
    birthday DATE COMMENT '生日',
    bio VARCHAR(500) COMMENT '个人简介',
    
    -- 地理位置
    country VARCHAR(50),
    province VARCHAR(50),
    city VARCHAR(50),
    timezone VARCHAR(50) DEFAULT 'Asia/Shanghai',
    
    -- 社交信息
    wechat_unionid VARCHAR(100) UNIQUE,
    wechat_openid VARCHAR(100),
    github_id VARCHAR(100),
    
    -- 时间戳
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- 外键
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- 索引
    INDEX idx_nickname (nickname),
    INDEX idx_wechat (wechat_unionid, wechat_openid)
) ENGINE=InnoDB COMMENT='用户资料表 - 用户画像';

-- 2.3 用户设置表 (user_settings)
CREATE TABLE user_settings (
    user_id VARCHAR(32) PRIMARY KEY,
    
    -- 通知设置
    notify_email JSON COMMENT '邮件通知设置',
    notify_push JSON COMMENT '推送通知设置',
    notify_sms JSON COMMENT '短信通知设置',
    
    -- 隐私设置
    privacy_profile TINYINT DEFAULT 1 COMMENT '资料可见性 1:公开 2:好友 3:私密',
    privacy_contact TINYINT DEFAULT 3 COMMENT '联系方式可见性',
    
    -- 界面设置
    ui_theme VARCHAR(20) DEFAULT 'light',
    ui_language VARCHAR(10) DEFAULT 'zh-CN',
    
    -- 业务相关设置
    custom_settings JSON COMMENT '自定义业务设置',
    
    -- 时间戳
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- 外键
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='用户设置表 - 低频更新';

-- =============================================================================
-- 3. 用户安全与状态 (User Security & State)
-- =============================================================================

-- 3.1 用户安全表 (user_security)
CREATE TABLE user_security (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(32) NOT NULL,
    
    -- 验证状态
    email_verified BOOLEAN DEFAULT FALSE,
    mobile_verified BOOLEAN DEFAULT FALSE,
    id_verified BOOLEAN DEFAULT FALSE,
    
    -- 验证时间
    email_verified_at TIMESTAMP NULL,
    mobile_verified_at TIMESTAMP NULL,
    
    -- 二次验证
    mfa_enabled BOOLEAN DEFAULT FALSE,
    mfa_type VARCHAR(20) COMMENT 'totp,sms,email',
    mfa_secret VARCHAR(100),
    
    -- 登录安全
    last_login_at TIMESTAMP NULL,
    last_login_ip VARCHAR(45),
    last_password_changed_at TIMESTAMP NULL,
    
    -- 索引
    INDEX idx_user (user_id),
    INDEX idx_verified (email_verified, mobile_verified),
    
    -- 外键
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='用户安全表 - 安全审计';

-- 3.2 用户统计表 (user_statistics)
CREATE TABLE user_statistics (
    user_id VARCHAR(32) PRIMARY KEY,
    
    -- 登录统计
    login_count INT DEFAULT 0,
    continuous_login_days INT DEFAULT 0,
    
    -- 活跃度统计
    last_active_at TIMESTAMP NULL,
    total_online_seconds BIGINT DEFAULT 0,
    
    -- 内容统计
    post_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    like_count INT DEFAULT 0,
    follower_count INT DEFAULT 0,
    following_count INT DEFAULT 0,
    
    -- 业务统计
    order_count INT DEFAULT 0,
    total_spent DECIMAL(12, 2) DEFAULT 0.00,
    
    -- 更新时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- 外键
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- 索引
    INDEX idx_follower (follower_count),
    INDEX idx_total_spent (total_spent)
) ENGINE=InnoDB COMMENT='用户统计表 - 高频更新';

-- 3.3 用户会话表 (user_sessions)  放置在redis当中进行管理
CREATE TABLE user_sessions (
    id VARCHAR(64) PRIMARY KEY COMMENT 'Session ID',
    user_id VARCHAR(32) NOT NULL,
    
    -- 设备信息
    device_id VARCHAR(100) COMMENT '设备唯一标识',
    device_type VARCHAR(20) COMMENT 'web,ios,android',
    user_agent VARCHAR(500),
    
    -- 令牌信息
    refresh_token_hash VARCHAR(255) UNIQUE,
    access_token_expires_at TIMESTAMP,
    refresh_token_expires_at TIMESTAMP,
    
    -- 会话状态
    status TINYINT DEFAULT 1 COMMENT '1:活跃 2:已退出 3:过期',
    
    -- 位置信息
    ip_address VARCHAR(45),
    location VARCHAR(100),
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_activity_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- 索引
    INDEX idx_user_status (user_id, status),
    INDEX idx_refresh_token (refresh_token_hash),
    INDEX idx_expires (refresh_token_expires_at),
    
    -- 外键
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='用户会话表';

-- 3.4 第三方登录绑定表 (user_oauth_bindings)
CREATE TABLE user_oauth_bindings (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    
    -- 第三方平台信息
    provider VARCHAR(20) NOT NULL,                              -- wechat, qq, weibo, github, google, apple
    provider_user_id VARCHAR(255) NOT NULL,                     -- 第三方用户ID
    union_id VARCHAR(255),                                      -- 微信union_id
    
    -- 第三方用户信息（缓存）
    open_id VARCHAR(255),
    nickname VARCHAR(100),
    avatar_url TEXT,
    raw_user_info JSONB,                                        -- 原始用户信息
    
    -- 认证信息
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP,
    scope VARCHAR(255),
    
    -- 绑定状态
    is_primary BOOLEAN DEFAULT FALSE,                           -- 是否主绑定
    is_valid BOOLEAN DEFAULT TRUE,                              -- 是否有效
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP,                                     -- 最后使用时间
    
    -- 索引
    INDEX idx_user_provider (user_id, provider),
    INDEX idx_provider_user (provider, provider_user_id),
    INDEX idx_union_id (union_id),
    UNIQUE KEY uk_provider_user (provider, provider_user_id),
    
    -- 外键
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- 约束
    CONSTRAINT chk_provider CHECK (provider IN ('wechat', 'qq', 'weibo', 'github', 'google', 'apple', 'facebook', 'twitter'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='第三方登录绑定表';

-- =============================================================================
-- 4. 模型网关配置 (Model Gateway Config)
-- =============================================================================

-- 4.1 AI模型路由配置表 (ai_models)
CREATE TABLE ai_models (
    id VARCHAR(50) PRIMARY KEY,                                 -- 网关内模型标识，如 'gpt-4-turbo'
    description VARCHAR(255),
    default_provider_id VARCHAR(36),                            -- 默认供应商KEY
    routing_strategy ENUM('load_balance', 'priority', 'lowest_cost') DEFAULT 'load_balance',
    is_active BOOLEAN DEFAULT TRUE,
    config_json JSON,                                           -- 扩展配置：温度、最大token等默认值
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) COMMENT = 'AI模型路由配置表';

-- 4.2 供应商API密钥池 (api_provider_keys)
CREATE TABLE api_provider_keys (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    model_id VARCHAR(50) NOT NULL,                              -- 对应的网关模型标识
    provider_name VARCHAR(50) NOT NULL,                         -- 供应商，如 'openai', 'azure', 'claude'
    api_key_encrypted TEXT NOT NULL,                            -- 加密存储的API Key
    endpoint_base_url VARCHAR(500),                             -- 供应商API端点（如有）
    priority INT DEFAULT 1,                                     -- 优先级，数字越小优先级越高
    max_rpm INT,                                                -- 供应商限速：每分钟最大请求数
    max_tpm INT,                                                -- 供应商限速：每分钟最大Token数
    cost_per_million_input DECIMAL(12, 6),                      -- 输入Token成本（美元/百万）
    cost_per_million_output DECIMAL(12, 6),                     -- 输出Token成本（美元/百万）
    is_enabled BOOLEAN DEFAULT TRUE,                            -- 是否启用
    health_status ENUM('healthy', 'slow', 'down') DEFAULT 'healthy',
    last_used_at TIMESTAMP NULL,
    failure_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (model_id) REFERENCES ai_models(id) ON DELETE CASCADE,
    INDEX idx_model_enabled (model_id, is_enabled, priority),
    INDEX idx_health_status (health_status)
) COMMENT = '供应商API密钥池，支持同一模型多KEY负载均衡与故障转移';

-- =============================================================================
-- 5. 配额与成本 (Quotas & Costs)
-- =============================================================================

-- 5.1 多维度Token配额表 (token_quotas)
CREATE TABLE token_quotas (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    department_id VARCHAR(36) NOT NULL,                         -- 配额所属部门
    model_id VARCHAR(50),                                       -- NULL表示该部门对所有模型的总配额
    quota_type ENUM('total', 'daily', 'monthly') DEFAULT 'monthly',
    max_tokens BIGINT NOT NULL,                                 -- 周期内最大Token数
    used_tokens BIGINT DEFAULT 0,                               -- 当前周期已用（需要定期重置或累加）
    reset_at TIMESTAMP,                                         -- 下次重置时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE,
    FOREIGN KEY (model_id) REFERENCES ai_models(id) ON DELETE CASCADE,
    UNIQUE KEY uk_department_model_type (department_id, model_id, quota_type),
    INDEX idx_reset_at (reset_at)
) COMMENT = '多维度Token配额表，支持部门、模型、周期的组合限额';

-- 5.2 月度成本聚合表 (monthly_costs)
CREATE TABLE monthly_costs (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    department_id VARCHAR(36) NOT NULL,
    `year_month` CHAR(7) NOT NULL,                              -- 年月，格式 '2024-05'
    model_id VARCHAR(50) NOT NULL,
    total_requests INT DEFAULT 0,
    total_prompt_tokens BIGINT DEFAULT 0,
    total_completion_tokens BIGINT DEFAULT 0,
    total_tokens BIGINT DEFAULT 0,
    estimated_cost_usd DECIMAL(14, 4) DEFAULT 0,                -- 月度总成本
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (model_id) REFERENCES ai_models(id),
    UNIQUE KEY uk_department_month_model (department_id, `year_month`, model_id)
) COMMENT = '月度成本聚合表，用于快速生成部门账单和报告';

-- =============================================================================
-- 6. 审计与安全日志 (Audit & Logs)
-- =============================================================================

-- 6.1 请求审计日志表 (request_logs)
CREATE TABLE request_logs (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),                -- 唯一请求ID，用于全链路追踪
    request_id VARCHAR(100),                                    -- 网关生成的唯一标识
    user_id VARCHAR(36) NOT NULL,                               -- 请求用户
    department_id VARCHAR(36) NOT NULL,                         -- 请求部门（冗余存储，避免关联查询）
    model_id VARCHAR(50) NOT NULL,                              -- 请求的模型
    provider_key_id VARCHAR(36),                                -- 实际使用的供应商KEY
    provider_name VARCHAR(50),                                  -- 供应商名（冗余）
    
    -- 请求详情
    prompt TEXT,                                                -- 原始请求提示词（考虑安全可脱敏存储）
    request_body JSON,                                          -- 完整请求体
    request_model VARCHAR(100),                                 -- 用户实际请求的模型名
    
    -- 响应详情
    response_body LONGTEXT,                                     -- 完整响应体（非流式）
    response_status INT,                                        -- HTTP状态码
    error_message TEXT,
    
    -- Token与成本核算
    prompt_tokens INT DEFAULT 0,
    completion_tokens INT DEFAULT 0,
    total_tokens INT DEFAULT 0,
    estimated_cost_usd DECIMAL(12, 6) DEFAULT 0,
    
    -- 性能指标
    gateway_latency_ms INT,                                     -- 网关内部处理耗时
    total_latency_ms INT,                                       -- 总耗时（用户感受到的）
    upstream_latency_ms INT,                                    -- 上游供应商处理耗时
    
    -- 时间戳
    requested_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6),     -- 高精度时间戳
    responded_at TIMESTAMP(6) NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (model_id) REFERENCES ai_models(id),
    FOREIGN KEY (provider_key_id) REFERENCES api_provider_keys(id),
    INDEX idx_time_range (requested_at),                        -- 按时间范围查询（最重要）
    INDEX idx_user_department (user_id, department_id),
    INDEX idx_model_provider (model_id, provider_name),
    INDEX idx_request_id (request_id)
) COMMENT = '请求审计日志表，所有AI调用的黄金数据源，数据量巨大，考虑分区';

-- 6.2 安全日志表 (security_logs)
CREATE TABLE security_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(36),                                        -- 可为空（未登录操作）
    
    -- 操作信息
    action VARCHAR(50) NOT NULL,                                -- login_success, login_failed, logout, password_change, profile_update, etc.
    action_type VARCHAR(20) NOT NULL,                           -- authentication, authorization, profile, security
    ip_address VARCHAR(45) NOT NULL,                            -- [Fixed] INET -> VARCHAR(45) for MySQL compatibility
    user_agent TEXT,
    
    -- 详情
    details JSONB,                                              -- 操作详情
    risk_level VARCHAR(20) DEFAULT 'low',                       -- low, medium, high, critical
    is_suspicious BOOLEAN DEFAULT FALSE,                        -- 是否可疑
    
    -- 位置信息
    country VARCHAR(100),
    city VARCHAR(100),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    
    -- 设备信息
    device_type VARCHAR(20),
    device_id VARCHAR(100),
    
    -- 关联信息
    session_id VARCHAR(64),
    request_id VARCHAR(100),                                    -- 请求唯一标识
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- 索引
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at),
    INDEX idx_ip_action (ip_address, action),
    INDEX idx_risk_level (risk_level),
    
    -- 外键
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (session_id) REFERENCES user_sessions(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='安全日志表';
