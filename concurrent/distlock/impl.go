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

type DistLockImpl struct {
	store     Store
	namespace string
	uuid      string
	expire    time.Duration
	reentry   bool
}

// NewMutex returns a non-reentry distributed lock
//	namespace is used to separate different projects
//	expire indicates the expiration of an active lock and it will be removed if no {Keep} and {Unlock} was invoked during this.
//	store decides which storage it uses
func NewMutex(namespace string, expire time.Duration, store Store) DistLock {
	return &DistLockImpl{
		store:     store,
		namespace: namespace,
		uuid:      strutils.RandString(20),
		expire:    expire,
		reentry:   false,
	}
}

// NewMutex returns a reentry distributed lock
//	namespace is used to separate different projects
//	expire indicates the expiration of an active lock and it will be removed if no {Keep} and {Unlock} was invoked during this.
//	store decides which storage it uses
func NewReentry(namespace string, expire time.Duration, store Store) DistLock {
	return &DistLockImpl{
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

func (l *DistLockImpl) key(target interface{}) *LockKey {
	ns := l.namespace
	if ns == "" {
		ns = "distributed-lock"
	}
	return &LockKey{
		Namespace: ns,
		Key:       fmt.Sprintf("%v", target),
	}
}

func (l *DistLockImpl) Close() {
	l.store.Close()
}

func (l *DistLockImpl) Keep(target interface{}) {
	lockKey := l.key(target)
	valid, myself := l.verify(lockKey)
	if valid && myself {
		val := fmt.Sprintf("%s|%d", l.uuid, time.Now().UnixNano()/1e6)
		l.store.Set(lockKey, val, l.expire)
	}
}

func (l *DistLockImpl) Lock(target interface{}, wait time.Duration) error {
	timer := time.NewTimer(wait)
	for {
		succ := l.TryLock(target)
		if succ {
			timer.Stop()
			return nil
		}
		select {
		case <-timer.C:
			// timeout
			return LockFailed
		default:
			time.Sleep(TRY_INTERVAL)
		}
	}
}

// verify an existed lock data structure and return true when valid
//	valid indicates whether the lock is valid
//	myself indicates whether the owner is myself
func (l *DistLockImpl) verify(lockKey *LockKey) (valid bool, myself bool) {
	val := l.store.Get(lockKey)
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

func (l *DistLockImpl) TryLock(target interface{}) bool {
	lockKey := l.key(target)
	if l.store.Exists(lockKey) {
		// verify the lock
		if valid, myself := l.verify(lockKey); valid {
			// valid lock
			if l.reentry && myself {
				// allow reentry, check whether already locked and update it
				val := fmt.Sprintf("%s|%d", l.uuid, time.Now().UnixNano()/1e6)
				l.store.Set(lockKey, val, l.expire)
				return true
			}
			return false
		}
		log.Warnf("Force release an invalid lock for %v", target)
		l.store.Delete(lockKey)
	}
	// try to lock
	val := fmt.Sprintf("%s|%d", l.uuid, time.Now().UnixNano()/1e6)
	succ := l.store.SetIfAbsent(lockKey, val, l.expire)
	return succ
}

func (l *DistLockImpl) UnLock(target interface{}) bool {
	lockKey := l.key(target)
	uuid, _ := parseLockData(l.store.Get(lockKey))
	if uuid != l.uuid {
		// only the lock who locked it can unlock
		return false
	}
	l.store.Delete(lockKey)
	return true
}
