package main

import "fmt"

type Student struct {
	Name string
	Age  int
}

func main() {
	var i interface{} = new(Student)
	s, ok := i.(Student)
	if ok {
		fmt.Println(s)
	}
}
