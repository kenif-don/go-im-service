package main

import "fmt"

func main() {
	i := test()
	fmt.Println(i)
}
func test() int {
	defer func() {
		fmt.Println(2)
	}()
	return 1
}
