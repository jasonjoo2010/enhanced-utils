package distlock

import (
	"fmt"
	"testing"
	"time"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock/mock"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	uuid, created := parseLockData("")
	assert.Empty(t, uuid)

	uuid, created = parseLockData("uuid")
	assert.Empty(t, uuid)

	uuid, created = parseLockData("uuid|")
	assert.Empty(t, uuid)

	uuid, created = parseLockData("|123333")
	assert.Empty(t, uuid)

	uuid, created = parseLockData("uuid|123333")
	assert.Equal(t, "uuid", uuid)
	assert.Equal(t, int64(123333), created)
}

func TestKey(t *testing.T) {
	lock := NewMutex("", 0, &mock.MockLocker{})
	assert.Equal(t, "lock::distributed-lock::asdf", lock.key("asdf"))

	lock = NewMutex("test", 0, &mock.MockLocker{})
	assert.Equal(t, "lock::test::asdf", lock.key("asdf"))
}

func TestVerify(t *testing.T) {
	store := mock.New()
	lock := NewMutex("test", 500*time.Second, store)
	key := lock.key("demo")

	store.Delete(key)
	valid, myself := lock.verify(key)
	assert.False(t, valid)

	store.Set(key, "", 500*time.Second)
	valid, myself = lock.verify(key)
	assert.False(t, valid)

	store.Set(key, fmt.Sprintf("%s|%d", lock.uuid, time.Now().UnixNano()/1e6-501*1e3), 500*time.Second)
	valid, myself = lock.verify(key)
	assert.False(t, valid)

	store.Set(key, fmt.Sprintf("%s|%d", lock.uuid, time.Now().UnixNano()/1e6-499*1e3), 500*time.Second)
	valid, myself = lock.verify(key)
	assert.True(t, valid)
	assert.True(t, myself)
}
