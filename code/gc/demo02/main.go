package main

func allocate() {
	_ = make([]byte, 1<<20)
}
func main() {
	for i := 1; i < 10000; i++ {
		allocate()
	}
}
