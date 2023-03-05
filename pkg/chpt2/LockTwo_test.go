package chpt2

import "testing"

type LockTwo struct {
	victim int
}

func (locktwo *LockTwo) Lock(id int) {
	locktwo.victim = id
	for locktwo.victim == id {
	}
}

func (locktwo *LockTwo) Unlock() {
	// do nothing
}

// locktwo 经不起测试，当一个线程完全先
// 于另外一个线程,就非常容易出现死锁的情况
func TestLockTwo(t *testing.T) {

}
