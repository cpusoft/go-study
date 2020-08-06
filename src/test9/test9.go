package main

import (
	"fmt"
)

func main() {
	// max
	var a uint8 = 1
	var b byte
	fmt.Println("从1开始，左移并填充1:", a)
	for i := 0; i < 3; i++ {
		a = a | a<<1
		//	fmt.Println(a)
		b = byte(a)
		//	fmt.Printf("%x\n", b)
	}
	fmt.Printf("%x\n\n\n", b)
	var z uint8 = 0x80
	z = z | b
	fmt.Printf("%x,%d,%b\n\n\n", z, z, z)

	// min
	a = 0xff
	fmt.Println("从0xFF开始，左移后填充0:", a)
	for i := 0; i < 10; i++ {
		a = a << 1
		b = byte(a)
	}
	z = 0xA8
	z = z | a
	fmt.Printf("%x\n\n\n", z)
}
