

1. [1. for and range](#1-for-and-range)
2. [2. select](#2-select)
3. [3. defer](#3-defer)
4. [4. panic and recover](#4-panic-and-recover)
5. [5. make and new](#5-make-and-new)



## 1. for and range

```go
package main

import "fmt"

func main() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
}
```

```go
$ go tool compile -S main.go
"".main STEXT size=121 args=0x0 locals=0x48 funcid=0x0 align=0x0
	0x0000 00000 (main.go:5)	TEXT	"".main(SB), ABIInternal, $72-0

	0x0014 00020 (main.go:5)	XORL	AX, AX            // i:=0

	0x0016 00022 (main.go:6)	JMP	98
	0x0054 00084 ($GOROOT/src/fmt/print.go:274)	CALLfmt.Fprintln(SB)

	0x005e 00094 (main.go:6)	LEAQ	1(CX), AX            // i++
	0x0062 00098 (main.go:6)	CMPQ	AX, $10              // 比较变量 i 和 10
	0x0066 00102 (main.go:6)	JLT	24                       // 跳转到 24 行,如果 i < 10
```




## 2. select




## 3. defer

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

```go

```


















## 4. panic and recover




## 5. make and new