package main

import (
	"fmt"

	"github.com/cpusoft/goutil/convert"
)

func main() {
	var a uint16
	var b uint16

	a = uint16(1)
	b = uint16(2)
	aa, _ := convert.IntToBytes(a)
	bb, _ := convert.IntToBytes(b)
	fmt.Println(aa, len(aa))
	fmt.Println(bb, len(bb))

}
