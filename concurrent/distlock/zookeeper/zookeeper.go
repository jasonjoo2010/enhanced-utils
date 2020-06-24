// Copyright 2020 The enhanced-utils Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package zookeeper

import (
	"strings"
	"time"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock"
	"github.com/jasonjoo2010/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
)

// Structure: /lock/<namespace>/[sharding]/<md5(key)>

type ACL struct {
	Username, Password string
}

type Option func(l *LockerOption)

type ZookeeperLocker struct {
	conn      *zk.Conn
	namespace string
	ttl       int64
	prefix    string
	shards    int
	acl       []zk.ACL
	stopped   bool
}

type LockerOption struct {
	root         string
	acl          *ACL
	shardingBits int
}

func WithShardingBits(bits int) Option {
	return func(o *LockerOption) {
		if bits > 0 && bits < 8 {
			o.shardingBits = bits
		} else {
			o.shardingBits = 1
		}
	}
}

func WithRoot(rootPath string) Option {
	rootPath = strings.TrimRight(rootPath, "/")
	return func(o *LockerOption) {
		o.root = rootPath
	}
}

func WithAcl(username, password string) Option {
	return func(o *LockerOption) {
		if username != "" && password != "" {
			o.acl = &ACL{username, password}
		} else {
			o.acl = nil
		}
	}
}

// New create a zookeeper locker based on specific namespace
//	Due to extra initialization should be done before using so you cannot
//	invoke it with different namespaces after creation.
func New(namespace string, ttl int64, addrs []string) *ZookeeperLocker {
	return NewWithRoot("/lock", namespace, ttl, addrs)
}

func NewWithOptions(namespace string, ttl int64, addrs []string, opts ...Option) *ZookeeperLocker {
	opt := &LockerOption{
		root:         "/lock",
		shardingBits: 1,
	}
	for _, fn := range opts {
		fn(opt)
	}
	conn, eventC, err := zk.Connect(
		addrs,
		60*time.Second,
		zk.WithLogger(logrus.StandardLogger()),
		zk.WithLogInfo(true),
	)
	if err != nil {
		logrus.Error("Can't initialize zookeeper locker: ", err.Error())
		return nil
	}
	timeout := time.NewTimer(10 * time.Second)
LOOP_CHECK:
	for {
		select {
		case event := <-eventC:
			if event.State == zk.StateHasSession {
				break LOOP_CHECK
			}
		case <-timeout.C:
			logrus.Error("Can't connect to zookeeper server: timeout")
			return nil
		}
	}
	// initial structure
	instance := &ZookeeperLocker{
		conn:      conn,
		namespace: namespace,
		ttl:       ttl,
		shards:    1 << opt.shardingBits,
		prefix:    opt.root + "/" + strings.ReplaceAll(namespace, "/", "_"),
		acl:       zk.WorldACL(zk.PermAll),
	}
	if opt.acl != nil && opt.acl.Username != "" {
		instance.acl = []zk.ACL{}
		instance.acl = append(instance.acl, zk.WorldACL(zk.PermRead)...)
		instance.acl = append(instance.acl, zk.DigestACL(zk.PermAll, opt.acl.Username, opt.acl.Password)...)
		conn.AddAuth("digest", []byte(opt.acl.Username+":"+opt.acl.Password))
	}
	instance.initialize()
	return instance
}

func NewWithRoot(root, namespace string, ttl int64, addrs []string) *ZookeeperLocker {
	return NewWithRootAcl(root, namespace, ttl, addrs, nil)
}

func NewWithRootAcl(root, namespace string, ttl int64, addrs []string, acl *ACL) *ZookeeperLocker {
	var opts []Option
	opts = append(opts, WithRoot(root))
	if acl != nil {
		opts = append(opts, WithAcl(acl.Username, acl.Password))
	}
	return NewWithOptions(namespace, ttl, addrs, opts...)
}

func (z *ZookeeperLocker) Keep(lockKey *distlock.LockKey, val string, expire time.Duration) {
	z.check(lockKey)
	z.conn.Set(z.key(lockKey), []byte(val), -1)
}

func (z *ZookeeperLocker) Exists(lockKey *distlock.LockKey) bool {
	z.check(lockKey)
	key := z.key(lockKey)
	_, stat, err := z.conn.Get(key)
	if err != nil {
		return false
	}
	now := time.Now().UnixNano() / 1e6
	if now-stat.Mtime > z.ttl {
		z.conn.Delete(key, -1)
		return false
	}
	return true
}

func (z *ZookeeperLocker) Get(lockKey *distlock.LockKey) string {
	z.check(lockKey)
	data, _, _ := z.conn.Get(z.key(lockKey))
	return string(data)
}

func (z *ZookeeperLocker) Set(lockKey *distlock.LockKey, val string, expire time.Duration) {
	z.check(lockKey)
	if !z.SetIfAbsent(lockKey, val, expire) {
		_, err := z.conn.Set(z.key(lockKey), []byte(val), -1)
		if err != nil {
			logrus.Warn("Write zookeeper failed: ", err.Error())
		}
	}
}

func (z *ZookeeperLocker) SetIfAbsent(lockKey *distlock.LockKey, val string, expire time.Duration) bool {
	z.check(lockKey)
	key := z.key(lockKey)
	if z.Exists(lockKey) {
		return false
	}
	_, err := z.conn.Create(key, []byte(val), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return false
	}
	return true
}

func (z *ZookeeperLocker) Delete(lockKey *distlock.LockKey) {
	z.check(lockKey)
	err := z.conn.Delete(z.key(lockKey), -1)
	if err != nil {
		logrus.Warn("Delete from zookeeper failed: ", err.Error())
	}
}

func (z *ZookeeperLocker) Close() {
	z.stopped = true
	z.conn.Close()
}
