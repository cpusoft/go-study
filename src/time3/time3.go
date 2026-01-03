package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	time.Sleep(10 * time.Second)
	end := time.Now()
	diff := time.Duration(int64(end.Sub(start)/time.Second)) * time.Second
	fmt.Println("diff:", diff)
}
