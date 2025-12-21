package handler

import (
	"net/http"
	"time"

	"next-ai-gateway/internal/dto"
	"next-ai-gateway/internal/pkg/errors"
	"next-ai-gateway/internal/service"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// RegisterAccount 账号密码注册
func (h *AuthHandler) RegisterAccount(c echo.Context) error {
	req := new(dto.RegisterAccountRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(400, "Bad Request", "无效的请求参数").KV("error", err.Error()))
	}

	// Basic validation (manual for now as I don't see validator package setup)
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return c.JSON(http.StatusBadRequest, errors.New(400, "Bad Request", "必填参数缺失"))
	}
	if req.Password != req.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, errors.New(400, "Bad Request", "两次密码输入不一致"))
	}

	resp, err := h.service.RegisterAccount(c.Request().Context(), req)
	if err != nil {
		if errX, ok := err.(*errors.ErrorX); ok {
			return c.JSON(errX.Code, errX)
		}
		return c.JSON(http.StatusInternalServerError, errors.New(500, "Internal Server Error", err.Error()))
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"code":      201,
		"message":   "注册成功",
		"data":      resp,
		"timestamp": time.Now(),
	})
}

// LoginAccount 账号密码登录
func (h *AuthHandler) LoginAccount(c echo.Context) error {
	req := new(dto.LoginAccountRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(400, "Bad Request", "无效的请求参数").KV("error", err.Error()))
	}

	if req.Account == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, errors.New(400, "Bad Request", "账号或密码不能为空"))
	}

	resp, err := h.service.LoginAccount(c.Request().Context(), req)
	if err != nil {
		if errX, ok := err.(*errors.ErrorX); ok {
			return c.JSON(errX.Code, errX)
		}
		return c.JSON(http.StatusInternalServerError, errors.New(500, "Internal Server Error", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":      200,
		"message":   "登录成功",
		"data":      resp,
		"timestamp": time.Now(),
	})
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	req := new(dto.RefreshTokenRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New(400, "Bad Request", "无效的请求参数").KV("error", err.Error()))
	}

	if req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, errors.New(400, "Bad Request", "刷新令牌不能为空"))
	}

	resp, err := h.service.RefreshToken(c.Request().Context(), req)
	if err != nil {
		if errX, ok := err.(*errors.ErrorX); ok {
			return c.JSON(errX.Code, errX)
		}
		return c.JSON(http.StatusInternalServerError, errors.New(500, "Internal Server Error", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":      200,
		"message":   "令牌刷新成功",
		"data":      resp,
		"timestamp": time.Now(),
	})
}

// GetProfile 获取用户信息
func (h *AuthHandler) GetProfile(c echo.Context) error {
	// Need to get userID from context (set by middleware)
	// For now assume middleware sets "user_id"
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, errors.New(401, "Unauthorized", "未登录"))
	}

	user, err := h.service.GetProfile(c.Request().Context(), userID)
	if err != nil {
		if errX, ok := err.(*errors.ErrorX); ok {
			return c.JSON(errX.Code, errX)
		}
		return c.JSON(http.StatusInternalServerError, errors.New(500, "Internal Server Error", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":      200,
		"message":   "获取成功",
		"data":      map[string]interface{}{"user": user},
		"timestamp": time.Now(),
	})
}
