package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64
	p := reflect.ValueOf(&x)
	fmt.Println("type of p: ", p.Type())         // type of p:  *float64
	fmt.Println("settability of p:", p.CanSet()) // settability of p: false

	v := p.Elem()
	v.SetFloat(5.2)
	fmt.Println(v.Interface()) // 5.2
	fmt.Println(x)             // 5.2
}
