package main

import "fmt"

func main() {
	// 1. 初始化 Student 对象指针
	// 2. 将 Student 对象指针转换成 interface
	var s Person = &Student{name: "mojo"}
	// 3. 调用 interface 的方法
	s.sayHello("everyone")
}

type Person interface {
	sayHello(name string) string
	sayGoodbye(name string) string
}

type Student struct {
	name string
}

// go:noinline
func (s *Student) sayHello(name string) string {
	return fmt.Sprintf("%v: Hello! %v, Nice to meet you.\n", s.name, name)
}

// go:noinline
func (s *Student) sayGoodbye(name string) string {
	return fmt.Sprintf("%v: Hi! %v, see you next time.\n", s.name, name)
}

// go tool compile -S -N -l main.go >main.s1 2>&1
