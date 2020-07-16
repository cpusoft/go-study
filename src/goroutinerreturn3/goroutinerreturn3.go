package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type FuncResponse struct {
	Error    error
	Response interface{}
}
type FuncProcess func(funcParam interface{}) FuncResponse

func GoFuncProcess(funcProcess FuncProcess, funcParam interface{}, funcResponseCh chan FuncResponse, limit chan int, wg *sync.WaitGroup) {
	// 函数执行完毕时 计数器-1
	defer wg.Done()
	// 将拿到的结果, 发送到参数中传递过来的channel中

	funcResponse := funcProcess(funcParam)
	fmt.Println(funcParam, funcResponse)
	funcResponseCh <- funcResponse
	// 释放一个坑位
	<-limit
}

func GoFuncProcesses(funcProcess FuncProcess, funcParams []interface{}, limitCount int) (funcResponses []FuncResponse) {
	funcResponses = make([]FuncResponse, 0)

	// 设置返回值channel，和控制的wg
	funcResponseCh := make(chan FuncResponse, len(funcParams))
	wgResponse := &sync.WaitGroup{}
	go func() {
		wgResponse.Add(1)
		for r := range funcResponseCh {
			funcResponses = append(funcResponses, r)
		}
		wgResponse.Done()
	}()

	// 设置运行协程的个数、控制的wg， 结果放到返回值channel中
	limit := make(chan int, limitCount)
	defer close(limit)
	wg := &sync.WaitGroup{}
	for i, funcParam := range funcParams {
		wg.Add(1)
		limit <- 1
		fmt.Println("funcParam", i, funcParam)
		go GoFuncProcess(funcProcess, funcParam, funcResponseCh, limit, wg)
	}
	wg.Wait()
	fmt.Println("all funcParams ", len(funcParams), " end")

	close(funcResponseCh)
	wgResponse.Wait()

	return funcResponses

}

func main() {
	urls := []interface{}{"http://www.baidu.com", "http://www.google.com", "http://www.sina.com.cn", "http://localtest/", "http://www.gov.cn",
		"http://www.163.com", "http://www.126.com", "http://www.douyu.com"}

	funcResponses := GoFuncProcesses(HttpGet, urls, 2)
	fmt.Println(funcResponses)
}

func HttpGet(funcParam interface{}) (funcResponse FuncResponse) {
	start := time.Now()
	url := funcParam.(string)
	resp, err := http.Get(url)
	response := "http get " + url
	if err != nil {
		response += " fail. "
		if resp != nil {
			response += (" status is " + resp.Status)
		}
		funcResponse = FuncResponse{Error: err, Response: response}
	} else {
		response += " ok. "
		funcResponse = FuncResponse{Error: nil, Response: response}
	}
	fmt.Println("url:", url, "  funcResponse:", funcResponse, "  time(s):", time.Now().Sub(start).Seconds())
	return funcResponse
}
