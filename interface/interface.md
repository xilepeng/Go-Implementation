
1. [iface 和 eface 的区别是什么?](#iface-和-eface-的区别是什么)
2. [Go 接口与 C++ 接口有何异同？](#go-接口与-c-接口有何异同)
3. [如何用 interface 实现多态 ?](#如何用-interface-实现多态-)
4. [接口转换的原理 ?](#接口转换的原理-)

## iface 和 eface 的区别是什么?

iface 和 eface 都是 Go 中描述接口的底层结构体，区别在于 iface 描述的接口包含方法，而 eface 则是不包含任何方法的空接口：interface{}。

## Go 接口与 C++ 接口有何异同？

接口定义了一种规范，描述了类的行为和功能，而不做具体实现。

C++ 的接口是使用抽象类来实现的，如果类中至少有一个函数被声明为纯虚函数，则这个类就是抽象类。纯虚函数是通过在声明中使用 "= 0" 来指定的。

```c++
class Shape
{
   public:
      // 纯虚函数
      virtual double getArea() = 0;
   private:
      string name;      // 名称
};
```

设计抽象类的目的，是为了给其他类提供一个可以继承的适当的基类。抽象类不能被用于实例化对象，它只能作为接口使用。
派生类需要明确地声明它继承自基类，并且需要实现基类中所有的纯虚函数。

**C++ 定义接口的方式称为“侵入式”，而 Go 采用的是 “非侵入式”，不需要显式声明，只需要实现接口定义的函数，编译器自动会识别。**

C++ 和 Go 在定义接口方式上的不同，也导致了底层实现上的不同。C++ 通过虚函数表来实现基类调用派生类的函数；而 Go 通过 itab 中的 fun 字段来实现接口变量调用实体类型的函数。C++ 中的虚函数表是在编译期生成的；而 Go 的 itab 中的 fun 字段是在运行期间动态生成的。原因在于，Go 中实体类型可能会无意中实现 N 多接口，很多接口并不是本来需要的，所以不能为类型实现的所有接口都生成一个 itab， 这也是“非侵入式”带来的影响；这在 C++ 中是不存在的，因为派生需要显示声明它继承自哪个基类。


## 如何用 interface 实现多态 ?


```go
package main

import "fmt"

// 定义 Person 接口
type Person interface {
	job()
	growUp()
}

// 定义结构体
type Student struct {
	age int
}
type Coder struct {
	age int
}

func main() {
	// 定义 Student 对象
	mojo := Student{age: 18}
	WhatJob(&mojo)
	growUp(&mojo)
	fmt.Println(mojo)

	// 定义 Coder 对象
	hfbpw := Coder{age: 24}
	WhatJob(hfbpw)
	growUp(hfbpw)
	fmt.Println(hfbpw)

}

// 定义函数参数为 Person 的2个函数
func WhatJob(p Person) {
	p.job()
}
func growUp(p Person) {
	p.growUp()
}

// Student 类型没有实现接口
func (p Student) job() {
	fmt.Println("I am a Student")
	return
}

// *Student 类型实现了接口
func (p *Student) growUp() {
	p.age += 1
	return
}

// Coder 类型实现了接口
func (p Coder) job() {
	fmt.Println("I am a Coder")
	return
}
func (p Coder) growUp() {
	p.age += 10
}
```

```go
I am a Student
{19}
I am a Coder
{24}
```


main 函数里先生成 Student 和 Programmer 的对象，再将它们分别传入到函数 whatJob 和 growUp。函数中，直接调用接口函数，实际执行的时候是看最终传入的实体类型是什么，调用的是实体类型实现的函数。于是，不同对象针对同一消息就有多种表现，多态就实现了。
更深入一点来说的话，在函数 whatJob() 或者 growUp() 内部，接口 person 绑定了实体类型 *Student 或者 Programmer。根据前面分析的 iface 源码，这里会直接调用 fun 里保存的函数，类似于： s.tab->fun[0]，而因为 fun 数组里保存的是实体类型实现的函数，所以当函数传入不同的实体类型时，调用的实际上是不同的函数实现，从而实现多态。



## 接口转换的原理 ?








**nil 和 non-nil** 

我们可以通过一个例子理解Go **语言的接口类型不是任意类型** 这一句话，下面的代码在 main 函数中初始化了一个 *TestStruct 类型的变量，由于指针的零值是 nil，所以变量 s 在初始化之后也是 nil：

```go
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
```


我们简单总结一下上述代码执行的结果：

- 将上述变量与 nil 比较会返回 true；
- 将上述变量传入 NilOrNot 方法并与 nil 比较会返回 false；
  
出现上述现象的原因是 —— 调用 NilOrNot 函数时发生了隐式的类型转换，除了向方法传入参数之外，变量的赋值也会触发隐式类型转换。在类型转换时，*TestStruct 类型会转换成 interface{} 类型，转换后的变量不仅包含转换前的变量，还包含变量的类型信息 TestStruct，所以转换后的变量与 nil 不相等。




以下代码基于 Go 1.18.2


src/runtime/runtime2.go


**一. 数据结构**

 1. 非空 interface 数据结构

非空的 interface 初始化的底层数据结构是 iface，稍后在汇编代码中能验证这一点。
```go
type iface struct {
	tab  *itab
	data unsafe.Pointer
}
```

tab 中存放的是类型、方法等信息。data 指针指向的 iface 绑定对象的原始数据的副本。这里同样遵循 Go 的统一规则，值传递。tab 是 itab 类型的指针。

```go
// layout of Itab known to compilers
// allocated in non-garbage-collected memory
// Needs to be in sync with
// ../cmd/compile/internal/reflectdata/reflect.go:/^func.WriteTabs.
type itab struct {
	inter *interfacetype // inner 存的是 interface 自己的静态类型
	_type *_type        // _type 存的是 interface 对应具体对象的类型。
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}
```

itab 中包含 5 个字段。inner 存的是 interface 自己的静态类型。_type 存的是 interface 对应具体对象的类型。itab 中的 _type 和 iface 中的 data 能简要描述一个变量。_type 是这个变量对应的类型，data 是这个变量的值。这里的 hash 字段和 _type 中存的 hash 字段是完全一致的，这么做的目的是为了类型断言(下文会提到)。fun 是一个函数指针，它指向的是具体类型的函数方法。虽然这里只有一个函数指针，但是它可以调用很多方法。在这个指针对应内存地址的后面依次存储了多个方法，利用指针偏移便可以找到它们。

由于 Go 语言是强类型语言，编译时对每个变量的类型信息做强校验，所以每个类型的元信息要用一个结构体描述。再者 Go 的反射也是基于类型的元信息实现的。_type 就是所有类型最原始的元信息。

```go
// Needs to be in sync with ../cmd/link/internal/ld/decodesym.go:/^func.commonsize,
// ../cmd/compile/internal/gc/reflect.go:/^func.dcommontype and
// ../reflect/type.go:/^type.rtype.
// ../internal/reflectlite/type.go:/^type.rtype.
type _type struct {
	size       uintptr // 类型占用内存大小
	ptrdata    uintptr // 包含所有指针的内存前缀大小
	hash       uint32  // 类型 hash
	tflag      tflag   // 标记位，主要用于反射
	align      uint8   // 对齐字节信息
	fieldAlign uint8   // 当前结构字段的对齐字节数
	kind       uint8   // 基础类型枚举值
	equal func(unsafe.Pointer, unsafe.Pointer) bool // 比较两个形参对应对象的类型是否相等
	gcdata    *byte    // GC 类型的数据
	str       nameOff  // 类型名称字符串在二进制文件段中的偏移量
	ptrToThis typeOff  // 类型元信息指针在二进制文件段中的偏移量
}
```

- size 字段存储了类型占用的内存空间，为内存空间的分配提供信息；
- hash 字段能够帮助我们快速确定类型是否相等；
- equal 字段用于判断当前类型的多个对象是否相等，该字段是为了减少 Go 语言二进制包大小从 typeAlg 结构体中迁移过来的4；
  
我们只需要对 runtime._type 结构体中的字段有一个大体的概念，不需要详细理解所有字段的作用和意义。




/src/runtime/runtime2.go 

Go 语言根据接口类型是否包含一组方法将接口类型分成了两类：

- 使用 runtime.iface 结构体表示包含方法的接口
- 使用 runtime.eface 结构体表示不包含任何方法的 interface{} 类型；



runtime.eface 结构体在 Go 语言中的定义是这样的：

```go
type eface struct {
	_type *_type
	data  unsafe.Pointer
}
```

由于 interface{} 类型不包含任何方法，所以它的结构也相对来说比较简单，只包含指向底层数据和类型的两个指针。从上述结构我们也能推断出 — Go 语言的任意类型都可以转换成 interface{}。

另一个用于表示接口的结构体是 runtime.iface，这个结构体中有指向原始数据的指针 data，不过更重要的是 runtime.itab 类型的 tab 字段。



```go
type iface struct {
	tab  *itab
	data unsafe.Pointer
}
```

接下来我们将详细分析 Go 语言接口中的这两个类型，即 runtime._type 和 runtime.itab。



**类型结构体** 

runtime._type 是 Go 语言类型的运行时表示。下面是运行时包中的结构体，其中包含了很多类型的元信息，例如：类型的大小、哈希、对齐以及种类等。



