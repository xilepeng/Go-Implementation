


1. [for and range](#for-and-range)



## for and range

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













