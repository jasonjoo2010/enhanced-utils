package mock

import (
	"sync"
	"time"
)

type item struct {
	val   string
	dueTo int64 // in nano
}

type MockLocker struct {
	sync.Mutex
	store map[string]*item
}

func New() *MockLocker {
	return &MockLocker{
		store: make(map[string]*item),
	}
}

func (m *MockLocker) Keep(key, val string, expire time.Duration) {
	m.Lock()
	defer m.Unlock()
	t, ok := m.store[key]
	if ok {
		t.dueTo = time.Now().UnixNano() + expire.Nanoseconds()
		t.val = val
	}
}

func (m *MockLocker) Exists(key string) bool {
	m.Lock()
	defer m.Unlock()
	t, ok := m.store[key]
	if ok && t.dueTo < time.Now().UnixNano() {
		delete(m.store, key)
		ok = false
	}
	return ok
}

func (m *MockLocker) Get(key string) string {
	m.Lock()
	defer m.Unlock()
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

func (m *MockLocker) Set(key, val string, expire time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.store[key] = &item{
		val:   val,
		dueTo: time.Now().UnixNano() + expire.Nanoseconds(),
	}
}

func (m *MockLocker) SetIfAbsent(key, val string, expire time.Duration) bool {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.store[key]; ok {
		return false
	}
	m.store[key] = &item{
		val:   val,
		dueTo: time.Now().UnixNano() + expire.Nanoseconds(),
	}
	return true
}

func (m *MockLocker) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.store, key)
}
