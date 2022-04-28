package main

import "fmt"

func main() {
	s := "Hello"
	// s[0] = 'h' // cannot assign to s[0] (value of type byte)compilerUnassignableOperand

	s = "hello" // 1. 给变量整体赋新值，地址指向新内容，并没有修改原来的内存
	fmt.Printf("%c\n", s[0])
	for i := 0; i < len(s); i++ {
		fmt.Printf("%x\n", s[i])
	}

	// 2. 将变量强制转变为字节slice,这样会为slice变量重新分配一段内存
	b := ([]byte)(s)
	b[0] = 'H'
	fmt.Printf("%c\n", b[0])
}
