package storetest

import (
	"testing"
	"time"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock"
	"github.com/stretchr/testify/assert"
)

func DoTest(t *testing.T, s distlock.Store) {
	key := "demo"
	expire := time.Second * 2
	s.Delete(key)

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
}
