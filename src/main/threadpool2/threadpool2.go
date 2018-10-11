package main

import (
	"fmt"
)

func worker(id int, jobs chan int, results chan int) {
	for j := range jobs {
		fmt.Println("worker", id, "procesing job", j)
		results <- j * 2
	}
}

func main() {

	job := make(chan int, 100)
	result := make(chan int, 200)

	for w := 1; w <= 3; w++ {
		go worker(w, job, result)
	}

	for j := 1; j <= 9; j++ {
		job <- j
	}

	for a := 1; a <= 9; a++ {
		fmt.Println(<-result)
	}
}
