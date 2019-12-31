package main

import (
	_ "errors"
	"fmt"
	_ "strings"
	"sync"
	"time"
)

func ParseCer(threadindex int, filenamesCh chan string, resultChan chan string, wg *sync.WaitGroup) {
	// Decreasing internal counter for wait-group as soon as goroutine finishes
	defer wg.Done()

	for filenameCh := range filenamesCh {
		time.Sleep(1 * time.Second)
		fmt.Printf("Done processing filenameCh #%s\n", filenameCh)
		resultChan <- filenameCh + "  ok"

	}

}
func main() {
	filenames := make([]string, 8)
	filenames[0] = `E:\Go\parse-cert\data\ROUTER-0000FBF0_new.cer`
	filenames[1] = `E:\Go\parse-cert\data\ROUTER-00010000_new.cer`
	filenames[2] = `E:\Go\parse-cert\data\err1.cer`
	filenames[3] = `E:\Go\parse-cert\data\H.cer`
	filenames[4] = `E:\Go\parse-cert\data\1.cer`
	filenames[5] = `E:\Go\parse-cert\data\range_ipv6.cer`
	filenames[6] = `E:\Go\parse-cert\data\41870XBX5RmmOBSWl-AwgOrYdys_test.cer`
	filenames[7] = `E:\Go\parse-cert\data\test1.crl`

	filenameChan := make(chan string)
	resultChan := make(chan string)

	wg := new(sync.WaitGroup)

	// Adding routines to workgroup and running then
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go ParseCer(i, filenameChan, resultChan, wg)
	}

	// Processing all links by spreading them to `free` goroutines
	for _, filename := range filenames {
		filenameChan <- filename
	}
	wg.Wait()
	close(filenameChan)
	// Waiting for all goroutines to finish (otherwise they die as main routine dies)

	var reslut string
	for _, filename := range filenames {
		reslut = <-resultChan
		fmt.Println(filename + "   " + reslut)
	}

	// Closing channel (waiting in goroutines won't continue any more)

	close(resultChan)

}
