package store

import (
	"context"
	"time"
	"todo_api/internal/config"

	"github.com/redis/go-redis/v9"
)

type Redis struct{ Client *redis.Client }

func NewRedis(cfg *config.Config) *Redis {
	addr := cfg.RedisAddr
	if addr == "" {
		addr = "localhost:6379"
	}

	//remember to defer close.Client when calling this func
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return &Redis{Client: rdb}
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
