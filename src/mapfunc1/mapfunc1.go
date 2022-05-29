package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/cpusoft/goutil/jsonutil"
)

//定义函数类型

type Msg func(name string) string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := &sync.WaitGroup{}
	c := make(chan os.Signal, 1)
	handleMap := make(map[int]Msg)
	handleMap[1] = handle1
	handleMap[2] = handle2
	handleMap[3] = handle3
	fmt.Println(handleMap)
	fmt.Println(len(handleMap))
	fmt.Println(jsonutil.MarshalJson(handleMap))
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-c
		_ = sig
		s := handleMap[3]
		s("gggggg")
		wg.Done()
	}()
	wg.Add(1)
	fmt.Println("run.......")
	wg.Wait()
	fmt.Printf("end")
}

func handle1(name string) string {
	fmt.Println("handle1:", name)
	return "handle1"

}
func handle2(name string) string {
	fmt.Println("handle2:", name)
	return "handle2"

}
func handle3(name string) string {
	fmt.Println("handle3:", name)
	return "handle3"

}
