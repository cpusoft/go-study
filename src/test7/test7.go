package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {
	recvByte := []byte{0x00, 0x00, 0x00, 0x14}
	buf := bytes.NewReader(recvByte)
	var length uint32
	err := binary.Read(buf, binary.BigEndian, &length)
	fmt.Println(length, err, buf)
}
