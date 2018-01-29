package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 5000; i++ {
		go printHello(i)
	}
	time.Sleep(time.Microsecond)
}

func printHello(i int) {
	fmt.Printf("hello world %d\n", i)
}
