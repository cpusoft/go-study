package main

import (
	"fmt"
	"math/big"
)

func main() {
	s := `07b2Z`
	n, ok := new(big.Int).SetString(s, 16)
	fmt.Println(n, ok)
}
