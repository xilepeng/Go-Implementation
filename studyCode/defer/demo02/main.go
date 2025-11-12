package main

import "fmt"

func a() {
	i := 0
	defer fmt.Println(i) // 0
	i++
	return
}

// func b() {
// 	var j *int
// 	*j = 0
// 	defer fmt.Println(*j) // 0
// 	*j++
// 	return
// }

func main() {
	a()
	// b()
}
