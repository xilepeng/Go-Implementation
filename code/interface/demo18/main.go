package main

import "fmt"

type MyError struct{}

func (i MyError) Error() string {
	return "MyError"
}

func process() error {
	var err *MyError = nil
	return err
}

func main() {
	err := process()
	fmt.Println(err)
	fmt.Println(err == nil)
	fmt.Printf("err: %T, %v\n", err, err)
}

// <nil>
// false
// err: *main.MyError, <nil>