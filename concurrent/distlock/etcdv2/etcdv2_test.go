// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package etcdv2

import (
	"testing"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock/storetest"
)

func TestEtcdv2(t *testing.T) {
	storetest.DoTest(t, New([]string{"http://127.0.0.1:2379"}))
}
