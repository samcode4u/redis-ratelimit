package redisratelimit

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRateLimit struct {
	ctx context.Context
	rdb *redis.Client
}

func (rrl *RedisRateLimit) InitClient() {
	rrl.ctx = context.Background()
	rrl.rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (rrl *RedisRateLimit) ResetRateLimit(key string) bool {
	rrl.rdb.Del(rrl.ctx, key)
	return true
}

func (rrl *RedisRateLimit) CheckRateLimit(key string, limit int64, isincr bool, duration time.Duration) (bool, int64, time.Duration) {
	var durationTTL time.Duration
	val, err := rrl.rdb.Get(rrl.ctx, key).Int64()
	if err == redis.Nil {
		rrl.rdb.Incr(rrl.ctx, key)
		rrl.rdb.Expire(rrl.ctx, key, duration)
	} else if err != nil {
		return false, -1, -1
	} else if val > limit {
		durationTTL = rrl.rdb.TTL(rrl.ctx, key).Val()
		return false, limit - val, durationTTL
	} else if isincr {
		rrl.rdb.Incr(rrl.ctx, key)
	}
	durationTTL = rrl.rdb.TTL(rrl.ctx, key).Val()
	return true, limit - val, durationTTL
}
