package main

import (
	"fmt"
)

func main() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)
	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case a, b := <-c1:
			c1Count++
			fmt.Println("c1", a, b)
		case a, b := <-c2:
			c2Count++
			fmt.Println("c2", a, b)
		}
	}
	fmt.Println(c1Count, c2Count)
}
