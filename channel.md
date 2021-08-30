
# 深入 Go 并发原语 — Channel 底层实现

1. 什么是 CSP
2. channel 底层的数据结构是什么 
3. channel 发送和接收元素的本质是什么 
4. 从 channel 接收数据的过程是怎样的 
5. channel 在什么情况下会引起资源泄漏 
6. channel 有哪些应用 
7. 从一个关闭的 channel 仍然能读出数据吗 
8. 关于 channel 的 happened-before 有哪些 
9. 关闭一个 channel 的过程是怎样的 
10. 向 channel 发送数据的过程是怎样的 
11. 如何优雅地关闭 channel 
12. 操作 channel 的情况总结
13. 手写 worker pool（goroutine池）

## 1. 什么是 CSP ？

Go 语言的并发同步模型来自一个叫作通信顺序进程（Communicating Sequential Processes，CSP） 的范型（paradigm）。
CSP 是一种消息传递模型，通过在 goroutine 之间传递数据来传递消息，而不是对数据进行加锁来实现同步访问。
用于在 goroutine 之间同步和传递数据的关键数据类型叫作通道（channel）。


Do not communicate by sharing memory; instead, share memory by communicating.
## 不要通过共享内存进行通信。建议，通过通信来共享内存。

### 当一个资源需要在 goroutine 之间共享时，通道在 goroutine 之间架起了一个管道，并提供了确保同步交换数据的机制。

声明通道时，需要指定将要被共享的数据的类型。可以通过通道共享内置类型、命名类型、结构类型和引用类型的值或者指针。

#### 不要使用 sync 同步包组件实现并发编程、而是使用 channel 实现并发编程

Go 语言中，要传递某个数据给另一个goroutine(协程)，可以把这个数据封装成一个对象，然后把这个对象的指针传入某个 channel 中，另一个 goroutine 从这个 channel 中读出这个指针，并处理其指向的内存对象。Go从语言层面保证同一时间只有一个 goroutine 能访问 channel 中的数据，为开发者提供一种优雅简单的工具，所以 Go 的做法就是使用 channel 来通信，通过通信来传递内存数据，使得内存数据在不同 goroutine 中传递，而不是使用共享内存来通信。





### 1. 推荐使用 sync 包的 2 种情况：

- 对性能要求极高的临界区
- 保护某个结构内部状态和完整性
关于保护某个结构内部的状态和完整性。例如 Go 源码中如下代码：

```go
var sum struct {
	sync.Mutex
	i int
}

//export Add
func Add(x int) {
	defer func() {
		recover()
	}()
	sum.Lock()
	sum.i += x
	sum.Unlock()
	var p *int
	*p = 2
}
```
sum 这个结构体不想将内部的变量暴露在结构体之外，所以使用 sync.Mutex 来保护线程安全。

### 推荐使用 channel 的 2 种情况：

- 输出数据给其他使用方
- 组合多个逻辑
输出数据给其他使用方的目的是转移数据的使用权。并发安全的实质是保证同时只有一个并发上下文拥有数据的所有权。channel 可以很方便的将数据所有权转给其他使用方。另一个优势是组合型。如果使用 sync 里面的锁，想实现组合多个逻辑并且保证并发安全，是比较困难的。但是使用 channel + select 实现组合逻辑实在太方便了。以上就是 CSP 的基本概念和何时选择 channel 的时机。下一章从 channel 基本数据结构开始详细分析 channel 底层源码实现。


# 以下代码基于 Go 1.17

## 2. channel 底层的数据结构是什么 ？
### 二. 基本数据结构
channel 的底层源码和相关实现在 src/runtime/chan.go 中。

```go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}
```

## 3. channel 发送和接收元素的本质是什么 

## 4. 从 channel 接收数据的过程是怎样的 

## 5. channel 在什么情况下会引起资源泄漏 

## 6. channel 有哪些应用 

## 7. 从一个关闭的 channel 仍然能读出数据吗 

## 8. 关于 channel 的 happened-before 有哪些 

## 9. 关闭一个 channel 的过程是怎样的 

## 10. 向 channel 发送数据的过程是怎样的 

## 11. 如何优雅地关闭 channel 

## 12. 操作 channel 的情况总结

## 13. 手写 worker pool（goroutine池）



