package main

import (
	"fmt"
	"time"
)

func main() {
	var d time.Duration = 500000 * time.Millisecond
	fmt.Println("d:", fmt.Sprintf("%v", d))
	var a interface{}
	a = d
	fmt.Println("v:", fmt.Sprintf("%v", a))
	if v, p := a.(time.Duration); p {
		fmt.Println("v,p:", fmt.Sprintf("%v", v))
	}
}
