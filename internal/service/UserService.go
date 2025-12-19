package service

import (
	"context"
	"next-ai-gateway/internal/repository"
)

type UserService struct {
	dao repository.UserDAO
}

// NewUserService 创建用户服务
func NewUserService(dao repository.UserDAO) *UserService {
	return &UserService{dao: dao}
}

// GetUserByID 获取用户信息
func (s *UserService) GetUserByID(ctx context.Context, id string) (*repository.User, error) {
	return nil, nil
}
