package chpt1

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
经典互斥问题实现2:

读者写者问题:
一个容量为2的字符串切片
实现5个读线程和5个写线程的并发访问该切片

怎么不用go语言提供的mutex包来实现互斥(automic等这些语言提供的方案不使用)?
*/

type Buffer struct {
	mutex  sync.RWMutex
	buffer []string
}

type Reader struct {
	ch <-chan int
	id int
}

type Writer struct {
	ch <-chan int
	id int
}

var data *Buffer

func init() {
	data = &Buffer{
		buffer: make([]string, 0),
	}
}

func run() {
	writers := make([]*Writer, 0)
	readers := make([]*Reader, 0)
	chans := make([]chan int, 0)
	for i := 0; i < 5; i++ {
		chans = append(chans, make(chan int))
		chans = append(chans, make(chan int))
		writers = append(writers, &Writer{
			id: i,
			ch: chans[i*2],
		})
		readers = append(readers, &Reader{
			id: i,
			ch: chans[i*2+1],
		})
	}

	for i := 0; i < 5; i++ {
		go func(idx int) {
			for {
				select {
				case <-writers[idx].ch:
					fmt.Println("writer ", idx, " exit ")
					return
				default:
					writers[idx].Write()
				}
			}
		}(i)

		go func(idx int) {
			for {
				select {
				case <-readers[idx].ch:
					fmt.Println("reader ", idx, " exit ")
					return
				default:
					readers[idx].Read()
				}

			}
		}(i)
	}
	time.Sleep(2 * time.Millisecond)
	for i := 0; i < 10; i++ {
		chans[i] <- 1
	}
	time.Sleep(2 * time.Second)
}

func TestReadWrite(t *testing.T) {
	run()
}

func (writer *Writer) Write() {
	data.mutex.Lock()
	defer data.mutex.Unlock()
	data.buffer = data.buffer[:0]
	data.buffer = append(data.buffer, "write")
	data.buffer = append(data.buffer, fmt.Sprintf("%d", writer.id))
	fmt.Println("Write ", writer.id)
}

func (reader *Reader) Read() {
	data.mutex.RLock()
	defer data.mutex.RUnlock()
	if len(data.buffer) > 0 {
		fmt.Println("reader", reader.id, " is reading ", data.buffer)
	}
}
