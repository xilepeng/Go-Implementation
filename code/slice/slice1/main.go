package main

import (
	"fmt"
)

func main() {
	var nil_slice []int
	new_slice := new([]int)
	empty_slice := make([]int, 0) // []int{}
	slice := []int{0, 1, 2}

	fmt.Printf("nil_slice: %p, %v, %T \n\t array: %p, len: %v, cap: %v\n\n",
		&nil_slice, nil_slice, nil_slice,
		nil_slice, len(nil_slice), cap(nil_slice))

	fmt.Printf("new_slice: %p, %v, %T \n\t array: %p, len: %v, cap: %v\n\n",
		&new_slice, new_slice, new_slice,
		new_slice, len(*new_slice), cap(*new_slice))

	fmt.Printf("empty_slice: %p, %v, %T \n\t array: %p, len: %v, cap: %v\n\n",
		&empty_slice, empty_slice, empty_slice,
		empty_slice, len(empty_slice), cap(empty_slice))

	fmt.Printf("slice: %p, %v, %T \n\t array: %p, len: %v, cap: %v\n\n",
		&slice, slice, slice,
		&slice[0], len(slice), cap(slice))

}
