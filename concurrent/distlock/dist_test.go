package distlock

import (
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
	lock := New("", &mock.MockLocker{})
	assert.Equal(t, "lock::distributed-lock::asdf", lock.key("asdf"))

	lock = New("test", &mock.MockLocker{})
	assert.Equal(t, "lock::test::asdf", lock.key("asdf"))
}

func TestVerify(t *testing.T) {
	store := &mock.MockLocker{}
	lock := New("test", store)
	key := lock.key("demo")

	store.SetIfAbsent(key, "", 500*time.Second)

}
