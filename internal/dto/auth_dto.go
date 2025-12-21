package dto

import "time"

// RegisterAccountRequest 账号密码注册请求
type RegisterAccountRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Email           string `json:"email" validate:"required,email"`
	Mobile          string `json:"mobile" validate:"omitempty,len=11"`
	Password        string `json:"password" validate:"required,min=8,max=32"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Nickname        string `json:"nickname"`
	Avatar          string `json:"avatar"`
}

// RegisterResponse 注册成功响应
type RegisterResponse struct {
	User   UserDTO  `json:"user"`
	Tokens TokenDTO `json:"tokens"`
}

// LoginAccountRequest 账号密码登录请求
type LoginAccountRequest struct {
	Account    string `json:"account" validate:"required"` // Email, Username, or Mobile
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me"`
}

// LoginResponse 登录成功响应
type LoginResponse struct {
	User   UserDTO  `json:"user"`
	Tokens TokenDTO `json:"tokens"`
}

// UserDTO 用户信息DTO
type UserDTO struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	IsActive  bool      `json:"is_active"`
	LastLogin time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenDTO 令牌DTO
type TokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
