package main

import "fmt"

func main() {
	slice := make([]int, 1e6)
	fmt.Printf("slice pointer = %p\n", &slice)
	slice = foo(slice)
	fmt.Printf("slice pointer = %p\n", &slice)
}
func foo(slice []int) []int {
	return slice
}
