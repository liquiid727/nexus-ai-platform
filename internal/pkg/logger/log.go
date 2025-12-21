package logger

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 定义了项目的日志接口
type Logger interface {
	// Debug 记录调试级别的日志 用于开发和调试阶段，包含详细的调试信息
	Debugw(msg string, kvs ...any)
	// Info 记录信息级别的日志 表示系统的正常运行状态
	Infow(msg string, kvs ...any)
	// Warn 记录警告级别的日志 表示系统可能出现问题
	Warnw(msg string, kvs ...any)
	// Error 记录错误级别的日志 表示系统发生了不可恢复的错误
	Errorw(msg string, kvs ...any)
	// Panic 记录紧急级别的日志 表示系统发生了不可恢复的错误，会触发程序崩溃
	Panicw(msg string, kvs ...any)
	// Fatal 记录致命级别的日志 表示系统发生了不可恢复的错误，会触发程序崩溃
	Fatalw(msg string, kvs ...any)
	// Sync 刷新日志缓冲区，确保所有日志都被写入文件
	Sync()
}

// zapLogger 实现了 Logger 接口，使用 zap 库作为底层日志记录器
type zapLogger struct {
	z *zap.Logger
}

// 确保 zapLogger 实现了 Logger 接口
var _ Logger = (*zapLogger)(nil)

var (
	mu sync.Mutex
	// std 定义默认全局的Logger
	std = New(NewOptions())
)

// Init 初始化默认全局的Logger
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}

func New(opts *Options) *zapLogger {
	// 如果opts为nil 则使用默认配置
	if opts == nil {
		opts = NewOptions()
	}
	// 解析日志级别
	level, err := zapcore.ParseLevel(opts.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// 创建encoder 配置，用于控制日志的输出格式
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.MessageKey = "msg"
	encoderCfg.TimeKey = "timestamp"

	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// 指定 time.Duration 序列化函数，将 time.Duration 序列化为经过的毫秒数的浮点数
	// 毫秒数比默认的秒数更精确
	encoderCfg.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}

	if opts.Format == "" {
		opts.Format = "console"
	}
	if len(opts.OutputPaths) == 0 {
		opts.OutputPaths = []string{"stdout"}
	}
	// 创建构建 zap.Logger 需要的配置
	cfg := &zap.Config{
		// 是否在日志中显示调用日志所在的文件和行号，例如：`"caller":"miniblog/miniblog.go:75"`
		DisableCaller: opts.DisableCaller,
		// 是否禁止在 panic 及以上级别打印堆栈信息
		DisableStacktrace: opts.DisableStackTrace,
		// 指定日志级别
		Level: zap.NewAtomicLevelAt(level),
		// 指定日志显示格式，可选值：console, json
		Encoding:      opts.Format,
		EncoderConfig: encoderCfg,
		// 指定日志输出位置
		OutputPaths: opts.OutputPaths,
		// 设置 zap 内部错误输出位置
		ErrorOutputPaths: []string{"stderr"},
	}
	// 创建 zap 日志记录器
	// 跳过2层调用
	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}
	// 重定向标准库的日志输出到 zap 日志记录器
	zap.RedirectStdLog(z)
	return &zapLogger{z: z}
}

// Sync 刷新日志缓冲区，确保所有日志都被写入文件
// 注意：主程序在退出前需要调用Sync
func Sync() {
	std.Sync()
}
func (l *zapLogger) Sync() {
	// 忽略返回值
	_ = l.z.Sync()
}

// Debugw implements Logger.
func (l *zapLogger) Debugw(msg string, kvs ...any) {
	l.z.Sugar().Debugw(msg, kvs...)
}
func Debugw(msg string, kvs ...any) {
	std.Debugw(msg, kvs...)
}

func Infow(msg string, kvs ...any) {
	std.Infow(msg, kvs...)
}
func (l *zapLogger) Infow(msg string, kvs ...any) {
	l.z.Sugar().Infow(msg, kvs...)
}

func Warnw(msg string, kvs ...any) {
	std.Warnw(msg, kvs...)
}
func (l *zapLogger) Warnw(msg string, kvs ...any) {
	l.z.Sugar().Warnw(msg, kvs...)
}

func Errorw(msg string, kvs ...any) {
	std.Errorw(msg, kvs...)
}
func (l *zapLogger) Errorw(msg string, kvs ...any) {
	l.z.Sugar().Errorw(msg, kvs...)
}
func (l *zapLogger) Panicw(msg string, kvs ...any) {
	l.z.Sugar().Panicw(msg, kvs...)
}
func Panicw(msg string, kvs ...any) {
	std.Panicw(msg, kvs...)
}
func (l *zapLogger) Fatalw(msg string, kvs ...any) {
	l.z.Sugar().Fatalw(msg, kvs...)
}
func Fatalw(msg string, kvs ...any) {
	std.Fatalw(msg, kvs...)
}
