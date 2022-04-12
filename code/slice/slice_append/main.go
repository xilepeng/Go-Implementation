package main

import "fmt"

func main() {
	// 创建一个整型切片
	// 其长度和容量都是 4 个元素
	slice := []int{10, 20, 30, 40}
	// 向切片追加一个新元素
	// 将新元素赋值为 50
	newSlice := append(slice, 50)
	// 当这个 append 操作完成后，newSlice 拥有一个全新的底层数组，这个数组的容量是原来的两倍
	fmt.Printf("slice=%v, slice_addr=%p, len=%d, cap=%d\n", slice, &slice, len(slice), cap(slice))
	fmt.Printf("newSlice=%v, newSlice_addr=%p, len=%d, cap=%d\n", newSlice, &newSlice, len(newSlice), cap(newSlice))


}


// ➜  slice_append git:(main) ✗ go run main.go
// slice=[10 20 30 40], slice_addr=0xc0000a4018, len=4, cap=4
// newSlice=[10 20 30 40 50], newSlice_addr=0xc0000a4030, len=5, cap=8

// 函数 append 会智能地处理底层数组的容量增长。
// 在切片的容量小于 256 个元素时，总是会成倍地增加容量。
// 一旦元素个数超过 256，容量的增长因子会设为 1.25，也就是会每次增加 25% 的容量。