package distlock

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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
}

func New(namespace string, store Store) *DistLock {
	return &DistLock{
		store:     store,
		namespace: namespace,
		uuid:      "", // XXX
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
func (l *DistLock) verify(key string) bool {
	val := l.store.Get(key)
	// XXX: Pay attention to the phantom reads of redis (double reading could solve it, but confirmed to do that)
	uuid, created := parseLockData(val)
	if uuid == "" {
		return false
	}
	diff := time.Duration(time.Now().UnixNano())/time.Millisecond - time.Duration(created)
	if diff > l.expire {
		return false
	}
	return true
}

func (l *DistLock) TryLock(target interface{}) bool {
	key := l.key(target)
	if l.store.Exists(key) {
		// verify the lock
		if l.verify(key) {
			// valid lock
			return false
		}
		log.Warnf("Force release an invalid lock for %v", target)
		l.store.Delete(key)
	}
	// try to lock
	val := fmt.Sprintf("%s|%d", l.uuid, time.Now().UnixNano()/1e6)
	return l.store.SetIfAbsent(key, val, l.expire)
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
