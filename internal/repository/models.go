package repository

import (
	"time"
)

// Department 组织架构表
type Department struct {
	ID                 string    `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())" json:"id"`
	Name               string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	ParentID           *string   `gorm:"column:parent_id;type:varchar(36);index:idx_parent_id" json:"parent_id,omitempty"`
	CostCenterCode     string    `gorm:"column:cost_center_code;type:varchar(50);index:idx_cost_center" json:"cost_center_code"`
	MonthlyTokenBudget int64     `gorm:"column:monthly_token_budget;type:bigint" json:"monthly_token_budget"`
	IsActive           bool      `gorm:"column:is_active;type:boolean;default:true" json:"is_active"`
	CreatedAt          time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName returns the table name
func (Department) TableName() string {
	return "departments"
}

// User 用户表
type User struct {
	ID           string    `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())" json:"id"`
	EmployeeID   string    `gorm:"column:employee_id;type:varchar(50);unique;not null" json:"employee_id"`
	Name         string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Email        string    `gorm:"column:email;type:varchar(255);not null;index:idx_email" json:"email"`
	DepartmentID string    `gorm:"column:department_id;type:varchar(36);not null;index:idx_department" json:"department_id"`
	Role         string    `gorm:"column:role;type:enum('admin','user','billing');default:'user'" json:"role"`
	IsActive     bool      `gorm:"column:is_active;type:boolean;default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`

	Department Department `gorm:"foreignKey:DepartmentID;references:ID" json:"department,omitempty"`
}

// TableName returns the table name
func (User) TableName() string {
	return "users"
}

// AiModel AI模型路由配置表
type AiModel struct {
	ID                string    `gorm:"column:id;type:varchar(50);primaryKey" json:"id"`
	Description       string    `gorm:"column:description;type:varchar(255)" json:"description"`
	DefaultProviderID string    `gorm:"column:default_provider_id;type:varchar(36)" json:"default_provider_id"`
	RoutingStrategy   string    `gorm:"column:routing_strategy;type:enum('load_balance','priority','lowest_cost');default:'load_balance'" json:"routing_strategy"`
	IsActive          bool      `gorm:"column:is_active;type:boolean;default:true" json:"is_active"`
	ConfigJSON        []byte    `gorm:"column:config_json;type:json" json:"config_json"`
	CreatedAt         time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName returns the table name
func (AiModel) TableName() string {
	return "ai_models"
}

// ApiProviderKey 供应商API密钥池
type ApiProviderKey struct {
	ID                   string     `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())" json:"id"`
	ModelID              string     `gorm:"column:model_id;type:varchar(50);not null" json:"model_id"`
	ProviderName         string     `gorm:"column:provider_name;type:varchar(50);not null" json:"provider_name"`
	ApiKeyEncrypted      string     `gorm:"column:api_key_encrypted;type:text;not null" json:"-"` // Don't expose key in JSON
	EndpointBaseURL      string     `gorm:"column:endpoint_base_url;type:varchar(500)" json:"endpoint_base_url"`
	Priority             int        `gorm:"column:priority;type:int;default:1" json:"priority"`
	MaxRPM               int        `gorm:"column:max_rpm;type:int" json:"max_rpm"`
	MaxTPM               int        `gorm:"column:max_tpm;type:int" json:"max_tpm"`
	CostPerMillionInput  float64    `gorm:"column:cost_per_million_input;type:decimal(12,6)" json:"cost_per_million_input"`
	CostPerMillionOutput float64    `gorm:"column:cost_per_million_output;type:decimal(12,6)" json:"cost_per_million_output"`
	IsEnabled            bool       `gorm:"column:is_enabled;type:boolean;default:true;index:idx_model_enabled" json:"is_enabled"`
	HealthStatus         string     `gorm:"column:health_status;type:enum('healthy','slow','down');default:'healthy';index:idx_health_status" json:"health_status"`
	LastUsedAt           *time.Time `gorm:"column:last_used_at;type:timestamp" json:"last_used_at"`
	FailureCount         int        `gorm:"column:failure_count;type:int;default:0" json:"failure_count"`
	CreatedAt            time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`

	Model AiModel `gorm:"foreignKey:ModelID;references:ID" json:"model,omitempty"`
}

// TableName returns the table name
func (ApiProviderKey) TableName() string {
	return "api_provider_keys"
}

// TokenQuota 配额管理表
type TokenQuota struct {
	ID           string     `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())" json:"id"`
	DepartmentID string     `gorm:"column:department_id;type:varchar(36);not null" json:"department_id"`
	ModelID      *string    `gorm:"column:model_id;type:varchar(50)" json:"model_id"`
	QuotaType    string     `gorm:"column:quota_type;type:enum('total','daily','monthly');default:'monthly'" json:"quota_type"`
	MaxTokens    int64      `gorm:"column:max_tokens;type:bigint;not null" json:"max_tokens"`
	UsedTokens   int64      `gorm:"column:used_tokens;type:bigint;default:0" json:"used_tokens"`
	ResetAt      *time.Time `gorm:"column:reset_at;type:timestamp;index:idx_reset_at" json:"reset_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`

	Department Department `gorm:"foreignKey:DepartmentID;references:ID" json:"department,omitempty"`
	Model      *AiModel   `gorm:"foreignKey:ModelID;references:ID" json:"model,omitempty"`
}

// TableName returns the table name
func (TokenQuota) TableName() string {
	return "token_quotas"
}

// RequestLog 审计与日志表
type RequestLog struct {
	ID                string     `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())" json:"id"`
	RequestID         string     `gorm:"column:request_id;type:varchar(100);index:idx_request_id" json:"request_id"`
	UserID            string     `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_department" json:"user_id"`
	DepartmentID      string     `gorm:"column:department_id;type:varchar(36);not null;index:idx_user_department" json:"department_id"`
	ModelID           string     `gorm:"column:model_id;type:varchar(50);not null;index:idx_model_provider" json:"model_id"`
	ProviderKeyID     *string    `gorm:"column:provider_key_id;type:varchar(36)" json:"provider_key_id"`
	ProviderName      string     `gorm:"column:provider_name;type:varchar(50);index:idx_model_provider" json:"provider_name"`
	Prompt            string     `gorm:"column:prompt;type:text" json:"prompt"`
	RequestBody       []byte     `gorm:"column:request_body;type:json" json:"request_body"`
	RequestModel      string     `gorm:"column:request_model;type:varchar(100)" json:"request_model"`
	ResponseBody      string     `gorm:"column:response_body;type:longtext" json:"response_body"`
	ResponseStatus    int        `gorm:"column:response_status;type:int" json:"response_status"`
	ErrorMessage      string     `gorm:"column:error_message;type:text" json:"error_message"`
	PromptTokens      int        `gorm:"column:prompt_tokens;type:int;default:0" json:"prompt_tokens"`
	CompletionTokens  int        `gorm:"column:completion_tokens;type:int;default:0" json:"completion_tokens"`
	TotalTokens       int        `gorm:"column:total_tokens;type:int;default:0" json:"total_tokens"`
	EstimatedCostUSD  float64    `gorm:"column:estimated_cost_usd;type:decimal(12,6);default:0" json:"estimated_cost_usd"`
	GatewayLatencyMS  int        `gorm:"column:gateway_latency_ms;type:int" json:"gateway_latency_ms"`
	TotalLatencyMS    int        `gorm:"column:total_latency_ms;type:int" json:"total_latency_ms"`
	UpstreamLatencyMS int        `gorm:"column:upstream_latency_ms;type:int" json:"upstream_latency_ms"`
	RequestedAt       time.Time  `gorm:"column:requested_at;type:timestamp(6);default:CURRENT_TIMESTAMP(6);index:idx_time_range" json:"requested_at"`
	RespondedAt       *time.Time `gorm:"column:responded_at;type:timestamp(6)" json:"responded_at"`
}

// TableName returns the table name
func (RequestLog) TableName() string {
	return "request_logs"
}

// MonthlyCost 月度成本聚合表
type MonthlyCost struct {
	ID                    string    `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())" json:"id"`
	DepartmentID          string    `gorm:"column:department_id;type:varchar(36);not null" json:"department_id"`
	YearMonth             string    `gorm:"column:year_month;type:char(7);not null" json:"year_month"`
	ModelID               string    `gorm:"column:model_id;type:varchar(50);not null" json:"model_id"`
	TotalRequests         int       `gorm:"column:total_requests;type:int;default:0" json:"total_requests"`
	TotalPromptTokens     int64     `gorm:"column:total_prompt_tokens;type:bigint;default:0" json:"total_prompt_tokens"`
	TotalCompletionTokens int64     `gorm:"column:total_completion_tokens;type:bigint;default:0" json:"total_completion_tokens"`
	TotalTokens           int64     `gorm:"column:total_tokens;type:bigint;default:0" json:"total_tokens"`
	EstimatedCostUSD      float64   `gorm:"column:estimated_cost_usd;type:decimal(14,4);default:0" json:"estimated_cost_usd"`
	CreatedAt             time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`

	Department Department `gorm:"foreignKey:DepartmentID;references:ID" json:"department,omitempty"`
	Model      AiModel    `gorm:"foreignKey:ModelID;references:ID" json:"model,omitempty"`
}

// TableName returns the table name
func (MonthlyCost) TableName() string {
	return "monthly_costs"
}
