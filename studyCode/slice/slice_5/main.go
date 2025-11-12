package main

import "fmt"

func main() {

	// 创建一个整型切片
	// 其长度和容量都是 5 个元素
	slice := []int{10, 20, 30, 40, 50}
	// 创建一个新切片
	// 其长度是 2 个元素，容量是 4 个元素
	newSlice := slice[1:3]
	// 修改 newSlice 索引为 1 的元素
	// 同时也修改了原来的 slice 的索引为 2 的元素
	fmt.Printf("slice = %v\n", slice)
	fmt.Printf("newSlice = %v\n", newSlice)
	newSlice[1] = 35
	fmt.Println("修改 newSlice 索引为 1 的元素,同时也修改了原来的slice")
	fmt.Printf("slice = %v\n", slice)
	fmt.Printf("newSlice = %v\n", newSlice)

	// 修改 newSlice 索引为 3 的元素
	// 这个元素对于 newSlice 来说并不存在
	// newSlice[3] = 45
	// panic: runtime error: index out of range [3] with length 2

	// 使用原有的容量来分配一个新元素
	// 将新元素赋值为 60
	newSlice = append(newSlice, 60)
	fmt.Printf("slice = %v\n", slice)
	fmt.Printf("newSlice = %v\n", newSlice)
}
