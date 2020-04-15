package distlock

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jasonjoo2010/enhanced-utils/strutils"
	"github.com/prometheus/common/log"
)

// lock value format: {uuid}|{locked timestamp in millisecond}

const TRY_INTERVAL time.Duration = 10 * time.Millisecond

var LockFailed = errors.New("Lock failed")

type DistLock struct {
	store     Store
	namespace string
	uuid      string
	expire    time.Duration
	reentry   bool
}

func NewMutex(namespace string, expire time.Duration, store Store) *DistLock {
	return &DistLock{
		store:     store,
		namespace: namespace,
		uuid:      strutils.RandString(20),
		expire:    expire,
		reentry:   false,
	}
}

func NewReentry(namespace string, expire time.Duration, store Store) *DistLock {
	return &DistLock{
		store:     store,
		namespace: namespace,
		uuid:      strutils.RandString(20),
		expire:    expire,
		reentry:   true,
	}
}

func parseLockData(data string) (uuid string, created int64) {
	if len(data) < 1 {
		return
	}

	pos := strings.IndexByte(data, '|')
	if pos < 0 {
		return
	}

	uuid = data[:pos]
	created, err := strconv.ParseInt(data[pos+1:], 10, 64)
	if err != nil {
		created = 0
		uuid = ""
		return
	}
	return
}

func (l *DistLock) key(target interface{}) string {
	ns := l.namespace
	if ns == "" {
		ns = "distributed-lock"
	}
	return fmt.Sprintf("lock::%s::%v", ns, target)
}

// Keep renew a lock held already for another {expire} time
func (l *DistLock) Keep(target interface{}) {
	key := l.key(target)
	valid, myself := l.verify(key)
	if valid && myself {
		val := fmt.Sprintf("%s|%d", l.uuid, time.Now().UnixNano()/1e6)
		l.store.Set(key, val, l.expire)
	}
}

// Lock try to lock the specified resource in {wait} time or return a LockFailed error
func (l *DistLock) Lock(target interface{}, wait time.Duration) error {
	for wait > 0 {
		succ := l.TryLock(target)
		if succ {
			return nil
		}
		time.Sleep(TRY_INTERVAL)
		wait -= TRY_INTERVAL
	}
	return LockFailed
}

// verify an existed lock data structure and return true when valid
//	valid indicates whether the lock is valid
//	myself indicates whether the owner is myself
func (l *DistLock) verify(key string) (valid bool, myself bool) {
	val := l.store.Get(key)
	// XXX: Pay attention to the phantom reads of redis (double reading could solve it, but confirmed to do that)
	uuid, created := parseLockData(val)
	if uuid == "" {
		return
	}
	diff := time.Duration(time.Now().UnixNano()-created*1e6) * time.Nanosecond
	if diff > l.expire {
		return
	}
	myself = uuid == l.uuid
	valid = true
	return
}

func (l *DistLock) TryLock(target interface{}) bool {
	key := l.key(target)
	if l.store.Exists(key) {
		// verify the lock
		if valid, myself := l.verify(key); valid {
			// valid lock
			if l.reentry && myself {
				// allow reentry, check whether already locked and update it
				val := fmt.Sprintf("%s|%d", l.uuid, time.Now().UnixNano()/1e6)
				l.store.Set(key, val, l.expire)
				return true
			}
			return false
		}
		log.Warnf("Force release an invalid lock for %v", target)
		l.store.Delete(key)
	}
	// try to lock
	val := fmt.Sprintf("%s|%d", l.uuid, time.Now().UnixNano()/1e6)
	succ := l.store.SetIfAbsent(key, val, l.expire)
	return succ
}

// UnLock releases the lock of specified resource id and return true for success
func (l *DistLock) UnLock(target interface{}) bool {
	key := l.key(target)
	uuid, _ := parseLockData(l.store.Get(key))
	if uuid != l.uuid {
		// only the lock who locked it can unlock
		return false
	}
	l.store.Delete(key)
	return true
}
