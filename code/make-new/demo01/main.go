package main

import "fmt"

func main() {
	slice := make([]int, 10, 100)
	hash := make(map[int]bool, 10)
	ch := make(chan int, 10)

	fmt.Printf("slice 类型：%T, 值：%v,\n", slice, slice)
	fmt.Printf("hash  类型：%T, 值：%v\n", hash, hash)
	fmt.Printf("ch    类型：%T, 值：%v\n", ch, ch)
}


// slice 类型：[]int, 值：[0 0 0 0 0 0 0 0 0 0],
// hash  类型：map[int]bool, 值：map[]
// ch    类型：chan int, 值：0xc0000b80b0