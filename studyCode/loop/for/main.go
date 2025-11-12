package main

import (
	"fmt"
)

func main() {
	for i := 0; i < 10; i++ {
		if i%2 == 0 { // 跳过
			continue
		}
		if i == 8 {
			break // 中断for循环
		}
		fmt.Print(i, " ")
	}
	fmt.Println()

	// 模拟while
	i := 10
	for {
		if i < 0 {
			break
		}
		fmt.Print(i, " ")
		i--
	}
	fmt.Println()
}
