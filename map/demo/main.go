package main

import "fmt"

func main() {
    m := make(map[string]int)

    fmt.Println(&m["qcrao"])
}

// ➜  map git:(main) ✗ go run main.go 
// # command-line-arguments
// ./main.go:8:18: invalid operation: cannot take address of m["qcrao"] (map index expression of type int)