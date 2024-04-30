package main

import (
	"fmt"

	"github.com/cpusoft/goutil/convert"
	"github.com/guregu/null"
)

// main
func main() {
	en := null.IntFrom(2147352576)
	encodedSubTree := [4]byte{0x00}
	b, _ := convert.IntToBytes(en.ValueOrZero())
	copy(encodedSubTree[:], b[4:])
	fmt.Println(convert.PrintBytesOneLine(b), encodedSubTree)
}
