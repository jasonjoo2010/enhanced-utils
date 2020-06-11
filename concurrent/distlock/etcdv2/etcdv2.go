// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package etcdv2

import (
	"context"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock"
	"github.com/sirupsen/logrus"
)

type etcdv2LockerConfig struct {
	prefix, username, password string
}

type Etcdv2Locker struct {
	client  etcd.Client
	keysApi etcd.KeysAPI
	prefix  string
	stopped bool
}

type Option func(cfg *etcdv2LockerConfig)

func WithCredential(username, password string) Option {
	return func(cfg *etcdv2LockerConfig) {
		cfg.username = username
		cfg.password = password
	}
}

func WithPrefix(prefix string) Option {
	return func(cfg *etcdv2LockerConfig) {
		cfg.prefix = prefix
	}
}

func New(addrs []string, opts ...Option) *Etcdv2Locker {
	lockerConfig := &etcdv2LockerConfig{}
	for _, fn := range opts {
		fn(lockerConfig)
	}
	cfg := etcd.Config{
		Endpoints: addrs,
	}
	if lockerConfig.username != "" {
		cfg.Username = lockerConfig.username
		cfg.Password = lockerConfig.password
	}
	if lockerConfig.prefix == "" {
		lockerConfig.prefix = "/lock"
	}
	c, err := etcd.New(cfg)
	if err != nil {
		logrus.Error("Failed to create locker store based on etcdv2: " + err.Error())
		return nil
	}
	return &Etcdv2Locker{
		client:  c,
		keysApi: etcd.NewKeysAPI(c),
		prefix:  lockerConfig.prefix,
	}
}

func (s *Etcdv2Locker) key(lockKey *distlock.LockKey) string {
	return s.prefix + "/" + lockKey.Namespace + "/" + lockKey.Key
}

func (s *Etcdv2Locker) Keep(lockKey *distlock.LockKey, val string, expire time.Duration) {
	_, err := s.keysApi.Set(context.Background(), s.key(lockKey), val, &etcd.SetOptions{
		TTL:       expire,
		PrevExist: etcd.PrevExist,
	})
	if err == nil {
		return
	}
	logrus.Warn("Keep failed: ", err.Error())
}

func (s *Etcdv2Locker) Exists(lockKey *distlock.LockKey) bool {
	resp, err := s.keysApi.Get(context.Background(), s.key(lockKey), nil)
	if err != nil {
		if errEtcd, ok := err.(etcd.Error); ok && errEtcd.Code == etcd.ErrorCodeKeyNotFound {
			return false
		}
		logrus.Warn("Fetch data from etcdv2 failed: ", err.Error())
	}
	return resp.Node.Value != ""
}

func (s *Etcdv2Locker) Get(lockKey *distlock.LockKey) string {
	resp, err := s.keysApi.Get(context.Background(), s.key(lockKey), nil)
	if err != nil {
		return ""
	}
	return resp.Node.Value
}

func (s *Etcdv2Locker) Set(lockKey *distlock.LockKey, val string, expire time.Duration) {
	s.keysApi.Set(context.Background(), s.key(lockKey), val, &etcd.SetOptions{
		TTL: expire,
	})
}

func (s *Etcdv2Locker) SetIfAbsent(lockKey *distlock.LockKey, val string, expire time.Duration) bool {
	_, err := s.keysApi.Set(context.Background(), s.key(lockKey), val, &etcd.SetOptions{
		TTL:       expire,
		PrevExist: etcd.PrevNoExist,
	})
	if err == nil {
		return true
	}
	if errEtcd, ok := err.(etcd.Error); ok && errEtcd.Code == etcd.ErrorCodeNodeExist {
		return false
	}
	logrus.Warn("SetIfAbsent failed: ", err.Error())
	return false
}

func (s *Etcdv2Locker) Delete(lockKey *distlock.LockKey) {
	s.keysApi.Delete(context.Background(), s.key(lockKey), nil)
}

func (s *Etcdv2Locker) Close() {
	// do nothing
}
