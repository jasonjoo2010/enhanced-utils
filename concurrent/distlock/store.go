package distlock

import "time"

type Store interface {
	// Keep ensures the lock data in storage won't disappear during the period and update the content
	Keep(key, val string, expire time.Duration)
	// Exists return the existence of specified key
	Exists(key string) bool
	Get(key string) string
	SetIfAbsent(key, val string, expire time.Duration) bool
	Set(key, val string, expire time.Duration)
	Delete(key string)
}
