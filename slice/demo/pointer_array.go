package main

import "fmt"

func main() {
	// 复制数组指针，只会复制指针的值，而不会复制指针所指向的值

	// 声明第一个包含 3 个元素的指向字符串的指针数组
	var array1 [3]*string
	// 声明第二个包含 3 个元素的指向字符串的指针数组
	// 使用字符串指针初始化这个数组
	array2 := [3]*string{new(string), new(string), new(string)}
	// 使用颜色为每个元素赋值
	*array2[0] = "Red"
	*array2[1] = "Blue"
	*array2[2] = "Green"
	// 将 array2 复制给 array1
	array1 = array2
	// 复制之后，两个数组指向同一组字符串
	for i := 0; i < len(array1); i++ {
		fmt.Printf("array1_addr = %x, array2_addr = %x \n", array1[i], array2[i])
	}
}
