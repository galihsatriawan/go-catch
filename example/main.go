package main

import (
	"errors"
	"fmt"

	"github.com/galihsatriawan/go-catch"
)

func main() {
	catch.Catch(&catch.CatchHandler{
		OnError: func(err interface{}) {
			fmt.Println(err)
		},
		OnSuccess: func() {
			fmt.Println("Success")
		},
		Finally: func() {
			fmt.Println("Finally")
		},
	}, errors.New("Test Error"), "Test")
}
