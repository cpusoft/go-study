package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
)

func main() {
	b := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}
	buf := bytes.NewReader(b)
	var arCount uint16
	err := binary.Read(buf, binary.BigEndian, &arCount)
	fmt.Println("arcount:", jsonutil.MarshalJson(*buf), arCount, err)

	//test1(buf)
	tmp := make([]byte, 4)
	err = binary.Read(buf, binary.BigEndian, &tmp)
	fmt.Println("tmp", jsonutil.MarshalJson(*buf), tmp, err)

	err = binary.Read(buf, binary.BigEndian, &arCount)
	fmt.Println("arcount2:", jsonutil.MarshalJson(*buf), arCount, err)
}
func test1(buf *bytes.Reader) {
	tmp := make([]byte, 4)
	err := binary.Read(buf, binary.BigEndian, &tmp)
	fmt.Println("tmp", tmp, err)
}
