package db

import "fmt"

type OnlyDb struct {
}

func (c OnlyDb) call() {
	fmt.Println("only db")
	CallMySQL1()
}
