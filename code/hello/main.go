package main

import "fmt"

func main() {
	/*
		一个 goroutine 打印数字、另一个 goroutine 打印字母，观察打印结果
	*/
	//1.创建并启动子 goroutine, 执行printNum()
	go printNum()
	//2.主 goroutine 打印字母
	for i := 0; i < 100; i++ {
		fmt.Printf("主 goroutine 打印字母 X %d\n", i)
	}
	fmt.Println("main is over...")

}

func printNum() {
	for i := 0; i < 100; i++ {
		fmt.Printf("\t子 goroutine 中打印数字%d\n", i)
	}
}
