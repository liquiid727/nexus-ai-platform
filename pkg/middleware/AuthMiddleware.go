package middleware

import (
	"net/http"
	"next-ai-gateway/pkg/logger"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware 是一个中间件，用于验证请求中的认证信息

type AuthMiddleware struct {
	// client ssoClient
	// tokenMgr manager.GatewayTokenManager
	logger logger.Logger
}

func (m *AuthMiddleware) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 从Header中提取X-Gateway-Token
		authHeader := c.Request().Header.Get("X-Gateway-Token")
		// 使用tokenMgr验证token
		// claims,err:= m.tokenMgr.ValidateAndParse(token)
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
		}

		// 验证认证信息
		// ...

		// 继续处理请求
		return next(c)
	}
}
