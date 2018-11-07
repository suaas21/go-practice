package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	file, err := os.Create("test.txt")
	if err != nil {
		// handle the error here
		return
	}
	defer file.Close()

	file.WriteString("test")

	bs, err := ioutil.ReadFile("test.txt")
	if err != nil {
		return
	}
	fmt.Println(bs)
	str := string(bs)
	fmt.Println(str)

}
