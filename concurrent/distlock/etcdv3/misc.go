// Copyright 2020 The enhanced-utils Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package etcdv3

import (
	"context"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock"
	"github.com/sirupsen/logrus"
)

func (s *Etcdv3Locker) check() {
	if s.stopped {
		panic("Locker has been stopped")
	}
}

func (s *Etcdv3Locker) notExisted(lockKey *distlock.LockKey) etcd.Cmp {
	return etcd.Compare(etcd.CreateRevision(s.key(lockKey)), "=", 0)
}

func (s *Etcdv3Locker) key(lockKey *distlock.LockKey) string {
	return s.prefix + "/" + lockKey.Namespace + "/" + lockKey.Key
}

func (s *Etcdv3Locker) lease(expire time.Duration) etcd.LeaseID {
	resp, err := s.leaseApi.Grant(context.Background(), leaseTTL(expire))
	if err != nil {
		logrus.Warn("Create lease failed: ", err.Error())
		return 0
	}
	return resp.ID
}

// leaseTTL returns the seconds of TTL at least(seem to ceil())
func leaseTTL(expire time.Duration) int64 {
	millis := int64(expire / time.Millisecond)
	second := millis / 1000
	if millis%1000 > 0 {
		second++
	}
	if second < 1 {
		second = 1
	}
	return second
}
