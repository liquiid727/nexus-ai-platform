package service

import (
	"context"
	"next-ai-gateway/internal/dao"
	"next-ai-gateway/internal/repository/entity"
)

type UserService struct {
	dao *dao.UserDAO
}

// NewUserService 创建用户服务
func NewUserService(dao *dao.UserDAO) *UserService {
	return &UserService{dao: dao}
}
func (s *UserService) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return s.dao.GetByID(ctx, id)
}
