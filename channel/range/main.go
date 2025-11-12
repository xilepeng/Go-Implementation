package main

import (
	"fmt"
	"time"
)

// 通道上持续监听消息，直到通道关闭。当通道关闭时，range 循环停止迭代
func listener(ch chan int) {
	for i := range ch {
		fmt.Println("listener: i=", i)
	}
}

// select 语句允许一个 goroutine 监听多个通道并响应第一个准备好的通道
func operater(ch1, ch2 chan int, quit <-chan struct{}) {
	for ch1 != nil || ch2 != nil {
		select {
		case msg := <-ch1:
			fmt.Println("ch1: ", msg)
		case msg := <-ch2:
			fmt.Println("ch2: ", msg)
		case <-quit:
			fmt.Println("quit")
			break
			// 仅当其他情况都不匹配时，才会选择在 select 语句中使用 default
			//default:
			//	fmt.Println("default: quit")
		}
	}
}
func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	quit := make(chan struct{})

	go listener(ch1)
	go operater(ch1, ch2, quit)

	ch1 <- 1
	ch1 <- 1
	quit <- struct{}{}
	ch2 <- 2

	time.Sleep(2 * time.Second) // 等待 goroutine 完成
}
