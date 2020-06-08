// Copyright 2020 The GoSchedule Authors. All rights reserved.
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

// Structure: /lock/<namespace>/[0-127]/<md5(key)>

type ZookeeperLocker struct {
	conn      *zk.Conn
	namespace string
	ttl       int64
	prefix    string
	acl       []zk.ACL
	stopped   bool
}

// New create a zookeeper locker based on specific namespace
//	Due to extra initialization should be done before using so you cannot
//	invoke it with different namespaces after creation.
func New(namespace string, ttl int64, addrs []string) *ZookeeperLocker {
	return NewWithRoot("/lock", namespace, ttl, addrs)
}

func NewWithRoot(root, namespace string, ttl int64, addrs []string) *ZookeeperLocker {
	root = strings.TrimRight(root, "/")
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
		prefix:    root + "/" + strings.ReplaceAll(namespace, "/", "_"),
		acl:       zk.WorldACL(zk.PermAll),
	}
	instance.initialize()
	return instance
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
