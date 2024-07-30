package main

import "fmt"

type OnlyDb struct {
}

func (c OnlyDb) call() {
	fmt.Println("only db")
	callMySQL1()
}
