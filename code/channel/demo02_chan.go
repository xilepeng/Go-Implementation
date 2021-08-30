package main

import "fmt"

func main() {
	var ch chan bool
	fmt.Printf("%T, %v\n", ch, ch)
	ch = make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("子 goroutine 中 i:", i)
		}
		ch <- true
		fmt.Println("\t子 goroutine 结束")
	}()
	data := <-ch
	fmt.Println("data--->", data)
	fmt.Println("main is over...")
}
