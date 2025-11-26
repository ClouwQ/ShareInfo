package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	_ "github.com/redis/go-redis/v9"
	"time"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

// Get получаем значение по ключу
func (r *CacheRepository) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("miss value for key: %s ", key)
	}
	if err != nil {
		return nil, fmt.Errorf("error with redis client: %s", err)
	}
	return val, nil
}

// Set задаем значение по ключу
func (r *CacheRepository) Set(ctx context.Context, key string, value []byte, ttl int) error {
	err := r.client.Set(ctx, key, value, time.Duration(ttl)).Err()
	if err != nil {
		return fmt.Errorf("error with redis client: %s", err)
	}
	return nil
}

// Delete удаляем значение по ключу
func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	if r == nil || r.client == nil {
		return nil
	}
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error with redis client: %s", err)
	}
	return nil
}

// GetJSON получить данные в JSON
func (r *CacheRepository) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Get(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %s", err)
	}
	return nil
}

// SetJSON записать данные из JSON
func (r *CacheRepository) SetJSON(ctx context.Context, key string, value interface{}, ttl int) error {
	data, err := json.Marshal(&value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %s", err)
	}
	return r.Set(ctx, key, data, ttl)
}

// Exists проверяет существование значения по ключу
func (r *CacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("error with redis client: %s", err)
	}
	return result > 0, nil
}

// Expire устанавливает ttl по ключу
func (r *CacheRepository) Expire(ctx context.Context, key string, ttl int) error {
	err := r.client.Expire(ctx, key, time.Duration(ttl)).Err()
	if err != nil {
		return fmt.Errorf("error with redis client: %s", err)
	}
	return nil
}

// Increment увеличивает значение счетчика
func (r *CacheRepository) Increment(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment: %w", err)
	}
	return val, nil
}

// IncrementBy увеличивает значение счетчика на указанную величину
func (r *CacheRepository) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	val, err := r.client.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment by: %w", err)
	}
	return val, nil
}
