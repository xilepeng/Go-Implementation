package main

import "fmt"

func main() {
	var myIint interface{} = 123
	k, ok := myIint.(int)
	if ok {
		fmt.Println("Success:", k)
	}
	v, ok := myIint.(float64)
	if ok {
		fmt.Println(v)
	} else {
		fmt.Println("Failed without panicking!")
	}

	i := myIint.(int)
	fmt.Println("No checking", i)
	j := myIint.(bool)
	fmt.Println(j)
}
