package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"grender/core/configReader"
)

type RedisUtil struct {
	Client *redis.Client
}

func (r *RedisUtil) Connect(cfg configReader.RedisCfg) {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Url,
		Password: cfg.Password,
		DB:       cfg.Db,
	})
}

func (r *RedisUtil) Get(key string) string {
	ctx := context.Background()
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func (r *RedisUtil) Lpop() string {
	ctx := context.Background()
	val, err := r.Client.LPop(ctx, configReader.Config.Redis.TaskKey).Result()
	if err != nil {
		return ""
	}
	return val
}

func (r *RedisUtil) Lpush(key, value string) bool {
	ctx := context.Background()
	err := r.Client.LPush(ctx, key, value).Err()
	if err != nil {
		return false
	}
	return true
}
