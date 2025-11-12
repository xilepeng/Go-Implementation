package main

import "fmt"

func main() {
	// 创建两个切片，并分别用两个整数进行初始化
	s1 := []int{1, 2}
	s2 := []int{3, 4}
	// 将两个切片追加在一起，并显示结果
	fmt.Printf("%v\n", append(s1, s2...)) // [1 2 3 4]
}

// pkgm 快速自动补全
