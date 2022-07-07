


## for

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





## for range

```go
package main

import "fmt"

func main() {
	slice := []int{10, 20, 30, 40}
	for index, copy_value := range slice {
		fmt.Printf("Value = %d , Value-Addr = %x , Elem-Addr = %x\n", copy_value, &copy_value, &slice[index])
	}
}
```
程序输出：

```go
Value = 10 , Value-Addr = c0000b2008 , Elem-Addr = c0000b6000
Value = 20 , Value-Addr = c0000b2008 , Elem-Addr = c0000b6008
Value = 30 , Value-Addr = c0000b2008 , Elem-Addr = c0000b6010
Value = 40 , Value-Addr = c0000b2008 , Elem-Addr = c0000b6018
```

当迭代切片时，关键字 range 会返回两个值。第一个值是当前迭代到的索引位置，第二个 值是该位置对应元素值的一份副本
因为迭代返回的变量是一个迭代过程中根据切片依次赋值的新变量，所以 value 的地址总 是相同的。
要想获取每个元素的地址，可以使用切片变量和索引值。
关键字 range 总是会从切片头部开始迭代。如果想对迭代做更多的控制，依旧可以使用传统的 for 循环


从上面结果我们可以看到，**如果用 range 的方式去遍历一个切片，拿到的 Value 其实是切片里面的值拷贝**。所以每次打印 Value 的地址都不变。

由于 Value 是值拷贝的，并非引用传递，所以直接改 Value 是达不到更改原切片值的目的的，需要通过 &slice[index] 获取真实的地址。


身体要挺起来，中央脊椎上提，状态挺拔向上，肩自然下垂，形成对抗，挺胸收腹提臀提挎，收下颚，眼神定住。










