1. [map 的底层如何实现](#map-的底层如何实现)



**map 的底层实现原理是什么**

map 是由 key-value 对组成的；key 只会出现一次。

map 的设计也被称为 “The dictionary problem”，它的任务是设计一种数据结构用来维护一个集合的数据，并且可以同时对集合进行增删查改的操作。
最主要的数据结构有两种：**哈希查找表（Hash table）、搜索树（Search tree）**。

**哈希查找表**用一个哈希函数将 key 分配到不同的桶（bucket，也就是数组的不同 index）。这样，开销主要在哈希函数的计算以及数组的常数访问时间。在很多场景下，哈希查找表的性能很高。

哈希查找表一般会存在“碰撞”的问题，就是说不同的 key 被哈希到了同一个 bucket。一般有两种应对方法：**链表法和开放地址法**。链表法将一个 bucket 实现成一个链表，落在同一个 bucket 中的 key 都会插入这个链表。开放地址法则是碰撞发生后，通过一定的规律，在数组的后面挑选“空位”，用来放置新的 key。

搜索树法一般采用自平衡搜索树，包括：AVL 树，红黑树。面试时经常会被问到，甚至被要求手写红黑树代码，很多时候，面试官自己都写不上来，非常过分。

自平衡搜索树法的最差搜索效率是 O(logN)，而哈希查找表最差是 O(N)。当然，哈希查找表的平均查找效率是 O(1)，如果哈希函数设计的很好，最坏的情况基本不会出现。还有一点，遍历自平衡搜索树，返回的 key 序列，一般会按照从小到大的顺序；而哈希查找表则是乱序的。


## map 的底层如何实现


**Go 语言采用的是哈希查找表，并且使用链表解决哈希冲突。**

代码基于 
go version go1.17 darwin/amd64


Go 的 map 实现在 src/runtime/map.go 这个文件中。

map 底层实质还是一个 hash table。

先来看看 Go 定义了一些常量。

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
// A header for a Go map.
type hmap struct {
	
	count     int 	 // map 的长度
	flags     uint8
	B         uint8  // log以2为底，桶个数的对数 (总共能存 6.5 * 2^B 个元素)
	noverflow uint16 // 近似溢出桶的个数
	hash0     uint32 // 哈希种子

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	extra *mapextra // optional fields
}
```



