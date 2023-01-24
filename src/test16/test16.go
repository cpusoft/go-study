package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cpusoft/goutil/convert"
)

func main() {
	m := []byte("www1")
	z := []byte{0x00}
	n := []byte("www2")
	wr := bytes.NewBuffer([]byte{})
	binary.Write(wr, binary.BigEndian, m)
	binary.Write(wr, binary.BigEndian, z)
	binary.Write(wr, binary.BigEndian, n)
	b := wr.Bytes()
	fmt.Println(convert.PrintBytesOneLine(b))

	buf := bytes.NewBuffer(b)
	line, err := buf.ReadBytes(0x00)
	fmt.Println(line, err)

	line, err = buf.ReadBytes(0x00)
	fmt.Println(line, err)

	pos := bytes.IndexByte(b, 0x00)
	fmt.Println(pos)
}
