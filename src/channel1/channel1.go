package main

import (
	"fmt"
	"strconv"
)

func main() {
	// parse and validate each certs in syncLogFileModes
	parseConcurrentCh := make(chan int, 5)
	parseFailFileCh := make(chan string, 300)
	parseFailFiles := make([]string, 0, 20)

	for i := 1; i < 30; i++ {
		fmt.Println("for i :", i)
		parseConcurrentCh <- 1
		go add(i, parseConcurrentCh, parseFailFileCh)
	}

	for i := 1; i < 30; i++ {
		parseFailFile := <-parseFailFileCh
		parseFailFiles = append(parseFailFiles, parseFailFile)
	}
	close(parseConcurrentCh)
	close(parseFailFileCh)
	fmt.Println("parseFailFiles:", parseFailFiles)
}
func add(i int,
	parseConcurrentCh chan int, parseFailFileCh chan string) {
	defer func() {

		<-parseConcurrentCh
	}()

	if i%2 == 0 {
		fmt.Println("add(), i :", i)
		parseFailFileCh <- strconv.Itoa(i)

	} else {
		parseFailFileCh <- ""
	}

}
