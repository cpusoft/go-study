package main

import (
	"fmt"
	jsonutil "github.com/cpusoft/goutil/jsonutil"
)

type A struct {
	A1 string
	A2 string
}
type B struct {
	A
	B1 string
}

func main() {
	b := B{B1: "b1"}
	b.A1 = "a1"
	b.A2 = "a2"
	fmt.Println(b)
	js := jsonutil.MarshalJson(b)
	fmt.Println(js)
	f(b.A)

	var a1 A
	jsonutil.UnmarshalJson(js, &a1)
	fmt.Print(a1)

	var b1 B
	jsonutil.UnmarshalJson(js, &b1)
	fmt.Print(b1)
}

func f(a A) {
	fmt.Println(a)
}
