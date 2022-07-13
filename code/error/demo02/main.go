package main

import (
	"errors"
	"fmt"
	"log"
)

func doError() (string, error) {
	return "哇塞计划", errors.New("This is my error")
}

func doNoError() (string, error) {
	return "My response", nil
}

func doFmtError() error {
	errCode := 401
	return fmt.Errorf("This my error code: %d", errCode)
}

func main() {
	resp, err := doError()
	if err != nil {
		log.Printf("There was an error: %v\n", err)
	}
	fmt.Println("My message:", resp)
	resp, err = doNoError()
	if err != nil {
		log.Printf("this should nor print")
	}
	fmt.Println("My response:", resp)
	err = doFmtError()
	if err != nil {
		log.Printf("There was an error: %v\n", err)
	}
}

// 2022/07/12 14:30:50 There was an error: This is my error
// My message: 哇塞计划
// My response: My response
// 2022/07/12 14:30:50 There was an error: This my error code: 401
