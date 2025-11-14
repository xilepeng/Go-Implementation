package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ch := make(chan os.Signal, 1)
	// 注册要由 os.Interrupt 信号通知的通道
	signal.Notify(ch, os.Interrupt)
	fmt.Println("awaiting signal ...")
	s := <-ch
	fmt.Println("Got signal:", s)
}

//awaiting signal ...
//Got signal: interrupt
