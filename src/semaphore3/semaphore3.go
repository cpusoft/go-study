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
	sems sync.Map
	ctx  context.Context
}

func Acquire(key string) error {
	fmt.Println("Acquire will:", key)

	v, ok := globalMapSemaphore.sems.Load(key)
	if !ok {
		fmt.Println("Acquire not found key:", key)
		sem := semaphore.NewWeighted(1)
		globalMapSemaphore.sems.Store(key, sem)
		fmt.Println("Acquire store new sem ok, key:", key, sem)
		return nil
	}
	fmt.Println("Acquire found key, will sem.Acquire:", key)
	m := v.(*semaphore.Weighted)
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

	v, ok := globalMapSemaphore.sems.Load(key)
	if !ok {
		fmt.Println("Release not found:", key)
		return
	}
	m := v.(*semaphore.Weighted)
	fmt.Println("Release found, will sem.Release:", key)
	m.Release(1)
	globalMapSemaphore.sems.Delete(key)
	fmt.Println("Release ok:", key)
}

var globalMapSemaphore mapSemaphore

func init() {
	globalMapSemaphore = mapSemaphore{

		ctx: context.Background(),
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
