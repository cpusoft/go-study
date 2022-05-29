package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	fmt.Println(len(args))
	fmt.Println(args)
}
