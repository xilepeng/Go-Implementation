package main

import "fmt"

func main() {
	ch := make(chan int)
	go sendData(ch)
	for {
		v, ok := <-ch
		if !ok {
			fmt.Println("数据读取结束")
			break
		}
		fmt.Println("读取数据v:", v)
	}
	fmt.Println("main is over...")
}

func sendData(c chan int) {
	for i := 0; i < 10; i++ {
		c <- i
	}
	fmt.Println("\t子goroutine结束")
	close(c)
}
