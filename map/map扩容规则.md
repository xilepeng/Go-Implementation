

## map 的扩容过程是怎样的 ?


loadFactor := count / (2^B) > 6.5   -> 翻倍扩容 hmap.B ++

使用了太多的溢出桶时（溢出桶使用的太多会导致map处理速度降低）。
loadFactor 没超标                    -> 等量扩容
noverflow  较多

B <= 15 noverflow >= 2^B
B >  15 noverflow >= 2^15

## 等量扩容有啥用？
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
