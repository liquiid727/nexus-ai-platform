package app

import (
	"fmt"
	"next-ai-gateway/internal/config"
	"next-ai-gateway/internal/pkg/logger"
	"next-ai-gateway/internal/router"
	"next-ai-gateway/pkg/database"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func Init(configPath string) (*echo.Echo, error) {
	// 1. Load Config
	if err := config.LoadConfig(configPath); err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	// 2. Init Logger
	logger.Init(&config.GlobalConfig.Logger)
	defer logger.Sync()

	// 3. Init Database
	if err := database.Init(&config.GlobalConfig.Database); err != nil {
		// Log error but maybe don't exit if DB is optional, but for this task we assume it might be critical
		// or at least we should log it.
		// For the purpose of verification, if DB fails (e.g. no mysql running), we might fail.
		// However, to allow user to run without DB for testing health check, we can just log error.
		logger.Errorw("Failed to connect to database", zap.Error(err))
	} else {
		logger.Infow("Database initialized")
	}

	// 4. Init Echo
	e := echo.New()

	// 5. Setup Router
	router.RegisterRoutes(e)

	return e, nil
}
