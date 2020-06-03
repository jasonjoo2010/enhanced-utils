// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package distlock

import (
	"fmt"
	"time"
)

type LockKey struct {
	Namespace, Key string
}

func (lk *LockKey) String() string {
	ns := lk.Namespace
	if ns == "" {
		ns = "distributed-lock"
	}
	return fmt.Sprintf("lock::%s::%v", ns, lk.Key)
}

type Store interface {
	// Keep ensures the lock data in storage won't disappear during the period and update the content
	Keep(lockKey *LockKey, val string, expire time.Duration)
	// Exists return the existence of specified key
	Exists(lockKey *LockKey) bool
	Get(lockKey *LockKey) string
	SetIfAbsent(lockKey *LockKey, val string, expire time.Duration) bool
	Set(lockKey *LockKey, val string, expire time.Duration)
	Delete(lockKey *LockKey)
	Close()
}
