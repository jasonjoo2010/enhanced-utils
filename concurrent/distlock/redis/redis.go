package redis

import (
	"time"

	goredis "github.com/go-redis/redis"
)

type RedisLocker struct {
	client goredis.UniversalClient
}

func New(addrs []string) *RedisLocker {
	return &RedisLocker{
		client: goredis.NewUniversalClient(&goredis.UniversalOptions{
			Addrs:        addrs,
			PoolSize:     2,
			DialTimeout:  2 * time.Second,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		}),
	}
}

func (r *RedisLocker) Keep(key, val string, expire time.Duration) {
	r.client.Set(key, val, expire)
}

func (r *RedisLocker) Exists(key string) bool {
	return r.client.Exists(key).Val() > 0
}

func (r *RedisLocker) Get(key string) string {
	return r.client.Get(key).Val()
}

func (r *RedisLocker) Set(key, val string, expire time.Duration) {
	r.client.Set(key, val, expire)
}

func (r *RedisLocker) SetIfAbsent(key, val string, expire time.Duration) bool {
	return r.client.SetNX(key, val, expire).Val()
}

func (r *RedisLocker) Delete(key string) {
	r.client.Del(key)
}
