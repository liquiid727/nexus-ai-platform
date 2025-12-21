package router

import (
	"net/http"

	"next-ai-gateway/internal/dao"
	handler "next-ai-gateway/internal/handler"
	"next-ai-gateway/internal/middleware"
	"next-ai-gateway/internal/service"
	"next-ai-gateway/pkg/database"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	// Dependencies
	userDAO := dao.NewUserDAO(database.DB)
	authService := service.NewAuthService(userDAO)
	authHandler := handler.NewAuthHandler(authService)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	e.GET("/user/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello user")
	})

	// Auth Routes
	v1 := e.Group("/api/v1/auth")
	v1.POST("/register/account", authHandler.RegisterAccount)
	v1.POST("/login/account", authHandler.LoginAccount)
	v1.POST("/refresh", authHandler.RefreshToken)

	// Protected routes
	v1.GET("/profile", authHandler.GetProfile, middleware.JWTAuth)
}
