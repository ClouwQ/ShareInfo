package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RateLimiterConfig struct {
	RequestsPerTTL int // Количество запросов в TTL
	TTL            time.Duration
	KeyPrefix      string // Префикс для ключей в redis
}

// RateLimiter — ограничение частоты запросов
func RateLimiter(redisClient *redis.Client, config RateLimiterConfig, logger *zap.Logger) gin.HandlerFunc {
	if config.KeyPrefix == "" {
		config.KeyPrefix = "rate_limit"
	}

	return func(c *gin.Context) {
		clientID := getClientIdentifier(c)
		key := string(config.KeyPrefix) + ":" + clientID

		ctx := context.Background()

		// Увеличиваем счетчик
		count, err := redisClient.Incr(ctx, key).Result()
		// При ошибке redis пропускаем дальше
		if err != nil {
			logger.Error("Failed to increment rate limit counter (have access)",
				zap.Error(err),
				zap.String("key", key),
			)
			c.Next()
			return
		}

		// Устанавливаем TTL при первом запросе
		if count == 1 {
			err = redisClient.Expire(ctx, key, config.TTL).Err()
			if err != nil {
				logger.Warn("failed to set rate limit TTL", zap.Error(err))
			}
		}

		// Получаем TTL для заголовков
		ttl, _ := redisClient.TTL(ctx, key).Result()

		// Устанавливаем заголовки
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.RequestsPerTTL))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, config.RequestsPerTTL-int(count))))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(ttl).Unix()))

		// Проверяем лимит
		if count > int64(config.RequestsPerTTL) {
			c.Header("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":     false,
				"error":       "rate_limit_exceeded",
				"message":     "Too many requests. Please try again later.",
				"retry_after": int(ttl.Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getClientIdentifier получает идентификатор клиента
func getClientIdentifier(c *gin.Context) string {
	// будем индицировать только по ip
	return fmt.Sprintf("ip:%s", c.ClientIP())
}
