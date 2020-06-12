// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package etcdv3

import (
	etcd "github.com/coreos/etcd/clientv3"
)

type etcdv3LockerConfig struct {
	prefix, username, password string
	ttl                        int
}

type Option func(cfg *etcdv3LockerConfig)

func WithCredential(username, password string) Option {
	return func(cfg *etcdv3LockerConfig) {
		cfg.username = username
		cfg.password = password
	}
}

func WithPrefix(prefix string) Option {
	return func(cfg *etcdv3LockerConfig) {
		cfg.prefix = prefix
	}
}

func WithTTL(sec int) Option {
	return func(cfg *etcdv3LockerConfig) {
		cfg.ttl = sec
	}
}

func New(addrs []string, opts ...Option) (*Etcdv3Locker, error) {
	lockerConfig := &etcdv3LockerConfig{}
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
	if lockerConfig.ttl < 1 {
		lockerConfig.ttl = 60
	}
	c, err := etcd.New(cfg)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &Etcdv3Locker{
		client:   c,
		kvApi:    etcd.NewKV(c),
		leaseApi: etcd.NewLease(c),
		prefix:   lockerConfig.prefix,
	}, nil
}
