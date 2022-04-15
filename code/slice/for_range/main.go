package main

import "fmt"

func main() {
	slice := []int{10, 20, 30, 40}
	for index, value := range slice {
		fmt.Printf("value = %d , value-addr = %x , slice-addr = %x\n", value, &value, &slice[index])
	}
}

// 当迭代切片时，关键字 range 会返回两个值。第一个值是当前迭代到的索引位置，第二个 值是该位置对应元素值的一份副本

// ➜  demo git:(main) ✗ go run main.go
// value = 10 , value-addr = c0000b2008 , slice-addr = c0000b4000
// value = 20 , value-addr = c0000b2008 , slice-addr = c0000b4008
// value = 30 , value-addr = c0000b2008 , slice-addr = c0000b4010
// value = 40 , value-addr = c0000b2008 , slice-addr = c0000b4018

// 因为迭代返回的变量是一个迭代过程中根据切片依次赋值的新变量，所以 value 的地址总 是相同的。
// 要想获取每个元素的地址，可以使用切片变量和索引值。
// 关键字 range 总是会从切片头部开始迭代。如果想对迭代做更多的控制，依旧可以使用传 统的 for 循环
