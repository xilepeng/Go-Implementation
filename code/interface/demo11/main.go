package main

import "fmt"

type Person interface {
	growUp()
}

type student struct {
	age int
}

func (p student) growUp() {
	p.age += 1
	return
}

// func (p *student) growUp() {
// 	p.age += 1
// 	return
// }

func main() {
	var mojo = Person(student{age: 18})
	// var mojo Person = &student{18}
	mojo.growUp()
	fmt.Println(mojo)
}
