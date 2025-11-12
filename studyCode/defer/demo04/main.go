package main

import "fmt"

func deferFuncReturn() (result int) {
	i := 0
	defer func() {
		result++
	}()
	return i // 1
}

func main() {
	i := deferFuncReturn()
	fmt.Println(i)
}
