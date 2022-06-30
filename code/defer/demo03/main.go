package main

import "fmt"

func d1() {
	for i := 3; i > 0; i-- {
		defer fmt.Print(i, " ") // 1 2 3
	}
	fmt.Println()
	return
}

// defer 执行没有带参数的匿名函数
func d2() {
	for i := 3; i > 0; i-- {
		defer func() {
			fmt.Print(i, " ") // 0 0 0
		}()
	}
	fmt.Println()
	return
}

// defer 执行带参数的匿名函数
func d3() {
	for i := 3; i > 0; i-- {
		defer func(n int) {
			fmt.Print(n, " ") // 1 2 3
		}(i)
	}
	fmt.Println()
}

func main() {
	d1()
	d2()
	d3()
}
