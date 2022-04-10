package main

import "fmt"

func main() {
	ss := make([]string, 0)
	ss = append(ss, "")
	for _, k := range ss {
		fmt.Println("k", k)
	}
}
