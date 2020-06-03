package redis

import (
	"time"

	goredis "github.com/go-redis/redis"
	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock"
)

type RedisLocker struct {
	client  goredis.UniversalClient
	stopped bool
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

func (r *RedisLocker) check() {
	if r.stopped {
		panic("Locker has been closed")
	}
}

func (r *RedisLocker) Keep(lockKey *distlock.LockKey, val string, expire time.Duration) {
	r.check()
	r.client.Set(lockKey.String(), val, expire)
}

func (r *RedisLocker) Exists(lockKey *distlock.LockKey) bool {
	r.check()
	return r.client.Exists(lockKey.String()).Val() > 0
}

func (r *RedisLocker) Get(lockKey *distlock.LockKey) string {
	r.check()
	return r.client.Get(lockKey.String()).Val()
}

func (r *RedisLocker) Set(lockKey *distlock.LockKey, val string, expire time.Duration) {
	r.check()
	r.client.Set(lockKey.String(), val, expire)
}

func (r *RedisLocker) SetIfAbsent(lockKey *distlock.LockKey, val string, expire time.Duration) bool {
	r.check()
	return r.client.SetNX(lockKey.String(), val, expire).Val()
}

func (r *RedisLocker) Delete(lockKey *distlock.LockKey) {
	r.check()
	r.client.Del(lockKey.String())
}

func (r *RedisLocker) Close() {
	r.check()
	r.stopped = true
	r.client.Close()
}
