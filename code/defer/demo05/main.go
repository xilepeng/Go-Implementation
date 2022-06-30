package main

import "fmt"

// 主函数拥有一个匿名的返回值，返回字面值
// return语句，直接把1写入栈中作为返回值，延迟函数无法操作该返回值，所以就无法影响返回值。
func foo() int {
	i := 0
	defer func() {
		i++
	}()
	return 1 // 1
}

//  主函数拥有匿名返回值，返回变量
// defer语句可以引用到返回值，但不会改变返回值。
func foo2() int {
	i := 0
	defer func() {
		i++
	}()
	return i // 0
}

func foo3() (result int) {
	i := 0
	defer func() {
		result++
	}()
	return i // 1
}

func main() {
	fmt.Println(foo())
	fmt.Println(foo2())
	fmt.Println(foo3())
}
