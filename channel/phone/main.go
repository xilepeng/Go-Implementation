package main

import "fmt"

func yiyi(ch chan string) {
	msg := <-ch
	fmt.Println("x said: ", msg)

	ch <- "hello x"
}

func main() {
	phone := make(chan string)
	defer close(phone)

	go yiyi(phone) // goroutine 与主函数并发执行

	phone <- "hello yiyi"
	msg := <-phone
	fmt.Println("yiyi said: ", msg)
}

//x said:  hello yiyi
//yiyi said:  hello x
