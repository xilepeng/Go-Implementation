
1. [iface 和 eface 的区别是什么?](#iface-和-eface-的区别是什么)
2. [Go 接口与 C++ 接口有何异同？](#go-接口与-c-接口有何异同)
3. [如何用 interface 实现多态 ?](#如何用-interface-实现多态-)
4. [接口转换的原理 ?](#接口转换的原理-)
5. [Go 语言与鸭子类型的关系 ?](#go-语言与鸭子类型的关系-)
6. [方法的值接收者和指针接收者的区别 ?](#方法的值接收者和指针接收者的区别-)
7. [接口的构造过程是怎样的 ？](#接口的构造过程是怎样的-)
8. [编译器自动检测类型是否实现接口 ？](#编译器自动检测类型是否实现接口-)
9. [类型转换和断言的区别 ?](#类型转换和断言的区别-)
10. [接口的动态类型和动态值 ?](#接口的动态类型和动态值-)

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

当判定一种类型是否满足某个接口时，Go 使用类型的方法集和接口所需要的方法集进行匹配，如果类型的方法集完全包含接口的方法集，则可认为该类型实现了该接口。

```go
package main

import "fmt"

type coder interface {
    code()
    run()
}

type runner interface {
    run()
}

type Gopher struct {
    language string
}

func (g Gopher) code() {
    return
}

func (g Gopher) run() {
    return
}

func main() {
    var c coder = Gopher{}

    var r runner
    r = c
    fmt.Println(c, r)
}
```

简单解释下上述代码：定义了两个 interface: coder 和 runner。定义了一个实体类型 Gopher，类型 Gopher 实现了两个方法，分别是 run() 和 code()。main 函数里定义了一个接口变量 c，绑定了一个 Gopher 对象，之后将 c 赋值给另外一个接口变量 r 。赋值成功的原因是 c 中包含 run() 方法。这样，两个接口变量完成了转换。



## Go 语言与鸭子类型的关系 ?

先直接来看维基百科里的定义：

```go
If it looks like a duck, swims like a duck, and quacks like a duck, then it probably is a duck.
```


翻译过来就是：如果某个东西长得像鸭子，像鸭子一样游泳，像鸭子一样嘎嘎叫，那它就可以被看成是一只鸭子。
Duck Typing，鸭子类型，是动态编程语言的一种对象推断策略，它更关注对象能如何被使用，而不是对象的类型本身。Go 语言作为一门静态语言，它通过通过接口的方式完美支持鸭子类型。

Go 语言作为一门现代静态语言，是有后发优势的。它引入了动态语言的便利，同时又会进行静态语言的类型检查，写起来是非常 Happy 的。Go 采用了折中的做法：不要求类型显示地声明实现了某个接口，只要实现了相关的方法即可，编译器就能检测到。

```go
package main

import "fmt"

// 定义一个接口，和使用此接口作为参数的函数
type IGreeting interface {
	sayHello()
}

func sayHello(i IGreeting) {
	i.sayHello()
}

// 定义 Go 类型
type Go struct{}

// sayHello 使用值接收者实现了一个方法
func (g Go) sayHello() {
	fmt.Println("Hi, I am Go")
}

// 方法能给用户定义的类型添加新的行为。方法实际上也是函数，只是在声明时，在关键字 func 和方法名之间增加了一个参数
type CPlus struct{}

func (c CPlus) sayHello() {
	fmt.Println("Hi, I am CPlus")
}

func main() {
	Golang := Go{}
	CPlus := CPlus{}
	sayHello(Golang)
	sayHello(CPlus)
}

// Hi, I am Go
// Hi, I am CPlus

```


在 main 函数中，调用调用 sayHello() 函数时，传入了 Golang、CPlus 对象，它们并没有显式地声明实现了 IGreeting 类型，只是实现了接口所规定的 sayHello() 函数。实际上，编译器在调用 sayHello() 函数时，会隐式地将 Golang、CPlus 对象转换成 IGreeting 类型，这也是静态语言的类型检查功能。

顺带再提一下动态语言的特点：

`变量绑定的类型是不确定的，在运行期间才能确定 函数和方法可以接收任何类型的参数，且调用时不检查参数类型 不需要实现接口`

总结一下，鸭子类型是一种动态语言的风格，在这种风格中，一个对象有效的语义，不是由继承自特定的类或实现特定的接口，而是由它"当前方法和属性的集合"决定。Go 作为一种静态语言，通过接口实现了 鸭子类型，实际上是 Go 的编译器在其中作了隐匿的转换工作。



## 方法的值接收者和指针接收者的区别 ?


方法能给用户自定义的类型添加新的行为。它和函数的区别在于方法有一个接收者，给一个函数添加一个接收者，那么它就变成了方法。接收者可以是值接收者，也可以是指针接收者。

在调用方法的时候，值类型既可以调用值接收者的方法，也可以调用指针接收者的方法；指针类型既可以调用指针接收者的方法，也可以调用值接收者的方法。

也就是说，不管方法的接收者是什么类型，该类型的值和指针都可以调用，不必严格符合接收者的类型。

```go
package main

import "fmt"

type Person struct {
	age int
}

func (p Person) HowOld() int {
	return p.age
}

func (p *Person) GrowUp() {
	p.age++
}

func main() {
	// mojo 是值类型
	mojo := Person{age: 18}
	// 值类型 调用接收者也是值类型的方法
	fmt.Println(mojo.HowOld())
	// 值类型 调用接收者是指针类型的方法
	mojo.GrowUp()
	fmt.Println(mojo.HowOld())

	// ------------------------
	// mojo 是指针类型
	pointer_mojo := &Person{age: 100}
	// 指针类型 调用接收者也是值类型的方法
	fmt.Println(pointer_mojo.HowOld())
	// 指针类型 调用接收者是指针类型的方法
	pointer_mojo.GrowUp()
	fmt.Println(pointer_mojo.HowOld())

}

// 18
// 19
// 100
// 101

```

调用了 growUp 函数后，不管调用者是值类型还是指针类型，它的 Age 值都改变了。
实际上，当类型和方法的接收者类型不同时，其实是编译器在背后做了一些工作，用一个表格来呈现：



|           |值接收者      |指针接收者   |
|:--------- |:-----------|:-----------|
|值类型调用者 | 方法会使用调用者的一个副本，类似于“传值” | 使用值的引用来调用方法，上例中，mojo.GrowUp() 实际上是 (&mojo).GrowUp()|
|指针类型调用者| 指针被解引用为值，上例中，pointer_mojo 实际上是 (*pointer_mojo).HowOld() | 实际上也是“传值”，方法里的操作会影响到调用者，类似于指针传参，拷贝了一份指针|


**值接收者和指针接收者**

前面说过，不管接收者类型是值类型还是指针类型，都可以通过值类型或指针类型调用，这里面实际上通过语法糖起作用的。

先说结论：**实现了接收者是值类型的方法，相当于自动实现了接收者是指针类型的方法；而实现了接收者是指针类型的方法，不会自动生成对应接收者是值类型的方法。**


```go
package main

import "fmt"

type coder interface {
	code()
	debug()
}

// 定义了一个结构体 Gopher，它实现了两个方法，一个值接收者，一个指针接收者。
type Gopher struct {
	language string
}

// 实现了接收者是值类型的方法，相当于自动实现了接收者是指针类型的方法
func (g Gopher) code() {
	fmt.Printf("I am Coding %s \n", g.language)
}

// 指针类型的接收者，不会自动生成对应值类型的方法
func (g *Gopher) debug() {
	fmt.Printf("I am Coding %s \n", g.language)
}

func main() {
	var c coder = &Gopher{"Go"}
	c.code()
	c.debug()
}

// var c coder = Gopher{"Go"}

// # command-line-arguments
// ./main.go:23:16: cannot use Gopher{…} (value of type Gopher) as type coder in variable declaration:
//         Gopher does not implement coder (debug method has pointer receiver)

```


**如果实现了接收者是值类型的方法，会隐含地也实现了接收者是指针类型的方法。**

- 如果方法的接收者是值类型，无论调用者是对象还是对象指针，修改的都是对象的副本，不影响调用者；
- 如果方法的接收者是指针类型，则调用者修改的是指针指向的对象本身。


**使用指针作为方法的接收者的理由：**
- 方法能够修改接收者指向的值。
- 避免在每次调用方法时复制该值，在值的类型为大型结构体时，这样做会更加高效。

是使用值接收者还是指针接收者，不是由该方法是否修改了调用者（也就是接收者）来决定，而是应该基于该类型的本质。
如果类型具备“原始的本质”，也就是说它的成员都是由 Go 语言里内置的原始类型，如字符串，整型值等，那就定义值接收者类型的方法。像内置的引用类型，如 slice，map，interface，channel，这些类型比较特殊，声明他们的时候，实际上是创建了一个 header， 对于他们也是直接定义值接收者类型的方法。这样，调用函数时，是直接 copy 了这些类型的 header，而 header 本身就是为复制设计的。
如果类型具备非原始的本质，不能被安全地复制，这种类型总是应该被共享，那就定义指针接收者的方法。比如 go 源码里的文件结构体（struct File）就不应该被复制，应该只有一份实体。





## 接口的构造过程是怎样的 ？



```go
package main

import "fmt"

type Person interface {
	growUp()
}

type student struct {
	age int
}

func (p student) growUp() {
	p.age += 1
	return
}

// func (p *student) growUp() {
// 	p.age += 1
// 	return
// }

func main() {
	var mojo = Person(student{age: 18})
	// var mojo Person = &student{18}
	mojo.growUp()
	fmt.Println(mojo)
}

```


```go
➜  demo11 git:(main) ✗ go tool compile -S main.go
"".student.growUp STEXT nosplit size=1 args=0x8 locals=0x0 funcid=0x0 align=0x0
	0x0000 00000 (main.go:13)	TEXT	"".student.growUp(SB), NOSPLIT|ABIInternal, $0-8
	0x0000 00000 (main.go:13)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (main.go:13)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (main.go:13)	FUNCDATA	$5, "".student.growUp.arginfo1(SB)
	0x0000 00000 (main.go:13)	FUNCDATA	$6, "".student.growUp.argliveinfo(SB)
	0x0000 00000 (main.go:13)	PCDATA	$3, $1
	0x0000 00000 (main.go:15)	RET
	0x0000 c3                                               .
"".main STEXT size=138 args=0x0 locals=0x50 funcid=0x0 align=0x0
	0x0000 00000 (main.go:23)	TEXT	"".main(SB), ABIInternal, $80-0
	0x0000 00000 (main.go:23)	CMPQ	SP, 16(R14)
	0x0004 00004 (main.go:23)	PCDATA	$0, $-2
	0x0004 00004 (main.go:23)	JLS	125
	0x0006 00006 (main.go:23)	PCDATA	$0, $-1
	0x0006 00006 (main.go:23)	SUBQ	$80, SP
	0x000a 00010 (main.go:23)	MOVQ	BP, 72(SP)
	0x000f 00015 (main.go:23)	LEAQ	72(SP), BP
	0x0014 00020 (main.go:23)	FUNCDATA	$0, gclocals·69c1753bd5f81501d95132d08af04464(SB)
	0x0014 00020 (main.go:23)	FUNCDATA	$1, gclocals·713abd6cdf5e052e4dcd3eb297c82601(SB)
	0x0014 00020 (main.go:23)	FUNCDATA	$2, "".main.stkobj(SB)
	0x0014 00020 (main.go:24)	MOVQ	$18, ""..autotmp_9+40(SP)
	0x001d 00029 (main.go:24)	MOVQ	""..autotmp_9+40(SP), AX
	0x0022 00034 (main.go:24)	PCDATA	$1, $0
	0x0022 00034 (main.go:24)	CALL	runtime.convT64(SB)
	0x0027 00039 (main.go:24)	MOVQ	AX, ""..autotmp_21+48(SP)
	0x002c 00044 (main.go:26)	MOVQ	(AX), CX
	0x002f 00047 (main.go:26)	MOVQ	CX, AX
	0x0032 00050 (main.go:26)	PCDATA	$1, $1
	0x0032 00050 (main.go:26)	CALL	"".student.growUp(SB)
	0x0037 00055 (main.go:27)	MOVUPS	X15, ""..autotmp_13+56(SP)
	0x003d 00061 (main.go:27)	MOVQ	go.itab."".student,"".Person+8(SB), CX
	0x0044 00068 (main.go:27)	MOVQ	CX, ""..autotmp_13+56(SP)
	0x0049 00073 (main.go:27)	MOVQ	""..autotmp_21+48(SP), CX
	0x004e 00078 (main.go:27)	MOVQ	CX, ""..autotmp_13+64(SP)
	0x0053 00083 (<unknown line number>)	NOP
	0x0053 00083 ($GOROOT/src/fmt/print.go:274)	MOVQos.Stdout(SB), BX
	0x005a 00090 ($GOROOT/src/fmt/print.go:274)	LEAQgo.itab.*os.File,io.Writer(SB), AX
	0x0061 00097 ($GOROOT/src/fmt/print.go:274)	LEAQ""..autotmp_13+56(SP), CX
	0x0066 00102 ($GOROOT/src/fmt/print.go:274)	MOVL$1, DI
	0x006b 00107 ($GOROOT/src/fmt/print.go:274)	MOVQDI, SI
	0x006e 00110 ($GOROOT/src/fmt/print.go:274)	PCDATA	$1, $0
	0x006e 00110 ($GOROOT/src/fmt/print.go:274)	CALLfmt.Fprintln(SB)
	0x0073 00115 (main.go:28)	MOVQ	72(SP), BP
	0x0078 00120 (main.go:28)	ADDQ	$80, SP
	0x007c 00124 (main.go:28)	RET
	0x007d 00125 (main.go:28)	NOP
	0x007d 00125 (main.go:23)	PCDATA	$1, $-1
	0x007d 00125 (main.go:23)	PCDATA	$0, $-2
	0x007d 00125 (main.go:23)	NOP
	0x0080 00128 (main.go:23)	CALL	runtime.morestack_noctxt(SB)
	0x0085 00133 (main.go:23)	PCDATA	$0, $-1
	0x0085 00133 (main.go:23)	JMP	0
	
```





## 编译器自动检测类型是否实现接口 ？

```go
var _ io.Writer = (*myWriter)(nil)
```

编译器会由此检查 *myWriter 类型是否实现了 io.Writer 接口。 

```go
package main

import "io"

type myWriter struct{}

// func (w myWriter) Write(p []byte) (n int, err error) {
// 	return
// }

func main() {
	// 检查 *myWriter 类型是否实现了 io.Writer 接口
	var _ io.Writer = (*myWriter)(nil)
	// 检查 myWriter 类型是否实现了 io.Writer 接口
	var _ io.Writer = myWriter{}
}

```

注释掉为 myWriter 定义的 Write 函数后，运行程序：

```go
# command-line-arguments
./main.go:13:20: cannot use (*myWriter)(nil) (value of type *myWriter) as type io.Writer in variable declaration:
        *myWriter does not implement io.Writer (missing Write method)
./main.go:15:20: cannot use myWriter{} (value of type myWriter) as type io.Writer in variable declaration:
        myWriter does not implement io.Writer (missing Write method)
```

报错信息：*myWriter/myWriter 未实现 io.Writer 接口，也就是未实现 Write 方法。

解除注释后，运行程序不报错。

实际上，上述赋值语句会发生隐式地类型转换，在转换的过程中，编译器会检测等号右边的类型是否实现了等号左边接口所规定的函数。

总结一下，可通过在代码中添加类似如下的代码，用来检测类型是否实现了接口：

```go
var _ io.Writer = (*myWriter)(nil)
var _ io.Writer = myWriter{}
```



## 类型转换和断言的区别 ?

我们知道，Go 语言中不允许隐式类型转换，也就是说 = 两边，不允许出现类型不相同的变量。
- 类型转换、类型断言本质都是把一个类型转换成另外一个类型。
- 不同之处在于，类型断言是对接口变量进行的操作。

**类型转换**

对于类型转换而言，转换前后的两个类型要相互兼容才行。类型转换的语法为：
<结果类型> := <目标类型> ( <表达式> )


```go
package main

import "fmt"

func main() {
	var i int = 2
	var f float64

	f = float64(i)
	fmt.Printf("类型：%T, 值：%v\n", f, f)

	f = 5.2
	a := int(f)
	fmt.Printf("类型：%T, 值：%v\n", a, a)

	// s := []int(a)
	// cannot convert a (variable of type int) to type []int
}

```

程序输出：
```go
类型：float64, 值：2
类型：int, 值：5
```

上面的代码里，我定义了一个 int 型和 float64 型的变量，尝试在它们之前相互转换，结果是成功的：int 型和 float64 是相互兼容的。
如果我把最后一行代码的注释去掉，编译器会报告类型不兼容的错误：


```go
cannot convert a (variable of type int) to type []int
```


**断言**

前面说过，因为空接口 interface{} 没有定义任何函数，因此 Go 中所有类型都实现了空接口。当一个函数的形参是 interface{}，那么在函数中，需要对形参进行断言，从而得到它的真实类型。

类型转换和类型断言有些相似，不同之处，在于类型断言是对接口进行的操作。

```go
package main

import "fmt"

type Student struct {
	Name string
	Age  int
}

func main() {
	var i interface{} = new(Student)
	s := i.(Student)
	fmt.Println(s)
}

// panic: interface conversion: interface {} is *main.Student, not main.Student
```

运行输出：

```go
panic: interface conversion: interface {} is *main.Student, not main.Student
```

直接 panic 了，这是因为 i 是 *Student 类型，并非 Student 类型，断言失败。这里直接发生了 panic，线上代码可能并不适合这样做，可以采用“安全断言”的语法：

```go
package main

import "fmt"

type Student struct {
	Name string
	Age  int
}

func main() {
	var i interface{} = new(Student)
	s, ok := i.(Student)
	if ok {
		fmt.Println(s)
	}
}
```

这样，即使断言失败也不会 panic。

断言其实还有另一种形式，就是用在利用 switch 语句判断接口的类型。每一个 case 会被顺序地考虑。当命中一个 case 时，就会执行 case 中的语句，因此 case 语句的顺序是很重要的，因为很有可能会有多个 case 匹配的情况。

代码示例如下：

```go
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

```

运行代码：
```go
// var i interface{}
0xc000096210, <nil>
0xc000096220 <nil>
nil type[<nil>] <nil>

// var i interface{} = new(Student)
0xc000096210, &{ 0}
0xc0000ac020 &{ 0}
*Student type[*main.Student] &{ 0}

// var i interface{} = (*Student)(nil)
0xc000096210, <nil>
0xc0000ac020 <nil>
*Student type[*main.Student] <nil>
```


```go
var i interface{}
```
i 是 nil 类型。


```go
var i interface{} = new(Student)
```
i 是一个 *Student 类型，匹配上第三个 case，从打印的三个地址来看，这三处的变量实际上都是不一样的。在 main 函数里有一个局部变量 i；调用函数时，实际上是复制了一份参数，因此函数里又有一个变量 v，它是 i 的拷贝；断言之后，又生成了一份新的拷贝。所以最终打印的三个变量的地址都不一样。


```go
var i interface{} = (*Student)(nil)
```

这里想说明的其实是 i 在这里动态类型是 (*Student), 数据为 nil，它的类型并不是 nil，它与 nil 作比较的时候，得到的结果也是 false。




- 【引申1】 fmt.Println 函数的参数是 interface。对于内置类型，函数内部会用穷举法，得出它的真实类型，然后转换为字符串打印。而对于自定义类型，首先确定该类型是否实现了 String() 方法，如果实现了，则直接打印输出 String() 方法的结果；否则，会通过反射来遍历对象的成员进行打印。


```go
package main

import "fmt"

type Student struct {
	Name string
	Age  int
}

func main() {
	var s = Student{
		Name: "mojo",
		Age:  18,
	}

	fmt.Println(s)
}
```

因为 Student 结构体没有实现 String() 方法，所以 fmt.Println 会利用反射挨个打印成员变量：

```go
{mojo 18}
```

增加一个 String() 方法的实现：
```go
func (s Student) String() string {
	return fmt.Sprintf("[Name: %s], [Age: %d]", s.Name, s.Age)
}
```

打印结果：
```go
[Name: mojo], [Age: 18]
```
按照我们自定义的方法来打印了。

- 【引申2】 针对上面的例子，如果改一下：


```go
func (s *Student) String() string {
    return fmt.Sprintf("[Name: %s], [Age: %d]", s.Name, s.Age)
}
```

注意看两个函数的接受者类型不同，现在 Student 结构体只有一个接受者类型为 指针类型 的 String() 函数，打印结果：

```go
{mojo 18}
```

为什么？

- 类型 T 只有接受者是 T 的方法；而类型 *T 拥有接受者是 T 和 *T 的方法。语法上 T 能直接调 *T 的方法仅仅是 Go 的语法糖。

所以， Student 结构体定义了接受者类型是值类型的 String() 方法时，通过

```go
fmt.Println(s)
fmt.Println(&s)

// [Name: mojo], [Age: 18]
// [Name: mojo], [Age: 18]
```

均可以按照自定义的格式来打印。

如果 Student 结构体定义了接受者类型是指针类型的 String() 方法时，只有通过

```go
fmt.Println(&s)

// [Name: mojo], [Age: 18]
```

才能按照自定义的格式打印。






## 接口的动态类型和动态值 ?


从源码里可以看到：iface包含两个字段：tab 是接口表指针，指向类型信息；data 是数据指针，则指向具体的数据。它们分别被称为动态类型和动态值。而接口值包括动态类型和动态值。

- 【引申1】接口类型和 nil 作比较

接口值的零值是指动态类型和动态值都为 nil。当仅且当这两部分的值都为 nil 的情况下，这个接口值就才会被认为 接口值 == nil。

来看个例子：

```go
package main

import "fmt"

type Coder interface {
	code()
}

type Gopher struct {
	name string
}

func (g Gopher) code() {
	fmt.Printf("%s is coding\n", g.name)
}

func main() {
	var c Coder
	fmt.Println(c == nil)
	fmt.Printf("c: %T, %v\n", c, c)

	var g *Gopher
	fmt.Println(g == nil)
	fmt.Printf("g: %T, %v\n", g, g)

	c = g
	fmt.Println(c == nil)
	fmt.Printf("c: %T, %v\n", c, c)
}

```

输出：
```go
true
c: <nil>, <nil>
true
g: *main.Gopher, <nil>
false
c: *main.Gopher, <nil>
```

一开始，c 的 动态类型和动态值都为 nil，g 也为 nil，当把 g 赋值给 c 后，c 的动态类型变成了 *main.Gopher，仅管 c 的动态值仍为 nil，但是当 c 和 nil 作比较的时候，结果就是 false 了。

- 【引申2】 来看一个例子，看一下它的输出：

```go
package main

import "fmt"

type MyError struct{}

func (i MyError) Error() string {
	return "MyError"
}

func process() error {
	var err *MyError = nil
	return err // 隐含类型转换 
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
```

这里先定义了一个 MyError 结构体，实现了 Error 函数，也就实现了 error 接口。Process 函数返回了一个 error 接口，这块隐含了类型转换。所以，虽然它的值是 nil，其实它的类型是 *MyError，最后和 nil 比较的时候，结果为 false。



- 【引申3】如何打印出接口的动态类型和值？

```go
package main

import (
    "unsafe"
    "fmt"
)

type iface struct {
    itab, data uintptr
}

func main() {
    var a interface{} = nil

    var b interface{} = (*int)(nil)

    x := 5
    var c interface{} = (*int)(&x)

    ia := *(*iface)(unsafe.Pointer(&a))
    ib := *(*iface)(unsafe.Pointer(&b))
    ic := *(*iface)(unsafe.Pointer(&c))

    fmt.Println(ia, ib, ic)

    fmt.Println(*(*int)(unsafe.Pointer(ic.data)))
}
```

代码里直接定义了一个 iface 结构体，用两个指针来描述 itab 和 data，之后将 a, b, c 在内存中的内容强制解释成我们自定义的 iface。最后就可以打印出动态类型和动态值的地址。

```go
{0 0} {17359136 0} {17359136 824634322632}
5
```

a 的动态类型和动态值的地址均为 0，也就是 nil；b 的动态类型和 c 的动态类型一致，都是 *int；最后，c 的动态值为 5。







































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



