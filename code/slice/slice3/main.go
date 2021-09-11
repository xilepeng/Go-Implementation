package main

import "fmt"

func main() {
	slice := []int{1, 2}
	slice = append(slice, 3, 4, 5)
	fmt.Printf("len = %d, cap = %d\n", len(slice), cap(slice))
}

// len = 5, cap = 6
