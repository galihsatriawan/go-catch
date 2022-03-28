package main

import (
	"fmt"

	"github.com/galihsatriawan/go-catch/v2"
)

func main() {
	var arr []int
	err := catch.Catch(func() error {
		var err error = fmt.Errorf("test")
		fmt.Println(arr[0])
		return err
	},
		catch.OnError(func(err interface{}) {
			fmt.Println(err)
		}),
		catch.OnSuccess(nil, func() interface{} {
			return nil
		}),
		catch.OnFailure(&arr, func(err interface{}) interface{} {
			var a []int
			fmt.Println("hello")
			return &a
		}),
		catch.Finally(&arr, func() interface{} {
			var a []int
			a = append(a, 1)
			return a
		}),
	)

	fmt.Println(err, arr)
}
