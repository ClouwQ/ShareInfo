package http

import (
	"ShareInfo/internal/config"
	middleware2 "ShareInfo/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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

	router.GET("/health", handler.HealthCheck)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		v1.POST("/links", handler.CreateLink)
		v1.GET("/links/:id", handler.GetLink)
		v1.POST("links/freeze")

		v1.POST("files/upload", handler.UploadFile)
		v1.DELETE("files/delete", handler.DeleteFile)
		v1.GET("files/:id", handler.GetZipFiles)
	}

	return router
}
