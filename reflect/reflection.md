

1. [Go 语言中反射有哪些应用 ?](#go-语言中反射有哪些应用-)
2. [什么情况下需要使用反射 ?](#什么情况下需要使用反射-)
3. [如何比较两个对象完全相同 ？](#如何比较两个对象完全相同-)
4. [什么是反射 ?](#什么是反射-)
5. [Go 语言如何实现反射 ?](#go-语言如何实现反射-)

## Go 语言中反射有哪些应用 ?

Go 语言中反射的应用非常广：IDE 中的代码自动补全功能、对象序列化（encoding/json）、fmt 相关函数的实现、ORM（全称是：Object Relational Mapping，对象关系映射）……

## 什么情况下需要使用反射 ?

使用反射的常见场景有以下两种：
1. 不能明确接口调用哪个函数，需要根据传入的参数在运行时决定。
2. 不能明确传入函数的参数类型，需要在运行时处理任意对象。


【引申1】不推荐使用反射的理由有哪些？

1. 与反射相关的代码，经常是难以阅读的。在软件工程中，代码可读性也是一个非常重要的指标。
2. Go 语言作为一门静态语言，编码过程中，编译器能提前发现一些类型错误，但是对于反射代码是无能为力的。所以包含反射相关的代码，很可能会运行很久，才会出错，这时候经常是直接 panic，可能会造成严重的后果。
3. 反射对性能影响还是比较大的，比正常代码运行速度慢一到两个数量级。所以，对于一个项目中处于运行效率关键位置的代码，尽量避免使用反射特性。


## 如何比较两个对象完全相同 ？

Go 语言中提供了一个函数可以完成此项功能：

```go
func DeepEqual(x, y interface{}) bool
```

DeepEqual 函数的参数是两个 interface，实际上也就是可以输入任意类型，输出 true 或者 flase 表示输入的两个变量是否是“深度”相等。
先明白一点，如果是不同的类型，即使是底层类型相同，相应的值也相同，那么两者也不是“深度”相等。

```go
type XInt int
type YInt int

func main() {
	x := XInt(1)
	y := YInt(1)
	fmt.Println(reflect.DeepEqual(x, y)) // false
}
```
上面的代码中，x, y 底层都是 int，而且值都是 1，但是两者静态类型不同，前者是 XInt，后者是 YInt，因此两者不是“深度”相等。




## 什么是反射 ?

在计算机学中，反射式编程 reflective programming 或反射 reflection，是指计算机程序在运行时 runtime 可以访问、检测和修改它本身状态或行为的一种能力。用比喻来说，反射就是程序在运行的时候能够“观察”并且修改自己的行为。

《Go 语言圣经》中是这样定义反射的：
Go 语言提供了一种机制在运行时更新变量和检查它们的值、调用它们的方法，但是在编译时并不知道这些变量的具体类型，这称为反射机制。



## Go 语言如何实现反射 ?





```go
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
```





参考内容：

[码农桃花源](https://qcrao91.gitbook.io/go/fan-she/shi-mo-qing-kuang-xia-xu-yao-shi-yong-fan-she)