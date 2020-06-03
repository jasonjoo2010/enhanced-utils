package distlock

import "time"

type DistLock interface {
	// Keep renew a lock held already for another {expire} time
	Keep(target interface{})
	// Lock try to lock the specified resource in {wait} time or return a LockFailed error
	Lock(target interface{}, wait time.Duration) error
	TryLock(target interface{}) bool
	// UnLock releases the lock of specified resource id and return true for success
	UnLock(target interface{}) bool
}
