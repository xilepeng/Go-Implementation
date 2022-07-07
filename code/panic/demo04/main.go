package main

import "fmt"

func main() {
	defer fmt.Println("in main")
	// 失效的崩溃恢复
	if err := recover(); err != nil {
		fmt.Println(err)
	}
	panic("unknow err") // panic 在 recover 后，
}

// 程序没有正常退出
// in main
// panic: unknow err
// goroutine 1 [running]:
