package main

import "fmt"

func main() {
	// 没有消息被打印（运行 goroutine 之前应用程序已经退出）
	go fmt.Println("hello world")
}
