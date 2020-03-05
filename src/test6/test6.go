package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	pm := make(map[string]*Person, 10)
	p1 := &Person{"aa", 1}
	p2 := &Person{"ab", 2}
	p3 := &Person{"ac", 3}
	p4 := &Person{"ad", 4}

	pm["aa"] = p1
	pm["ab"] = p2
	pm["ac"] = p3
	pm["ad"] = p4
	fmt.Println(pm)
	fmt.Println(pm["aa"])

	pp1 := pm["aa"]
	pp1.Age = 99

	fmt.Println(pm)
	fmt.Println(pm["aa"])
}
