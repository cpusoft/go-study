package main

import (
	"fmt"
	"net/http"
	"time"
)

// result封装错误，和其他需要附加的信息
type Result struct {
	Error    error
	Response *http.Response
}

// 通过一个函数，将go封装起来，
// done 用于终止函数
//  Result为返回值
func checkStatus(done <-chan interface{}, urls ...string) <-chan Result {
	resultCh := make(chan Result)
	go func() {
		defer close(resultCh)
		for _, url := range urls {
			resp, err := http.Get(url)
			result := Result{Error: err, Response: resp}
			fmt.Println("result:", url, result)
			select {
			case <-done:
				fmt.Println("<-done")
				return
			case resultCh <- result:
			}
		}
		fmt.Println("close(resultCh)")
	}()
	return resultCh
}

//如何返回go routine 的返回值
func main() {
	done := make(chan interface{})
	defer close(done)

	urls := []string{"http://badhost/", "http://www.baidu.com/"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)

			continue
		}
		fmt.Printf("Resonse :%v\n", result.Response.Status)
	}
	//close(done)
	fmt.Println("close(done)")
	time.Sleep(time.Duration(2) * time.Second)
	fmt.Println("end")
}
