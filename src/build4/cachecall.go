//go:build cache

package main

import "fmt"

func init() {
	Inf = CacheDb{}
}

type CacheDb struct {
}

func (c CacheDb) call() {
	fmt.Println("cache+db")
	callCache1()
	callMySQL1()
}
