package main

import (
	"fmt"
	"sync"
)

func httpGet(url int, response chan string, limiter chan bool, wg *sync.WaitGroup) {
	// 函数执行完毕时 计数器-1
	defer wg.Done()
	// 将拿到的结果, 发送到参数中传递过来的channel中
	response <- fmt.Sprintf("http get: %d", url)
	// 释放一个坑位
	<-limiter
}

// 将所有的返回结果, 以 []string 的形式返回
func collect(urls []int) []string {
	var result []string

	wg := &sync.WaitGroup{}
	// 控制并发数为10
	limiter := make(chan bool, 5)
	defer close(limiter)

	// 函数内的局部变量channel, 专门用来接收函数内所有goroutine的结果
	responseChannel := make(chan string, 20)
	// 为读取结果控制器创建新的WaitGroup, 需要保证控制器内的所有值都已经正确处理完毕, 才能结束
	wgResponse := &sync.WaitGroup{}
	// 启动读取结果的控制器
	go func() {
		// wgResponse计数器+1
		wgResponse.Add(1)
		// 读取结果
		for response := range responseChannel {
			// 处理结果
			result = append(result, response)
		}
		// 当 responseChannel被关闭时且channel中所有的值都已经被处理完毕后, 将执行到这一行
		wgResponse.Done()
	}()

	for _, url := range urls {
		// 计数器+1
		wg.Add(1)
		limiter <- true
		// 这里在启动goroutine时, 将用来收集结果的局部变量channel也传递进去
		go httpGet(url, responseChannel, limiter, wg)
	}

	// 等待所以协程执行完毕
	wg.Wait() // 当计数器为0时, 不再阻塞
	fmt.Println("所有协程已执行完毕")

	// 关闭接收结果channel
	close(responseChannel)

	// 等待wgResponse的计数器归零
	wgResponse.Wait()

	// 返回聚合后结果
	return result
}

func main() {
	urls := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := collect(urls)
	fmt.Println(result)
}
