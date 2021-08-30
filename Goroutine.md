

## Goroutine调度器



goroutine和线程的区别 
Runtime 和 OS 的关系

什么是M:N模型 
GPM 是什么 
g0 栈何用户栈如何切换 
goroutine 如何退出 
goroutine 调度时机有哪些 
什么是workstealing 

M 如何找工作 
mian gorutine 如何创建 
schedule 循环如何启动 
schedule 循环如何运转 
sysmon 后台监控线程做了什么 
一个调度相关的陷阱 
什么是 go shceduler 


描述 scheduler 的初始化过程




- G: goroutine
- M: OS thread (machine)
- P: processor （Go中定义的一个摡念，不是指CPU），包含运行Go代码的必要资源，也有调度goroutine的能力。



```go
package main

import "fmt"

func main() {
	/*
		一个 goroutine 打印数字、另一个 goroutine 打印字母，观察打印结果
	*/
	//1.创建并启动子 goroutine, 执行printNum()
	go printNum()
	//2.主 goroutine 打印字母
	for i := 0; i < 10; i++ {
		fmt.Printf("主 goroutine 打印字母 X %d\n", i)
	}
	fmt.Println("main is over...")

}

func printNum() {
	for i := 0; i < 10; i++ {
		fmt.Printf("\t子 goroutine 中打印数字%d\n", i)
	}
}

```

