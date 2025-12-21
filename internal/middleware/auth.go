package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"next-ai-gateway/internal/config"
	"next-ai-gateway/internal/pkg/errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// JWTAuth 简单的JWT认证中间件
func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, errors.New(401, "Unauthorized", "未提供认证令牌"))
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, errors.New(401, "Unauthorized", "令牌格式错误"))
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.GlobalConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, errors.New(401, "Unauthorized", "无效的令牌").KV("error", err.Error()))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, errors.New(401, "Unauthorized", "无效的令牌声明"))
		}

		// Check if token is access token
		if typeVal, ok := claims["type"]; ok && typeVal != "access" {
			return c.JSON(http.StatusUnauthorized, errors.New(401, "Unauthorized", "令牌类型错误"))
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, errors.New(401, "Unauthorized", "无效的用户ID"))
		}

		c.Set("user_id", userID)

		return next(c)
	}
}
