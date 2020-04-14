package redis

import (
	"testing"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock/storetest"
)

func TestRedis(t *testing.T) {
	storetest.DoTest(t, New([]string{"127.0.0.1:6379"}))
}
