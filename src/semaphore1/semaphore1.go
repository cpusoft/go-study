package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/semaphore"
)

func main() {
	// 创建一个初始值为 2 的信号量
	sem := semaphore.NewWeighted(2)
	ctx := context.Background()
	// 开启 5 个 goroutine，但只允许 2 个同时执行
	for i := 1; i <= 5; i++ {
		go func(i int) {
			// 请求信号量
			if err := sem.Acquire(ctx, 1); err != nil {
				fmt.Printf("goroutine %d acquire semaphore failed: %v\n", i, err)
				return
			}

			// 执行任务
			fmt.Printf("goroutine %d start running\n", i)
			time.Sleep(time.Second)
			fmt.Printf("goroutine %d stop running\n", i)

			// 释放信号量
			sem.Release(1)
		}(i)
	}

	// 等待所有 goroutine 执行完成
	time.Sleep(time.Second * 6)
}
