package main

import "reflect"

func Foo() {

}

func main() {
	// Without reflection
	f := Foo{}
	f.Hello()

	// With reflection

	fT := reflect.TypeOf(Foo{})
	fV := reflect.New(fT)

	m := fV.MethodByName("Hello")
	if m.IsValid() {
		m.Call(nil)
	}

}

