// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package zookeeper

import (
	"container/list"
	"fmt"
	"strings"
	"testing"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock"
	"github.com/jasonjoo2010/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func getChildren(conn *zk.Conn, path string, recursive bool) []string {
	var result []string
	queue := list.New()
	queue.PushBack(strings.TrimRight(path, "/"))
	for queue.Len() > 0 {
		len := queue.Len()
		for i := 0; i < len; i++ {
			item := queue.Front()
			p := item.Value.(string)
			queue.Remove(item)
			children, _, err := conn.Children(p)
			if err != nil {
				logrus.Warn("Fetch children failed for ", p, ": ", err.Error())
				continue
			}
			for _, child := range children {
				childPath := p + "/" + child
				result = append(result, childPath)
				queue.PushBack(childPath)
			}
		}
		if !recursive {
			break
		}
	}
	return result
}

func removePath(conn *zk.Conn, root string) error {
	pathList := getChildren(conn, root, true)
	for i := len(pathList) - 1; i >= 0; i-- {
		err := conn.Delete(pathList[i], 0)
		if err != nil {
			return err
		}
	}
	return conn.Delete(root, 0)
}

func TestInitialize(t *testing.T) {
	store := NewWithRoot("/demolock", "test", 60000, []string{"127.0.0.1:2181"})

	assert.True(t, store.exists("/demolock"))
	assert.True(t, store.exists("/demolock/test"))
	assert.True(t, store.exists("/demolock/test/0"))
	assert.True(t, store.exists("/demolock/test/1"))

	removePath(store.conn, "/demolock")
	store.Close()
}

func TestCheck(t *testing.T) {
	store := NewWithRoot("/demolock", "test", 60000, []string{"127.0.0.1:2181"})

	assert.NotPanics(t, func() { store.check(&distlock.LockKey{"test", "a"}) })
	assert.Panics(t, func() { store.check(&distlock.LockKey{"test1", "a"}) })

	removePath(store.conn, "/demolock")

	store.Close()
	assert.Panics(t, func() { store.check(&distlock.LockKey{"test", "a"}) })
}

func TestKey(t *testing.T) {
	store := NewWithRoot("/demolock", "test", 60000, []string{"127.0.0.1:2181"})

	assert.Equal(t, 0, strings.Index(store.key(&distlock.LockKey{"test", "a"}), "/demolock/test/"))

	_, err := store.conn.Create("/demolock/test/a", nil, zk.FlagEphemeral, store.acl)
	fmt.Println("Created:", err)
	result, _, err := store.conn.Exists("/demolock/test/a")
	fmt.Println("Exists:", result, " with err:", err)
	result, _, err = store.conn.Exists("/demolock/test/b")
	fmt.Println("Exists:", result, " with err:", err)

	removePath(store.conn, "/demolock")
	store.Close()
}
