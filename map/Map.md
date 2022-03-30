1. [map 的底层如何实现](#map-的底层如何实现)
2. [1.新建 Map](#1新建-map)
3. [2. 查找 Key](#2-查找-key)
4. [3. 插入 Key](#3-插入-key)
5. [4. 删除 Key](#4-删除-key)



**map 的底层实现原理是什么**

map 是由 key-value 对组成的；key 只会出现一次。

map 的设计也被称为 “The dictionary problem”，它的任务是设计一种数据结构用来维护一个集合的数据，并且可以同时对集合进行增删查改的操作。
最主要的数据结构有两种：**哈希查找表（Hash table）、搜索树（Search tree）**。

**哈希查找表**用一个哈希函数将 key 分配到不同的桶（bucket，也就是数组的不同 index）。这样，开销主要在哈希函数的计算以及数组的常数访问时间。在很多场景下，哈希查找表的性能很高。

哈希查找表一般会存在“碰撞”的问题，就是说不同的 key 被哈希到了同一个 bucket。一般有两种应对方法：**链表法和开放地址法**。链表法将一个 bucket 实现成一个链表，落在同一个 bucket 中的 key 都会插入这个链表。开放地址法则是碰撞发生后，通过一定的规律，在数组的后面挑选“空位”，用来放置新的 key。

搜索树法一般采用自平衡搜索树，包括：AVL 树，红黑树。

自平衡搜索树法的最差搜索效率是 O(logN)，而哈希查找表最差是 O(N)。当然，哈希查找表的平均查找效率是 O(1)，如果哈希函数设计的很好，最坏的情况基本不会出现。还有一点，遍历自平衡搜索树，返回的 key 序列，一般会按照从小到大的顺序；而哈希查找表则是乱序的。


**为什么要用 map**

从 Go 语言官方博客摘录一段话：

"One of the most useful data structures in computer science is the hash table. Many hash table implementations exist with varying properties, but in general they offer fast lookups, adds, and deletes. Go provides a built-in map type that implements a hash table."

hash table 是计算机数据结构中一个最重要的设计。大部分 hash table 都实现了快速查找、添加、删除的功能。Go 语言内置的 map 实现了上述所有功能。

**因为它太强大了，各种增删查改的操作效率非常高。**



## map 的底层如何实现


**Go 语言 map 采用的是哈希查找表，并且使用链表（拉链法）解决哈希冲突。**

代码基于 GOVERSION="go1.18"

```go
➜  ~ go version
go version go1.18 darwin/amd64
```





Go 的 map 实现在 src/runtime/map.go 这个文件中。

map 底层实质还是一个 hash table。

先来看看 Go 定义了一些常量。

```shell
➜  ~ cd /usr/local/go/src
➜  src code .
或
➜  src atom .
```

```go

const (
	// 一个桶 bucket 里面最多可以装的键值对的个数，8对
	bucketCntBits = 3
	bucketCnt     = 1 << bucketCntBits	// 1<<3 == 2^3 == 8

	// 触发扩容操作的最大装载因子的临界值是 6.5
	// Represent as loadFactorNum/loadFactorDen, to allow integer math.
	loadFactorNum = 13
	loadFactorDen = 2

	// 为了保持内联，键 和 值 的最大长度都是128字节，如果超过了128个字节，就存储它的指针
	maxKeySize  = 128
	maxElemSize = 128

	// 数据偏移应该是 bmap 的整数倍，但是需要正确的对齐。
	dataOffset = unsafe.Offsetof(struct {
		b bmap
		v int64
	}{}.v)

	// tophash 的一些值
	emptyRest      = 0 // cell 是空的（没有键值对），并且在更高的索引或溢出处不再有非空 cell 单元格.
	emptyOne       = 1 // cell 是空的
	evacuatedX     = 2 // 键值对有效，并且已经迁移了一个表的前半段
	evacuatedY     = 3 // 键值对有效，并且已经迁移了一个表的后半段
	evacuatedEmpty = 4 // cell是空的，并且桶内的键值被迁移走了
	minTopHash     = 5 // 最小的 tophash

	// flags 标记
	iterator     = 1 // 当前桶的迭代子
	oldIterator  = 2 // 旧桶的迭代子
	hashWriting  = 4 // 一个goroutine正在写入map
	sameSizeGrow = 8 // 当前 map 增长到新 map 相同尺寸

	// 迭代子检查桶ID的哨兵
	noCheck = 1<<(8*sys.PtrSize) - 1
)
```


这里值得说明的一点是触发扩容操作的临界值6.5是怎么得来的。这个值太大会导致overflow buckets过多，查找效率降低，过小会浪费存储空间。

据 Google 开发人员称，这个值是一个测试的程序，测量出来选择的一个经验值。

**loadFactor = loadFactorNum / loadFactorDen = 13 / 2 = 6.5**

```go
loadFactor := count / (2^B)
```

loadFactorNum：map 的元素个数 count；

loadFactorDen：2^B 表示 bucket 数量。


```go
// Picking loadFactor: too large and we have lots of overflow
// buckets, too small and we waste a lot of space. I wrote
// a simple program to check some stats for different loads:
// (64-bit, 8 byte keys and elems)
//  loadFactor    %overflow  bytes/entry     hitprobe    missprobe
//        4.00         2.13        20.77         3.00         4.00
//        4.50         4.05        17.30         3.25         4.50
//        5.00         6.85        14.77         3.50         5.00
//        5.50        10.55        12.94         3.75         5.50
//        6.00        15.27        11.67         4.00         6.00
          6.50        20.90        10.79         4.25         6.50
//        7.00        27.14        10.15         4.50         7.00
//        7.50        34.03         9.73         4.75         7.50
//        8.00        41.10         9.40         5.00         8.00
```

%overflow ：
溢出率，平均一个 bucket 有多少个 键值kv 的时候会溢出。

bytes/entry ：
平均存一个 键值kv 需要额外存储多少字节的数据。

hitprobe ：
查找一个存在的 key 平均查找次数。

missprobe ：
查找一个不存在的 key 平均查找次数。


**经过这几组测试数据，最终选定 6.5 作为临界的装载因子。**

渐进式扩容：键值对迁移的时间分摊到多次哈希表操作中的方式，可避免一次性扩容带来的性能瞬时抖动


选择桶时用的是 ”与“ 运算的方法

**Go 中 map header 的定义：**

```go
// A header for a Go map.

type hmap struct {
	count     int // map 的长度(键值对数目)
	flags     uint8
	B         uint8  // B = log_2 buckets  log以2为底，桶个数的对数 (总共能存 6.5 * 2^B 个元素)
	noverflow uint16 // 近似溢出桶的个数
	hash0     uint32 // 哈希种子

	buckets    unsafe.Pointer // 有 buckets = 2^B 个桶的数组. count==0 的时候，这个数组为 nil.
	oldbuckets unsafe.Pointer // 旧的桶数组一半的元素
	nevacuate  uintptr        // 扩容增长过程中的计数器(即将迁移的旧桶编号)

	extra *mapextra // 可选字段
}

```

1. count 表示当前哈希表中的元素数量；
2. B 表示当前哈希表持有的 buckets 数量，但是因为哈希表中桶的数量都 2 的倍数，所以该字段会存储对数，也就是 len(buckets) == 2^B；
3. hash0 是哈希的种子，它能为哈希函数的结果引入随机性，这个值在创建哈希表时确定，并在调用哈希函数时作为参数传入；
4. oldbuckets 是哈希在扩容时用于保存之前 buckets 的字段，它的大小是当前 buckets 的一半；

B 是 buckets 数组的长度的对数，也就是说 buckets 数组的长度就是 2^B。bucket 里面存储了 key 和 value

```go
buckets = 2^B   B = log_2 buckets
```
解释：如果2的 B 次方等于 buckets，那么 B 叫做以2为底 buckets 的对数。





hmap 的最后一个字段是一个指向 mapextra 结构的指针，它的定义如下：

```go
type mapextra struct {
	overflow    *[]*bmap
	oldoverflow *[]*bmap

	nextOverflow *bmap
}
```
如果一个键值对没有找到对应的指针，那么就会把它们先存到溢出桶
overflow 里面。在 mapextra 中还有一个指向下一个可用的溢出桶的指针。

溢出桶 overflow 是一个数组指针（是一个指针变量，占有内存中一个指针的存储空间），里面存了指向 *bmap 数组的指针。overflow[0] 里面装的是 hmap.buckets 。overflow[1] 里面装的是 hmap.oldbuckets。






再看看桶的数据结构的定义，bmap 就是 Go 中 map 里面桶对应的结构体类型。

```go
// A bucket for a Go map.
type bmap struct {
	tophash [bucketCnt]uint8
}
```

在运行期间，runtime.bmap 结构体其实不止包含 tophash 字段，因为哈希表中可能存储不同类型的键值对，而且 Go 语言也不支持泛型，所以键值对占据的内存空间大小只能在编译时进行推导。runtime.bmap 中的其他字段在运行时也都是通过计算内存地址的方式访问的，所以它的定义中就不包含这些字段，不过我们能根据编译期间的
= [cmd/compile/internal/gc.bmap](https://github.com/golang/go/blob/ac0ba6707c1655ea4316b41d06571a0303cc60eb/src/cmd/compile/internal/gc/reflect.go#L83)
函数重建它的结构：

```go
type bmap struct {
    topbits  [8]uint8
    keys     [8]keytype
    values   [8]valuetype
    pad      uintptr
    overflow uintptr
}
```


桶的定义比较简单，里面就只是包含了一个 uint8 类型的数组，里面包含8个元素。这8个元素存储的是 hash 值的高8位。

在 tophash 之后的内存布局里还有2块内容。紧接着 tophash 之后的是8对 键值 key- value 对。并且排列方式是 8个 key 和 8个 value 放在一起。

8对 键值 key- value 对结束以后紧接着一个 overflow 指针，指向下一个 bmap。从此也可以看出 Go 中 map是用链表的方式处理 hash 冲突的。


为何 Go 存储键值对的方式不是普通的 key/value、key/value、key/value……这样存储的呢？它是键 key 都存储在一起，然后紧接着是 值value 都存储在一起，为什么会这样呢？



在 Redis 中，当使用 REDIS_ENCODING_ZIPLIST 编码哈希表时， 程序通过将键和值一同推入压缩列表， 从而形成保存哈希表所需的键-值对结构，如上图。新添加的 key-value 对会被添加到压缩列表的表尾。

这种结构有一个弊端，如果存储的键和值的类型不同，在内存中布局中所占字节不同的话，就需要对齐。比如说存储一个 map[int64]int8 类型的字典。

Go 为了节约内存对齐的内存消耗，于是把它设计成上图所示那样。

如果 map 里面存储了上万亿的大数据，这里节约出来的内存空间还是比较可观的。






bmap 就是我们常说的“桶”，桶里面会最多装 8 个 key，这些 key 之所以会落入同一个桶，是因为它们经过哈希计算后，哈希结果是“一类”的。
在桶内，又会根据 key 计算出来的 hash 值的高 8 位来决定 key 到底落入桶内的哪个位置（一个桶内最多有8个位置）。

来一个整体的图：



![hashmap-bmap](images/hashmap-bmap.png)






hmap 的最后一个字段是一个指向 mapextra 结构的指针，它的定义如下：

当 map 的 key 和 value 都不是指针，并且 size 都小于 128 字节的情况下，会把 bmap 标记为不含指针，这样可以避免 gc 时扫描整个 hmap。
但是，我们看 bmap 其实有一个 overflow 的字段，是指针类型的，破坏了 bmap 不含指针的设想，这时会把 overflow 移动到 extra 字段来。
```go
type mapextra struct {
	overflow    *[]*bmap
	oldoverflow *[]*bmap
  // nextOverflow 包含空闲的 overflow bucket，这是预分配的 bucket
	nextOverflow *bmap
}
```



## 1.新建 Map


```go
func makemap(t *maptype, hint int, h *hmap) *hmap {
	// 1. 计算哈希占用的内存是否溢出或者超出能分配的最大值；
	mem, overflow := math.MulUintptr(uintptr(hint), t.bucket.size)
	if overflow || mem > maxAlloc {
		hint = 0
	}

	// 初始化 hmap
	if h == nil {
		h = new(hmap)
	}
	// 2. 调用 runtime.fastrand 获取一个随机的哈希种子；
	h.hash0 = fastrand()

	// Find the size parameter B which will hold the requested # of elements.
	// For hint < 0 overLoadFactor returns false since hint < bucketCnt.
	// 3. 根据传入的 hint 计算出需要的最小需要的桶的数量；
	B := uint8(0)
	for overLoadFactor(hint, B) {
		B++
	}
	h.B = B

	// 分配内存并初始化哈希表
	// 如果此时 B = 0，那么 hmap 中的 buckets 字段稍后分配
	// 如果 hint 值很大，初始化这块内存需要一段时间。
	if h.B != 0 {
		var nextOverflow *bmap
		// 初始化 bucket 和 nextOverflow
		// 4. 使用 runtime.makeBucketArray 创建用于保存桶的数组；
		h.buckets, nextOverflow = makeBucketArray(t, h.B, nil)
		if nextOverflow != nil {
			h.extra = new(mapextra)
			h.extra.nextOverflow = nextOverflow
		}
	}

	return h
}
```

注意，这个函数返回的结果：*hmap，它是一个指针，而我们之前讲过的 makeslice 函数返回的是 Slice 结构体：

```go
// runtime/slice.go
type slice struct {
 array unsafe.Pointer
 len   int
 cap   int
}
```

```go
func makeslice(et *_type, len, cap int) unsafe.Pointer {
	return mallocgc(mem, et, true)
}
```
结构体内部包含底层的数据指针。

makemap 和 makeslice 的区别，带来一个不同点：当 map 和 slice 作为函数参数时，在函数参数内部对 map 的操作会影响 map 自身；而对 slice 却不会（之前讲 slice 的文章里有讲过）。

主要原因：一个是指针（*hmap），一个是结构体（slice）。Go 语言中的函数传参都是值传递，在函数内部，参数会被 copy 到本地。*hmap指针 copy 完之后，仍然指向同一个 map，因此函数内部对 map 的操作会影响实参。而 slice 被 copy 后，会成为一个新的 slice，对它进行的操作不会影响到实参。






[runtime.makeBucketArray](https://github.com/golang/go/blob/ac0ba6707c1655ea4316b41d06571a0303cc60eb/src/runtime/map.go#L344)
会根据传入的 B 计算出的需要创建的桶数量并在内存中分配一片连续的空间用于存储数据：
```go
// makeBucketArray initializes a backing array for map buckets.
// 1<<b is the minimum number of buckets to allocate.
// dirtyalloc should either be nil or a bucket array previously
// allocated by makeBucketArray with the same t and b parameters.
// If dirtyalloc is nil a new backing array will be alloced and
// otherwise dirtyalloc will be cleared and reused as backing array.

// 4. 使用 runtime.makeBucketArray 创建用于保存桶的数组；
func makeBucketArray(t *maptype, b uint8, dirtyalloc unsafe.Pointer) (buckets unsafe.Pointer, nextOverflow *bmap) {
	base := bucketShift(b)
	nbuckets := base
	// For small b, overflow buckets are unlikely.
	// Avoid the overhead of the calculation.
	// 当桶的数量小于 2^4时，由于数据较少、使用溢出桶的可能性较低，会省略创建的过程以减少额外开销
	if b >= 4 {
		// Add on the estimated number of overflow buckets
		// required to insert the median number of elements
		// used with this value of b.
		// 当桶的数量多于 2^4 时，会额外创建 2^𝐵−4 个溢出桶；
		nbuckets += bucketShift(b - 4)
		sz := t.bucket.size * nbuckets
		up := roundupsize(sz)
		// 如果申请 sz 大小的桶，系统只能返回 up 大小的内存空间，那么桶的个数为 up / t.bucket.size
		if up != sz {
			nbuckets = up / t.bucket.size
		}
	}

	if dirtyalloc == nil {
		buckets = newarray(t.bucket, int(nbuckets))
	} else {
		// dirtyalloc was previously generated by
		// the above newarray(t.bucket, int(nbuckets))
		// but may not be empty.
		buckets = dirtyalloc
		size := t.bucket.size * nbuckets
		if t.bucket.ptrdata != 0 {
			memclrHasPointers(buckets, size)
		} else {
			memclrNoHeapPointers(buckets, size)
		}
	}
	// 当 b > 4 并且计算出来桶的个数与 1 << b 个数不等的时候，
	if base != nbuckets {
		// We preallocated some overflow buckets.
		// To keep the overhead of tracking these overflow buckets to a minimum,
		// we use the convention that if a preallocated overflow bucket's overflow
		// pointer is nil, then there are more available by bumping the pointer.
		// We need a safe non-nil pointer for the last overflow bucket; just use buckets.
		// 此时 nbuckets 比 base 大，那么会预先分配 nbuckets - base 个 nextOverflow 桶
		nextOverflow = (*bmap)(add(buckets, base*uintptr(t.bucketsize)))
		last := (*bmap)(add(buckets, (nbuckets-1)*uintptr(t.bucketsize)))
		last.setoverflow(t, (*bmap)(buckets))
	}
	return buckets, nextOverflow
}
```
- 当桶的数量小于 2^4时，由于数据较少、使用溢出桶的可能性较低，会省略创建的过程以减少额外开销;
- 当桶的数量多于 2^4 时，会额外创建 2^𝐵−4 个溢出桶；

根据上述代码，我们能确定在正常情况下，正常桶和溢出桶在内存中的存储空间是连续的，只是被 [runtime.hmap](https://github.com/golang/go/blob/41d8e61a6b9d8f9db912626eb2bbc535e929fefc/src/runtime/map.go#L115) 中的不同字段引用，当溢出桶数量较多时会通过 [runtime.newobject](https://github.com/golang/go/blob/41d8e61a6b9d8f9db912626eb2bbc535e929fefc/src/runtime/malloc.go#L1176) 创建新的溢出桶。


这里的 newarray 就已经是 mallocgc 了。

从上述代码里面可以看出，只有当 B >=4 的时候，makeBucketArray 才会生成 nextOverflow 指针指向 bmap，从而在 Map 生成 hmap 的时候才会生成 mapextra 。

- 当 B = 3 ( B < 4 ) 的时候，初始化 hmap 只会生成8个桶。
- 当 B = 4 ( B >= 4 ) 的时候，初始化 hmap 的时候还会额外生成 mapextra ，并初始化 nextOverflow。mapextra 的 nextOverflow 指针会指向第16个桶结束，第17个桶的首地址。第17个桶（从0开始，也就是下标为16的桶）的 bucketsize - sys.PtrSize 地址开始存一个指针，这个指针指向当前整个桶的首地址。这个指针就是 bmap 的 overflow 指针。






**哈希函数**

map 的一个关键点在于，哈希函数的选择。在程序启动时，会检测 cpu 是否支持 aes，如果支持，则使用 aes hash，否则使用 memhash。这是在函数 alginit() 中完成，位于路径：src/runtime/alg.go 下。

hash 函数，有加密型和非加密型。
加密型的一般用于加密数据、数字摘要等，典型代表就是 md5、sha1、sha256、aes256 这种；
非加密型的一般就是查找。在 map 的应用场景中，用的是查找。
选择 hash 函数主要考察的是两点：性能、碰撞概率。












## 2. 查找 Key


```go



// bucketShift returns 1<<b, optimized for code generation.
func bucketShift(b uint8) uintptr {
	// Masking the shift amount allows overflow checks to be elided.
	return uintptr(1) << (b & (goarch.PtrSize*8 - 1))
}

// bucketMask returns 1<<b - 1, optimized for code generation.
func bucketMask(b uint8) uintptr {
	return bucketShift(b) - 1
}
// hash & (1<<B - 1) 求出 key 在哪个桶
// hash & m 求出 key 在哪个桶


    // 比如 B=5，那 m 就是 2^5=31，二进制是全 1
    // 求 bucket 索引时，将 hash 与 m 相与，
    // 达到 bucket 位置下标由 hash 的低 8 位决定的效果


// tophash calculates the tophash value for hash.
func tophash(hash uintptr) uint8 {
	top := uint8(hash >> (goarch.PtrSize*8 - 8))
	// 如果 top 小于 minTopHash，就让它加上 minTopHash 的偏移。
	// 因为 0 - minTopHash 这区间的数都已经用来作为标记位了
	if top < minTopHash {
		top += minTopHash
	}
	return top
}
```


在 Go 中，如果字典里面查找一个不存在的 key ，查找不到并不会返回一个 nil ，而是返回当前类型的零值。比如，字符串就返回空字符串，int 类型就返回 0 。



```go

// mapaccess1 returns a pointer to h[key].  Never returns nil, instead
// it will return a reference to the zero object for the elem type if
// the key is not in the map.
// NOTE: The returned pointer may keep the whole map live, so don't
// hold onto it for very long.
func mapaccess1(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer {
	if raceenabled && h != nil {
		// 获取 caller 的 程序计数器 program counter
		callerpc := getcallerpc()
		pc := abi.FuncPCABIInternal(mapaccess1)
		racereadpc(unsafe.Pointer(h), callerpc, pc)
		raceReadObjectPC(t.key, key, callerpc, pc)
	}
	if msanenabled && h != nil {
		msanread(key, t.key.size)
	}
	if asanenabled && h != nil {
		asanread(key, t.key.size)
	}
	// 如果 h 什么都没有，返回零值
	if h == nil || h.count == 0 {
		if t.hashMightPanic() {
			t.hasher(key, 0) // see issue 23734
		}
		return unsafe.Pointer(&zeroVal[0])
	}
	// 如果多线程读写，直接抛出异常
	// 并发检查 go hashmap 不支持并发访问
	if h.flags&hashWriting != 0 {
		throw("concurrent map read and map write")
	}
	// 不同类型 key 使用的 hash 算法在编译期确定
	// 计算 key 的 hash 值, 加入 hash0 引入随机性
	hash := t.hasher(key, uintptr(h.hash0))
	m := bucketMask(h.B)
	// hash & (1<<B - 1) 求出 key 在哪个桶
	// hash & m 求出 key 在哪个桶
	// b 就是 bucket 的地址
	b := (*bmap)(add(h.buckets, (hash&m)*uintptr(t.bucketsize)))
	// oldbuckets 不为 nil，说明发生了扩容
	if c := h.oldbuckets; c != nil {
		// 如果不是等量扩容
		if !h.sameSizeGrow() {
			// There used to be half as many buckets; mask down one more power of two.
			// 如果 oldbuckets 未迁移完成 则找找 oldbuckets 中对应的 bucket(低 B-1 位)
			// 否则为 buckets 中的 bucket(低 B 位)
			// 把 mask 缩小 1 倍
			m >>= 1
		}
		// 求出 key 在老的 map 中的 bucket 位置
		oldb := (*bmap)(add(c, (hash&m)*uintptr(t.bucketsize)))
		if !evacuated(oldb) {
			// 如果 oldbuckets 桶存在，并且还没有扩容迁移，就在老的桶里面查找 key
			b = oldb
		}
	}
	// 取出 hash 值的高 8 位	
	top := tophash(hash)
bucketloop:
	for ; b != nil; b = b.overflow(t) {
		for i := uintptr(0); i < bucketCnt; i++ {
			// 如果 hash 的高8位和当前 key 记录的不一样，就找下一个
			// 这样比较很高效，因为只用比较高8位，不用比较所有的 hash 值
			// 如果高8位都不相同，hash 值肯定不同，但是高8位如果相同，那么就要比较整个 hash 值了
			if b.tophash[i] != top {
				if b.tophash[i] == emptyRest {
					break bucketloop
				}
				continue
			}
			// 取出 key 值的方式是用偏移量，bmap 首地址 + i 个 key 值大小的偏移量
			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
			// 比较 key 值是否相等
			if t.indirectkey() {
				k = *((*unsafe.Pointer)(k))
			}
			if t.key.equal(key, k) {
				// 如果找到了 key，那么取出 value 值
				// 取出 value 值的方式是用偏移量，bmap 首地址 + 8 个 key 值大小的偏移量 + i 个 value 值大小的偏移量
				e := add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.elemsize))
				if t.indirectelem() {
					e = *((*unsafe.Pointer)(e))
				}
				return e
			}
		}
	}
	return unsafe.Pointer(&zeroVal[0])
}
```
图片引用[一缕殇流化隐半边冰霜](https://halfrost.com/go_map_chapter_one/)
![](images/select_key.png)

如上图，这是一个查找 key 的全过程。

首先计算出 key 对应的 hash 值，hash 值对 B 取余。

这里有一个优化点。m % n 这步计算，如果 n 是2的倍数，那么可以省去这一步取余操作。

```go
m % n = m & ( n - 1 )
```

这样优化就可以省去耗时的取余操作了。这里例子中计算完取出来是 0010 ，也就是2，于是对应的是桶数组里面的第3个桶。为什么是第3个桶呢？首地址指向第0个桶，往下偏移2个桶的大小，于是偏移到了第3个桶的首地址了，具体实现可以看上述代码。

- *hash 的低 B 位决定了桶数组里面的第几个桶*
- *hash 值的高8位决定了这个桶数组 bmap 里面 key 存在 tophash 数组的第几位了。*

如上图，hash 的高8位用来和 tophash 数组里面的每个值进行对比，如果高8位和 tophash[i] 不等，就直接比下一个。如果相等，则取出 bmap 里面对应完整的 key，再比较一次，看是否完全一致。


整个查找过程优先在 oldbucket 里面找(如果存在 lodbucket 的话)，找完再去新 bmap 里面找。

有人可能会有疑问，为何这里要加入 tophash 多一次比较呢？

tophash 的引入是为了加速查找的。由于它只存了 hash 值的高8位，比查找完整的64位要快很多。通过比较高8位，迅速找到高8位一致hash 值的索引，接下来再进行一次完整的比较，如果还一致，那么就判定找到该 key 了。

如果找到了 key 就返回对应的 value。如果没有找到，还会继续去 overflow 桶继续寻找，直到找到最后一个桶，如果还没有找到就返回对应类型的零值。





图片引用[码农桃花源](https://qcrao91.gitbook.io/go/map/map-de-di-ceng-shi-xian-yuan-li-shi-shi-mo)
![](images/select_key2.png)


上图中，假定 B = 5，所以 bucket 总数就是 2^5 = 32。

1. 首先计算出待查找 key 的哈希，
2. 使用低 5 位 00110，找到对应的 6 号 bucket，
3. 使用高 8 位 10010111，对应十进制 151，在 6 号 bucket 中寻找 tophash 值（HOB hash）为 151 的 key，找到了 2 号槽位，这样整个查找过程就结束了。


如果在 bucket 中没找到，并且 overflow 不为空，还要继续去 overflow bucket 中寻找，直到找到或是所有的 key 槽位都找遍了，包括所有的 overflow bucket。


## 3. 插入 Key

插入 key 的过程和查找 key 的过程大体一致。

```go

// Like mapaccess, but allocates a slot for the key if it is not present in the map.
func mapassign(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer {
	if h == nil {
		panic(plainError("assignment to entry in nil map"))
	}
	if raceenabled {
		callerpc := getcallerpc()
		pc := abi.FuncPCABIInternal(mapassign)
		racewritepc(unsafe.Pointer(h), callerpc, pc)
		raceReadObjectPC(t.key, key, callerpc, pc)
	}
	if msanenabled {
		msanread(key, t.key.size)
	}
	if asanenabled {
		asanread(key, t.key.size)
	}
	if h.flags&hashWriting != 0 {
		throw("concurrent map writes")
	}
	hash := t.hasher(key, uintptr(h.hash0))

	// Set hashWriting after calling t.hasher, since t.hasher may panic,
	// in which case we have not actually done a write.
	// 在计算完 hash 值以后立即设置 hashWriting 变量的值，如果在计算 hash 值的过程中没有完全写完，可能会导致 panic
	h.flags ^= hashWriting
	// 如果 hmap 的桶的个数为0，那么就新建一个桶
	if h.buckets == nil {
		h.buckets = newobject(t.bucket) // newarray(t.bucket, 1)
	}

again:
	// hash 值对 B 取余，求得所在哪个桶
	bucket := hash & bucketMask(h.B)
	// 如果还在扩容中，继续扩容
	if h.growing() {
		growWork(t, h, bucket)
	}
	// 根据 hash 值的低 B 位找到位于哪个桶
	b := (*bmap)(add(h.buckets, bucket*uintptr(t.bucketsize)))
	// 计算 hash 值的高 8 位
	top := tophash(hash)

	var inserti *uint8
	var insertk unsafe.Pointer
	var elem unsafe.Pointer
bucketloop:
	for {
		// 遍历当前桶所有键值，查找 key 对应的 value
		for i := uintptr(0); i < bucketCnt; i++ {
			if b.tophash[i] != top {
				if isEmpty(b.tophash[i]) && inserti == nil {
					// 如果往后找都没有找到，这里先记录一个标记，方便找不到以后插入到这里
					inserti = &b.tophash[i]
					// 计算出偏移 i 个 key 值的位置
					insertk = add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
					// 计算出 val 所在的位置，当前桶的首地址 + 8个 key 值所占的大小 + i 个 value 值所占的大小
					elem = add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.elemsize))
				}
				if b.tophash[i] == emptyRest {
					break bucketloop
				}
				continue
			}
			// 依次取出 key 值
			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
			// 如果 key 值是一个指针，那么就取出改指针对应的 key 值
			if t.indirectkey() {
				k = *((*unsafe.Pointer)(k))
			}
			// 比较 key 值是否相等
			if !t.key.equal(key, k) {
				continue
			}
			// already have a mapping for key. Update it.
			// 如果需要更新，那么就把 t.key 拷贝从 k 拷贝到 key
			if t.needkeyupdate() {
				typedmemmove(t.key, k, key)
			}
			// 计算出 val 所在的位置，当前桶的首地址 + 8个 key 值所占的大小 + i 个 value 值所占的大小
			elem = add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.elemsize))
			goto done
		}
		ovf := b.overflow(t)
		if ovf == nil {
			break
		}
		b = ovf
	}

	// Did not find mapping for key. Allocate new cell & add entry.

	// If we hit the max load factor or we have too many overflow buckets,
	// and we're not already in the middle of growing, start growing.
	// 没有找到当前的 key 值，并且检查最大负载因子，如果达到了最大负载因子，或者存在很多溢出的桶
	if !h.growing() && (overLoadFactor(h.count+1, h.B) || tooManyOverflowBuckets(h.noverflow, h.B)) {
		// 开始扩容
		hashGrow(t, h)
		goto again // Growing the table invalidates everything, so try again
	}
    // 如果找不到一个空的位置可以插入 key 值
	if inserti == nil {
		// The current bucket and all the overflow buckets connected to it are full, allocate a new one.
		// 意味着当前桶已经全部满了，那么就生成一个新的
		newb := h.newoverflow(t, b)
		inserti = &newb.tophash[0]
		insertk = add(unsafe.Pointer(newb), dataOffset)
		elem = add(insertk, bucketCnt*uintptr(t.keysize))
	}

	// store new key/elem at insert position
	if t.indirectkey() {
		// 如果是存储 key 值的指针，这里就用 insertk 存储 key 值的地址
		kmem := newobject(t.key)
		*(*unsafe.Pointer)(insertk) = kmem
		insertk = kmem
	}
	if t.indirectelem() {
		vmem := newobject(t.elem)
		*(*unsafe.Pointer)(elem) = vmem
	}
	// 将 t.key 从 insertk 拷贝到 key 的位置
	typedmemmove(t.key, insertk, key)
	*inserti = top
	// hmap 中保存的总 key 值的数量 + 1
	h.count++

done:
	// 禁止并发写
	if h.flags&hashWriting == 0 {
		throw("concurrent map writes")
	}
	h.flags &^= hashWriting
	if t.indirectelem() {
		// 如果 value 里面存储的是指针，那么取值该指针指向的 value 值
		elem = *((*unsafe.Pointer)(elem))
	}
	return elem
}
```
插入 key 的过程中和查找 key 有几点不同，需要注意：

1. 如果找到要插入的 key ，只需要直接更新对应的 value 值就好了。
2. 如果没有在 bmap 中没有找到待插入的 key ，这么这时分几种情况。
情况一: bmap 中还有空位，在遍历 bmap 的时候预先标记空位，一旦查找结束也没有找到 key，就把 key 放到预先遍历时候标记的空位上。
情况二：bmap中已经没有空位了。这个时候 bmap 装的很满了。此时需要检查一次最大负载因子是否已经达到了。如果达到了，立即进行扩容操作。扩容以后在新桶里面插入 key，流程和上述的一致。如果没有达到最大负载因子，那么就在新生成一个 bmap，并把前一个 bmap 的 overflow 指针指向新的 bmap。
3. 在扩容过程中，oldbucke t是被冻结的，查找 key 时会在
oldbucket 中查找，但不会在 oldbucket 中插入数据。如果在
oldbucket 是找到了相应的key，做法是将它迁移到新 bmap 后加入 evalucated 标记。

其他流程和查找 key 基本一致，这里就不再赘述了。




## 4. 删除 Key


```go

func mapdelete(t *maptype, h *hmap, key unsafe.Pointer) {
	if raceenabled && h != nil {
		callerpc := getcallerpc()
		pc := abi.FuncPCABIInternal(mapdelete)
		racewritepc(unsafe.Pointer(h), callerpc, pc)
		raceReadObjectPC(t.key, key, callerpc, pc)
	}
	if msanenabled && h != nil {
		msanread(key, t.key.size)
	}
	if asanenabled && h != nil {
		asanread(key, t.key.size)
	}
	if h == nil || h.count == 0 {
		if t.hashMightPanic() {
			t.hasher(key, 0) // see issue 23734
		}
		return
	}
	if h.flags&hashWriting != 0 {
		throw("concurrent map writes")
	}

	hash := t.hasher(key, uintptr(h.hash0))

	// Set hashWriting after calling t.hasher, since t.hasher may panic,
	// in which case we have not actually done a write (delete).
	h.flags ^= hashWriting

	bucket := hash & bucketMask(h.B)
	if h.growing() {
		growWork(t, h, bucket)
	}
	b := (*bmap)(add(h.buckets, bucket*uintptr(t.bucketsize)))
	bOrig := b
	top := tophash(hash)
search:
	for ; b != nil; b = b.overflow(t) {
		for i := uintptr(0); i < bucketCnt; i++ {
			if b.tophash[i] != top {
				if b.tophash[i] == emptyRest {
					break search
				}
				continue
			}
			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
			k2 := k
			if t.indirectkey() {
				k2 = *((*unsafe.Pointer)(k2))
			}
			if !t.key.equal(key, k2) {
				continue
			}
			// Only clear key if there are pointers in it.
			if t.indirectkey() {
				*(*unsafe.Pointer)(k) = nil
			} else if t.key.ptrdata != 0 {
				memclrHasPointers(k, t.key.size)
			}
			e := add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.elemsize))
			if t.indirectelem() {
				*(*unsafe.Pointer)(e) = nil
			} else if t.elem.ptrdata != 0 {
				memclrHasPointers(e, t.elem.size)
			} else {
				memclrNoHeapPointers(e, t.elem.size)
			}
			b.tophash[i] = emptyOne
			// If the bucket now ends in a bunch of emptyOne states,
			// change those to emptyRest states.
			// It would be nice to make this a separate function, but
			// for loops are not currently inlineable.
			if i == bucketCnt-1 {
				if b.overflow(t) != nil && b.overflow(t).tophash[0] != emptyRest {
					goto notLast
				}
			} else {
				if b.tophash[i+1] != emptyRest {
					goto notLast
				}
			}
			for {
				b.tophash[i] = emptyRest
				if i == 0 {
					if b == bOrig {
						break // beginning of initial bucket, we're done.
					}
					// Find previous bucket, continue at its last entry.
					c := b
					for b = bOrig; b.overflow(t) != c; b = b.overflow(t) {
					}
					i = bucketCnt - 1
				} else {
					i--
				}
				if b.tophash[i] != emptyOne {
					break
				}
			}
		notLast:
			h.count--
			// Reset the hash seed to make it more difficult for attackers to
			// repeatedly trigger hash collisions. See issue 25237.
			if h.count == 0 {
				h.hash0 = fastrand()
			}
			break search
		}
	}

	if h.flags&hashWriting == 0 {
		throw("concurrent map writes")
	}
	h.flags &^= hashWriting
}

```




































参考博客：
[一缕殇流化隐半边冰霜](https://halfrost.com/go_map_chapter_one/)
[码农桃花源](https://qcrao91.gitbook.io/go/map/map-de-di-ceng-shi-xian-yuan-li-shi-shi-mo)




