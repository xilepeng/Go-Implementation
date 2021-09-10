package main

import "fmt"

func main() {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := slice[2:5]
	s2 := s1[2:6:7]
	fmt.Printf("s1: len=%v, cap=%v\n", len(s1), cap(s1))
	fmt.Printf("s2: len=%v, cap=%v\n", len(s2), cap(s2))

	s2 = append(s2, 100)
	s2 = append(s2, 200)
	fmt.Printf("\t扩容后 s2: len=%v, cap=%v\n", len(s2), cap(s2))

	s1[2] = 20

	fmt.Println("slice = ", slice)
	fmt.Println("s1 = ", s1)
	fmt.Println("s2 = ", s2)

}

// s1: len=3, cap=8
// s2: len=4, cap=5
//         扩容后 s2: len=6, cap=10
// slice =  [0 1 2 3 20 5 6 7 100 9]
// s1 =  [2 3 20]
// s2 =  [4 5 6 7 100 200]
