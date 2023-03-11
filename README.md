```text
    this is a repo of the art of multiprocess programming implemetations. We will
use go  program languages. There will be a rust repo for spec chapts.
```
## 第一章
```go
    实现了读者写者问题和生产者消费者的模拟程序
```
## 第二章
```go
    实现了LockOne和LockTwo以及两者的集合版本PeterSon算法,不过他们都是针对两个
线程的锁算法，前两者都是不完美的算法,而PeterSon则是满足无死锁,满足互斥,满足无饥
饿的完美算法.
    实现了FilterLock和BakeryLock算法,这两种算法都是针对多线程的锁算法,都是无死锁,
满足互斥,满足无饥饿.
    目前是实现有点问题，线程全部结束后,writer.num有时会出现比预期值小1的情况，这是
因为会出现指令重排序，可能出现类似的Lock()执行的时候跑到了writer.num的后面的情况,这
种情况下就导致并发写num了，这种是无法阻止的,结论就是软件是无法实现锁的，必须依赖硬件
    本章的其它内容也推荐好好读一读
```