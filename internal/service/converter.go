package service

import (
	"time"

	"next-ai-gateway/internal/dto"
	"next-ai-gateway/internal/repository/entity"
)

// ToUserDTO converts User entity to UserDTO
func ToUserDTO(user *entity.User) dto.UserDTO {
	lastLogin := time.Time{}
	if user.Security != nil && user.Security.LastLoginAt != nil {
		lastLogin = *user.Security.LastLoginAt
	}

	var username, email, mobile, nickname, avatar string
	if user.Username != nil {
		username = *user.Username
	}
	if user.Email != nil {
		email = *user.Email
	}
	if user.Mobile != nil {
		mobile = *user.Mobile
	}
	if user.Profile != nil {
		if user.Profile.Nickname != nil {
			nickname = *user.Profile.Nickname
		}
		if user.Profile.AvatarURL != nil {
			avatar = *user.Profile.AvatarURL
		}
	}

	return dto.UserDTO{
		ID:        user.ID,
		Username:  username,
		Email:     email,
		Mobile:    mobile,
		Nickname:  nickname,
		Avatar:    avatar,
		IsActive:  user.IsActive,
		LastLogin: lastLogin,
		CreatedAt: user.CreatedAt,
	}
}
