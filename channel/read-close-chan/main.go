package main

import (
	"fmt"
	"time"
)

func listener(ch <-chan int) {
	for {
		i, ok := <-ch
		if !ok {
			fmt.Println("listener: channel closed")
			return
		}
		fmt.Println("listener: channel received:", i)

	}
}
func main() {
	c := make(chan int)
	c <- 0
	close(c)

	go listener(c)

	time.Sleep(5 * time.Second)
}
