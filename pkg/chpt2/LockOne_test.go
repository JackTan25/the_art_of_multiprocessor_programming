package chpt2

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type LockOne struct {
	flags []bool
}

func (lockone *LockOne) Lock(id int) {
	lockone.flags[id] = true
	for lockone.flags[1-id] {
	}
	// fmt.Println("break ", lockone.flags[1-id], "id: ", 1-id)
}

func (lockone *LockOne) Unlock(id int) {
	lockone.flags[id] = false
}

type Writer struct {
	lock   *LockOne
	num    int
	waiter sync.WaitGroup
}

func (writer *Writer) increment(id int) {
	writer.lock.Lock(id)
	writer.num++
	// fmt.Println(writer.num, writer.lock.flags[id], "id: ", id)
	writer.lock.Unlock(id)
	writer.waiter.Done()
}

// LockOne程序的缺陷在于
// 会发送死锁的情况当两
// 个线程的flags更新为true
// 都发生在for循环这一步的
// 读这里的话,就会发生死锁
func TestLockOne(t *testing.T) {
	for i := 0; i < 500001; i++ {
		writer := &Writer{
			num: 0,
			lock: &LockOne{
				flags: make([]bool, 2),
			},
		}
		writer.waiter.Add(2)
		go writer.increment(0)
		go writer.increment(1)
		writer.waiter.Wait()
		// fmt.Println("--------------------------------")
		require.Equal(t, writer.num, 2)
	}
}
