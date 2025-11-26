package http

import (
	"ShareInfo/internal/config"
	middleware2 "ShareInfo/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type RouterConfig struct {
	RedisClient *redis.Client
	Logger      *zap.Logger
	Config      *config.Config
}

func SetupRouter(handler *Handler, cfg RouterConfig) *gin.Engine {
	if cfg.Config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware2.Logger(cfg.Logger))
	router.Use(middleware2.Recovery(cfg.Logger))
	rateLimiterConfig := middleware2.RateLimiterConfig{
		RequestsPerTTL: cfg.Config.RateLimit.RequestsPerMinute,
		TTL:            time.Second * time.Duration(cfg.Config.RateLimit.BurstSize),
	}
	router.Use(middleware2.RateLimiter(cfg.RedisClient, rateLimiterConfig, cfg.Logger))

	// router.GET("/health", handler.HealthCheck)

	// v1 := router.Group("/api/v1")
	{

	}

	return router
}
