package main

import "fmt"

type coder interface {
	code()
	debug()
}

type Gopher struct {
	language string
}

// 实现了接收者是值类型的方法，相当于自动实现了接收者是指针类型的方法
func (g Gopher) code() {
	fmt.Printf("I am Coding %s \n", g.language)
}

// 指针类型的接收者，不会自动生成对应值类型的方法
func (g *Gopher) debug() {
	fmt.Printf("I am Coding %s \n", g.language)
}

func main() {
	var c coder = &Gopher{"Go"}
	c.code()
	c.debug()
}

// var c coder = Gopher{"Go"}

// # command-line-arguments
// ./main.go:23:16: cannot use Gopher{…} (value of type Gopher) as type coder in variable declaration:
//         Gopher does not implement coder (debug method has pointer receiver)
