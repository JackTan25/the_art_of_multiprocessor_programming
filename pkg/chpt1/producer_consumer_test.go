package chpt1

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
经典互斥问题实现1:

第一章提到的生成者-消费者问题的程序实现
一个容量为10的string切片:
5个生产者和5个消费者一共10个线程,完成对
这个切片的并发访问

怎么不用go语言提供的mutex包来实现互斥(automic等这些语言提供的方案不使用)?
下面这种大锁对整个buffers加锁的方式其实对并发性能影响很大,粒度应该小到对每一个Bucket

*/

type Buffers struct {
	mutex    sync.RWMutex
	buffers  []string
	capcaity int
}

type Consumer struct {
	ch <-chan int
	id int
}

type Producer struct {
	ch <-chan int
	id int
}

var datas *Buffers

func init() {
	datas = &Buffers{
		buffers:  make([]string, 0),
		capcaity: 10,
	}
}

func run2() {
	producers := make([]*Producer, 0)
	consumers := make([]*Consumer, 0)
	chans := make([]chan int, 0)
	for i := 0; i < 5; i++ {
		chans = append(chans, make(chan int))
		chans = append(chans, make(chan int))
		producers = append(producers, &Producer{
			id: i,
			ch: chans[i*2],
		})
		consumers = append(consumers, &Consumer{
			id: i,
			ch: chans[i*2+1],
		})
	}

	for i := 0; i < 5; i++ {
		go func(idx int) {
			for {
				select {
				case <-producers[idx].ch:
					fmt.Println("producer ", idx, " exit ")
					return
				default:
					producers[idx].Write()
				}
			}
		}(i)

		go func(idx int) {
			for {
				select {
				case <-consumers[idx].ch:
					fmt.Println("consumer ", idx, " exit ")
					return
				default:
					consumers[idx].Read()
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

func TestProducerAndConsumer(t *testing.T) {
	run2()
}

func (producer *Producer) Write() {
	datas.mutex.Lock()
	defer datas.mutex.Unlock()
	if len(datas.buffers) < datas.capcaity {
		datas.buffers = append(datas.buffers, fmt.Sprintf("write%d", producer.id))
		fmt.Println("Write ", producer.id)
	}
}

func (consumer *Consumer) Read() {
	datas.mutex.Lock()
	defer datas.mutex.Unlock()
	if len(datas.buffers) > 0 {
		d := len(datas.buffers)
		fmt.Println("consumer", consumer.id, " is reading ", datas.buffers[d-1])
		datas.buffers = datas.buffers[:d-1]
	}
}
