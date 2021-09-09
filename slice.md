
### 以下代码基于 Go 1.17


## 2. 基本数据结构

slice 的底层源码和相关实现在 src/runtime/slice.go

```go
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```

 

