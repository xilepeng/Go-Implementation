package main

import "fmt"

func main() {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)

	for v := range ch {
		fmt.Println(v)
	}
	// 向已关闭的channel写入数据导致panic
	ch <- 42
}

/*
0
1
2
3
4
panic: send on closed channel
*/
