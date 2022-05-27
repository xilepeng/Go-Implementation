package main

import "fmt"

type TestStruct struct{}

func NilOrNot(v interface{}) bool {
	return v == nil
}

func main() {
	var s *TestStruct
	fmt.Println(s == nil)
	fmt.Println(NilOrNot(s))
}

// ➜  demo01 git:(main) ✗ go run main.go
// true
// false