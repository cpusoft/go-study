package main

import (
	"fmt"
	"strings"
)

func main() {
	ff := ""
	split := strings.Split(ff, ",")
	fmt.Println(split, len(split))

	zz := make([]string, 0)
	fmt.Println(zz, len(zz))
	zz = append(zz, split...)
	fmt.Println(zz, len(zz))
}
