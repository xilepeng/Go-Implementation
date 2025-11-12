package main

import "fmt"

func a() {
	fmt.Println("Inside a()")
	defer func() {
		if c := recover(); c != nil { // 处理 b 函数 panic 情况
			fmt.Println("Recover inside a()!")
		}
	}()
	fmt.Println("About to call b()")
	b()
	// recover 可以中止 panic 造成的程序崩溃。它是一个只能在 defer 中发挥作用的函数，在其他作用域中调用不会发挥作用；
	fmt.Println("b() exited!") // 未执行
	fmt.Println("Exiting a()") // 未执行
}

func b() {
	fmt.Println("Inside b()")
	// panic 能够改变程序的控制流，调用 panic 后会立刻停止执行当前函数的剩余代码，
	// 并在当前 Goroutine 中递归执行调用方的 defer；
	panic("Panic in b()!")     // 未执行
	fmt.Println("Exiting b()") // 未执行
}

func main() {
	a()
	fmt.Println("main() ended")
}

// Inside a()
// About to call b()
// Inside b()
// Recover inside a()!
// main() ended
