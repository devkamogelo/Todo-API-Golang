package store

import (
	"context"
	"fmt"
	"time"
	"todo_api/internal/config"

	"github.com/redis/go-redis/v9"
)

type Redis struct{ Client *redis.Client }

func NewRedis(cfg *config.Config) (*Redis, error) {
	url := cfg.RedisUrl
	if url == "" {
		return nil, fmt.Errorf("no redis Url provided")
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("invalid redis Url")
	}

	//remember to defer close.Client when calling this func
	rdb := redis.NewClient(opts)

	return &Redis{Client: rdb}, nil
}

func (r *Redis) SetJTI(ctx context.Context, key string, userId string, exp time.Time) error {
	return r.Client.Set(ctx, key, userId, time.Until(exp)).Err()
}

func (r *Redis) DelJTI(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *Redis) GetUserByJTI(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}
