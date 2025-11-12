package main

import (
	"fmt"
)

type Student struct {
	Name string
	Age  int
}

func judge(v interface{}) {
	fmt.Printf("%p, %v\n", &v, v)
	switch v := v.(type) {
	case nil:
		fmt.Printf("%p %v\n", &v, v)
		fmt.Printf("nil type[%T] %v\n", v, v)
	case Student:
		fmt.Printf("%p %v\n", &v, v)
		fmt.Printf("Student type[%T] %v\n", v, v)

	case *Student:
		fmt.Printf("%p %v\n", &v, v)
		fmt.Printf("*Student type[%T] %v\n", v, v)

	default:
		fmt.Printf("%p %v\n", &v, v)
		fmt.Printf("unknow\n")
	}
}

func main() {
	// var i interface{}
	// var i interface{} = new(Student)
	var i interface{} = (*Student)(nil)
	judge(i)
}

// 0xc000096210, <nil>
// 0xc000096220 <nil>
// nil type[<nil>] <nil>

// 0xc000096210, &{ 0}
// 0xc0000ac020 &{ 0}
// *Student type[*main.Student] &{ 0}

// 0xc000096210, <nil>
// 0xc0000ac020 <nil>
// *Student type[*main.Student] <nil>
