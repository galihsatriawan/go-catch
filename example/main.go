package main

import (
	"fmt"

	"github.com/galihsatriawan/go-catch/v2"
)

func main() {
	var arr []int64
	catch.Catch(func() error {
		var err error = fmt.Errorf("test")
		fmt.Println(arr[0])
		return err
	}, &catch.CatchHandler{
		OnFailure: &catch.OnFailureHandler{
			Dst: &arr,
			Callback: func(err interface{}) (dst interface{}) {
				var a []int64
				a = append(a, 1)
				return &a
			},
		},
	})
	fmt.Println(arr)
}
