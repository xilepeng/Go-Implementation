package main

import "fmt"

func main() {
	// 创建字符串切片
	// 其长度和容量都是 5 个元素
	source := []string{"Apple", "Orange", "Plum", "Banana", "Grape"}
	// 对第三个元素做切片，并限制容量 // 其长度和容量都是 1 个元素
	slice := source[2:3:3]
	fmt.Printf("slice=%v, slice_addr=%p, len=%d, cap=%d\n", slice, &slice, len(slice), cap(slice))
	// 向 slice 追加新字符串
	newSlice := append(slice, "Kiwi")
	fmt.Printf("newSlice=%v, newSlice_addr=%p, len=%d, cap=%d\n", newSlice, &newSlice, len(newSlice), cap(newSlice))

}

// slice=[Plum], slice_addr=0xc00011a000, len=1, cap=1
// newSlice=[Plum Kiwi], newSlice_addr=0xc00011a030, len=2, cap=2
