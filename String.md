

1. [String 标准概念](#string-标准概念)
2. [为什么字符串不允许修改？](#为什么字符串不允许修改)
3. [[]byte转换成string一定会拷贝内存吗？](#byte转换成string一定会拷贝内存吗)
4. [string和[]byte如何取舍](#string和byte如何取舍)

## String 标准概念


Go标准库builtin给出了所有内置类型的定义。 源代码位于src/builtin/builtin.go，其中关于string的描述如下:

```go
// string is the set of all strings of 8-bit bytes, conventionally but not
// necessarily representing UTF-8-encoded text. A string may be empty, but
// not nil. Values of string type are immutable.
type string string
```

所以string是8比特字节的集合，通常但并不一定是UTF-8编码的文本。

另外，还提到了两点，非常重要：

- string可以为空（长度为0），但不会是nil；
- string对象不可以修改。


**string 数据结构**

源码包src/runtime/string.go:stringStruct定义了string的数据结构：

```go
type stringStruct struct {
	str unsafe.Pointer
	len int
}
```
其数据结构很简单：
- stringStruct.str：字符串的首地址；
- stringStruct.len：字符串的长度；

string数据结构跟切片有些类似，只不过切片还有一个表示容量的成员，事实上string和切片，准确的说是byte切片经常发生转换。这个后面再详细介绍。


**string操作**

**声明**

如下代码所示，可以声明一个string变量变赋予初值：



## 为什么字符串不允许修改？

- 像C++语言中的string，其本身拥有内存空间，修改string是支持的。
- 但Go的实现中，string不包含内存空间，只有一个内存的指针，这样做的好处是string变得非常轻量，可以很方便的进行传递而不用担心内存拷贝。

因为string通常指向字符串字面量，而字符串字面量存储位置是只读段，而不是堆或栈上，所以才有了string不可修改的约定。

## []byte转换成string一定会拷贝内存吗？

byte切片转换成string的场景很多，为了性能上的考虑，有时候只是临时需要字符串的场景下，byte切片转换成string时并不会拷贝内存，而是直接返回一个string，这个string的指针(string.str)指向切片的内存。

比如，编译器会识别如下临时场景：

- 使用m[string(b)]来查找map（map是string为key，临时把切片b转成string）；
- 字符串拼接，如"<" + "string(b)" + ">"；
- 字符串比较：string(b) == "foo"
因为是临时把byte切片转换成string，也就避免了因byte切片同容改成而导致string引用失败的情况，所以此时可以不必拷贝内存新建一个string。


## string和[]byte如何取舍

string和[]byte都可以表示字符串，但因数据结构不同，其衍生出来的方法也不同，要跟据实际应用场景来选择。

**string 擅长的场景：**

- 需要字符串比较的场景；
- 不需要nil字符串的场景；

**[]byte擅长的场景：**

- 修改字符串的场景，尤其是修改粒度为1个字节；
- 函数返回值，需要用nil表示含义的场景；
- 需要切片操作的场景；

虽然看起来string适用的场景不如[]byte多，但因为string直观，在实际应用中还是大量存在，在偏底层的实现中[]byte使用更多。





参考：

[Go专家编程](https://books.studygolang.com/GoExpertProgramming/chapter01/1.6-string.html)





   
