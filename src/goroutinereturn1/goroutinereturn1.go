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

// 通过一个函数，将go封装起来，
// chan Result为返回值，接收go中的错误，然后在外面的函数中返回
func checkStatus(urls ...string) <-chan Result {
	resultCh := make(chan Result)
	go func() {
		defer close(resultCh)
		for _, url := range urls {
			fmt.Println(url)
			resp, err := http.Get(url)
			result := Result{Error: err, Response: resp}
			fmt.Println("get url:", url, err)
			select {

			case resultCh <- result:
			}
		}
		fmt.Println("finish all urls")
	}()
	return resultCh
}

//如何返回go routine 的返回值
func main() {

	urls := []string{"http://www.baidu.com/", "http://badhost/"}
	for result := range checkStatus(urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			continue
		}
		fmt.Printf("Resonse: %v\n", result.Response.Status)
	}
}

/*
http://www.baidu.com/
get url: http://www.baidu.com/ <nil>
http://badhost/
Resonse: 200 OK
get url: http://badhost/ Get "http://badhost/": dial tcp: lookup badhost: no such host
finish all urls
error: Get "http://badhost/": dial tcp: lookup badhost: no such host
*/
