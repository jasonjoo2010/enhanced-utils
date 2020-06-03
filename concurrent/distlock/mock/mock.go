package mock

import (
	"sync"
	"time"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock"
)

type item struct {
	val   string
	dueTo int64 // in nano
}

type MockLocker struct {
	sync.Mutex
	store   map[string]*item
	stopped bool
}

func New() *MockLocker {
	return &MockLocker{
		store: make(map[string]*item),
	}
}

func (m *MockLocker) Close() {
	m.Lock()
	defer m.Unlock()
	m.stopped = true
}

func (m *MockLocker) check() {
	if m.stopped {
		panic("Locker has been closed")
	}
}

func (m *MockLocker) Keep(lockKey *distlock.LockKey, val string, expire time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.check()
	key := lockKey.String()
	t, ok := m.store[key]
	if ok {
		t.dueTo = time.Now().UnixNano() + expire.Nanoseconds()
		t.val = val
	}
}

func (m *MockLocker) Exists(lockKey *distlock.LockKey) bool {
	m.Lock()
	defer m.Unlock()
	m.check()
	key := lockKey.String()
	t, ok := m.store[key]
	if ok && t.dueTo < time.Now().UnixNano() {
		delete(m.store, key)
		ok = false
	}
	return ok
}

func (m *MockLocker) Get(lockKey *distlock.LockKey) string {
	m.Lock()
	defer m.Unlock()
	m.check()
	key := lockKey.String()
	t, ok := m.store[key]
	if ok && t.dueTo < time.Now().UnixNano() {
		delete(m.store, key)
		ok = false
	}
	if ok {
		return t.val
	} else {
		return ""
	}
}

func (m *MockLocker) Set(lockKey *distlock.LockKey, val string, expire time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.check()
	key := lockKey.String()
	m.store[key] = &item{
		val:   val,
		dueTo: time.Now().UnixNano() + expire.Nanoseconds(),
	}
}

func (m *MockLocker) SetIfAbsent(lockKey *distlock.LockKey, val string, expire time.Duration) bool {
	m.Lock()
	defer m.Unlock()
	m.check()
	key := lockKey.String()
	if _, ok := m.store[key]; ok {
		return false
	}
	m.store[key] = &item{
		val:   val,
		dueTo: time.Now().UnixNano() + expire.Nanoseconds(),
	}
	return true
}

func (m *MockLocker) Delete(lockKey *distlock.LockKey) {
	m.Lock()
	defer m.Unlock()
	m.check()
	delete(m.store, lockKey.String())
}
