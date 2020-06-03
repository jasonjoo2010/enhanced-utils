package storetest

import (
	"testing"
	"time"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock"
	"github.com/stretchr/testify/assert"
)

func DoTest(t *testing.T, s distlock.Store) {
	key := &distlock.LockKey{"testns", "demo"}
	expire := time.Second * 2
	s.Delete(key)

	s.Set(key, "t1", expire)
	assert.True(t, s.Exists(key))
	assert.Equal(t, "t1", s.Get(key))
	s.Set(key, "t2", expire)
	assert.True(t, s.Exists(key))
	assert.Equal(t, "t2", s.Get(key))
	s.Delete(key)
	assert.False(t, s.Exists(key))

	assert.False(t, s.Exists(key))
	assert.True(t, s.SetIfAbsent(key, "test", expire))
	assert.True(t, s.Exists(key))
	assert.False(t, s.SetIfAbsent(key, "test", expire))
	assert.Equal(t, "test", s.Get(key))
	s.Keep(key, "test1", expire)
	time.Sleep(time.Second)
	assert.True(t, s.Exists(key))
	assert.Equal(t, "test1", s.Get(key))
	s.Keep(key, "test2", expire)
	time.Sleep(time.Second)
	assert.True(t, s.Exists(key))
	assert.Equal(t, "test2", s.Get(key))
	time.Sleep(time.Second)
	time.Sleep(time.Second)
	assert.False(t, s.Exists(key))

	DoTestMutex(t, s)
	DoTestReentry(t, s)
}

func DoTestMutex(t *testing.T, s distlock.Store) {
	lock := distlock.NewMutex("test", 2*time.Second, s)
	id := 3333
	id1 := 4444

	assert.True(t, lock.TryLock(id))
	assert.True(t, lock.TryLock(id1))
	assert.False(t, lock.TryLock(id))
	assert.False(t, lock.TryLock(id1))
	assert.True(t, lock.UnLock(id))
	assert.True(t, lock.UnLock(id1))
	assert.True(t, lock.TryLock(id))
	assert.True(t, lock.UnLock(id))

	assert.NoError(t, lock.Lock(id, 1*time.Second))
	assert.Error(t, lock.Lock(id, 1*time.Second))
	assert.NoError(t, lock.Lock(id, 2*time.Second))
	assert.Error(t, lock.Lock(id, 1*time.Second))
	assert.Error(t, lock.Lock(id, 500*time.Millisecond))
	lock.Keep(id)
	assert.Error(t, lock.Lock(id, 1500*time.Millisecond))
	assert.True(t, lock.UnLock(id))
}

func DoTestReentry(t *testing.T, s distlock.Store) {
	lock := distlock.NewReentry("test", 2*time.Second, s)
	lock1 := distlock.NewReentry("test", 2*time.Second, s)
	id := 3333
	id1 := 4444

	assert.True(t, lock.TryLock(id))
	assert.True(t, lock.TryLock(id1))
	assert.True(t, lock.TryLock(id))
	assert.True(t, lock.TryLock(id1))
	assert.True(t, lock.UnLock(id))
	assert.True(t, lock.UnLock(id1))
	assert.False(t, lock.UnLock(id))
	assert.False(t, lock.UnLock(id1))
	assert.True(t, lock.TryLock(id))
	assert.True(t, lock.UnLock(id))

	assert.NoError(t, lock1.Lock(id, 1*time.Second))
	assert.Error(t, lock.Lock(id, 1*time.Second))
	assert.NoError(t, lock.Lock(id, 2*time.Second))
	assert.Error(t, lock1.Lock(id, 1*time.Second))
	assert.Error(t, lock1.Lock(id, 500*time.Millisecond))
	lock.Keep(id)
	assert.Error(t, lock1.Lock(id, 1500*time.Millisecond))
	assert.True(t, lock.UnLock(id))
}
