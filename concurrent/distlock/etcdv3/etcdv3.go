// Copyright 2020 The GoSchedule Authors. All rights reserved.
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

type Etcdv3Locker struct {
	distlock.Store
	client   *etcd.Client
	kvApi    etcd.KV
	leaseApi etcd.Lease
	prefix   string
	stopped  bool
}

func (s *Etcdv3Locker) Keep(lockKey *distlock.LockKey, val string, expire time.Duration) {
	s.check()
	key := s.key(lockKey)
	_, err := s.kvApi.Txn(context.Background()).
		If(s.notExisted(lockKey)).
		Else(etcd.OpPut(key, val, etcd.WithLease(s.lease(expire)))).
		Commit()
	if err != nil {
		logrus.Warn("Keep failed: ", err.Error())
		return
	}
}

func (s *Etcdv3Locker) Exists(lockKey *distlock.LockKey) bool {
	s.check()
	resp, err := s.kvApi.Get(context.Background(), s.key(lockKey))
	if err != nil {
		logrus.Warn("Fetch from etcdv3 error: ", err.Error())
		return false
	}
	return resp.Count == 1 && len(resp.Kvs[0].Value) > 0
}

func (s *Etcdv3Locker) Get(lockKey *distlock.LockKey) string {
	s.check()
	resp, err := s.kvApi.Get(context.Background(), s.key(lockKey))
	if err != nil {
		logrus.Warn("Fetch from etcdv3 error: ", err.Error())
		return ""
	}
	if resp.Count < 1 {
		return ""
	}
	return string(resp.Kvs[0].Value)
}

func (s *Etcdv3Locker) SetIfAbsent(lockKey *distlock.LockKey, val string, expire time.Duration) bool {
	s.check()
	key := s.key(lockKey)
	resp, err := s.kvApi.Txn(context.Background()).
		If(s.notExisted(lockKey)).
		Then(etcd.OpPut(key, val, etcd.WithLease(s.lease(expire)))).
		Commit()
	if err != nil {
		logrus.Warn("SetIfAbsent failed: ", err.Error())
		return false
	}
	return resp.Succeeded
}

func (s *Etcdv3Locker) Set(lockKey *distlock.LockKey, val string, expire time.Duration) {
	s.check()
	s.kvApi.Put(context.Background(), s.key(lockKey), val, etcd.WithLease(s.lease(expire)))
}

func (s *Etcdv3Locker) Delete(lockKey *distlock.LockKey) {
	s.check()
	s.kvApi.Delete(context.Background(), s.key(lockKey))
}

func (s *Etcdv3Locker) Close() {
	s.check()
	s.stopped = true
	s.client.Close()
}
