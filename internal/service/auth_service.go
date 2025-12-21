package service

import (
	"context"
	"time"

	"next-ai-gateway/internal/config"
	"next-ai-gateway/internal/dao"
	"next-ai-gateway/internal/dto"
	"next-ai-gateway/internal/pkg/errors"
	"next-ai-gateway/internal/repository/entity"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userDAO *dao.UserDAO
}

func NewAuthService(userDAO *dao.UserDAO) *AuthService {
	return &AuthService{userDAO: userDAO}
}

// RegisterAccount 账号密码注册
func (s *AuthService) RegisterAccount(ctx context.Context, req *dto.RegisterAccountRequest) (*dto.RegisterResponse, error) {
	// 1. Check if user exists (username or email)
	// Check username
	_, err := s.userDAO.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, errors.New(409, "User already exists", "用户名已存在").KV("field", "username").KV("value", req.Username)
	} else if err != gorm.ErrRecordNotFound {
		return nil, errors.New(500, "Internal Server Error", "数据库查询失败").KV("error", err.Error())
	}

	// Check email
	_, err = s.userDAO.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New(409, "User already exists", "邮箱已存在").KV("field", "email").KV("value", req.Email)
	} else if err != gorm.ErrRecordNotFound {
		return nil, errors.New(500, "Internal Server Error", "数据库查询失败").KV("error", err.Error())
	}

	// Check mobile if provided
	if req.Mobile != "" {
		_, err = s.userDAO.GetByMobile(ctx, req.Mobile)
		if err == nil {
			return nil, errors.New(409, "User already exists", "手机号已存在").KV("field", "mobile").KV("value", req.Mobile)
		} else if err != gorm.ErrRecordNotFound {
			return nil, errors.New(500, "Internal Server Error", "数据库查询失败").KV("error", err.Error())
		}
	}

	// 2. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(500, "Internal Server Error", "密码加密失败").KV("error", err.Error())
	}

	// 3. Create user
	user := &entity.User{
		Username:     &req.Username,
		Email:        &req.Email,
		Mobile:       &req.Mobile,
		IsActive:     true,
		DepartmentID: "dept_default", // TODO: Determine correct department
		Password: &entity.UserPassword{
			PasswordHash: string(hashedPassword),
			PasswordSalt: "", // bcrypt handles salt
		},
		Profile: &entity.UserProfile{
			Nickname:  &req.Nickname,
			AvatarURL: &req.Avatar,
		},
		Security: &entity.UserSecurity{},
	}

	if err := s.userDAO.Create(ctx, user); err != nil {
		return nil, errors.New(500, "Internal Server Error", "创建用户失败").KV("error", err.Error())
	}

	// 4. Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, errors.New(500, "Internal Server Error", "生成令牌失败").KV("error", err.Error())
	}

	return &dto.RegisterResponse{
		User:   ToUserDTO(user),
		Tokens: tokens,
	}, nil
}

// LoginAccount 账号密码登录
func (s *AuthService) LoginAccount(ctx context.Context, req *dto.LoginAccountRequest) (*dto.LoginResponse, error) {
	// 1. Find user by Account (Email, Username, or Mobile)
	var user *entity.User
	var err error

	// Try by username
	user, err = s.userDAO.GetByUsername(ctx, req.Account)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New(500, "Internal Server Error", "查询用户失败").KV("error", err.Error())
	}

	// Try by email
	if user == nil {
		user, err = s.userDAO.GetByEmail(ctx, req.Account)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, errors.New(500, "Internal Server Error", "查询用户失败").KV("error", err.Error())
		}
	}

	// Try by mobile
	if user == nil {
		user, err = s.userDAO.GetByMobile(ctx, req.Account)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, errors.New(500, "Internal Server Error", "查询用户失败").KV("error", err.Error())
		}
	}

	if user == nil {
		return nil, errors.New(401, "Unauthorized", "账号或密码错误")
	}

	// 2. Verify password
	if user.Password == nil {
		return nil, errors.New(500, "Internal Server Error", "用户密码数据缺失")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New(401, "Unauthorized", "账号或密码错误")
	}

	// 3. Check if active
	if !user.IsActive {
		return nil, errors.New(403, "Forbidden", "账号已被禁用")
	}

	// 4. Update last login
	now := time.Now()
	if user.Security == nil {
		user.Security = &entity.UserSecurity{UserID: user.ID}
	}
	user.Security.LastLoginAt = &now
	if err := s.userDAO.Update(ctx, user); err != nil {
		// Log error but proceed? Or fail? Best to log and proceed.
		// For now, we ignore update error or log it if we had a logger here.
	}

	// 5. Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, errors.New(500, "Internal Server Error", "生成令牌失败").KV("error", err.Error())
	}

	// Adjust expiry based on RememberMe if needed (simplified here)

	return &dto.LoginResponse{
		User:   ToUserDTO(user),
		Tokens: tokens,
	}, nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.TokenDTO, error) {
	// Parse refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New(401, "Unauthorized", "无效的刷新令牌")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New(401, "Unauthorized", "无效的刷新令牌")
	}

	// Check type
	if typeVal, ok := claims["type"]; !ok || typeVal != "refresh" {
		return nil, errors.New(401, "Unauthorized", "令牌类型错误")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New(401, "Unauthorized", "无效的用户ID")
	}

	// Check if user exists/active
	user, err := s.userDAO.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New(401, "Unauthorized", "用户不存在")
	}
	if !user.IsActive {
		return nil, errors.New(403, "Forbidden", "账号已被禁用")
	}

	// Generate new tokens
	// Note: generateTokens returns dto.TokenDTO and error.
	// We need to return *dto.TokenDTO
	tokens, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}
	return &tokens, nil
}

// GetProfile 获取用户信息
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*dto.UserDTO, error) {
	user, err := s.userDAO.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(404, "Not Found", "用户不存在")
		}
		return nil, errors.New(500, "Internal Server Error", "数据库查询失败").KV("error", err.Error())
	}

	userDTO := ToUserDTO(user)
	return &userDTO, nil
}

func (s *AuthService) generateTokens(user *entity.User) (dto.TokenDTO, error) {
	accessExpiry := time.Hour * 2
	refreshExpiry := time.Hour * 24 * 30

	// Access Token
	accessTokenClaims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(accessExpiry).Unix(),
		"type":     "access",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessStr, err := accessToken.SignedString([]byte(config.GlobalConfig.JWTSecret))
	if err != nil {
		return dto.TokenDTO{}, err
	}

	// Refresh Token
	refreshTokenClaims := jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(refreshExpiry).Unix(),
		"type": "refresh",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshStr, err := refreshToken.SignedString([]byte(config.GlobalConfig.JWTSecret))
	if err != nil {
		return dto.TokenDTO{}, err
	}

	return dto.TokenDTO{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessExpiry.Seconds()),
	}, nil
}
