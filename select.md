

select 是操作系统中的系统调用，我们经常会使用 select、poll 和 epoll 等函数构建 I/O 多路复用模型提升程序的性能。Go 语言的 select 与操作系统中的 select 比较相似，本节会介绍 Go 语言 select 关键字常见的现象、数据结构以及实现原理。

C 语言的 select 系统调用可以同时监听多个文件描述符的可读或者可写的状态，Go 语言中的 select 也能够让 Goroutine 同时等待多个 Channel 可读或者可写，在多个文件或者 Channel状态改变之前，select 会一直阻塞当前线程或者 Goroutine。




我们简单总结一下 select 结构的执行过程与实现原理，首先在编译期间，Go 语言会对 select 语句进行优化，它会根据 select 中 case 的不同选择不同的优化路径：

1. 空的 select 语句会被转换成调用 runtime.block 直接挂起当前 Goroutine；
2. 如果 select 语句中只包含一个 case，编译器会将其转换成 if ch == nil { block }; n; 表达式；
    - 首先判断操作的 Channel 是不是空的；
    - 然后执行 case 结构中的内容；
3. 如果 select 语句中只包含两个 case 并且其中一个是 default，那么会使用 runtime.selectnbrecv 和 runtime.selectnbsend 非阻塞地执行收发操作；
4. 在默认情况下会通过 runtime.selectgo 获取执行 case 的索引，并通过多个 if 语句执行对应 case 中的代码；


在编译器已经对 select 语句进行优化之后，Go 语言会在运行时执行编译期间展开的 runtime.selectgo 函数，该函数会按照以下的流程执行：

1. 随机生成一个遍历的轮询顺序 pollOrder 并根据 Channel 地址生成锁定顺序 lockOrder；
2. 根据 pollOrder 遍历所有的 case 查看是否有可以立刻处理的 Channel；
    - 如果存在，直接获取 case 对应的索引并返回；
    - 如果不存在，创建 runtime.sudog 结构体，将当前 Goroutine 加入到所有相关 Channel 的收发队列，并调用 runtime.gopark 挂起当前 Goroutine 等待调度器的唤醒；
3. 当调度器唤醒当前 Goroutine 时，会再次按照 lockOrder 遍历所有的 case，从中查找需要被处理的 runtime.sudog 对应的索引；
select 关键字是 Go 语言特有的控制结构，它的实现原理比较复杂，需要编译器和运行时函数的通力合作。




参考：

[《Go 语言设计与实现》](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-select/#52-select)




- select语句中除default外，每个case操作一个channel，要么读要么写
- select语句中除default外，各case执行顺序是随机的
- select语句中如果没有default语句，则会阻塞等待任一case
- select语句中读操作要判断是否成功读取，关闭的channel也可以读取

