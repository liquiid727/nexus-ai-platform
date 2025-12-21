package config

import (
	"fmt"
	"next-ai-gateway/pkg/database"
	"next-ai-gateway/internal/pkg/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  database.Config `mapstructure:"database"`
	Logger    logger.Options  `mapstructure:"logger"`
	JWTSecret string          `mapstructure:"jwt_secret"`
	JWTExpiry int64           `mapstructure:"jwt_expiry"`
	RedisURL  string          `mapstructure:"redis_url"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

var GlobalConfig Config

func LoadConfig(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	// Map environment variables
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.username", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.database", "DB_NAME")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("jwt_secret", "JWT_SECRET")
	viper.BindEnv("jwt_expiry", "JWT_EXPIRY")
	viper.BindEnv("redis_url", "REDIS_URL")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		// Config file not found is okay if we have env vars,
		// but typically we want at least defaults.
		// For now we allow missing file if envs are set.
	}

	// Try to read .env file if it exists, for local dev convenience
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.MergeInConfig() // Merge .env into existing config

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if err := viper.Unmarshal(&GlobalConfig); err != nil {
			fmt.Println("Error unmarshaling config:", err)
		}
	})

	return nil
}
