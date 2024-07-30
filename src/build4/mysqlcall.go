//go:build !cache

package main

import "fmt"

func init() {
	Inf = OnlyDb{}
}

type OnlyDb struct {
}

func (c OnlyDb) call() {
	fmt.Println("only db")
	callMySQL1()
}
