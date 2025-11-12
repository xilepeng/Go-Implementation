package main

import (
	"fmt"
)

func main() {
	var nil_slice []int
	new_slice := new([]int)

	pointer_slice := new([]string)
	*pointer_slice = append(*pointer_slice, "x")

	fmt.Printf("*pointer_slice = %v\n\n ", *pointer_slice) // *pointer_slice = [x]

	empty_slice := make([]int, 0) // empty_slice := []int{}
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

// nil_slice: 0xc00000c030, [], []int
// 	 array: 0x0, len: 0, cap: 0

// new_slice: 0xc00000e028, &[], *[]int
// 	 array: 0xc00000c048, len: 0, cap: 0

// empty_slice: 0xc00000c060, [], []int
// 	 array: 0x116ce80, len: 0, cap: 0

// slice: 0xc00000c078, [0 1 2], []int
// 	 array: 0xc00001e0d8, len: 3, cap: 3
