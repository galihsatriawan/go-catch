package main

import (
	"errors"
	"fmt"

	"github.com/galihsatriawan/go-catch"
)

func main() {
	catch.Catch(&catch.CatchHandler{
		ErrorHandling: func(err interface{}) {
			fmt.Println(err)
		},
		SuccessHandling: func() {
			fmt.Println("Success")
		},
		FinallyHandling: func() {
			fmt.Println("Finally")
		},
	}, errors.New("Test Error"), "Test")
}
