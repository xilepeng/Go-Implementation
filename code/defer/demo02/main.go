package main

import "fmt"

func a() {
	i := 0
	defer fmt.Println(i) // 0
	i++
	return
}

func main() {
	a()
}
