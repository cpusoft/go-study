package main

import (
	"encoding/asn1"
	"encoding/hex"
	"fmt"

	jsonutil "github.com/cpusoft/goutil/jsonutil"
)

// asID as in rfc6482
type RouteOriginAttestation struct {
	AsID         ASID                 `json:"asID"`
	IpAddrBlocks []ROAIPAddressFamily `json:"ipAddrBlocks"`
}
type ASID int64
type ROAIPAddressFamily struct {
	AddressFamily []byte         `json:"addressFamily"`
	Addresses     []ROAIPAddress `json:"addresses"`
}
type ROAIPAddress struct {
	Address   asn1.BitString `json:"address"`
	MaxLength int64          `asn1:"optional" json:"maxLength"`
}

type IPAddress asn1.BitString

func main() {
	roaHex := `301D0203020FEE3016301404020002300E300C030700240469800043020130`
	roaByte, err := hex.DecodeString(roaHex)
	if err != nil {
		fmt.Println(
			"err:", err)
		return
	}
	roa := RouteOriginAttestation{}
	asn1.Unmarshal(roaByte, &roa)
	fmt.Println("ParseRoaModelByOpensslResults(): roa:", jsonutil.MarshalJson(roa))
}
