package filter

import (
	"multiprocesser/pkg/chpt2"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// 过滤锁是对peterson锁算法的更一般化的表现
// 假设有n个线程,level[i]代表第i个线程正在尝试
// 进入level[i]层，而victim[j]则表示第第j层淘
// 汰了线程victim[j]
type FilterLock struct {
	level  []int
	victim []int
	n      int
}

func NewFilterLock(n int) *FilterLock {
	return &FilterLock{
		level:  make([]int, n),
		victim: make([]int, n),
		n:      n,
	}
}

func (lock *FilterLock) getLevelNum(t, id int) bool {
	for i := 0; i < len(lock.level); i++ {
		if lock.level[i] >= t && i != id {
			return true
		}
	}
	return false
}

func (lock *FilterLock) Lock(id int) {
	// 没循环一次,对应那一层就会淘汰至少一个
	// 循环n-1次后淘汰了n-1个线程,那么最后剩
	// 下的拿到锁
	for i := 1; i < lock.n; i++ {
		// 尝试进入第i层
		lock.level[id] = i
		// 假设第i层淘汰了id线程
		lock.victim[i] = id
		// 当第i层没有人进入过或者说第i层已经又人被淘汰了
		// 就结束循环
		for lock.getLevelNum(i, id) && lock.victim[i] == id {
		}
	}
}

func (lock *FilterLock) Unlock(id int) {
	// 第id线程释放锁,下面这里可以理解为
	// id线程被挡在所有层外面,一个也没进去了
	lock.level[id] = 0
}

type Writer struct {
	lock   *FilterLock
	num    int
	waiter sync.WaitGroup
}

func (writer *Writer) increment(id int) {
	writer.lock.Lock(id)
	// fmt.Println(writer.num, id)
	writer.num++
	// fmt.Println(writer.num, id)
	// fmt.Println(writer.num, writer.lock.flags[id], "id: ", id)
	writer.lock.Unlock(id)
	writer.waiter.Done()
}

// 过滤锁算法满足互斥,无饥饿,无死锁
func TestFilterLock(t *testing.T) {
	chpt2.Mfence()
	var n int = 10
	for i := 0; i < 10000; i++ {
		writer := &Writer{
			lock: NewFilterLock(n),
			num:  0,
		}
		writer.waiter.Add(n)
		for i := 0; i < n; i++ {
			go writer.increment(i)
		}
		writer.waiter.Wait()
		// fmt.Println("----------")
		require.Equal(t, n, writer.num)
	}
}
