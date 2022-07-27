package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
)

func printStats(mem runtime.MemStats) {
	runtime.ReadMemStats(&mem)
	fmt.Println("mem.Alloc:", mem.Alloc)
	fmt.Println("mem.TotalAlloc:", mem.TotalAlloc)
	fmt.Println("mem.HeapAlloc:", mem.HeapAlloc)
	fmt.Println("mem.NumGC:", mem.NumGC)
	fmt.Println("-----")
}
func main() {
	f, err := os.Create("/tmp/traceFile.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = trace.Start(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer trace.Stop()
	var mem runtime.MemStats
	printStats(mem)
	for i := 0; i < 10; i++ {
		s := make([]byte, 50000000)
		if s == nil {
			fmt.Println("Operation failed!")
		}
	}
	printStats(mem)

}

// go run gColl.go

// mem.Alloc: 84096
// mem.TotalAlloc: 84096
// mem.HeapAlloc: 84096
// mem.NumGC: 0
// -----
// mem.Alloc: 50077000
// mem.TotalAlloc: 500128648
// mem.HeapAlloc: 50077000
// mem.NumGC: 9
// -----



// GODEBUG=gctrace=1 go run gColl.go

// gc 1 @0.026s 0%: 0.024+5.2+0.007 ms clock, 0.096+0.61/0.28/0+0.030 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 2 @0.054s 0%: 0.023+4.1+0.005 ms clock, 0.095+0.051/0.39/0+0.021 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 3 @0.081s 0%: 0.027+3.5+0.005 ms clock, 0.11+0.54/0/0+0.020 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 4 @0.198s 0%: 0.074+0.36+0.002 ms clock, 0.29+0.33/0.30/0.29+0.009 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 5 @0.265s 2%: 6.1+3.3+0.054 ms clock, 24+0/1.2/0.96+0.21 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 6 @0.316s 10%: 28+4.4+0.003 ms clock, 115+0.33/0.71/0+0.014 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 7 @0.353s 10%: 0.22+0.60+0.003 ms clock, 0.90+0.38/0.49/0+0.012 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 8 @0.363s 10%: 0.040+0.94+0.004 ms clock, 0.16+0.78/0/0+0.019 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 9 @0.370s 10%: 0.094+0.95+0.002 ms clock, 0.37+0.77/0/0+0.011 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 10 @0.381s 9%: 0.73+1.5+0.11 ms clock, 2.9+0.36/0.036/0.72+0.47 ms cpu, 4->4->1 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
// # command-line-arguments
// gc 1 @0.005s 9%: 0.015+3.1+0.054 ms clock, 0.062+0.26/2.8/0.097+0.21 ms cpu, 4->4->3 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
// # command-line-arguments
// gc 1 @0.003s 31%: 0.011+6.0+0.006 ms clock, 0.045+6.0/5.8/0.39+0.026 ms cpu, 4->4->3 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 2 @0.018s 16%: 0.009+2.9+0.004 ms clock, 0.038+0/2.8/0.23+0.016 ms cpu, 6->6->6 MB, 7 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 3 @0.048s 8%: 0.024+3.7+0.081 ms clock, 0.099+0.47/2.7/0.70+0.32 ms cpu, 10->10->9 MB, 12 MB goal, 0 MB stacks, 0 MB globals, 4 P
// mem.Alloc: 84352
// mem.TotalAlloc: 84352
// mem.HeapAlloc: 84352
// mem.NumGC: 0
// -----
// gc 1 @0.005s 0%: 0.019+2.8+0.004 ms clock, 0.077+0.19/0/0+0.018 ms cpu, 47->47->0 MB, 47 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 2 @0.059s 1%: 0.86+3.4+0.007 ms clock, 3.4+0.13/0/0+0.031 ms cpu, 47->47->0 MB, 47 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 3 @0.074s 4%: 2.2+3.3+0.006 ms clock, 9.0+0.58/0/0+0.027 ms cpu, 47->47->0 MB, 47 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 4 @0.087s 6%: 2.7+2.8+0.005 ms clock, 11+0.13/0/0+0.020 ms cpu, 47->47->0 MB, 47 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 5 @0.106s 8%: 2.8+1.9+0.004 ms clock, 11+0.12/0/0+0.018 ms cpu, 47->47->0 MB, 47 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 6 @0.118s 7%: 0.049+0.31+0.037 ms clock, 0.19+0.10/0.039/0+0.15 ms cpu, 47->47->0 MB, 47 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 7 @0.125s 7%: 0.15+0.20+0.027 ms clock, 0.60+0.14/0/0+0.11 ms cpu, 47->47->0 MB, 47 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 8 @0.131s 7%: 0.21+0.31+0.020 ms clock, 0.84+0.10/0/0+0.080 ms cpu, 47->47->0 MB, 47 MB goal, 0 MB stacks, 0 MB globals, 4 P
// gc 9 @0.138s 7%: 0.11+0.46+0.050 ms clock, 0.44+0.11/0/0+0.20 ms cpu, 47->47->0 MB, 47 MB goal, 0 MB stacks, 0 MB globals, 4 P
// mem.Alloc: 50077120
// mem.TotalAlloc: 500129016
// mem.HeapAlloc: 50077120
// mem.NumGC: 9
// -----
// gc 10 @0.146s 6%: 0.27+2.1+0.009 ms clock, 1.1+0.11/0.056/0+0.038 ms cpu, 47->47->0 MB, 47 MB goal, 0 MB stacks, 0 MB globals, 4 P