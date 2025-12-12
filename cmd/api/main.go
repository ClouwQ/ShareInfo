package main

import (
	"ShareInfo/internal/config"
	httpDelivery "ShareInfo/internal/delivery/http"
	"ShareInfo/internal/infrastructure/cache"
	"ShareInfo/internal/infrastructure/database"
	postgresRepo "ShareInfo/internal/repository/database"
	fileUC "ShareInfo/internal/usecase/file"
	linkUC "ShareInfo/internal/usecase/link"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "ShareInfo/docs"
	_ "go.uber.org/zap/zapcore"
)

func main() {
	// Загружаем config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to create config: %v", err))
	}

	logger := initLogger(cfg.Logger)
	defer logger.Sync()

	logger.Info("launching server",
		zap.String("environment", cfg.Environment),
		zap.Int("port", cfg.Server.Port))

	// postgres
	db, err := database.NewPostgresConnection(&cfg.Database, logger)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	logger.Info("successfully connected to postgres")

	// redis
	var redisClient *redis.Client
	if cfg.Redis.Host != "" {
		redisClient, err = cache.NewRedisClient(&cfg.Redis, logger)
		if err != nil {
			logger.Warn("failed to connect to redis, rate limiting will use memory fallback", zap.Error(err))
		} else {
			defer cache.Close(context.Background(), redisClient, logger)
			logger.Info("successfully connected to redis")
		}
	} else {
		logger.Warn("redis host not configured, using memory cache")
	}

	// postgres репозитории
	linkRepo := postgresRepo.NewLinkRepository(db)
	linkMessageRepo := postgresRepo.NewLinkMessageRepository(db)
	fileRepo := postgresRepo.NewUploadedFileRepository(db)

	logger.Info("repositories initialized")

	// use cases
	linkUsecase := linkUC.NewService(logger, *linkRepo, *fileRepo)
	fileUsecase := fileUC.NewService(logger, *linkRepo, *fileRepo)

	logger.Info("use cases initialized")

	// handler
	handler := httpDelivery.NewHandler(logger, linkUsecase, fileUsecase, linkMessageRepo, *linkRepo, *fileRepo)

	// роутер
	routerConfig := httpDelivery.RouterConfig{
		Logger:      logger,
		RedisClient: redisClient,
		Config:      cfg,
	}

	router := httpDelivery.SetupRouter(handler, routerConfig)
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// launch
	go func() {
		logger.Info("server started",
			zap.String("address", srv.Addr),
			zap.String("environment", cfg.Environment),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server stopped gracefully")

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
