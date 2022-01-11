package main

import (
	"encoding/asn1"
	"encoding/hex"
	"fmt"

	"github.com/cpusoft/goutil/asn1util/asn1node"
	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/jsonutil"
)

func main() {
	func1()
}

func func3() {
	hexStr := `300B0609608648016503040201`
	sigBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	node, err := asn1node.ParseBytes(sigBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(jsonutil.MarshalJson(node))

	oids := make([]asn1.ObjectIdentifier, 0)
	_, err = asn1.Unmarshal(sigBytes, &oids)
	fmt.Println("ParseRoaModelByOpensslResults(): roa:", jsonutil.MarshalJson(oids), err)
}

// asID as in rfc6482
type ROAIpAddressFamilyssss []ROAIpAddressFamilysss
type ROAIpAddressFamilysss []ROAIpAddressFamilyss
type ROAIpAddressFamilyss []ROAIpAddressFamilys
type ROAIpAddressFamilys struct {
	Rs []ROAIpAddressFamily
}

type ROAIpAddressFamily struct {
	AddressFamily []byte `json:"addressFamily"`
	//Addresses     []Address `json:"addresses"`
	IPAddressRange Address `json:"IPAddressRange"`
}
type Address struct {
	Address   asn1.BitString `json:"address"`
	MaxLength int64          `asn1:"optional" json:"maxLength"`
}

func func2() {
	hexStr := `3014A1123010300E04010230090307002001067C208C`
	sigBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	node, err := asn1node.ParseBytes(sigBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(jsonutil.MarshalJson(node))
	/*
		value1 := node.Nodes[0].Nodes[0].Value
		v1, _ := value1.([]byte)
		fmt.Println(convert.PrintBytesOneLine(v1))
		value2 := node.Nodes[0].Nodes[1].Nodes[0].Value
		v2, _ := value2.([]byte)
		fmt.Println(convert.PrintBytesOneLine(v2))
	*/
	rbs := make([][][][]ROAIpAddressFamily, 0)
	_, err = asn1.Unmarshal(sigBytes, &rbs)
	fmt.Println("ParseRoaModelByOpensslResults(): roa:", jsonutil.MarshalJson(rbs), err)

	rb := ROAIpAddressFamilys{}
	_, err = asn1.Unmarshal(sigBytes, &rb)
	fmt.Println("ParseRoaModelByOpensslResults(): roa:", jsonutil.MarshalJson(rb), err)
}

func func1() {
	hexStr := `3010300E04010230090307002001067C208C`
	sigBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	node, err := asn1node.ParseBytes(sigBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(jsonutil.MarshalJson(node))

	value1 := node.Nodes[0].Nodes[0].Value
	v1, _ := value1.([]byte)
	fmt.Println(convert.PrintBytesOneLine(v1))
	value2 := node.Nodes[0].Nodes[1].Nodes[0].Value
	v2, _ := value2.([]byte)
	fmt.Println(convert.PrintBytesOneLine(v2))

	rb := make([]ROAIpAddressFamily, 0)
	_, err = asn1.Unmarshal(sigBytes, &rb)
	fmt.Println("ParseRoaModelByOpensslResults(): roa:", jsonutil.MarshalJson(rb), err)
}
