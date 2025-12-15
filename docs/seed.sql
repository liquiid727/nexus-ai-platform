-- ----------------------------
-- Data Seeding Script
-- Generated for Next AI Gateway
-- ----------------------------

-- Disable foreign key checks temporarily to allow bulk insertion order flexibility if needed
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 1. Departments (20 records)
-- ----------------------------
INSERT INTO departments (id, name, parent_id, cost_center_code, monthly_token_budget, is_active) VALUES
('dept_001', 'Headquarters', NULL, 'CC001', 100000000, 1),
('dept_002', 'R&D Division', 'dept_001', 'CC002', 50000000, 1),
('dept_003', 'Sales Division', 'dept_001', 'CC003', 20000000, 1),
('dept_004', 'Marketing', 'dept_003', 'CC004', 10000000, 1),
('dept_005', 'HR', 'dept_001', 'CC005', 5000000, 1),
('dept_006', 'Finance', 'dept_001', 'CC006', 5000000, 1),
('dept_007', 'IT Support', 'dept_002', 'CC007', 8000000, 1),
('dept_008', 'AI Research', 'dept_002', 'CC008', 30000000, 1),
('dept_009', 'Product Management', 'dept_002', 'CC009', 5000000, 1),
('dept_010', 'Legal', 'dept_001', 'CC010', 2000000, 1),
('dept_011', 'Customer Success', 'dept_003', 'CC011', 5000000, 1),
('dept_012', 'Operations', 'dept_001', 'CC012', 4000000, 1),
('dept_013', 'Security', 'dept_002', 'CC013', 5000000, 1),
('dept_014', 'Data Science', 'dept_002', 'CC014', 15000000, 1),
('dept_015', 'DevOps', 'dept_002', 'CC015', 10000000, 1),
('dept_016', 'QA', 'dept_002', 'CC016', 5000000, 1),
('dept_017', 'Design', 'dept_009', 'CC017', 3000000, 1),
('dept_018', 'Content', 'dept_004', 'CC018', 4000000, 1),
('dept_019', 'Events', 'dept_004', 'CC019', 2000000, 1),
('dept_020', 'Strategy', 'dept_001', 'CC020', 1000000, 1);

-- ----------------------------
-- 2. Users (20 records)
-- ----------------------------
INSERT INTO users (id, employee_id, name, email, department_id, role, is_active) VALUES
('user_001', 'EMP001', 'Alice Smith', 'alice@example.com', 'dept_001', 'admin', 1),
('user_002', 'EMP002', 'Bob Jones', 'bob@example.com', 'dept_002', 'user', 1),
('user_003', 'EMP003', 'Charlie Brown', 'charlie@example.com', 'dept_002', 'user', 1),
('user_004', 'EMP004', 'David Wilson', 'david@example.com', 'dept_008', 'user', 1),
('user_005', 'EMP005', 'Eve Davis', 'eve@example.com', 'dept_008', 'user', 1),
('user_006', 'EMP006', 'Frank Miller', 'frank@example.com', 'dept_003', 'user', 1),
('user_007', 'EMP007', 'Grace Lee', 'grace@example.com', 'dept_004', 'user', 1),
('user_008', 'EMP008', 'Hank Green', 'hank@example.com', 'dept_005', 'user', 1),
('user_009', 'EMP009', 'Ivy White', 'ivy@example.com', 'dept_006', 'billing', 1),
('user_010', 'EMP010', 'Jack Black', 'jack@example.com', 'dept_007', 'user', 1),
('user_011', 'EMP011', 'Kevin King', 'kevin@example.com', 'dept_014', 'user', 1),
('user_012', 'EMP012', 'Laura Scott', 'laura@example.com', 'dept_014', 'user', 1),
('user_013', 'EMP013', 'Mike Ross', 'mike@example.com', 'dept_015', 'user', 1),
('user_014', 'EMP014', 'Nina Hall', 'nina@example.com', 'dept_015', 'user', 1),
('user_015', 'EMP015', 'Oscar Young', 'oscar@example.com', 'dept_016', 'user', 1),
('user_016', 'EMP016', 'Paul Adams', 'paul@example.com', 'dept_017', 'user', 1),
('user_017', 'EMP017', 'Quinn Baker', 'quinn@example.com', 'dept_018', 'user', 1),
('user_018', 'EMP018', 'Rachel Clark', 'rachel@example.com', 'dept_011', 'user', 1),
('user_019', 'EMP019', 'Steve Wright', 'steve@example.com', 'dept_013', 'user', 1),
('user_020', 'EMP020', 'Tina Evans', 'tina@example.com', 'dept_009', 'user', 1);

-- ----------------------------
-- 3. AI Models (20 records)
-- ----------------------------
INSERT INTO ai_models (id, description, default_provider_id, routing_strategy, is_active, config_json) VALUES
('gpt-4', 'GPT-4 High Intelligence', NULL, 'load_balance', 1, '{"temperature": 0.7}'),
('gpt-4-turbo', 'GPT-4 Turbo Fast', NULL, 'priority', 1, '{"temperature": 0.7}'),
('gpt-3.5-turbo', 'GPT-3.5 Cost Effective', NULL, 'lowest_cost', 1, '{"temperature": 0.5}'),
('claude-3-opus', 'Claude 3 Opus', NULL, 'load_balance', 1, '{"max_tokens": 4096}'),
('claude-3-sonnet', 'Claude 3 Sonnet', NULL, 'priority', 1, '{"max_tokens": 4096}'),
('claude-3-haiku', 'Claude 3 Haiku', NULL, 'lowest_cost', 1, '{"max_tokens": 2048}'),
('gemini-pro', 'Gemini Pro', NULL, 'load_balance', 1, '{}'),
('gemini-ultra', 'Gemini Ultra', NULL, 'priority', 1, '{}'),
('llama-3-70b', 'Llama 3 70B Open Source', NULL, 'load_balance', 1, '{}'),
('llama-3-8b', 'Llama 3 8B Fast', NULL, 'load_balance', 1, '{}'),
('mistral-large', 'Mistral Large', NULL, 'priority', 1, '{}'),
('mistral-medium', 'Mistral Medium', NULL, 'lowest_cost', 1, '{}'),
('azure-gpt-4', 'Azure GPT-4', NULL, 'load_balance', 1, '{}'),
('azure-gpt-35', 'Azure GPT-3.5', NULL, 'load_balance', 1, '{}'),
('text-embedding-3-small', 'Embedding Small', NULL, 'load_balance', 1, '{}'),
('text-embedding-3-large', 'Embedding Large', NULL, 'load_balance', 1, '{}'),
('dall-e-3', 'Image Generation', NULL, 'load_balance', 1, '{}'),
('tts-1', 'Text to Speech', NULL, 'load_balance', 1, '{}'),
('whisper-1', 'Speech to Text', NULL, 'load_balance', 1, '{}'),
('custom-fin-model', 'Finance Fine-tuned', NULL, 'priority', 1, '{}');

-- ----------------------------
-- 4. API Provider Keys (20 records)
-- ----------------------------
INSERT INTO api_provider_keys (id, model_id, provider_name, api_key_encrypted, endpoint_base_url, priority, max_rpm, max_tpm, cost_per_million_input, cost_per_million_output, is_enabled) VALUES
('key_001', 'gpt-4', 'openai', 'enc_key_1', 'https://api.openai.com/v1', 1, 10000, 1000000, 30.00, 60.00, 1),
('key_002', 'gpt-4', 'azure', 'enc_key_2', 'https://azure.microsoft.com', 2, 20000, 2000000, 30.00, 60.00, 1),
('key_003', 'gpt-4-turbo', 'openai', 'enc_key_3', 'https://api.openai.com/v1', 1, 15000, 1500000, 10.00, 30.00, 1),
('key_004', 'gpt-3.5-turbo', 'openai', 'enc_key_4', 'https://api.openai.com/v1', 1, 50000, 5000000, 0.50, 1.50, 1),
('key_005', 'claude-3-opus', 'anthropic', 'enc_key_5', 'https://api.anthropic.com', 1, 5000, 500000, 15.00, 75.00, 1),
('key_006', 'claude-3-sonnet', 'anthropic', 'enc_key_6', 'https://api.anthropic.com', 1, 10000, 1000000, 3.00, 15.00, 1),
('key_007', 'claude-3-haiku', 'anthropic', 'enc_key_7', 'https://api.anthropic.com', 1, 20000, 2000000, 0.25, 1.25, 1),
('key_008', 'gemini-pro', 'google', 'enc_key_8', 'https://generativelanguage.googleapis.com', 1, 10000, 1000000, 0.50, 1.50, 1),
('key_009', 'llama-3-70b', 'groq', 'enc_key_9', 'https://api.groq.com', 1, 5000, 500000, 0.70, 0.80, 1),
('key_010', 'llama-3-8b', 'groq', 'enc_key_10', 'https://api.groq.com', 1, 10000, 1000000, 0.10, 0.10, 1),
('key_011', 'mistral-large', 'mistral', 'enc_key_11', 'https://api.mistral.ai', 1, 5000, 500000, 8.00, 24.00, 1),
('key_012', 'azure-gpt-4', 'azure', 'enc_key_12', 'https://azure.microsoft.com', 1, 10000, 1000000, 30.00, 60.00, 1),
('key_013', 'text-embedding-3-small', 'openai', 'enc_key_13', 'https://api.openai.com/v1', 1, 100000, 10000000, 0.02, 0.00, 1),
('key_014', 'dall-e-3', 'openai', 'enc_key_14', 'https://api.openai.com/v1', 1, 500, 0, 40.00, 0.00, 1),
('key_015', 'gpt-4', 'openai', 'enc_key_15', 'https://api.openai.com/v1', 2, 5000, 500000, 30.00, 60.00, 1),
('key_016', 'gpt-4-turbo', 'azure', 'enc_key_16', 'https://azure.microsoft.com', 2, 10000, 1000000, 10.00, 30.00, 1),
('key_017', 'claude-3-opus', 'aws-bedrock', 'enc_key_17', 'https://aws.amazon.com', 2, 5000, 500000, 15.00, 75.00, 1),
('key_018', 'gemini-ultra', 'google', 'enc_key_18', 'https://generativelanguage.googleapis.com', 1, 5000, 500000, 5.00, 10.00, 1),
('key_019', 'mistral-medium', 'mistral', 'enc_key_19', 'https://api.mistral.ai', 1, 8000, 800000, 2.70, 8.10, 1),
('key_020', 'custom-fin-model', 'local', 'enc_key_20', 'http://10.0.0.1:8000', 1, 1000, 100000, 0.00, 0.00, 1);

-- ----------------------------
-- 5. Token Quotas (20 records)
-- ----------------------------
INSERT INTO token_quotas (id, department_id, model_id, quota_type, max_tokens, used_tokens) VALUES
('quota_001', 'dept_001', NULL, 'monthly', 100000000, 5000000),
('quota_002', 'dept_002', NULL, 'monthly', 50000000, 2000000),
('quota_003', 'dept_008', 'gpt-4', 'monthly', 10000000, 1500000),
('quota_004', 'dept_008', 'claude-3-opus', 'monthly', 5000000, 200000),
('quota_005', 'dept_003', NULL, 'monthly', 20000000, 1000000),
('quota_006', 'dept_004', NULL, 'monthly', 10000000, 500000),
('quota_007', 'dept_014', 'gpt-4-turbo', 'monthly', 5000000, 100000),
('quota_008', 'dept_014', 'llama-3-70b', 'monthly', 10000000, 50000),
('quota_009', 'dept_001', 'gpt-4', 'daily', 100000, 2000),
('quota_010', 'dept_002', 'gpt-4', 'daily', 50000, 1000),
('quota_011', 'dept_007', NULL, 'monthly', 8000000, 400000),
('quota_012', 'dept_005', NULL, 'monthly', 5000000, 100000),
('quota_013', 'dept_006', NULL, 'monthly', 5000000, 50000),
('quota_014', 'dept_018', NULL, 'monthly', 4000000, 200000),
('quota_015', 'dept_017', 'dall-e-3', 'monthly', 100000, 500),
('quota_016', 'dept_002', 'gpt-3.5-turbo', 'monthly', 20000000, 500000),
('quota_017', 'dept_015', NULL, 'monthly', 10000000, 300000),
('quota_018', 'dept_016', NULL, 'monthly', 5000000, 10000),
('quota_019', 'dept_008', 'gemini-pro', 'monthly', 5000000, 50000),
('quota_020', 'dept_020', NULL, 'monthly', 1000000, 1000);

-- ----------------------------
-- 6. Request Logs (20 records)
-- ----------------------------
INSERT INTO request_logs (id, request_id, user_id, department_id, model_id, provider_key_id, provider_name, prompt_tokens, completion_tokens, total_tokens, estimated_cost_usd, gateway_latency_ms, total_latency_ms) VALUES
('log_001', 'req_001', 'user_001', 'dept_001', 'gpt-4', 'key_001', 'openai', 100, 50, 150, 0.006, 10, 1500),
('log_002', 'req_002', 'user_002', 'dept_002', 'gpt-3.5-turbo', 'key_004', 'openai', 500, 100, 600, 0.0004, 5, 500),
('log_003', 'req_003', 'user_004', 'dept_008', 'gpt-4-turbo', 'key_003', 'openai', 200, 200, 400, 0.008, 8, 2000),
('log_004', 'req_004', 'user_004', 'dept_008', 'claude-3-opus', 'key_005', 'anthropic', 1000, 500, 1500, 0.0525, 12, 4000),
('log_005', 'req_005', 'user_006', 'dept_003', 'gpt-3.5-turbo', 'key_004', 'openai', 100, 20, 120, 0.00008, 5, 300),
('log_006', 'req_006', 'user_011', 'dept_014', 'llama-3-70b', 'key_009', 'groq', 800, 200, 1000, 0.00072, 6, 400),
('log_007', 'req_007', 'user_001', 'dept_001', 'gpt-4', 'key_002', 'azure', 150, 150, 300, 0.0135, 11, 1600),
('log_008', 'req_008', 'user_016', 'dept_017', 'dall-e-3', 'key_014', 'openai', 100, 0, 100, 0.040, 20, 5000),
('log_009', 'req_009', 'user_005', 'dept_008', 'gemini-pro', 'key_008', 'google', 300, 300, 600, 0.0006, 7, 800),
('log_010', 'req_010', 'user_002', 'dept_002', 'gpt-4', 'key_001', 'openai', 50, 50, 100, 0.0045, 9, 1200),
('log_011', 'req_011', 'user_003', 'dept_002', 'gpt-3.5-turbo', 'key_004', 'openai', 1000, 200, 1200, 0.0008, 6, 700),
('log_012', 'req_012', 'user_007', 'dept_004', 'gpt-3.5-turbo', 'key_004', 'openai', 200, 50, 250, 0.000175, 5, 350),
('log_013', 'req_013', 'user_008', 'dept_005', 'gpt-4', 'key_001', 'openai', 100, 100, 200, 0.009, 10, 1800),
('log_014', 'req_014', 'user_012', 'dept_014', 'claude-3-sonnet', 'key_006', 'anthropic', 500, 500, 1000, 0.009, 8, 2500),
('log_015', 'req_015', 'user_013', 'dept_015', 'gpt-3.5-turbo', 'key_004', 'openai', 100, 10, 110, 0.000065, 4, 200),
('log_016', 'req_016', 'user_014', 'dept_015', 'gpt-3.5-turbo', 'key_004', 'openai', 200, 20, 220, 0.00013, 4, 250),
('log_017', 'req_017', 'user_015', 'dept_016', 'gpt-4-turbo', 'key_003', 'openai', 100, 100, 200, 0.004, 8, 1500),
('log_018', 'req_018', 'user_018', 'dept_011', 'gpt-3.5-turbo', 'key_004', 'openai', 300, 50, 350, 0.000225, 5, 400),
('log_019', 'req_019', 'user_019', 'dept_013', 'gpt-4', 'key_001', 'openai', 50, 50, 100, 0.0045, 9, 1300),
('log_020', 'req_020', 'user_020', 'dept_009', 'gpt-3.5-turbo', 'key_004', 'openai', 100, 10, 110, 0.000065, 4, 220);

-- ----------------------------
-- 7. Monthly Costs (20 records)
-- ----------------------------
INSERT INTO monthly_costs (id, department_id, `year_month`, model_id, total_requests, total_tokens, estimated_cost_usd) VALUES
('cost_001', 'dept_001', '2024-05', 'gpt-4', 100, 50000, 2.50),
('cost_002', 'dept_002', '2024-05', 'gpt-3.5-turbo', 500, 100000, 0.15),
('cost_003', 'dept_008', '2024-05', 'gpt-4-turbo', 200, 80000, 1.20),
('cost_004', 'dept_008', '2024-05', 'claude-3-opus', 50, 50000, 1.50),
('cost_005', 'dept_003', '2024-05', 'gpt-3.5-turbo', 300, 60000, 0.09),
('cost_006', 'dept_014', '2024-05', 'llama-3-70b', 100, 40000, 0.03),
('cost_007', 'dept_017', '2024-05', 'dall-e-3', 20, 2000, 0.80),
('cost_008', 'dept_008', '2024-05', 'gemini-pro', 100, 50000, 0.05),
('cost_009', 'dept_004', '2024-05', 'gpt-3.5-turbo', 150, 30000, 0.045),
('cost_010', 'dept_005', '2024-05', 'gpt-4', 50, 10000, 0.45),
('cost_011', 'dept_014', '2024-05', 'claude-3-sonnet', 80, 80000, 0.60),
('cost_012', 'dept_015', '2024-05', 'gpt-3.5-turbo', 200, 40000, 0.06),
('cost_013', 'dept_016', '2024-05', 'gpt-4-turbo', 30, 10000, 0.20),
('cost_014', 'dept_011', '2024-05', 'gpt-3.5-turbo', 100, 20000, 0.03),
('cost_015', 'dept_013', '2024-05', 'gpt-4', 20, 5000, 0.225),
('cost_016', 'dept_009', '2024-05', 'gpt-3.5-turbo', 50, 10000, 0.015),
('cost_017', 'dept_001', '2024-04', 'gpt-4', 80, 40000, 2.00),
('cost_018', 'dept_002', '2024-04', 'gpt-3.5-turbo', 400, 80000, 0.12),
('cost_019', 'dept_008', '2024-04', 'gpt-4', 150, 75000, 3.375),
('cost_020', 'dept_003', '2024-04', 'gpt-3.5-turbo', 200, 40000, 0.06);

SET FOREIGN_KEY_CHECKS = 1;
