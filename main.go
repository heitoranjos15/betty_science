package main

import "fmt"

func printCount(n int, c chan<- string) {
	c <- fmt.Sprintf("print counting to %d", n)
}

func main() {
	c := make(chan string)
	d := make(chan string)
	go printCount(5, c)
	go printCount(10, d)

	msg, msg1 := <-c, <-d
	fmt.Println(msg)
	fmt.Println(msg1)
}
