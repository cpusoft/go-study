package main

import (
	"fmt"
	"strings"
)

func main() {
	domain := "c.bbb.com."
	origin := "bbb.com."
	suffix := strings.HasSuffix(domain, origin)
	fmt.Println(suffix)

	rrName := strings.TrimSuffix(domain, origin)
	fmt.Println(rrName)

	rrName = strings.TrimSuffix(rrName, ".")
	fmt.Println(rrName)

	test := ""
	t := []byte(test)
	var tt []byte
	fmt.Println(t, tt)
}
