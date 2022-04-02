package main

import "fmt"

func main() {
	var m map[string]int
	var n map[string]int

	fmt.Println(m == nil)
	fmt.Println(n == nil)

	// 不能通过编译
	// fmt.Println(m == n)
    // ./main2.go:13:14: invalid operation: m == n (map can only be compared to nil)
}
