package sync

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMutex(t *testing.T) {
	var mu Mutex
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			mu.Lock()
			time.Sleep(time.Second)
			mu.Unlock()
			wg.Done()
		}()
	}
	time.Sleep(time.Second)
	fmt.Printf("waitings: %d, isLocked: %t, IsWoken: %v, IsStarving: %v\n", mu.Count(), mu.IsLocked(), mu.IsWoken(), mu.IsStarving())
	wg.Wait()
	fmt.Printf("waitings: %d, isLocked: %t, IsWoken: %v, IsStarving: %v\n", mu.Count(), mu.IsLocked(), mu.IsWoken(), mu.IsStarving())
}

func TestRecursiveMutexNoToken(t *testing.T) {
	mtx1 := RecursiveMutex{}
	mtx1.Lock()
	mtx1.Lock()
	t.Log("twice lock")
	mtx1.Unlock()
	mtx1.Unlock()
}

func TestRecursiveMutexByToken(t *testing.T) {
	mtx1 := RecursiveMutexByToken{}
	token := int64(uuid.New().ID())
	mtx1.Lock(token)
	mtx1.Lock(token)
	t.Logf("token %d: twice lock", token) //token 1472495600: twice lock
	mtx1.Unlock(token)
	mtx1.Unlock(token)
}
