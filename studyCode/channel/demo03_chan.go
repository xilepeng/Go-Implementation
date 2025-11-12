package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	done := make(chan bool)
	go func() {
		fmt.Println("\t子goroutine开始...")
		time.Sleep(3 * time.Second)
		data := <-ch
		fmt.Println("data:", data)
		done <- true
	}()
	time.Sleep(5 * time.Second)
	ch <- 100
	<-done
	fmt.Println("main is over...")
}
