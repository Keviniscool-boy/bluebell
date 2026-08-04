package logger

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lg *zap.Logger

// Init 初始化日志
func Init() error {
	// 读取配置
	logLevel := viper.GetString("log.level")
	logMode := viper.GetString("log.mode") // dev: 控制台+文件 | prod: 仅文件
	logFilename := viper.GetString("log.filename")
	maxSize := viper.GetInt("log.max_size")
	maxAge := viper.GetInt("log.max_age")
	maxBackups := viper.GetInt("log.max_backups")

	// 解析日志级别
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 确保日志目录存在
	logDir := filepath.Dir(logFilename)
	if logDir != "." {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
	}

	// 配置文件输出（JSON 格式）
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		LocalTime:  true,
		Compress:   true,
	})
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(fileWriter),
		level,
	)

	// 根据 mode 决定是否同时输出到控制台
	var core zapcore.Core
	if logMode == "dev" || logMode == "debug" {
		// dev/debug 模式：文件 + 控制台（Console 格式，易读）
		consoleWriter := zapcore.AddSync(os.Stdout)
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(consoleWriter),
			level,
		)
		core = zapcore.NewTee(fileCore, consoleCore)
	} else {
		// prod/release 模式：仅输出到文件
		core = fileCore
	}

	lg = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// 替换全局 logger
	zap.ReplaceGlobals(lg)

	return nil
}

// Debug 输出 debug 级别日志
func Debug(msg string, fields ...zap.Field) {
	lg.Debug(msg, fields...)
}

// Info 输出 info 级别日志
func Info(msg string, fields ...zap.Field) {
	lg.Info(msg, fields...)
}

// Warn 输出 warn 级别日志
func Warn(msg string, fields ...zap.Field) {
	lg.Warn(msg, fields...)
}

// Error 输出 error 级别日志
func Error(msg string, fields ...zap.Field) {
	lg.Error(msg, fields...)
}

// Fatal 输出 fatal 级别日志并退出
func Fatal(msg string, fields ...zap.Field) {
	lg.Fatal(msg, fields...)
}

// With 创建带有额外字段的 logger
func With(fields ...zap.Field) *zap.Logger {
	return lg.With(fields...)
}

// GinLogger 返回 Gin 中间件使用的日志记录器
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		zap.L().Info("HTTP Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", time.Since(start)),
		)
	}
}

// GinRecovery 返回 Gin 中间件，用于捕获 panic 并记录日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				zap.L().Error("HTTP Request panic recovered",
					zap.Any("error", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
				)

				if stack {
					zap.L().Error("panic stacktrace", zap.Stack("stack"))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// Close 关闭 logger，刷新缓冲区
func Close() {
	if lg != nil {
		_ = lg.Sync()
	}
}
