package main

import "fmt"

func main() {
	ch1 := make(chan int)
	// ch2 := make(chan<- int) //单向，只能写不能读
	// ch3 := make(<- chan int)//单向，只能读不能写
	go write(ch1)
	data := <-ch1
	fmt.Println("write 写入的数据是", data)
}

func write(ch chan<- int) {
	ch <- 100
	fmt.Println("\t子goroutine写入100")
}
