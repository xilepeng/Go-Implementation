package main

import "fmt"

func main() {
	i := new(int)

	var v int
	j := &v
	fmt.Printf("i j 是相同类型: %v\n", i == j)
	fmt.Printf("i: %T    j: %T\n", i, j)
}

// i j 是相同类型: false
// i: *int    j: *int
