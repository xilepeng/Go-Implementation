package main

import (
	"errors"
	"fmt"
)
// errors.Unwrap 将嵌套的 error 解析出来，多层嵌套需要调用 Unwrap 函数多次，才能获取最里层的 error。
func main() {
	err1 := errors.New("error1")
	err2 := fmt.Errorf("error2: [%w]", err1)
	fmt.Println(err2)
	fmt.Println(errors.Unwrap(err2))
}

// Output
// error2: [error1]
// error1
