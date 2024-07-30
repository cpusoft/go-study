package main

import "fmt"

type CacheDb struct {
}

func (c CacheDb) call() {
	fmt.Println("cache+db")
	callCache1()
	callMySQL1()
}
