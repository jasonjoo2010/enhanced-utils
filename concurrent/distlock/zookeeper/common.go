// Copyright 2020 The enhanced-utils Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package zookeeper

import (
	"crypto/md5"
	"strconv"
	"strings"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock"
	"github.com/jasonjoo2010/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
)

func (z *ZookeeperLocker) exists(pathPlain string) bool {
	result, _, err := z.conn.Exists(pathPlain)
	if err != nil {
		logrus.Warn("Failed to execute Exists(", pathPlain, "): ", err.Error())
		return false
	}
	return result
}

func (z *ZookeeperLocker) createPath(pathPlain string, createParent bool) error {
	if !createParent {
		_, err := z.conn.Create(pathPlain, nil, 0, z.acl)
		return err
	}
	b := strings.Builder{}
	for _, str := range splitPath(pathPlain) {
		b.WriteString("/")
		b.WriteString(str)
		path := b.String()
		if z.exists(path) {
			continue
		}
		_, err := z.conn.Create(path, nil, 0, z.acl)
		if err != nil {
			logrus.Warn("Failed to create path ", path, ": ", err.Error())
			return err
		}
	}
	return nil
}

func (z *ZookeeperLocker) initialize() {
	if !z.exists(z.prefix) {
		z.createPath(z.prefix, true)
	}
	for i := 0; i < z.shards; i++ {
		path := z.prefix + "/" + strconv.Itoa(i)
		if z.exists(path) {
			continue
		}
		err := z.createPath(path, false)
		if err != nil && err.Error() != zk.ErrNodeExists.Error() {
			logrus.Error("Initialize distlock failed: ", err.Error())
		}
	}
}

func (z *ZookeeperLocker) key(lockKey *distlock.LockKey) string {
	b := strings.Builder{}
	hash := md5.Sum([]byte(lockKey.Key))
	b.WriteString(z.prefix)
	b.WriteString("/")
	b.WriteString(strconv.Itoa(int(hash[0]) % z.shards))
	b.WriteString("/")
	b.WriteString(lockKey.Key)
	return b.String()
}

func (z *ZookeeperLocker) check(lockKey *distlock.LockKey) {
	if z.stopped {
		panic("Locker has been stopped")
	}
	if lockKey.Namespace != z.namespace {
		panic("Specific namespace is not allowed: " + lockKey.Namespace)
	}
}
