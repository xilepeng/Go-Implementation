package main

import "fmt"

func deferFuncParameter() {
	var aInt = 1
	defer fmt.Println(aInt)
	aInt = 2
	return
}

func main() {
	deferFuncParameter() // 1
}
