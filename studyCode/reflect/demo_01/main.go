package main

import (
	"fmt"
	"reflect"
)

type XInt int
type YInt int

func main() {
	x := XInt(1)
	y := YInt(1)
	fmt.Println(reflect.DeepEqual(x, y)) // false
}
