package entity

import (
	"time"

	"gorm.io/gorm"
)

// Department 部门表
type Department struct {
	ID                 string    `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())"`
	Name               string    `gorm:"column:name;type:varchar(100);not null"`
	ParentID           *string   `gorm:"column:parent_id;type:varchar(36);index"`
	CostCenterCode     string    `gorm:"column:cost_center_code;type:varchar(50);index"`
	MonthlyTokenBudget int64     `gorm:"column:monthly_token_budget;type:bigint"`
	IsActive           bool      `gorm:"column:is_active;type:boolean;default:true"`
	CreatedAt          time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (Department) TableName() string {
	return "departments"
}

// User 用户表
type User struct {
	ID           string         `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())"`
	InternalID   uint64         `gorm:"column:internal_id;type:bigint unsigned;autoIncrement;unique"`
	Username     *string        `gorm:"column:username;type:varchar(50);unique"`
	Email        *string        `gorm:"column:email;type:varchar(255);unique;index"`
	Mobile       *string        `gorm:"column:mobile;type:varchar(20);unique;index"`
	Status       int            `gorm:"column:status;type:tinyint;default:1;index:idx_status_created,priority:1"` // 1:正常 2:禁用 3:注销
	Type         int            `gorm:"column:type;type:tinyint;default:1"`                                       // 1:普通用户 2:VIP 3:管理员
	EmployeeID   string         `gorm:"column:employee_id;type:varchar(50);unique;not null"`
	Name         string         `gorm:"column:name;type:varchar(100);not null"`
	DepartmentID string         `gorm:"column:department_id;type:varchar(36);not null;index"`
	Role         string         `gorm:"column:role;type:enum('admin','user','billing');default:'user'"`
	IsActive     bool           `gorm:"column:is_active;type:boolean;default:true"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;index:idx_status_created,priority:2"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime;index"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index"`

	// Relations
	Department *Department     `gorm:"foreignKey:DepartmentID;references:ID"`
	Password   *UserPassword   `gorm:"foreignKey:UserID;references:ID"`
	Profile    *UserProfile    `gorm:"foreignKey:UserID;references:ID"`
	Security   *UserSecurity   `gorm:"foreignKey:UserID;references:ID"`
	Settings   *UserSettings   `gorm:"foreignKey:UserID;references:ID"`
	Statistics *UserStatistics `gorm:"foreignKey:UserID;references:ID"`
}

func (User) TableName() string {
	return "users"
}

// UserPassword 用户密码表
type UserPassword struct {
	UserID         string     `gorm:"column:user_id;type:varchar(32);primaryKey"`
	PasswordHash   string     `gorm:"column:password_hash;type:varchar(255);not null"`
	PasswordSalt   string     `gorm:"column:password_salt;type:varchar(64);not null"`
	Algorithm      string     `gorm:"column:algorithm;type:varchar(20);default:'bcrypt'"`
	LastChangedAt  time.Time  `gorm:"column:last_changed_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	ChangeRequired bool       `gorm:"column:change_required;type:boolean;default:false"`
	FailedAttempts int        `gorm:"column:failed_attempts;type:int;default:0"`
	LockedUntil    *time.Time `gorm:"column:locked_until;type:timestamp;index"`
	PreviousHash   *string    `gorm:"column:previous_hash;type:varchar(255)"`
}

func (UserPassword) TableName() string {
	return "user_passwords"
}

// UserProfile 用户资料表
type UserProfile struct {
	UserID        string     `gorm:"column:user_id;type:varchar(32);primaryKey"`
	Nickname      *string    `gorm:"column:nickname;type:varchar(100);index"`
	AvatarURL     *string    `gorm:"column:avatar_url;type:varchar(500)"`
	Gender        int        `gorm:"column:gender;type:tinyint;default:0"` // 0:未知 1:男 2:女
	Birthday      *time.Time `gorm:"column:birthday;type:date"`
	Bio           *string    `gorm:"column:bio;type:varchar(500)"`
	Country       *string    `gorm:"column:country;type:varchar(50)"`
	Province      *string    `gorm:"column:province;type:varchar(50)"`
	City          *string    `gorm:"column:city;type:varchar(50)"`
	Timezone      string     `gorm:"column:timezone;type:varchar(50);default:'Asia/Shanghai'"`
	WechatUnionID *string    `gorm:"column:wechat_unionid;type:varchar(100);unique;index:idx_wechat,priority:1"`
	WechatOpenID  *string    `gorm:"column:wechat_openid;type:varchar(100);index:idx_wechat,priority:2"`
	GithubID      *string    `gorm:"column:github_id;type:varchar(100)"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

// UserSettings 用户设置表
type UserSettings struct {
	UserID         string    `gorm:"column:user_id;type:varchar(32);primaryKey"`
	NotifyEmail    []byte    `gorm:"column:notify_email;type:json"`
	NotifyPush     []byte    `gorm:"column:notify_push;type:json"`
	NotifySMS      []byte    `gorm:"column:notify_sms;type:json"`
	PrivacyProfile int       `gorm:"column:privacy_profile;type:tinyint;default:1"` // 1:公开 2:好友 3:私密
	PrivacyContact int       `gorm:"column:privacy_contact;type:tinyint;default:3"`
	UITheme        string    `gorm:"column:ui_theme;type:varchar(20);default:'light'"`
	UILanguage     string    `gorm:"column:ui_language;type:varchar(10);default:'zh-CN'"`
	CustomSettings []byte    `gorm:"column:custom_settings;type:json"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (UserSettings) TableName() string {
	return "user_settings"
}

// UserSecurity 用户安全表
type UserSecurity struct {
	ID                    uint64     `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement"`
	UserID                string     `gorm:"column:user_id;type:varchar(32);not null;index"`
	EmailVerified         bool       `gorm:"column:email_verified;type:boolean;default:false;index:idx_verified,priority:1"`
	MobileVerified        bool       `gorm:"column:mobile_verified;type:boolean;default:false;index:idx_verified,priority:2"`
	IDVerified            bool       `gorm:"column:id_verified;type:boolean;default:false"`
	EmailVerifiedAt       *time.Time `gorm:"column:email_verified_at;type:timestamp"`
	MobileVerifiedAt      *time.Time `gorm:"column:mobile_verified_at;type:timestamp"`
	MFAEnabled            bool       `gorm:"column:mfa_enabled;type:boolean;default:false"`
	MFAType               *string    `gorm:"column:mfa_type;type:varchar(20)"`
	MFASecret             *string    `gorm:"column:mfa_secret;type:varchar(100)"`
	LastLoginAt           *time.Time `gorm:"column:last_login_at;type:timestamp"`
	LastLoginIP           *string    `gorm:"column:last_login_ip;type:varchar(45)"`
	LastPasswordChangedAt *time.Time `gorm:"column:last_password_changed_at;type:timestamp"`
}

func (UserSecurity) TableName() string {
	return "user_security"
}

// UserStatistics 用户统计表
type UserStatistics struct {
	UserID              string     `gorm:"column:user_id;type:varchar(32);primaryKey"`
	LoginCount          int        `gorm:"column:login_count;type:int;default:0"`
	ContinuousLoginDays int        `gorm:"column:continuous_login_days;type:int;default:0"`
	LastActiveAt        *time.Time `gorm:"column:last_active_at;type:timestamp"`
	TotalOnlineSeconds  int64      `gorm:"column:total_online_seconds;type:bigint;default:0"`
	PostCount           int        `gorm:"column:post_count;type:int;default:0"`
	CommentCount        int        `gorm:"column:comment_count;type:int;default:0"`
	LikeCount           int        `gorm:"column:like_count;type:int;default:0"`
	FollowerCount       int        `gorm:"column:follower_count;type:int;default:0;index"`
	FollowingCount      int        `gorm:"column:following_count;type:int;default:0"`
	OrderCount          int        `gorm:"column:order_count;type:int;default:0"`
	TotalSpent          float64    `gorm:"column:total_spent;type:decimal(12,2);default:0.00;index"`
	UpdatedAt           time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (UserStatistics) TableName() string {
	return "user_statistics"
}

// UserSession 用户会话表
type UserSession struct {
	ID                    string     `gorm:"column:id;type:varchar(64);primaryKey"`
	UserID                string     `gorm:"column:user_id;type:varchar(32);not null;index:idx_user_status,priority:1"`
	DeviceID              *string    `gorm:"column:device_id;type:varchar(100)"`
	DeviceType            *string    `gorm:"column:device_type;type:varchar(20)"`
	UserAgent             *string    `gorm:"column:user_agent;type:varchar(500)"`
	RefreshTokenHash      *string    `gorm:"column:refresh_token_hash;type:varchar(255);unique;index"`
	AccessTokenExpiresAt  *time.Time `gorm:"column:access_token_expires_at;type:timestamp"`
	RefreshTokenExpiresAt *time.Time `gorm:"column:refresh_token_expires_at;type:timestamp;index"`
	Status                int        `gorm:"column:status;type:tinyint;default:1;index:idx_user_status,priority:2"` // 1:活跃 2:已退出 3:过期
	IPAddress             *string    `gorm:"column:ip_address;type:varchar(45)"`
	Location              *string    `gorm:"column:location;type:varchar(100)"`
	CreatedAt             time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	LastActivityAt        time.Time  `gorm:"column:last_activity_at;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}

// UserOAuthBinding 第三方登录绑定表
type UserOAuthBinding struct {
	ID             uint64     `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement"`
	UserID         string     `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_provider,priority:1"`
	Provider       string     `gorm:"column:provider;type:varchar(20);not null;index:idx_user_provider,priority:2;index:idx_provider_user,priority:1;uniqueIndex:uk_provider_user,priority:1"`
	ProviderUserID string     `gorm:"column:provider_user_id;type:varchar(255);not null;index:idx_provider_user,priority:2;uniqueIndex:uk_provider_user,priority:2"`
	UnionID        *string    `gorm:"column:union_id;type:varchar(255);index"`
	OpenID         *string    `gorm:"column:open_id;type:varchar(255)"`
	Nickname       *string    `gorm:"column:nickname;type:varchar(100)"`
	AvatarURL      *string    `gorm:"column:avatar_url;type:text"`
	RawUserInfo    []byte     `gorm:"column:raw_user_info;type:jsonb"`
	AccessToken    *string    `gorm:"column:access_token;type:text"`
	RefreshToken   *string    `gorm:"column:refresh_token;type:text"`
	TokenExpiresAt *time.Time `gorm:"column:token_expires_at;type:timestamp"`
	Scope          *string    `gorm:"column:scope;type:varchar(255)"`
	IsPrimary      bool       `gorm:"column:is_primary;type:boolean;default:false"`
	IsValid        bool       `gorm:"column:is_valid;type:boolean;default:true"`
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	LastUsedAt     *time.Time `gorm:"column:last_used_at;type:timestamp"`
}

func (UserOAuthBinding) TableName() string {
	return "user_oauth_bindings"
}

// AiModel AI模型路由配置表
type AiModel struct {
	ID                string    `gorm:"column:id;type:varchar(50);primaryKey"`
	Description       *string   `gorm:"column:description;type:varchar(255)"`
	DefaultProviderID *string   `gorm:"column:default_provider_id;type:varchar(36)"`
	RoutingStrategy   string    `gorm:"column:routing_strategy;type:enum('load_balance','priority','lowest_cost');default:'load_balance'"`
	IsActive          bool      `gorm:"column:is_active;type:boolean;default:true"`
	ConfigJSON        []byte    `gorm:"column:config_json;type:json"`
	CreatedAt         time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (AiModel) TableName() string {
	return "ai_models"
}

// APIProviderKey 供应商API密钥池
type APIProviderKey struct {
	ID                   string     `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())"`
	ModelID              string     `gorm:"column:model_id;type:varchar(50);not null;index:idx_model_enabled,priority:1"`
	ProviderName         string     `gorm:"column:provider_name;type:varchar(50);not null"`
	APIKeyEncrypted      string     `gorm:"column:api_key_encrypted;type:text;not null"`
	EndpointBaseURL      *string    `gorm:"column:endpoint_base_url;type:varchar(500)"`
	Priority             int        `gorm:"column:priority;type:int;default:1;index:idx_model_enabled,priority:3"`
	MaxRPM               *int       `gorm:"column:max_rpm;type:int"`
	MaxTPM               *int       `gorm:"column:max_tpm;type:int"`
	CostPerMillionInput  float64    `gorm:"column:cost_per_million_input;type:decimal(12,6)"`
	CostPerMillionOutput float64    `gorm:"column:cost_per_million_output;type:decimal(12,6)"`
	IsEnabled            bool       `gorm:"column:is_enabled;type:boolean;default:true;index:idx_model_enabled,priority:2"`
	HealthStatus         string     `gorm:"column:health_status;type:enum('healthy','slow','down');default:'healthy';index"`
	LastUsedAt           *time.Time `gorm:"column:last_used_at;type:timestamp"`
	FailureCount         int        `gorm:"column:failure_count;type:int;default:0"`
	CreatedAt            time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (APIProviderKey) TableName() string {
	return "api_provider_keys"
}

// TokenQuota 多维度Token配额表
type TokenQuota struct {
	ID           string     `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())"`
	DepartmentID string     `gorm:"column:department_id;type:varchar(36);not null;uniqueIndex:uk_department_model_type,priority:1"`
	ModelID      *string    `gorm:"column:model_id;type:varchar(50);uniqueIndex:uk_department_model_type,priority:2"`
	QuotaType    string     `gorm:"column:quota_type;type:enum('total','daily','monthly');default:'monthly';uniqueIndex:uk_department_model_type,priority:3"`
	MaxTokens    int64      `gorm:"column:max_tokens;type:bigint;not null"`
	UsedTokens   int64      `gorm:"column:used_tokens;type:bigint;default:0"`
	ResetAt      *time.Time `gorm:"column:reset_at;type:timestamp;index"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (TokenQuota) TableName() string {
	return "token_quotas"
}

// MonthlyCost 月度成本聚合表
type MonthlyCost struct {
	ID                    string    `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())"`
	DepartmentID          string    `gorm:"column:department_id;type:varchar(36);not null;uniqueIndex:uk_department_month_model,priority:1"`
	YearMonth             string    `gorm:"column:year_month;type:char(7);not null;uniqueIndex:uk_department_month_model,priority:2"`
	ModelID               string    `gorm:"column:model_id;type:varchar(50);not null;uniqueIndex:uk_department_month_model,priority:3"`
	TotalRequests         int       `gorm:"column:total_requests;type:int;default:0"`
	TotalPromptTokens     int64     `gorm:"column:total_prompt_tokens;type:bigint;default:0"`
	TotalCompletionTokens int64     `gorm:"column:total_completion_tokens;type:bigint;default:0"`
	TotalTokens           int64     `gorm:"column:total_tokens;type:bigint;default:0"`
	EstimatedCostUSD      float64   `gorm:"column:estimated_cost_usd;type:decimal(14,4);default:0"`
	CreatedAt             time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt             time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (MonthlyCost) TableName() string {
	return "monthly_costs"
}

// RequestLog 请求审计日志表
type RequestLog struct {
	ID                string     `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())"`
	RequestID         *string    `gorm:"column:request_id;type:varchar(100);index"`
	UserID            string     `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_department,priority:1"`
	DepartmentID      string     `gorm:"column:department_id;type:varchar(36);not null;index:idx_user_department,priority:2"`
	ModelID           string     `gorm:"column:model_id;type:varchar(50);not null;index:idx_model_provider,priority:1"`
	ProviderKeyID     *string    `gorm:"column:provider_key_id;type:varchar(36)"`
	ProviderName      *string    `gorm:"column:provider_name;type:varchar(50);index:idx_model_provider,priority:2"`
	Prompt            *string    `gorm:"column:prompt;type:text"`
	RequestBody       []byte     `gorm:"column:request_body;type:json"`
	RequestModel      *string    `gorm:"column:request_model;type:varchar(100)"`
	ResponseBody      *string    `gorm:"column:response_body;type:longtext"`
	ResponseStatus    *int       `gorm:"column:response_status;type:int"`
	ErrorMessage      *string    `gorm:"column:error_message;type:text"`
	PromptTokens      int        `gorm:"column:prompt_tokens;type:int;default:0"`
	CompletionTokens  int        `gorm:"column:completion_tokens;type:int;default:0"`
	TotalTokens       int        `gorm:"column:total_tokens;type:int;default:0"`
	EstimatedCostUSD  float64    `gorm:"column:estimated_cost_usd;type:decimal(12,6);default:0"`
	GatewayLatencyMS  *int       `gorm:"column:gateway_latency_ms;type:int"`
	TotalLatencyMS    *int       `gorm:"column:total_latency_ms;type:int"`
	UpstreamLatencyMS *int       `gorm:"column:upstream_latency_ms;type:int"`
	RequestedAt       time.Time  `gorm:"column:requested_at;type:timestamp(6);default:CURRENT_TIMESTAMP(6);index"`
	RespondedAt       *time.Time `gorm:"column:responded_at;type:timestamp(6)"`
}

func (RequestLog) TableName() string {
	return "request_logs"
}

// SecurityLog 安全日志表
type SecurityLog struct {
	ID           uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement"`
	UserID       *string   `gorm:"column:user_id;type:varchar(36);index"`
	Action       string    `gorm:"column:action;type:varchar(50);not null;index:idx_action;index:idx_ip_action,priority:2"`
	ActionType   string    `gorm:"column:action_type;type:varchar(20);not null"`
	IPAddress    string    `gorm:"column:ip_address;type:varchar(45);not null;index:idx_ip_action,priority:1"`
	UserAgent    *string   `gorm:"column:user_agent;type:text"`
	Details      []byte    `gorm:"column:details;type:jsonb"`
	RiskLevel    string    `gorm:"column:risk_level;type:varchar(20);default:'low';index"`
	IsSuspicious bool      `gorm:"column:is_suspicious;type:boolean;default:false"`
	Country      *string   `gorm:"column:country;type:varchar(100)"`
	City         *string   `gorm:"column:city;type:varchar(100)"`
	Latitude     *float64  `gorm:"column:latitude;type:decimal(10,8)"`
	Longitude    *float64  `gorm:"column:longitude;type:decimal(11,8)"`
	DeviceType   *string   `gorm:"column:device_type;type:varchar(20)"`
	DeviceID     *string   `gorm:"column:device_id;type:varchar(100)"`
	SessionID    *string   `gorm:"column:session_id;type:varchar(64)"`
	RequestID    *string   `gorm:"column:request_id;type:varchar(100)"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;index"`
}

func (SecurityLog) TableName() string {
	return "security_logs"
}
