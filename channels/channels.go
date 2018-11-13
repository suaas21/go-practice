package main

import (
	"fmt"
)

func pinger(c chan<- string) {
	for i := 0; ; i++ {
		c <- "ping"
		// time.Sleep(time.Second * 2)
	}
}
func ponger(c chan<- string) {
	for i := 0; ; i++ {
		c <- "pong"
		// time.Sleep(time.Second * 2)
	}
}

func printer(c <-chan string) {
	for {
		msg := <-c
		fmt.Println(msg)
		// time.Sleep(time.Second * 1)
	}
}

func main() {
	c := make(chan string)

	go pinger(c)
	go ponger(c)
	go printer(c)

	var input string
	fmt.Scanln(&input)
	//fmt.Scanln(&input)
}
