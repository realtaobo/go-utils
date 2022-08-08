package mutex

import (
	"testing"

	"github.com/google/uuid"
)

func TestRecursiveMutex(t *testing.T) {
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
