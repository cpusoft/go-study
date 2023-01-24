package main

import "fmt"

func main() {
	m := make(map[string]int)
	m["1"] = 1

	if _, ok := m["2"]; !ok {
		m["2"] = 0
	}
	m["2"] += 1

	fmt.Println(m)
}
