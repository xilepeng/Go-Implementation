package main

import "fmt"

func main() {
	/*
		channel 通道
	*/
	var c chan int
	fmt.Printf("%T, %v\n", c, c)
	if c == nil {
		fmt.Println("channel 是nil的,不能使用，需要先创建通道... ")
		c = make(chan int)
		fmt.Println(c)
	}
	test(c) //channel 是引用类型数据、传地址
}

func test(ch chan int) {
	fmt.Printf("%T, %v\n", ch, ch)
}
