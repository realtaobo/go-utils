package mutex

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/petermattis/goid"
)

// TryLock 方法已在go1.18之后实现，此处不在添加
// https://cs.opensource.google/go/go/+/refs/tags/go1.19:src/sync/mutex.go;l=98;bpv=0;bpt=1

// RecursiveMutex包装一个Mutex,实现可重入
type RecursiveMutex struct {
	sync.Mutex
	owner     int64 //当前持有锁的goroutine id
	recursion int32 // 当前goroutine 重入的次数
}

func (m *RecursiveMutex) Lock() {
	gid := goid.Get() // 获取到当前goroutine的id
	//如果当前持有锁的goroutine就是这次调用的goroutine,说明是重入
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	// 获得锁的goroutine第一次调用，记录下它的goroutine id
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

func (m *RecursiveMutex) Unlock() {
	gid := goid.Get()
	//非持有锁的goroutine尝试释放锁，错误的使用
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	m.recursion--
	if m.recursion != 0 {
		return
	}
	// 此goroutine最后一次调用，需要释放锁
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}

func (m *RecursiveMutex) TryLock() bool {
	gid := goid.Get()
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return true
	}
	if m.Mutex.TryLock() {
		atomic.StoreInt64(&m.owner, gid)
		m.recursion = 1
		return true
	}
	return false
}

// RecursiveMutexByToken包装一个Mutex,实现可重入
type RecursiveMutexByToken struct {
	sync.Mutex
	token     int64 //存储加锁成功时goroutine传入的token值
	recursion int32 // 当前goroutine 重入的次数
}

func (m *RecursiveMutexByToken) Lock(token int64) {
	if atomic.LoadInt64(&m.token) == token {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	atomic.StoreInt64(&m.token, token)
	m.recursion = 1
}

func (m *RecursiveMutexByToken) Unlock(token int64) {
	if atomic.LoadInt64(&m.token) != token {
		panic(fmt.Sprintf("wrong the token(%d): %d!", m.token, token))
	}
	m.recursion--
	if m.recursion != 0 {
		return
	}
	atomic.StoreInt64(&m.token, 0)
	m.Mutex.Unlock()
}

func (m *RecursiveMutexByToken) TryLock(token int64) bool {
	if atomic.LoadInt64(&m.token) == token {
		m.recursion++
		return true
	}
	if m.Mutex.TryLock() {
		atomic.StoreInt64(&m.token, token)
		m.recursion = 1
		return true
	}
	return false
}
