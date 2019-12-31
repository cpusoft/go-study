package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

//http://baixiaoustc.com/2016/06/23/golang-e7-9a-84-e5-bb-ba-e8-ae-aegoroutine-e7-ba-bf-e7-a8-8b-e6-b1-a0/
//https://golangtc.com/t/559e97d6b09ecc22f6000053
//https://godoc.org/golang.org/x/sync/semaphore#example-package--WorkerPool
//http://www.cnblogs.com/luckcs/articles/2588200.html
/////////////////////////////////////////////////////////////////////////////////////////
func usingselect() {
	var urls = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"}
	//最多4个同时运行
	const MAX int = 4
	urlChan := make(chan string)
	lens := make(chan bool, MAX)
	go func() {
		for {
			select {
			case url := <-urlChan:
				go func() {
					fmt.Println(url)
					//模拟下载...
					time.Sleep(time.Second * 2)
					lens <- true
				}()
			}
		}
	}()
	for k, v := range urls {
		if k < MAX {
			lens <- true
		}
		<-lens
		urlChan <- v
	}
}

/////////////////////////////////////////////////////////////////////////////////////////
//还是会生成大量的go进程，阻塞而已
type Pool struct {
	queue chan int
	wg    *sync.WaitGroup
}

// 创建并发控制池, 设置并发数量与总数量
func NewPool(cap, total int) *Pool {
	if cap < 1 {
		cap = 1
	}
	p := &Pool{
		queue: make(chan int, cap),
		wg:    new(sync.WaitGroup),
	}
	p.wg.Add(total)
	return p
}

// 向并发队列中添加一个
func (p *Pool) AddOne() {
	p.queue <- 1
}

// 并发队列中释放一个, 并从总数量中减去一个
func (p *Pool) DelOne() {
	<-p.queue
	p.wg.Done()
}
func Download(s string) error {
	// do download logic
	time.Sleep(time.Second * 2)
	println(s)
	return nil
}
func usingPool() {
	urls := []string{"a", "b", "c", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"}
	pool := NewPool(4, len(urls)) // 初始化一个容量为20的并发控制池
	for _, v := range urls {
		go func(url string) {
			pool.AddOne() // 向并发控制池中添加一个, 一旦池满则此处阻塞
			err := Download(url)
			if nil != err {
				println(err)
			}
			pool.DelOne() // 从并发控制池中释放一个, 之后其他被阻塞的可以进入池中
		}(v)
	}
	pool.wg.Wait() // 等待所有下载全部完成
}

///////////////////////////////////////////////////////////////////
//还是会生成大量的go进程，阻塞而已
func forPool() {
	const (
		GOROUTINE_COUNT = 4
		TASK_COUNT      = 100
	)
	chReq := make(chan string, GOROUTINE_COUNT)
	chRes := make(chan int, GOROUTINE_COUNT)
	for i := 0; i < GOROUTINE_COUNT; i++ {
		go func() {
			for {
				url := <-chReq
				time.Sleep(time.Second * 2)
				fmt.Println(url)
				chRes <- 0
			}
		}()
	}
	go func() {
		urls := make([]string, TASK_COUNT)
		for i := 0; i < TASK_COUNT; i++ {
			urls[i] = fmt.Sprintf("http://www.%d.com", i)
		}
		// got urls
		for i := 0; i < TASK_COUNT; i++ {
			chReq <- urls[i]
		}
	}()
	for i := 0; i < TASK_COUNT; i++ {
		d := <-chRes
		// check error
		_ = d
	}
}

/////////////////////////////////////////////////////////////////////////////////////////
// Example_workerPool demonstrates how to use a semaphore to limit the number of
// goroutines working on parallel tasks.
//
// This use of a semaphore mimics a typical “worker pool” pattern, but without
// the need to explicitly shut down idle workers when the work is done.
func semaphorePool() {
	ctx := context.TODO()

	var (
		maxWorkers = runtime.GOMAXPROCS(0)
		sem        = semaphore.NewWeighted(int64(maxWorkers))
		out        = make([]int, 32)
	)
	fmt.Println(maxWorkers)

	// Compute the output using up to maxWorkers goroutines at a time.
	for i := range out {
		// When maxWorkers goroutines are in flight, Acquire blocks until one of the
		// workers finishes.
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}

		go func(i int) {
			defer sem.Release(1)
			out[i] = collatzSteps(i + 1)
		}(i)
	}

	// Acquire all of the tokens to wait for any remaining workers to finish.
	//
	// If you are already waiting for the workers by some other means (such as an
	// errgroup.Group), you can omit this final Acquire call.
	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
	}

	fmt.Println(out)

}

// collatzSteps computes the number of steps to reach 1 under the Collatz
// conjecture. (See https://en.wikipedia.org/wiki/Collatz_conjecture.)
func collatzSteps(n int) (steps int) {
	if n <= 0 {
		panic("nonpositive input")
	}

	for ; n > 1; steps++ {
		if steps < 0 {
			panic("too many steps")
		}

		if n%2 == 0 {
			n /= 2
			continue
		}

		const maxInt = int(^uint(0) >> 1)
		if n > (maxInt-1)/3 {
			panic("overflow")
		}
		n = 3*n + 1
	}

	return steps
}

/////////////////////////////////////////////////////////////////////////////////////////
func main() {
	//usingselect()
	//usingPool()
	//forPool()
	semaphorePool()
}
