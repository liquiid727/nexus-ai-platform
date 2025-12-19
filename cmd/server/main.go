package main

import (
	"flag"
	"fmt"
	"next-ai-gateway/internal/app"
	"next-ai-gateway/internal/config"
	"next-ai-gateway/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/config.yaml", "path to config file")
	flag.Parse()

	e, err := app.Init(configPath)
	if err != nil {
		panic(err)
	}

	port := config.GlobalConfig.Server.Port
	if port == 0 {
		port = 8080
	}

	logger.Infow("Starting server", zap.Int("port", port))
	if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
		logger.Fatalw("Server start failed", zap.Error(err))
	}
}
