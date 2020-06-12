// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package etcdv3

import (
	"testing"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock/storetest"
)

func TestEtcdv3(t *testing.T) {
	store, _ := New([]string{"http://127.0.0.1:2379"})
	storetest.DoTest(t, store)
}
