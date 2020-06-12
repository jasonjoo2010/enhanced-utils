// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package etcdv3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLeaseTTL(t *testing.T) {
	assert.Equal(t, int64(1), leaseTTL(0))
	assert.Equal(t, int64(1), leaseTTL(time.Second))
	assert.Equal(t, int64(1), leaseTTL(1*time.Millisecond))
	assert.Equal(t, int64(1), leaseTTL(500*time.Millisecond))
	assert.Equal(t, int64(1), leaseTTL(999*time.Millisecond))
	assert.Equal(t, int64(1), leaseTTL(1000*time.Millisecond))
	assert.Equal(t, int64(2), leaseTTL(1001*time.Millisecond))
	assert.Equal(t, int64(1), leaseTTL(-1*time.Millisecond))
}
