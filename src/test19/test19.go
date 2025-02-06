package main

import (
	"fmt"
	_ "go-study/test19/t1"
	_ "go-study/test19/t1/t2"
	_ "go-study/test19/t1/t2/t3"
	"strings"
)

func main() {
	fmt.Println("main")

	s := `|     root:Rpstir-123 - Valid credentials`
	f := strings.Contains(s, "Valid credentials")
	fmt.Println(f)
}
