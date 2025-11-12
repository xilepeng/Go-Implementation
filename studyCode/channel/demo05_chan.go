package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)
	go sendData1(ch)
	for i := range ch { // i := <- ch
		fmt.Println("\t取出数据：", i)
	}
	fmt.Println("main is over...")
}
func sendData1(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)
}
