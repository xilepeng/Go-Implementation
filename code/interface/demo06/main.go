package main

import "fmt"

// 定义 Person 接口
type Person interface {
	job()
	growUp()
}

// 定义结构体
type Student struct {
	age int
}
type Coder struct {
	age int
}

func main() {
	// 定义 Student 对象
	mojo := Student{age: 18}
	WhatJob(&mojo)
	growUp(&mojo)
	fmt.Println(mojo)

	// 定义 Coder 对象
	hfbpw := Coder{age: 24}
	WhatJob(hfbpw)
	growUp(hfbpw)
	fmt.Println(hfbpw)

}

// 定义函数参数为 Person 的2个函数
func WhatJob(p Person) {
	p.job()
}
func growUp(p Person) {
	p.growUp()
}

// Student 类型没有实现接口
func (p Student) job() {
	fmt.Println("I am a Student")
	return
}

// *Student 类型实现了接口
func (p *Student) growUp() {
	p.age += 1
	return
}

// Coder 类型实现了接口
func (p Coder) job() {
	fmt.Println("I am a Coder")
	return
}
func (p Coder) growUp() {
	p.age += 10
}
