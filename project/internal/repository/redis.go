package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"

	"currency-converter/internal/config"
)

const logStreamKey = "conversion_logs"

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(ctx context.Context, cfg *config.Config) (*RedisRepository, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	if _, err := rdb.Ping().Result(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisRepository{client: rdb}, nil
}

func (r *RedisRepository) LogConversion(ctx context.Context, from, to string, amount, result float64) error {
	logMessage := fmt.Sprintf("Converted %.2f %s to %.2f %s", amount, from, result, to)

	err := r.client.LPush(logStreamKey, logMessage).Err()

	if err != nil {
		return fmt.Errorf("failed to log conversion to redis: %w", err)
	}

	return nil
}
