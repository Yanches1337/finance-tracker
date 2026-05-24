package services

import (
	"backend/internal/adapters/redis"
	"backend/internal/utils"
	"context"
	"fmt"
	redis1 "github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type TokenService struct {
	cfg *utils.Config
}

func NewTokenService(cfg *utils.Config) *TokenService {
	return &TokenService{
		cfg: cfg,
	}
}

func (s *TokenService) SaveRefreshToken(ctx context.Context, userID int64, refreshToken, ip, userAgent string) error {
	ttlStr := s.cfg.Redis.RefreshTokenTTL
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		ttl = 7 * 24 * time.Hour
	}

	value := fmt.Sprintf("%d|%s|%s|%s", userID, time.Now().UTC().Format(time.RFC3339), ip, userAgent)

	key := "refresh:" + refreshToken
	return redis.Client.Set(ctx, key, value, ttl).Err()
}

func (s *TokenService) ValidateRefreshToken(ctx context.Context, refreshToken string) (int64, error) {
	key := "refresh:" + refreshToken
	val, err := redis.Client.Get(ctx, key).Result()
	if err == redis1.Nil {
		return 0, fmt.Errorf("refresh token not found or expired")
	}
	if err != nil {
		return 0, err
	}

	parts := strings.Split(val, "|")
	if len(parts) < 1 {
		return 0, fmt.Errorf("invalid token data")
	}

	var userID int64
	_, err = fmt.Sscanf(parts[0], "%d", &userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *TokenService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	key := "refresh:" + refreshToken
	return redis.Client.Del(ctx, key).Err()
}
