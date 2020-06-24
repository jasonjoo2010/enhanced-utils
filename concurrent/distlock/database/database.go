// Copyright 2020 The enhanced-utils Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock"
	"github.com/jasonjoo2010/godao"
	"github.com/jasonjoo2010/godao/options"
	"github.com/jasonjoo2010/godao/types"
	"github.com/sirupsen/logrus"
)

// Lock table structure:
// CREATE TABLE `lock` (
//   `id` bigint NOT NULL AUTO_INCREMENT,
//   `key` varchar(100) NOT NULL DEFAULT '',
//   `value` varchar(100) NOT NULL DEFAULT '',
//   `version` bigint NOT NULL DEFAULT '0',
//   `created` bigint NOT NULL DEFAULT '0',
//   `expire` bigint NOT NULL DEFAULT '0',
//   PRIMARY KEY (`id`),
//   UNIQUE KEY `key` (`key`)
// ) ENGINE=InnoDB;

const (
	millis = 1e6
)

type lockStruct struct {
	Id      int64 `dao:"primary;auto_increment"`
	Key     string
	Value   string
	Version int64 // reserved
	Created int64
	Expire  int64
}

type databaseLockerConfig struct {
	prefix string
	table  string
}

type DatabaseLocker struct {
	db      *sql.DB
	dao     *godao.Dao
	table   string
	prefix  string
	stopped bool
}

type Option func(cfg *databaseLockerConfig)

func WithPrefix(prefix string) Option {
	return func(cfg *databaseLockerConfig) {
		cfg.prefix = prefix
	}
}

func WithTable(table string) Option {
	return func(cfg *databaseLockerConfig) {
		cfg.table = table
	}
}

func New(db *sql.DB, opts ...Option) *DatabaseLocker {
	lockerConfig := &databaseLockerConfig{}
	for _, fn := range opts {
		fn(lockerConfig)
	}
	if lockerConfig.prefix == "" {
		lockerConfig.prefix = "/lock"
	}
	if lockerConfig.table == "" {
		lockerConfig.table = "lock"
	}
	return &DatabaseLocker{
		db:     db,
		dao:    godao.NewDao(lockStruct{}, db, options.WithTable(lockerConfig.table)),
		table:  lockerConfig.table,
		prefix: lockerConfig.prefix,
	}
}

func (s *DatabaseLocker) key(lockKey *distlock.LockKey) string {
	return s.prefix + "/" + lockKey.Namespace + "/" + lockKey.Key
}

func (s *DatabaseLocker) Keep(lockKey *distlock.LockKey, val string, expire time.Duration) {
	affected, err := s.dao.UpdateBy(context.Background(), (&godao.Query{}).
		Equal("Key", s.key(lockKey)).
		Data(),
		&types.UpdateEntry{
			Field: "Value",
			Value: val,
		},
		&types.UpdateEntry{
			Field: "Expire",
			Value: time.Now().Add(expire).UnixNano() / millis,
		},
	)
	if err != nil {
		logrus.Warn("Keep failed: ", err.Error())
		return
	}
	if affected < 1 {
		logrus.Warn("Keep failed: no lock exists")
		return
	}
}

func (s *DatabaseLocker) Exists(lockKey *distlock.LockKey) bool {
	return s.Get(lockKey) != ""
}

func (s *DatabaseLocker) Get(lockKey *distlock.LockKey) string {
	obj, err := s.dao.SelectOneBy(context.Background(), "Key", s.key(lockKey))
	if err != nil {
		logrus.Warn("Fetch data from database failed: ", err.Error())
		return ""
	}
	if obj == nil {
		return ""
	}
	lockObj := obj.(*lockStruct)
	if lockObj.Expire < time.Now().UnixNano()/millis {
		// expired
		logrus.Info("Release an expired lock: ", lockKey.Key, "@", lockKey.Namespace, " ", lockObj.Expire)
		s.dao.Delete(context.Background(), lockObj.Id)
		return ""
	}
	return lockObj.Value
}

func (s *DatabaseLocker) newLock(lockKey *distlock.LockKey, val string, expire time.Duration) *lockStruct {
	return &lockStruct{
		Key:     s.key(lockKey),
		Value:   val,
		Version: 1,
		Created: time.Now().UnixNano() / millis,
		Expire:  time.Now().Add(expire).UnixNano() / millis,
	}
}

func (s *DatabaseLocker) Set(lockKey *distlock.LockKey, val string, expire time.Duration) {
	affected, _, err := s.dao.Insert(context.Background(), s.newLock(lockKey, val, expire), options.WithReplace())
	if err != nil {
		logrus.Warn("Failed update: ", err.Error())
		return
	}
	if affected < 1 {
		logrus.Warn("Failed update: no effected row")
		return
	}
}

func (s *DatabaseLocker) SetIfAbsent(lockKey *distlock.LockKey, val string, expire time.Duration) bool {
	affected, _, err := s.dao.Insert(context.Background(), s.newLock(lockKey, val, expire), options.WithInsertIgnore())
	if err == nil && affected > 0 {
		return true
	}
	if err != nil {
		logrus.Warn("SetIfAbsent failed: ", err.Error())
		return false
	}
	return false
}

func (s *DatabaseLocker) Delete(lockKey *distlock.LockKey) {
	s.dao.DeleteRange(context.Background(), (&godao.Query{}).
		Equal("Key", s.key(lockKey)).
		Data())
}

func (s *DatabaseLocker) Close() {
	// do nothing
}
