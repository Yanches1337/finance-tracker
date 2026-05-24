package redis

import (
	"backend/internal/utils"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var Client *redis.Client

func InitRedis(cfg *utils.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	utils.Log.Info("Successfully connected to Redis")
	return nil
}
