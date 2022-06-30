



##  defer

defer语句用于延迟函数的调用，每次defer都会把一个函数压入栈中，函数返回前再把延迟的函数取出并执行。

为了方便描述，我们把创建defer的函数称为主函数，defer语句后面的函数称为延迟函数。

延迟函数可能有输入参数，这些参数可能来源于定义defer的函数，延迟函数也可能引用主函数用于返回的变量，也就是说延迟函数可能会影响主函数的一些行为，这些场景下，如果不了解defer的规则很容易出错。

其实官方说明的defer的三个原则很清楚，本节试图汇总defer的使用场景并做简单说明。

**热身**

按照惯例，我们看几个有意思的题目，用于检验对defer的了解程度。

**题目**

下面函数输出结果是什么？

```go
func deferFuncParameter() {
	var aInt = 1
	defer fmt.Println(aInt)
	aInt = 2
	return
}
```
题目说明：
函数deferFuncParameter()定义一个整型变量并初始化为1，然后使用defer语句打印出变量值，最后修改变量值为2.

参考答案：
输出1。延迟函数fmt.Println(aInt)的参数在defer语句出现时就已经确定了，所以无论后面如何修改aInt变量都不会影响延迟函数。



 **defer规则**

 **3.1 规则一：延迟函数的参数在defer语句出现时就已经确定下来了**

官方给出一个例子，如下所示：
```go
func a() {
	i := 0
	defer fmt.Println(i) // 0
	i++
	return
}
```
defer语句中的fmt.Println()参数i值在defer出现时就已经确定下来，实际上是拷贝了一份。后面对变量i的修改不会影响fmt.Println()函数的执行，仍然打印"0"。

注意：对于指针类型参数，规则仍然适用，只不过延迟函数的参数是一个地址值，这种情况下，defer后面的语句对变量的修改可能会影响延迟函数。



**3.2 规则二：延迟函数执行按后进先出顺序执行，即先出现的defer最后执行**

这个规则很好理解，定义defer类似于入栈操作，执行defer类似于出栈操作。

defer 在外围函数返回之后，以后进先出 (LIFO)的原则执行。简单点说，在一个外围函数中有3个defer函数: f1() 最先出现，然后 f2() ，最后 f3() ，当外围函数执行返回之后， f3() 最先被执行，接着是 f2() ，最后是 f1() 。

```go
package main

import "fmt"

func d1() {
	for i := 3; i > 0; i-- {
		defer fmt.Print(i, " ") // 1 2 3
	}
	fmt.Println()
	return
}

// defer 执行没有带参数的匿名函数
// 循环结束后，i的值为0，因为没有参数，这意味着，将值为0的i进行三次打印输出
func d2() {
	for i := 3; i > 0; i-- {
		defer func() {
			fmt.Print(i, " ") // 0 0 0
		}()
	}
	fmt.Println()
	return
}

// defer 执行带参数的匿名函数
func d3() {
	for i := 3; i > 0; i-- {
		defer func(n int) {
			fmt.Print(n, " ") // 1 2 3
		}(i)
	}
	fmt.Println()
}

func main() {
	d1()
	d2()
	d3()
}

```






**3.3 规则三：延迟函数可能操作主函数的具名返回值**


定义defer的函数，即主函数可能有返回值，返回值有没有名字没有关系，defer所作用的函数，即延迟函数可能会影响到返回值。

若要理解延迟函数是如何影响主函数返回值的，只要明白函数是如何返回的就足够了。

- 3.3.1 函数返回过程

有一个事实必须要了解，关键字return不是一个原子操作，实际上return只代理汇编指令ret，即将跳转程序执行。比如语句**return i，实际上分两步进行，即将i值存入栈中作为返回值，然后执行跳转，而defer的执行时机正是跳转前，所以说defer执行时还是有机会操作返回值的。**

```go
func deferFuncReturn() (result int) {
	i := 0
	defer func() {
		result++
	}()
	return i // 1
}
```

该函数的return语句可以拆分成下面两行：

```go
result = i
return
```

**而延迟函数的执行正是在return之前，即加入defer后的执行过程如下：**

```go
result = i
result++
return
```

所以上面函数实际返回i++值。

关于主函数有不同的返回方式，但返回机制就如上机介绍所说，只要把return语句拆开都可以很好的理解，下面分别举例说明



- 主函数拥有匿名返回值，返回字面值



```go
package main

import "fmt"

// 主函数拥有一个匿名的返回值，返回字面值
// return语句，直接把1写入栈中作为返回值，延迟函数无法操作该返回值，所以就无法影响返回值。
func foo() int {
	i := 0
	defer func() {
		i++
	}()
	return 1 // 1
}

//  主函数拥有匿名返回值，返回变量
// defer语句可以引用到返回值，但不会改变返回值。
func foo2() int {
	i := 0
	defer func() {
		i++
	}()
	return i // 0
}

func main() {
	fmt.Println(foo())
	fmt.Println(foo2())
}

```

foo2() 函数，返回一个局部变量，同时defer函数也会操作这个局部变量。对于匿名返回值来说，可以假定仍然有一个变量存储返回值，假定返回值变量为"copy_result"，上面的返回语句可以拆分成以下过程：

```go
copy_result = i
i++
return
```

由于i是整型，会将值拷贝给copy_result，所以defer语句中修改i值，对函数返回值不造成影响。


foo3() 函数拆解出来，如下所示：
```go
result = 0
result ++
return 
```

函数真正返回前，在defer中对返回值做了+1操作，所以函数最终返回1。



**4.1 defer数据结构**

源码包src/src/runtime/runtime2.go:_defer定义了defer的数据结构：

```go
// A _defer holds an entry on the list of deferred calls.
// If you add a field here, add code to clear it in deferProcStack.
// This struct must match the code in cmd/compile/internal/ssagen/ssa.go:deferstruct
// and cmd/compile/internal/ssagen/ssa.go:(*state).call.
// Some defers will be allocated on the stack and some on the heap.
// All defers are logically part of the stack, so write barriers to
// initialize them are not required. All defers must be manually scanned,
// and for heap defers, marked.
type _defer struct {
	started bool
	heap    bool
	// openDefer indicates that this _defer is for a frame with open-coded
	// defers. We have only one defer record for the entire frame (which may
	// currently have 0, 1, or more defers active).
	openDefer bool
	sp        uintptr // sp at time of defer
	pc        uintptr // pc at time of defer
	fn        func()  // can be nil for open-coded defers
	_panic    *_panic // panic that is running defer
	link      *_defer // next defer on G; can point to either heap or stack!

	// If openDefer is true, the fields below record values about the stack
	// frame and associated function that has the open-coded defer(s). sp
	// above will be the sp for the frame, and pc will be address of the
	// deferreturn call in the function.
	fd   unsafe.Pointer // funcdata for the function associated with the frame
	varp uintptr        // value of varp for the stack frame
	// framepc is the current pc associated with the stack frame. Together,
	// with sp above (which is the sp associated with the stack frame),
	// framepc/sp can be used as pc/sp pair to continue a stack trace via
	// gentraceback().
	framepc uintptr
}
```

我们知道defer后面一定要接一个函数的，所以defer的数据结构跟一般函数类似，也有栈地址、程序计数器、函数地址等等。

与函数不同的一点是它含有一个指针，可用于指向另一个defer，每个goroutine数据结构中实际上也有一个defer指针，该指针指向一个defer的单链表，每次声明一个defer时就将defer插入到单链表表头，每次执行defer时就从单链表表头取出一个defer执行。



**defer的创建和执行**

源码包src/runtime/panic.go定义了两个方法分别用于创建defer和执行defer。

- deferproc()： 在声明defer处调用，其将defer函数存入goroutine的链表中；
- deferreturn()：在return指令，准确的讲是在ret指令前调用，其将defer从goroutine链表中取出并执行。
可以简单这么理解，在编译在阶段，声明defer处插入了函数deferproc()，在函数return前插入了函数deferreturn()。














