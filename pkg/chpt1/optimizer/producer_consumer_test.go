package chpt1

import (
	"fmt"
	"math/rand"
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
这里是生产者消费者的优化,粒度更低的方式
*/

type Bucket struct {
	data  string
	empty bool
	mutex sync.Mutex
}

type Buffers struct {
	mutex    sync.Mutex
	buffers  []*Bucket
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
		buffers:  make([]*Bucket, 10),
		capcaity: 10,
	}
	for i := 0; i < 10; i++ {
		datas.buffers[i] = &Bucket{
			empty: true,
		}
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
					// fmt.Println("producer ", idx, " exit ")
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
					// fmt.Println("consumer ", idx, " exit ")
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
	var flag bool = false
	for !flag {
		idx := rand.Int31() % int32(datas.capcaity)
		datas.buffers[idx].mutex.Lock()
		if datas.buffers[idx].empty {
			datas.buffers[idx].empty = false
			datas.buffers[idx].data = fmt.Sprintf("produce%d", producer.id)
			// 0 代表写
			fmt.Println(0)
			flag = true
		}
		datas.buffers[idx].mutex.Unlock()
	}
}

func (consumer *Consumer) Read() {
	var flag bool = false
	for !flag {
		idx := rand.Int31() % int32(datas.capcaity)
		datas.buffers[idx].mutex.Lock()
		if !datas.buffers[idx].empty {
			datas.buffers[idx].empty = true
			// 1 代表读
			fmt.Println(1)
			flag = true
		}
		datas.buffers[idx].mutex.Unlock()
	}
}
