package main

import (
	asn1 "encoding/asn1"
	"fmt"
)

/*

   id-pe-ipAddrBlocks      OBJECT IDENTIFIER ::= { id-pe 7 }

   IPAddrBlocks        ::= SEQUENCE OF IPAddressFamily

   IPAddressFamily     ::= SEQUENCE {    -- AFI & optional SAFI --
      addressFamily        OCTET STRING (SIZE (2..3)),
      ipAddressChoice      IPAddressChoice }

   IPAddressChoice     ::= CHOICE {
      inherit              NULL, -- inherit from issuer --
      addressesOrRanges    SEQUENCE OF IPAddressOrRange }

   IPAddressOrRange    ::= CHOICE {
      addressPrefix        IPAddress,
      addressRange         IPAddressRange }

   IPAddressRange      ::= SEQUENCE {
      min                  IPAddress,
      max                  IPAddress }

   IPAddress           ::= BIT STRING

*/

type IPAddress asn1.RawValue
type IPAddressRange struct {
	min IPAddress
	max IPAddress
}
type IPAddressOrRange struct {
	addressPrefix IPAddress      `asn1:"choice:prefix,optional"`
	addressRange  IPAddressRange `asn1:"choice:range,optional"`
}
type IPAddressChoice struct {
	//inherit           nil
	addressesOrRanges []IPAddressOrRange
}

func TestParseBigInt() {
	mdata, err := asn1.Marshal(13)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, md := range mdata {
		fmt.Print(fmt.Sprintf("0x%02x ", md))
	}
	var n int
	_, err1 := asn1.Unmarshal(mdata, &n)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	fmt.Println("After marshal/unmarshal: ", n)
}
func main() {
	TestParseBigInt()
	var ia IPAddressChoice
	var ar = []byte{0x30, 0x0c, 0x04, 0x02, 0x00, 0x01, 0x30, 0x06, 0x03, 0x04, 0x02, 0xb9, 0xa6, 0xfc, 0x30, 0x0c, 0x04, 0x02, 0x00, 0x01, 0x30, 0x06, 0x03, 0x04, 0x02, 0xb9, 0xa6, 0xfc}
	fmt.Println(fmt.Sprintf("%v", ar))
	addrs, err1 := asn1.Unmarshal(ar, &ia)
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println(fmt.Sprintf("%v", ia), fmt.Sprint("%v"), addrs)

}
