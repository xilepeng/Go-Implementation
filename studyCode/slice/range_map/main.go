package main

import "fmt"

func main() {
	// 创建一个映射，存储颜色以及颜色对应的十六进制代码
	colors := map[string]string{
		"AliceBlue":   "#f0f8ff",
		"Coral":       "#ff7F50",
		"DarkGray":    "#a9a9a9",
		"ForestGreen": "#228b22",
	}

	// 删除键为 Coral 的键值对
	delete(colors, "Coral")

	// 显示映射里的所有颜色
	for key, value := range colors {
		fmt.Printf("Key: %s Value: %s\n", key, value)
	}
}

// 对映射来说，range 返回的不是索引和值，而是键值对。

// Key: DarkGray Value: #a9a9a9
// Key: ForestGreen Value: #228b22
// Key: AliceBlue Value: #f0f8ff
