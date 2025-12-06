package cache

import (
	"ShareInfo/internal/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

func NewRedisClient(cfg *config.RedisConfig, logger *zap.Logger) (*redis.Client, error) {
	logger.Info("connecting to redis",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Int("db", cfg.DB),
	)

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		Password: cfg.Password,
		DB:       cfg.DB,

		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		PoolTimeout:     4 * time.Second,
		DialTimeout:     cfg.DialTimeout,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		MaxRetries:      cfg.MaxRetries,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			err := cn.Ping(ctx).Err()
			if err != nil {
				return err
			}
			return nil
		},
	})

	// healthcheck
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		err := client.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	logger.Info("successfully connected to redis",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Int("db", cfg.DB),
	)

	return client, nil
}

func Close(ctx context.Context, client *redis.Client, logger *zap.Logger) error {
	logger.Info("closing redis client")

	if err := client.Close(); err != nil {
		logger.Error("failed to close redis client", zap.Error(err))
		return err
	}

	logger.Info("successfully closed redis client")
	return nil
}
