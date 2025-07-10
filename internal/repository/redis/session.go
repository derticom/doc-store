package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Sessions struct {
	client *redis.Client
}

func NewSessions(ctx context.Context, addr string, numDB int) (*Sessions, error) {
	opts, err := redis.ParseURL(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to redis.ParseURL: %w", err)
	}

	opts.DB = numDB

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp := client.Ping(ctx)
	err = resp.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to client.Ping: %w", err)
	}

	return &Sessions{client: client}, nil
}

func (s *Sessions) Save(ctx context.Context, token, userID string, ttl time.Duration) error {
	return s.client.Set(ctx, "auth:"+token, userID, ttl).Err()
}

func (s *Sessions) GetUserID(ctx context.Context, token string) (string, error) {
	return s.client.Get(ctx, "auth:"+token).Result()
}

func (s *Sessions) Delete(ctx context.Context, token string) error {
	return s.client.Del(ctx, "auth:"+token).Err()
}
