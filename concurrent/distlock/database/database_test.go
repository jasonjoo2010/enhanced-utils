// Copyright 2020 The GoSchedule Authors. All rights reserved.
// Use of this source code is governed by BSD
// license that can be found in the LICENSE file.

package database

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonjoo2010/enhanced-utils/concurrent/distlock/storetest"
	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/test")
	assert.Nil(t, err)
	storetest.DoTest(t, New(db))
}
