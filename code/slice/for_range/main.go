package main

import "fmt"

func main() {
	slice := []int{10, 20, 30, 40}
	for index, value := range slice {
		fmt.Printf("value = %d , value-addr = %x , slice-addr = %x\n", value, &value, &slice[index])
	}
}

// ➜  demo git:(main) ✗ go run main.go
// value = 10 , value-addr = c0000b2008 , slice-addr = c0000b4000
// value = 20 , value-addr = c0000b2008 , slice-addr = c0000b4008
// value = 30 , value-addr = c0000b2008 , slice-addr = c0000b4010
// value = 40 , value-addr = c0000b2008 , slice-addr = c0000b4018
