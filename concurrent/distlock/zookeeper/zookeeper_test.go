// Copyright 2020 The enhanced-utils Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package zookeeper

import (
	"testing"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock/storetest"
)

func TestZookeeper(t *testing.T) {
	store := NewWithRoot("/demo", "testns", 2000, []string{"127.0.0.1:2181"})
	storetest.DoTest(t, store)
	removePath(store.conn, "/demo")
	store.Close()
}
