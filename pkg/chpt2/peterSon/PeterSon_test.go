package chpt2

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type PeterSon struct {
	flags  []bool
	victim int
}

func (peterson *PeterSon) Lock(id int) {
	peterson.flags[id] = true
	peterson.victim = id
	for peterson.flags[1-id] && peterson.victim == id {
	}
}

func (peterson *PeterSon) Unlock(id int) {
	peterson.flags[id] = false
}

type Writer struct {
	lock   *PeterSon
	num    int
	waiter sync.WaitGroup
}

func (writer *Writer) increment(id int) {
	writer.lock.Lock(id)
	writer.num++
	writer.lock.Unlock(id)
	writer.waiter.Done()
}

// PeterSon 是结合了LockOne和LockTwo算法的
// 简洁完美的双线程算法
// PeterSon锁是无死锁无饥饿的
func TestPeterSon(t *testing.T) {
	for i := 0; i < 500000; i++ {
		writer := &Writer{
			num: 0,
			lock: &PeterSon{
				flags: make([]bool, 2),
			},
		}
		writer.waiter.Add(2)
		go writer.increment(0)
		go writer.increment(1)
		writer.waiter.Wait()
		require.Equal(t, writer.num, 2)
	}
}
