package main

import (
	"fmt"
	"time"
)

func main() {
	chan1 := make(chan int)
	chan2 := make(chan int)

	writeFlag := false
	go func() {
		for {
			if writeFlag {
				chan1 <- 1
			}
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			if writeFlag {
				chan2 <- 1
			}
			time.Sleep(time.Second)
		}
	}()

	select {
	case <-chan1:
		fmt.Println("chan1 ready.")
	case <-chan2:
		fmt.Println("chan2 ready.")
	}

	fmt.Println("main exit.")
}

// 阻塞，无输出
