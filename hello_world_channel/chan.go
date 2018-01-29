package main

import "fmt"

func main() {
	var msg = make(chan string)
	for i := 0; i < 5000; i++ {
		go getHello(i, msg)
	}
	for {
		fmt.Println(<-msg)
	}
}

func getHello(i int, msg chan string) {
	for {
		msg <- fmt.Sprintf("hello world %d\n", i)
	}
}
