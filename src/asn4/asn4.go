package main

import (
	"encoding/asn1"
	"encoding/hex"
	"fmt"

	jsonutil "github.com/cpusoft/goutil/jsonutil"
)

type ResourceBlock struct {
	//AsIds        []AsId               `asn1:"optional" json:"asIds"`
	IpAddrBlocks []ROAIpAddressFamily `asn1:"optional" json:"ipAddrBlocks"`
}

type ResourceBlockAsId struct {
	AsIds []AsId `json:"asIds"`
}
type AsId int64

// asID as in rfc6482
type ResourceBlockIpAddress struct {
	IpAddrBlocks []ROAIpAddressFamily `json:"ipAddrBlocks"`
}

type ROAIpAddressFamily struct {
	AddressFamily []byte    `json:"addressFamily"`
	Addresses     []Address `json:"addresses"`
}
type Address struct {
	Address   asn1.BitString `json:"address"`
	MaxLength int64          `asn1:"optional" json:"maxLength"`
}

func main() {
	hexStr := `300E04010230090307002001067C208C`
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	rb := ResourceBlock{}
	_, err = asn1.Unmarshal(b, &rb)
	fmt.Println("ParseRoaModelByOpensslResults(): roa:", jsonutil.MarshalJson(rb), err)
}
