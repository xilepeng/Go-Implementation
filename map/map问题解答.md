
1. [map 是线程安全的吗?](#map-是线程安全的吗)
2. [可以边遍历边删除吗?](#可以边遍历边删除吗)
3. [map 的删除过程是怎样的?](#map-的删除过程是怎样的)
4. [可以对 map 的元素取地址吗?](#可以对-map-的元素取地址吗)
5. [map 的扩容过程是怎样的 ?](#map-的扩容过程是怎样的-)


## map 是线程安全的吗?

map 不是线程安全的。
在查找、赋值、遍历、删除的过程中都会检测写标志，一旦发现写标志置位（等于1），则直接 panic。赋值和删除函数在检测完写标志是复位之后，先将写标志位置位，才会进行之后的操作。
检测写标志：

```go
	// 如果多线程读写，直接抛出异常
	// 并发检查 go hashmap 不支持并发访问
	if h.flags&hashWriting != 0 {
		throw("concurrent map read and map write")
	}
```

设置写标志：

```go
	// Set hashWriting after calling t.hasher, since t.hasher may panic,
	// in which case we have not actually done a write.
	// 在计算完 hash 值以后立即设置 hashWriting 变量的值，如果在计算 hash 值的过程中没有完全写完，可能会导致 panic
	h.flags ^= hashWriting
```




## 可以边遍历边删除吗?

map 并不是一个线程安全的数据结构。同时读写一个 map 是未定义的行为，如果被检测到，会直接 panic。
上面说的是发生在多个协程同时读写同一个 map 的情况下。 如果在同一个协程内边遍历边删除，并不会检测到同时读写，理论上是可以这样做的。但是，遍历的结果就可能不会是相同的了，有可能结果遍历结果集中包含了删除的 key，也有可能不包含，这取决于删除 key 的时间：是在遍历到 key 所在的 bucket 时刻前或者后。
一般而言，这可以通过读写锁来解决：sync.RWMutex。
读之前调用 RLock() 函数，读完之后调用 RUnlock() 函数解锁；写之前调用 Lock() 函数，写完之后，调用 Unlock() 解锁。
另外，sync.Map 是线程安全的 map，也可以使用。



## map 的删除过程是怎样的?

map 的删除操作底层的执行函数是 mapdelete：

```go
func mapdelete(t *maptype, h *hmap, key unsafe.Pointer) { }
```
首先会检查 h.flags 标志，如果发现写标位是 1，直接 panic，因为这表明有其他协程同时在进行写操作。
计算 key 的哈希，找到落入的 bucket。检查此 map 如果正在扩容的过程中，直接触发一次搬迁操作。
删除操作同样是两层循环，核心还是找到 key 的具体位置。寻找过程都是类似的，在 bucket 中挨个 cell 寻找。
找到对应位置后，对 key 或者 value 进行“清零”操作：

```go
// Only clear key if there are pointers in it.
// 对 key 清零
			if t.indirectkey() {
				// key 的指针置空
				*(*unsafe.Pointer)(k) = nil
			} else if t.key.ptrdata != 0 {
				// 清除 key 的内存
				memclrHasPointers(k, t.key.size)
			}
// 对 value 清零			
			if t.indirectelem() {
				// value 的指针置空
				*(*unsafe.Pointer)(e) = nil
			} else if t.elem.ptrdata != 0 {
				// 清除 value 的内存
				memclrHasPointers(e, t.elem.size)
			} else {
				memclrNoHeapPointers(e, t.elem.size)
			}

// 清空 tophash 里面的值
			b.tophash[i] = emptyOne

// map 里面 key 的总个数减1
			h.count--
```
最后，将 count 值减 1，将对应位置的 tophash 值置成 Empty。

## 可以对 map 的元素取地址吗?

无法对 map 的 key 或 value 进行取址。以下代码不能通过编译：

```go
package main

import "fmt"

func main() {
    m := make(map[string]int)

    fmt.Println(&m["qcrao"])
}
```
编译报错：
```
➜  map git:(main) ✗ go run main.go 
# command-line-arguments
./main.go:8:18: invalid operation: cannot take address of m["qcrao"] (map index expression of type int)
```

如果通过其他 hack 的方式，例如 unsafe.Pointer 等获取到了 key 或 value 的地址，也不能长期持有，因为一旦发生扩容，key 和 value 的位置就会改变，之前保存的地址也就失效了。


























## map 的扩容过程是怎样的 ?


loadFactor := count / (2^B) > 6.5   -> 翻倍扩容 hmap.B ++

使用了太多的溢出桶时（溢出桶使用的太多会导致map处理速度降低）。
loadFactor 没超标                    -> 等量扩容
noverflow  较多

B <= 15 noverflow >= 2^B
B >  15 noverflow >= 2^15

**等量扩容有啥用？**
迁移到新桶，排列的更加紧凑，从而减少溢出桶的使用，这就是等量扩容的意义。



低B位决定哪个桶，高8位决定桶里哪个曹cell。

```go
loadFactor := count / (2^B) > 6.5
```

count 就是 map 的元素个数，2^B 表示 bucket 数量。

再来说触发 map 扩容的时机：在向 map 插入新 key 的时候，会进行条件检测，符合下面这 2 个条件，就会触发扩容：

1. 装载因子超过阈值，源码里定义的阈值是 6.5。
2. overflow 的 bucket 数量过多：当 B 小于 15，也就是 bucket 总数 2^B 小于 2^15 时，
如果 overflow 的 bucket 数量超过 2^B；当 B >= 15，也就是 bucket 总数 2^B 大于等于 2^15，
如果 overflow 的 bucket 数量超过 2^15。（key 太分散或太稀疏）


