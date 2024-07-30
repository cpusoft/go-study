package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type mapSemaphore struct {
	sems map[string]*semaphore.Weighted
	mu   sync.Mutex
	ctx  context.Context
}

func Acquire(key string) error {
	fmt.Println("Acquire will:", key)
	globalMapSemaphore.mu.Lock()
	fmt.Println("Acquire Lock:", key)
	m, ok := globalMapSemaphore.sems[key]
	if !ok {
		sem := semaphore.NewWeighted(1)
		globalMapSemaphore.sems[key] = sem
		fmt.Println("Acquire not found, added, before unlock:", key)
		globalMapSemaphore.mu.Unlock()
		fmt.Println("Acquire add ok:", key)
		return nil
	}
	fmt.Println("Acquire found before unlock:", key)
	globalMapSemaphore.mu.Unlock()
	fmt.Println("Acquire wait:", key)
	err := m.Acquire(globalMapSemaphore.ctx, 1)
	if err != nil {
		fmt.Println("Acquire fail:", key, err)
		return err
	}
	fmt.Println("Acquire ok:", key)
	return nil
}
func Release(key string) {
	fmt.Println("Release will:", key)
	globalMapSemaphore.mu.Lock()

	m, ok := globalMapSemaphore.sems[key]
	if !ok {
		fmt.Println("Release not found before unlock:", key)
		globalMapSemaphore.mu.Unlock()
		fmt.Println("Release not found:", key)
		return
	}
	m.Release(1)
	delete(globalMapSemaphore.sems, key)
	fmt.Println("Release found,deleted, before unlock:", key)
	globalMapSemaphore.mu.Unlock()
	fmt.Println("Release ok:", key)
}

var globalMapSemaphore mapSemaphore

func init() {
	globalMapSemaphore = mapSemaphore{
		sems: make(map[string]*semaphore.Weighted, 10000),
		ctx:  context.Background(),
	}
}

func main() {

	for i := 1; i <= 20; i++ {
		go func(i int) {
			fmt.Println("will ", i)
			if i%2 == 0 {
				fmt.Println("main.Acquire will aaa", i)
				err := Acquire("aaa")
				fmt.Println("main.Acquire aaa", i, err)
				time.Sleep(time.Duration(rand.Intn(2)))
				Release("aaa")
				fmt.Println("main.Acquire have release aaa", i)
			} else {
				fmt.Println("main.Acquire will bbb", i)
				err := Acquire("bbb")
				fmt.Println("main.Acquire bbb", i, err)
				time.Sleep(time.Duration(rand.Intn(2)))
				Release("bbb")
				fmt.Println("main.Acquire have release bbb", i)
			}

		}(i)
	}

	// 等待所有 goroutine 执行完成
	time.Sleep(time.Second * 199)
}
