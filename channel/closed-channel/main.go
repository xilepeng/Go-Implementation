package main

import (
	"fmt"
	"time"
)

func listener(i int, quit <-chan struct{}) {
	fmt.Printf("listener %d is waiting\n", i)

	<-quit // 阻塞通道直到关闭

	fmt.Printf("listener %d is exiting\n", i)
}

func main() {
	quit := make(chan struct{})

	for i := 0; i < 5; i++ {
		go listener(i, quit)

	}
	time.Sleep(2 * time.Second) // 允许监听器开始
	fmt.Println("closing the quit channel")
	close(quit)                 // 关闭通道，向监听器发出退出信号
	time.Sleep(1 * time.Second) // 允许监听器退出
}

//listener 1 is waiting
//listener 4 is waiting
//listener 2 is waiting
//listener 3 is waiting
//listener 0 is waiting
//closing the quit channel
//listener 1 is exiting
//listener 0 is exiting
//listener 3 is exiting
//listener 2 is exiting
//listener 4 is exiting
