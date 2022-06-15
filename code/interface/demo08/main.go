package main

import "fmt"

// 定义一个接口，和使用此接口作为参数的函数
type IGreeting interface {
	sayHello()
}

func sayHello(i IGreeting) {
	i.sayHello()
}

// 定义 Go 类型
type Go struct{}

// sayHello 使用值接收者实现了一个方法
func (g Go) sayHello() {
	fmt.Println("Hi, I am Go")
}

// 方法能给用户定义的类型添加新的行为。方法实际上也是函数，只是在声明时，在关键字 func 和方法名之间增加了一个参数
type CPlus struct{}

func (c CPlus) sayHello() {
	fmt.Println("Hi, I am CPlus")
}

func main() {
	Golang := Go{}
	CPlus := CPlus{}
	sayHello(Golang)
	sayHello(CPlus)
}

// Hi, I am Go
// Hi, I am CPlus
