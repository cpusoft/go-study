package cache

import (
	"fmt"

	db "build3/db"
)

type CacheDb struct {
}

func (c CacheDb) call() {
	fmt.Println("cache+db")
	callCache1()
	db.CallMySQL1()
}
