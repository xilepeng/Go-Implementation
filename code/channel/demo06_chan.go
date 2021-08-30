package main

import (
	"fmt"
)

func main() {
	ch := make(chan int, 5)
	go sendData2(ch)
	for i := range ch { // i := <- ch
		fmt.Println("取出数据：", i)
	}
	fmt.Println("main is over...")
}
func sendData2(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
		fmt.Println("\t子goroutine写入数据i:", i)
	}
	close(ch)
}
