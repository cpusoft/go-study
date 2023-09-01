package main

import (
	"encoding/asn1"
	"errors"
	"fmt"
)

func main() {
	fmt.Println("")
	b1 := []byte{0x03, 0x03, 0x04, 0x0a, 0x20}
	as1 := asn1.BitString{}
	_, err := asn1.Unmarshal(b1, &as1)
	fmt.Println("b1:", b1, "as1:", as1, err)

	as2, err := parseBitString(b1[2:])
	fmt.Println("b1:", b1, "as2:", as2, err)
}
func parseBitString(bytes []byte) (ret asn1.BitString, err error) {
	if len(bytes) == 0 {
		err = errors.New("zero length BIT STRING")
		return
	}
	fmt.Println("bytes:", bytes)
	paddingBits := int(bytes[0])
	fmt.Println("paddingBits:", paddingBits)
	if paddingBits > 7 ||
		len(bytes) == 1 && paddingBits > 0 ||
		bytes[len(bytes)-1]&((1<<bytes[0])-1) != 0 {
		err = errors.New("invalid padding bits in BIT STRING")
		return
	}
	ret.BitLength = (len(bytes)-1)*8 - paddingBits
	ret.Bytes = bytes[1:]
	return
}
