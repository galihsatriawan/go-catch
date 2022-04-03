package main

import (
	"fmt"
	"time"

	"github.com/galihsatriawan/go-catch/v2"
)

func main() {
	var arr []int
	err := catch.Catch(func() error {
		var err error = fmt.Errorf("test error")
		fmt.Println("test panic 1")
		fmt.Println(arr[0])
		return err
	},
		catch.OnPanic(func(err interface{}) {
			return
		}),
		catch.OnSuccess(nil, func() interface{} {
			catch.Catch(func() error {
				fmt.Println(arr[0])
				return nil
			})
			fmt.Println("[success] test panic 2")
			return nil
		}),
		catch.OnFailure(nil, func(err interface{}) interface{} {
			fmt.Println("[failure] test panic 2")
			// fmt.Println(arr[0])
			return nil
		}),
		catch.Finally(&arr, func() interface{} {
			var a []int
			a = append(a, 1)
			fmt.Println("test panic 3")
			// fmt.Println(arr[0])
			time.Sleep(3 * time.Second)
			return a
		}),
	)

	fmt.Println(err, arr)
}
