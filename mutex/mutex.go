package mutex

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	// copy from /usr/local/go/src/sync/mutex.go
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexStarving
	mutexWaiterShift = iota
)

type Mutex struct {
	sync.Mutex
}

// get mutex waiters at now
func (m *Mutex) Count() int {
	// get sync.Mutex.state filed
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	waiters := state >> mutexWaiterShift
	waiters = waiters + (state & mutexLocked)
	return int(waiters)
}

// IsLocked 锁是否被持有
func (m *Mutex) IsLocked() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexLocked == mutexLocked
}

// IsWoken 是否有等待者被唤醒
func (m *Mutex) IsWoken() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexWoken == mutexWoken
}

// IsStarving 锁是否处于饥饿状态
func (m *Mutex) IsStarving() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexStarving == mutexStarving
}
