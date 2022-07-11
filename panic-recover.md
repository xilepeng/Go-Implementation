

## panic 和 recover 

- panic 只会触发当前 Goroutine 的 defer；
- recover 只有在 defer 中调用才会生效；
- panic 允许在 defer 中嵌套多次调用；


```go
package main

import "fmt"

func a() {
	fmt.Println("Inside a()")
	defer func() {
		if c := recover(); c != nil { // 处理 b 函数 panic 情况
			fmt.Println("Recover inside a()!")
		}
	}()
	fmt.Println("About to call b()")
	b()
	// recover 可以中止 panic 造成的程序崩溃。它是一个只能在 defer 中发挥作用的函数，在其他作用域中调用不会发挥作用；
	fmt.Println("b() exited!") // 未执行
	fmt.Println("Exiting a()") // 未执行
}

func b() {
	fmt.Println("Inside b()")
	// panic 能够改变程序的控制流，调用 panic 后会立刻停止执行当前函数的剩余代码，
	// 并在当前 Goroutine 中递归执行调用方的 defer；
	panic("Panic in b()!")     // 未执行
	fmt.Println("Exiting b()") // 未执行
}

func main() {
	a()
	fmt.Println("main() ended")
}

// Inside a()
// About to call b()
// Inside b()
// Recover inside a()!
// main() ended

```



**跨协程失效**

```go
func main() {
	defer fmt.Println("in main")    // 未执行
	go func() {
		defer fmt.Println("in goroutine")
		panic("")
	}()
	time.Sleep(1 * time.Second)
}

// in goroutine
// panic:
```



**正确的崩溃恢复**

```go
func main() {
	defer fmt.Println("in main")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	panic("unknow err") // panic 在 recover 前
}

// unknow err
// in main
```

**失效的崩溃恢复**

```go
func main() {
	defer fmt.Println("in main")
	// 失效的崩溃恢复
	if err := recover(); err != nil {
		fmt.Println(err)
	}
	panic("unknow err") // panic 在 recover 后，
}

// 程序没有正常退出
// in main
// panic: unknow err
// goroutine 1 [running]:

```






## 数据结构 


/usr/local/Cellar/go/1.18.3/libexec/src/runtime/runtime2.go



```go
// A _panic holds information about an active panic.
//
// A _panic value must only ever live on the stack.
//
// The argp and link fields are stack pointers, but don't need special
// handling during stack growth: because they are pointer-typed and
// _panic values only live on the stack, regular stack pointer
// adjustment takes care of them.
type _panic struct {
	argp      unsafe.Pointer // pointer to arguments of deferred call run during panic; cannot move - known to liblink
	arg       any            // argument to panic
	link      *_panic        // link to earlier panic
	pc        uintptr        // where to return to in runtime if this panic is bypassed
	sp        unsafe.Pointer // where to return to in runtime if this panic is bypassed
	recovered bool           // whether this panic is over
	aborted   bool           // the panic was aborted
	goexit    bool
}
```

- argp 是指向 defer 调用时参数的指针；
- arg 是调用 panic 时传入的参数；
- link 指向了更早调用的 runtime._panic 结构；
- recovered 表示当前 runtime._panic 是否被 recover 恢复；
- aborted 表示当前的 panic 是否被强行终止；

从数据结构中的 link 字段我们就可以推测出以下的结论：panic 函数可以被连续多次调用，它们之间通过 link 可以组成链表。



## 程序崩溃 
这里先介绍分析 panic 函数是终止程序的实现原理。编译器会将关键字 panic 转换成 runtime.gopanic，该函数的执行过程包含以下几个步骤：

- 创建新的 runtime._panic 并添加到所在 Goroutine 的 _panic 链表的最前面；
- 在循环中不断从当前 Goroutine 的 _defer 中链表获取 runtime._defer 并调用 runtime.reflectcall 运行延迟调用函数；
- 调用 runtime.fatalpanic 中止整个程序；

/usr/local/Cellar/go/1.18.3/libexec/src/runtime/panic.go

```go
func gopanic(e interface{}) {
	gp := getg()
	...
	var p _panic
	p.arg = e
	p.link = gp._panic
	gp._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

	for {
		d := gp._defer
		if d == nil {
			break
		}

		d._panic = (*_panic)(noescape(unsafe.Pointer(&p)))

		reflectcall(nil, unsafe.Pointer(d.fn), deferArgs(d), uint32(d.siz), uint32(d.siz))

		d._panic = nil
		d.fn = nil
		gp._defer = d.link

		freedefer(d)
		if p.recovered {
			...
		}
	}

	fatalpanic(gp._panic)
	*(*int)(nil) = 0
}
```











panic 方法实际上就是处理当前 Goroutine(g) 上所挂载的 ._panic 链表（所以无法对其他 Goroutine 的异常事件响应），然后对其所属的 defer 链表和 recover 进行检测并处理，最后调用退出命令中止应用程序

