package main

import (
	"fmt"
)

func main() {
	x := []int{
		2, 3, 4,
	}
	y := &x
	fmt.Println(y)
}
