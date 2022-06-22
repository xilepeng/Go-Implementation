package main

import "fmt"

func main() {
	var i int = 2
	var f float64

	f = float64(i)
	fmt.Printf("类型：%T, 值：%v\n", f, f)

	f = 5.2
	a := int(f)
	fmt.Printf("类型：%T, 值：%v\n", a, a)

	// cannot convert a (variable of type int) to type []int
	// s := []int(a)
}
