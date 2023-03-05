package bakery

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// Bakery是无死锁,先来先服务(必然不会饥饿),且互斥的
type BakeryLock struct {
	// flags[i]表示线程i打算进入临界区
	flags []bool
	// label[i]表示线程i是第label[i]个
	// 进入临界区
	labels []int
}

func NewBakeryLock(n int) *BakeryLock {
	return &BakeryLock{
		flags:  make([]bool, n),
		labels: make([]int, n),
	}
}

type Writer struct {
	lock   *BakeryLock
	num    int
	waiter sync.WaitGroup
}

func (writer *Writer) increment(id int) {
	writer.lock.Lock(id)
	writer.num++
	writer.lock.Unlock(id)
	writer.waiter.Done()
}

func (lock *BakeryLock) getLabel() int {
	var Max int = -1
	for i := 0; i < len(lock.labels); i++ {
		if Max < lock.labels[i] {
			Max = lock.labels[i]
		}
	}
	return Max + 1
}

// 为了保证在所有正在打算进入临界区的线程当中id线程
// 的(label,id)是最小的
func (lock *BakeryLock) comp(k, id int) bool {
	return lock.labels[k] < lock.labels[id] || (lock.labels[k] == lock.labels[id] && k < id)
}

// 采取字典序排序而不是使用简单的比较label的方式,
// 是为了防止两个线程在门廊区拿到的label是一样的
func (lock *BakeryLock) getLevelNum(id int) bool {
	for i := 0; i < len(lock.flags); i++ {
		if lock.flags[i] && i != id && lock.comp(i, id) {
			return true
		}
	}
	return false
}

func (lock *BakeryLock) Lock(id int) {
	lock.labels[id] = lock.getLabel()
	lock.flags[id] = true
	for lock.getLevelNum(id) {
	}
}

func (lock *BakeryLock) Unlock(id int) {
	lock.flags[id] = false
}

func TestBakery(t *testing.T) {
	var n int = 10
	for i := 0; i < 10000; i++ {
		writer := &Writer{
			lock: NewBakeryLock(n),
			num:  0,
		}
		writer.waiter.Add(n)
		for i := 0; i < n; i++ {
			go writer.increment(i)
		}
		writer.waiter.Wait()
		require.Equal(t, n, writer.num)
	}
}
