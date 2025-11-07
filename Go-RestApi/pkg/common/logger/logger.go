package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

type LogConfig struct {
	Mode         string
	Level        string
	EnableFile   bool
	FilePath     string
	MaxSize      int
	MaxBackups   int
	MaxAge       int
	Compress     bool
	EnableStdout bool
}

func DefaultConfig() *LogConfig {
	return &LogConfig{
		Mode:         "development",
		Level:        "debug",
		EnableFile:   false,
		FilePath:     "./logs/app.log",
		MaxSize:      100,
		MaxBackups:   3,
		MaxAge:       7,
		Compress:     true,
		EnableStdout: true,
	}
}

func ProductionConfig() *LogConfig {
	return &LogConfig{
		Mode:         "production",
		Level:        "info",
		EnableFile:   true,
		FilePath:     "./logs/app.log",
		MaxSize:      100,
		MaxBackups:   10,
		MaxAge:       14,
		Compress:     true,
		EnableStdout: false,
	}
}

func InitLogger(mode string) error {
	var config *LogConfig

	if mode == "production" {
		config = ProductionConfig()
	} else {
		config = DefaultConfig()
		config.Mode = mode
	}

	return InitLoggerWithConfig(config)
}

func InitLoggerWithFile(mode string, logFilePath string) error {
	config := DefaultConfig()
	config.Mode = mode
	config.EnableFile = true
	config.FilePath = logFilePath
	config.EnableStdout = true

	if mode == "production" {
		config = ProductionConfig()
		config.FilePath = logFilePath
	}

	return InitLoggerWithConfig(config)
}

func InitLoggerWithConfig(config *LogConfig) error {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.CallerKey = "caller"
	encoderConfig.StacktraceKey = "stacktrace"

	if config.Mode == "development" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}

	var cores []zapcore.Core

	if config.EnableStdout {
		var consoleEncoder zapcore.Encoder
		if config.Mode == "production" {
			consoleEncoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			consoleEncoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		stdoutCore := zapcore.NewCore(
			consoleEncoder, zapcore.AddSync(os.Stdout), level,
		)

		cores = append(cores, stdoutCore)
	}

	if config.EnableFile {
		if err := os.MkdirAll("./logs", 0755); err != nil {
			return err
		}

		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		})

		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileCore := zapcore.NewCore(
			fileEncoder, fileWriter, level,
		)

		cores = append(cores, fileCore)
	}

	core := zapcore.NewTee(cores...)

	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

func Sync() {
	_ = Logger.Sync()
}

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	if err := InitLogger(env); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}
