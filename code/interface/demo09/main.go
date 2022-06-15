package main

import "fmt"

type Person struct {
	age int
}

func (p Person) HowOld() int {
	return p.age
}

func (p *Person) GrowUp() {
	p.age++
}

func main() {
	// mojo 是值类型
	mojo := Person{age: 18}
	// 值类型 调用接收者也是值类型的方法
	fmt.Println(mojo.HowOld())
	// 值类型 调用接收者是指针类型的方法
	mojo.GrowUp()
	fmt.Println(mojo.HowOld())

	// ------------------------
	// mojo 是指针类型
	pointer_mojo := &Person{age: 100}
	// 指针类型 调用接收者也是值类型的方法
	fmt.Println(pointer_mojo.HowOld())
	// 指针类型 调用接收者是指针类型的方法
	pointer_mojo.GrowUp()
	fmt.Println(pointer_mojo.HowOld())

}

// 18
// 19
// 100
// 101
