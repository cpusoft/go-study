package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(linkChan chan string, wg *sync.WaitGroup) {
	// Decreasing internal counter for wait-group as soon as goroutine finishes
	defer wg.Done()

	for url := range linkChan {
		time.Sleep(1 * time.Second)
		fmt.Printf("Done processing link #%s\n", url)
	}

}
func main() {
	//https://play.golang.org/p/fruJiGBWjn
	yourLinksSlice := make([]string, 50)
	for i := 0; i < 50; i++ {
		yourLinksSlice[i] = fmt.Sprintf("%d", i+1)
	}

	lCh := make(chan string)
	wg := new(sync.WaitGroup)

	// Adding routines to workgroup and running then
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(lCh, wg)
	}

	// Processing all links by spreading them to `free` goroutines
	for _, link := range yourLinksSlice {
		lCh <- link
	}

	// Closing channel (waiting in goroutines won't continue any more)
	close(lCh)

	// Waiting for all goroutines to finish (otherwise they die as main routine dies)
	wg.Wait()
}
