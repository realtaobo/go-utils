package sync

import (
	"runtime"
	"sync/atomic"
)

type SpinLock struct {
	lock uintptr
}

// Lock locks l.
func (l *SpinLock) Lock() {
	for !atomic.CompareAndSwapUintptr(&l.lock, 0, 1) {
		runtime.Gosched()
	}
}

// Unlock unlocks l.
func (l *SpinLock) Unlock() {
	atomic.StoreUintptr(&l.lock, 0)
}

// TryLock
func (l *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapUintptr(&l.lock, 0, 1)
}
