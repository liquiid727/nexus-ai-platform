package logger

import "go.uber.org/zap/zapcore"

type Options struct {
	Level             string   `mapstructure:"level"`               // 日志级别 (debug, info, warn, error)
	MaxSize           int      `mapstructure:"max_size"`            // 单个日志文件的最大大小，单位：MB
	MaxAge            int      `mapstructure:"max_age"`             // 旧日志文件保留的最大天数
	MaxBackups        int      `mapstructure:"max_backups"`         // 保留旧日志文件的最大数量
	Compress          bool     `mapstructure:"compress"`            // 是否压缩旧日志文件
	Console           bool     `mapstructure:"console"`             // 是否同时输出到控制台 (stdout)
	DisableCaller     bool     `mapstructure:"disable_caller"`      // 是否禁用调用者信息 (文件名、行号等)
	DisableStackTrace bool     `mapstructure:"disable_stack_trace"` // 是否禁用栈跟踪信息
	Format            string   `mapstructure:"format"`              // 日志格式 (console, json)
	OutputPaths       []string `mapstructure:"output_paths"`        // 日志输出路径 (stdout, stderr, 文件路径)

	// 其他自定义配置项...
}

func NewOptions() *Options {
	return &Options{
		Level:             zapcore.InfoLevel.String(),
		MaxSize:           100,
		MaxAge:            7,
		MaxBackups:        5,
		Compress:          true,
		Console:           true,
		DisableCaller:     false,
		DisableStackTrace: false,
		Format:            "console",
		OutputPaths:       []string{"stdout"},
	}
}
