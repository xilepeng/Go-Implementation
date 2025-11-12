package main

import "fmt"

func main() {
	// 分配包含 100 万个整型值的切片
	slice := make([]int, 1e6)
	fmt.Printf("slice pointer = %p\n", &slice)
	// 将 slice 传递到函数 foo
	slice = foo(slice)
	fmt.Printf("slice pointer = %p\n", &slice)
}

// 函数 foo 接收一个整型切片，并返回这个切片
func foo(slice []int) []int {
	return slice
}
