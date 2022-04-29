package main

func main() {
	str := "hello"
	println([]byte(str))
}

➜  GOOS=linux GOARCH=amd64 go tool compile -S main.go
go.string."hello" SRODATA dupok size=5
        0x0000 68 65 6c 6c 6f                                   hello