-- -----------------------------------------------------------------------------
-- Next AI Gateway Data Seeding Script
-- Generated based on init.sql schema
-- -----------------------------------------------------------------------------

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =============================================================================
-- 1. 组织架构 (Organization)
-- =============================================================================
INSERT INTO departments (id, name, parent_id, cost_center_code, monthly_token_budget, is_active) VALUES
('dept_root',   '总部 (Headquarters)',      NULL,          'CC-001', 100000000, 1),
('dept_tech',   '技术中心 (Tech Center)',   'dept_root',   'CC-101', 50000000,  1),
('dept_mkt',    '市场部 (Marketing)',       'dept_root',   'CC-201', 20000000,  1),
('dept_sales',  '销售部 (Sales)',           'dept_root',   'CC-301', 15000000,  1),
('dept_hr',     '人力资源 (HR)',            'dept_root',   'CC-401', 5000000,   1),
('dept_ai_lab', 'AI研究院 (AI Lab)',        'dept_tech',   'CC-102', 30000000,  1),
('dept_devops', '运维部 (DevOps)',          'dept_tech',   'CC-103', 10000000,  1);

-- =============================================================================
-- 2. 用户基础数据 (Users)
-- =============================================================================
INSERT INTO users (id, username, email, mobile, employee_id, name, department_id, role, status, type) VALUES
('usr_admin',    'admin',      'admin@nextai.com',      '13800000001', 'EMP001', '系统管理员', 'dept_root',    'admin',   1, 3),
('usr_cto',      'cto',        'cto@nextai.com',        '13800000002', 'EMP002', '技术总监',   'dept_tech',    'admin',   1, 2),
('usr_dev_lead', 'dev_lead',   'lead@nextai.com',       '13800000003', 'EMP003', '研发组长',   'dept_ai_lab',  'user',    1, 2),
('usr_dev_01',   'dev01',      'dev01@nextai.com',      '13800000004', 'EMP004', '研发工程师A','dept_ai_lab',  'user',    1, 1),
('usr_dev_02',   'dev02',      'dev02@nextai.com',      '13800000005', 'EMP005', '研发工程师B','dept_ai_lab',  'user',    1, 1),
('usr_mkt_mgr',  'mkt_mgr',    'mkt@nextai.com',        '13800000006', 'EMP006', '市场经理',   'dept_mkt',     'user',    1, 1),
('usr_sales_01', 'sales01',    'sales01@nextai.com',    '13800000007', 'EMP007', '销售代表',   'dept_sales',   'user',    1, 1),
('usr_finance',  'finance',    'finance@nextai.com',    '13800000008', 'EMP008', '财务专员',   'dept_root',    'billing', 1, 1);

-- =============================================================================
-- 3. 用户扩展数据 (User Extensions)
-- =============================================================================

-- 3.1 密码 (默认密码: Password123!)
INSERT INTO user_passwords (user_id, password_hash, password_salt, algorithm) VALUES
('usr_admin',    '$2a$10$X7.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6', 'salt_admin', 'bcrypt'),
('usr_cto',      '$2a$10$X7.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6', 'salt_cto',   'bcrypt'),
('usr_dev_lead', '$2a$10$X7.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6.G.6', 'salt_lead',  'bcrypt');

-- 3.2 资料
INSERT INTO user_profiles (user_id, nickname, country, city, bio) VALUES
('usr_admin',    'Root',       'China', 'Beijing',  'System Administrator'),
('usr_cto',      'TechMaster', 'China', 'Shanghai', 'CTO of NextAI'),
('usr_dev_lead', 'AI_Guru',    'China', 'Hangzhou', 'Focus on LLM');

-- 3.3 设置
INSERT INTO user_settings (user_id, ui_theme, ui_language, notify_email) VALUES
('usr_admin',    'dark',  'en-US', '{"system_alerts": true, "weekly_report": true}'),
('usr_dev_lead', 'light', 'zh-CN', '{"system_alerts": false, "weekly_report": true}');

-- 3.4 安全状态
INSERT INTO user_security (user_id, email_verified, mfa_enabled, last_login_ip) VALUES
('usr_admin',    1, 1, '192.168.1.10'),
('usr_cto',      1, 1, '192.168.1.11'),
('usr_dev_lead', 1, 0, '192.168.1.12');

-- 3.5 统计
INSERT INTO user_statistics (user_id, login_count, total_spent) VALUES
('usr_admin',    1520, 0.00),
('usr_dev_lead', 345,  150.50);

-- =============================================================================
-- 4. AI模型配置 (AI Models)
-- =============================================================================
INSERT INTO ai_models (id, description, default_provider_id, routing_strategy, config_json) VALUES
('gpt-4-turbo',    'OpenAI GPT-4 Turbo',      NULL, 'priority',    '{"temperature": 0.7, "max_tokens": 4096}'),
('gpt-3.5-turbo',  'OpenAI GPT-3.5 Turbo',    NULL, 'lowest_cost', '{"temperature": 0.5, "max_tokens": 2048}'),
('claude-3-opus',  'Anthropic Claude 3 Opus', NULL, 'load_balance','{"temperature": 0.7, "max_tokens": 4096}'),
('claude-3-sonnet','Anthropic Claude 3 Sonnet',NULL, 'priority',    '{"temperature": 0.7, "max_tokens": 4096}'),
('gemini-pro',     'Google Gemini Pro',       NULL, 'load_balance','{"temperature": 0.8}'),
('llama-3-70b',    'Meta Llama 3 70B',        NULL, 'lowest_cost', '{"temperature": 0.7}');

-- =============================================================================
-- 5. 供应商密钥 (Provider Keys)
-- =============================================================================
INSERT INTO api_provider_keys (id, model_id, provider_name, api_key_encrypted, priority, cost_per_million_input, cost_per_million_output) VALUES
('key_openai_1',  'gpt-4-turbo',    'openai',    'sk-enc-xxxxxxxx1', 1, 10.00, 30.00),
('key_azure_1',   'gpt-4-turbo',    'azure',     'az-enc-xxxxxxxx1', 2, 10.00, 30.00),
('key_openai_2',  'gpt-3.5-turbo',  'openai',    'sk-enc-xxxxxxxx2', 1, 0.50,  1.50),
('key_anth_1',    'claude-3-opus',  'anthropic', 'sk-ant-xxxxxxx1', 1, 15.00, 75.00),
('key_anth_2',    'claude-3-sonnet','anthropic', 'sk-ant-xxxxxxx2', 1, 3.00,  15.00),
('key_google_1',  'gemini-pro',     'google',    'ai-enc-xxxxxxxx1', 1, 0.50,  1.50),
('key_groq_1',    'llama-3-70b',    'groq',      'gsk-enc-xxxxxxx1', 1, 0.59,  0.79);

-- =============================================================================
-- 6. 配额管理 (Quotas)
-- =============================================================================
INSERT INTO token_quotas (department_id, model_id, quota_type, max_tokens, used_tokens) VALUES
-- AI研究院：总额度
('dept_ai_lab', NULL,           'monthly', 100000000, 5000000),
-- AI研究院：GPT-4 特别限额
('dept_ai_lab', 'gpt-4-turbo',  'monthly', 20000000,  1200000),
-- 市场部：总额度
('dept_mkt',    NULL,           'monthly', 10000000,  800000),
-- 研发组长：个人测试限额 (假设系统支持按用户配额，这里暂时演示部门级)
('dept_tech',   'claude-3-opus','daily',   1000000,   50000);

-- =============================================================================
-- 7. 模拟日志 (Logs - Samples)
-- =============================================================================
INSERT INTO request_logs (id, request_id, user_id, department_id, model_id, provider_name, prompt_tokens, completion_tokens, total_tokens, estimated_cost_usd, gateway_latency_ms, upstream_latency_ms, requested_at) VALUES
(UUID(), 'req_001', 'usr_dev_lead', 'dept_ai_lab', 'gpt-4-turbo',   'openai',    500,  200, 700, 0.011, 20, 1500, NOW() - INTERVAL 1 HOUR),
(UUID(), 'req_002', 'usr_dev_01',   'dept_ai_lab', 'gpt-3.5-turbo', 'openai',    1200, 500, 1700, 0.002, 15, 800,  NOW() - INTERVAL 50 MINUTE),
(UUID(), 'req_003', 'usr_mkt_mgr',  'dept_mkt',    'claude-3-sonnet','anthropic', 3000, 800, 3800, 0.021, 25, 2100, NOW() - INTERVAL 30 MINUTE),
(UUID(), 'req_004', 'usr_admin',    'dept_root',   'llama-3-70b',   'groq',      100,  50,  150,  0.0001, 10, 200,  NOW() - INTERVAL 10 MINUTE);

-- =============================================================================
-- 8. 成本聚合 (Monthly Costs)
-- =============================================================================
INSERT INTO monthly_costs (department_id, `year_month`, model_id, total_requests, total_tokens, estimated_cost_usd) VALUES
('dept_ai_lab', '2024-05', 'gpt-4-turbo',   1500, 5000000, 120.50),
('dept_ai_lab', '2024-05', 'gpt-3.5-turbo', 5000, 8000000, 12.00),
('dept_mkt',    '2024-05', 'claude-3-sonnet',200,  1000000, 6.00);

-- =============================================================================
-- 9. 安全日志 (Security Logs)
-- =============================================================================
INSERT INTO security_logs (user_id, action, action_type, ip_address, country, details) VALUES
('usr_admin', 'login_success', 'authentication', '192.168.1.10', 'China', '{"method": "password"}'),
('usr_dev_01', 'login_failed', 'authentication', '10.0.0.5',     'Unknown', '{"reason": "wrong_password"}');

SET FOREIGN_KEY_CHECKS = 1;
