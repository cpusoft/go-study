package main

import (
	"fmt"
	"net/http"
)

// result封装错误，和其他需要附加的信息
type Result struct {
	Error    error
	Response *http.Response
}

func checkStatus(done <-chan interface{}, urls ...string) <-chan Result {
	resultCh := make(chan Result)
	go func() {
		defer close(resultCh)
		for _, url := range urls {
			resp, err := http.Get(url)
			result := Result{Error: err, Response: resp}
			select {
			case <-done:
				return
			case resultCh <- result:
			}
		}
	}()
	return resultCh
}

//如何返回go routine 的返回值
func main() {
	done := make(chan interface{})
	defer close(done)

	urls := []string{"http://www.baidu.com/", "http://badhost/"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			continue
		}
		fmt.Printf("Resonse :%v\n", result.Response.Status)
	}
}
