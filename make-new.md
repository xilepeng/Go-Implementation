
## make 和 new 的区别：

- make 的作用是初始化内置的数据结构，也就是我们在前面提到的切片、哈希表和 Channel；
- new 的作用是根据传入的类型分配一片内存空间并返回指向这片内存空间的指针；




```go
package main

import "fmt"

func main() {
	slice := make([]int, 10, 100)
	hash := make(map[int]bool, 10)
	ch := make(chan int, 10)

	fmt.Printf("slice 类型：%T, 值：%v,\n", slice, slice)
	fmt.Printf("hash  类型：%T, 值：%v\n", hash, hash)
	fmt.Printf("ch    类型：%T, 值：%v\n", ch, ch)
}


// slice 类型：[]int, 值：[0 0 0 0 0 0 0 0 0 0],
// hash  类型：map[int]bool, 值：map[]
// ch    类型：chan int, 值：0xc0000b80b0
```

- slice 是一个包含 data、cap 和 len 的结构体 [reflect.SliceHeader](https://draveness.me/golang/tree/reflect.SliceHeader)；

/usr/local/Cellar/go/1.18.3/libexec/src/reflect/value.go
```go
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
```

- hash 是一个指向 [runtime.hmap](https://draveness.me/golang/tree/runtime.hmap) 结构体的指针；

```go
// A header for a Go map.
type hmap struct {
	// Note: the format of the hmap is also encoded in cmd/compile/internal/reflectdata/reflect.go.
	// Make sure this stays in sync with the compiler's definition.
	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	extra *mapextra // optional fields
}
```

- ch 是一个指向 [runtime.hchan](https://draveness.me/golang/tree/runtime.hchan) 结构体的指针；

```go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}
```



相比与复杂的 make 关键字，new 的功能就简单多了，它只能接收类型作为参数然后返回一个指向该类型的指针：

```go
package main

import "fmt"

func main() {
	i := new(int)

	var v int
	j := &v
	fmt.Printf("i j 是相同类型: %v\n", i == j)
	fmt.Printf("i: %T    j: %T\n", i, j)
}

// i j 是相同类型: false
// i: *int    j: *int
```

上述代码片段中的两种不同初始化方法是等价的，它们都会创建一个指向 int 零值的指针。



```go
// The make built-in function allocates and initializes an object of type
// slice, map, or chan (only). Like new, the first argument is a type, not a
// value. Unlike new, make's return type is the same as the type of its
// argument, not a pointer to it. The specification of the result depends on
// the type:
//	Slice: The size specifies the length. The capacity of the slice is
//	equal to its length. A second integer argument may be provided to
//	specify a different capacity; it must be no smaller than the
//	length. For example, make([]int, 0, 10) allocates an underlying array
//	of size 10 and returns a slice of length 0 and capacity 10 that is
//	backed by this underlying array.
//	Map: An empty map is allocated with enough space to hold the
//	specified number of elements. The size may be omitted, in which case
//	a small starting size is allocated.
//	Channel: The channel's buffer is initialized with the specified
//	buffer capacity. If zero, or the size is omitted, the channel is
//	unbuffered.
func make(t Type, size ...IntegerType) Type

// The new built-in function allocates memory. The first argument is a type,
// not a value, and the value returned is a pointer to a newly
// allocated zero value of that type.
func new(Type) *Type

// The complex built-in function constructs a complex value from two
// floating-point values. The real and imaginary parts must be of the same
// size, either float32 or float64 (or assignable to them), and the return
// value will be the corresponding complex type (complex64 for float32,
// complex128 for float64).
```

## make 

在编译期间的类型检查阶段，Go 语言会将代表 make 关键字的 OMAKE 节点根据参数类型的不同转换成了 OMAKESLICE、OMAKEMAP 和 OMAKECHAN 三种不同类型的节点，这些节点会调用不同的运行时函数来初始化相应的数据结构。

## new 
编译器会在中间代码生成阶段通过以下两个函数处理该关键字：

- cmd/compile/internal/gc.callnew 会将关键字转换成 ONEWOBJ 类型的节点2；
- cmd/compile/internal/gc.state.expr 会根据申请空间的大小分两种情况处理：
    - 如果申请的空间为 0，就会返回一个表示空指针的 zerobase 变量；
    - 在遇到其他情况时会将关键字转换成 runtime.newobject 函数：

runtime.newobject 函数会获取传入类型占用空间的大小，调用 runtime.mallocgc 在堆上申请一片内存空间并返回指向这片内存空间的指针：

```go
func newobject(typ *_type) unsafe.Pointer {
	return mallocgc(typ.size, typ, true)
}
```

