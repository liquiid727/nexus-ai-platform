package manager

import (
	"next-ai-gateway/internal/repository/entity"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type GatewayTokenManager struct {
	signingKey []byte
	logger     *zap.Logger
}

func NewGatewayTokenManager(signingKey []byte, logger *zap.Logger) *GatewayTokenManager {
	return &GatewayTokenManager{
		signingKey: signingKey,
		logger:     logger,
	}
}
func (m *GatewayTokenManager) GenerateToken(userInfo entity.User) (string, error) {
	// 生成token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userInfo.ID,
		"username": userInfo.Username,
		"email":    userInfo.Email,
	}).SignedString(m.signingKey)
	if err != nil {
		return "", err
	}
	//
	return token, nil
}

func (m *GatewayTokenManager) ValidateAndParse(token string) (string, error) {

	return token, nil
}
