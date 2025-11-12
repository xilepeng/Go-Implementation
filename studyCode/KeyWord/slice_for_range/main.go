package main

import "fmt"

func main() {
	slice := []int{10, 20, 30, 40}
	for index, copy_value := range slice {
		fmt.Printf("Value = %d , Value-Addr = %x , Elem-Addr = %x\n", copy_value, &copy_value, &slice[index])
	}
}

// Value = 10 , Value-Addr = c0000b2008 , Elem-Addr = c0000b6000
// Value = 20 , Value-Addr = c0000b2008 , Elem-Addr = c0000b6008
// Value = 30 , Value-Addr = c0000b2008 , Elem-Addr = c0000b6010
// Value = 40 , Value-Addr = c0000b2008 , Elem-Addr = c0000b6018

// 当迭代切片时，关键字 range 会返回两个值。第一个值是当前迭代到的索引位置，第二个 值是该位置对应元素值的一份副本
// 因为迭代返回的变量是一个迭代过程中根据切片依次赋值的新变量，所以 value 的地址总 是相同的。
// 要想获取每个元素的地址，可以使用切片变量和索引值。
// 关键字 range 总是会从切片头部开始迭代。如果想对迭代做更多的控制，依旧可以使用传统的 for 循环
