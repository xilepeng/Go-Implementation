
## Unwrap

```go
func (e *importError) Unwrap() error {
	// Don't return e.err directly, since we're only wrapping an error if %w
	// was passed to ImportErrorf.
	return errors.Unwrap(e.err)
}
```


Unwrap 将嵌套的 error 解析出来，多层嵌套需要调用 Unwrap 函数多次，才能获取最里层的 error。

```go
package main

import (
	"errors"
	"fmt"
)
// errors.Unwrap 将嵌套的 error 解析出来，多层嵌套需要调用 Unwrap 函数多次，才能获取最里层的 error。
func main() {
	err1 := errors.New("error1")
	err2 := fmt.Errorf("error2: [%w]", err1)
	fmt.Println(err2)
	fmt.Println(errors.Unwrap(err2))
}

// Output
// error2: [error1]
// error1
```

## Is

```go
func Is(err, target error) bool
```

判断 err 是否和 target 是同一类型，或者 err 嵌套的 error 有没有和 target 是同一类型的，如果是，则返回 true。

源码如下：

```go
func Is(err, target error) bool {
	if target == nil {
		return err == target
	}

	isComparable := reflectlite.TypeOf(target).Comparable()
    // 无限循环，比较 err 以及嵌套的 error
	for {
		if isComparable && err == target {
			return true
		}
        // 调用 error 的 Is 方法，这里可以自定义实现
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
        // 返回被嵌套的下一层的 error
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}
```

通过一个无限循环，使用 Unwrap 不断地将 err 里层嵌套的 error 解开，再看被解开的 error 是否实现了 Is 方法，并且调用它的 Is 方法，当两者都返回 true 的时候，整个函数返回 true。


## As

```go
func As(err error, target any) bool 
```

从 err 错误链里找到和 target 相等的并且设置 target 所指向的变量。

源码如下：

```go
func As(err error, target any) bool {
    // target 不能为 nil
	if target == nil {
		panic("errors: target cannot be nil")
	}
	val := reflectlite.ValueOf(target)
	typ := val.Type()
    // target 必须是一个非空指针
	if typ.Kind() != reflectlite.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	targetType := typ.Elem()
    // 保证 target 是一个接口类型或者实现了 Error 接口
	if targetType.Kind() != reflectlite.Interface && !targetType.Implements(errorType) {
		panic("errors: *target must be interface or implement error")
	}
	for err != nil {
        // 使用反射判断是否可被赋值，如果可以就赋值并且返回true
		if reflectlite.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflectlite.ValueOf(err))
			return true
		}
        // 调用 error 自定义的 As 方法，实现自己的类型断言代码
		if x, ok := err.(interface{ As(any) bool }); ok && x.As(target) {
			return true
		}
        // 不断地 Unwrap，一层层的获取嵌套的 error
		err = Unwrap(err)
	}
	return false
}
```

返回 true 的条件是错误链里的 err 能被赋值到 target 所指向的变量；或者 err 实现的 As(interface{}) bool 方法返回 true。

前者，会将 err 赋给 target 所指向的变量；后者，由 As 函数提供这个功能。

如果 target 不是一个指向“实现了 error 接口的类型或者其它接口类型”的非空的指针的时候，函数会 panic。





```go
package main

import (
	"errors"
	"fmt"
	"log"
)

func doError() (string, error) {
	return "哇塞计划", errors.New("This is my error")
}

func doNoError() (string, error) {
	return "My response", nil
}

func doFmtError() error {
	errCode := 401
	return fmt.Errorf("This my error code: %d", errCode)
}

func main() {
	resp, err := doError()
	if err != nil {
		log.Printf("There was an error: %v\n", err)
	}
	fmt.Println("My message:", resp)
	resp, err = doNoError()
	if err != nil {
		log.Printf("this should nor print")
	}
	fmt.Println("My response:", resp)
	err = doFmtError()
	if err != nil {
		log.Printf("There was an error: %v\n", err)
	}
}

// 2022/07/12 14:30:50 There was an error: This is my error
// My message: 哇塞计划
// My response: My response
// 2022/07/12 14:30:50 There was an error: This my error code: 401

```




