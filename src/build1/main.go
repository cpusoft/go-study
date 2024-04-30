package main

import "fmt"

var configArr []string

func main() {
	for _, conf := range configArr {
		fmt.Println(conf)
	}
}
