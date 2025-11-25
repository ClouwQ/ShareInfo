package main

import (
	"ShareInfo/internal/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "go.uber.org/zap/zapcore"
)

func main() {
	// Загружаем config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	logger := initLogger(cfg.Logger)
	defer logger.Sync()

	logger.Info("launching server",
		zap.String("environment", cfg.Environment),
		zap.Int("port", cfg.Server.Port))
}

func initLogger(cfg config.LoggerConfig) *zap.Logger {
	var level zapcore.Level
	switch cfg.Level {
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

	var zapConfig zap.Config
	if cfg.Format == "json" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapConfig.Level = zap.NewAtomicLevelAt(level)
	zapConfig.OutputPaths = []string{cfg.OutputPath}

	logger, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	return logger
}
