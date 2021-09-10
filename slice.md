
### 以下代码基于 Go 1.17


## ✏1. Array 数组

```go
package main

import "fmt"

func main() {
	Array1 := [5]int{}
	Array2 := [2]int{}
	fmt.Printf("Array1 的数据类型：%T\nArray2 的数据类型：%T\n", Array1, Array2)
}

// Array1 的数据类型：[5]int
// Array2 的数据类型：[2]int
```

### 📝Array 和 Slice 有什么区别？

1. **Slice 的底层数据是数组, Slice 是对底层数组进行了封装，它描述一个数组的片段。**
两者都可以通过下标来访问单个元素。

2. **数组是定长的**，长度定义好之后，不能再更改。在 Go 中，数组是不常见的，因为其长度是类型的一部分，限制了它的表达能力，比如 [3]int 和 [4]int 就是不同的类型。
**而切片则非常灵活，它可以动态地扩容。切片的类型和长度无关。**

3. **Array 是值类型**，赋值和函数传参操作都会复制整个数组数据。
**而 Slice 是引用类型。**

4. 数组是一片连续的内存， 而slice 实际上是一个结构体，包含三个字段：底层数组、长度、容量。


## 📝 2. Slice 基本数据结构

 slice 的底层源码和相关实现在 src/runtime/slice.go

```go
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```

 




练习

```go
package main

import "fmt"

func main() {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := slice[2:5]
	s2 := s1[2:6:7]
	fmt.Printf("s1: len=%v, cap=%v\n", len(s1), cap(s1))
	fmt.Printf("s2: len=%v, cap=%v\n", len(s2), cap(s2))

	s2 = append(s2, 100)
	s2 = append(s2, 200)
	fmt.Printf("\t扩容后 s2: len=%v, cap=%v\n", len(s2), cap(s2))

	s1[2] = 20

	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(slice)
}

// s1: len=3, cap=8
// s2: len=4, cap=5
//         扩容后 s2: len=6, cap=10
// [2 3 20]
// [4 5 6 7 100 200]
// [0 1 2 3 20 5 6 7 100 9]

```