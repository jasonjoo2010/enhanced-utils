package mock

import (
	"testing"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock/storetest"
)

func TestMock(t *testing.T) {
	storetest.DoTest(t, New())
}
