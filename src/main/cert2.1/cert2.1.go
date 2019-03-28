package main

import (
	"encoding/asn1"
	"fmt"
	"io/ioutil"
)

type Tsr struct {
	Result struct {
		Code   int
		Detail struct{ Message string }
	}
	SignedData struct {
		OID  asn1.ObjectIdentifier
		Zero struct {
			Seq struct {
				Int int
				Set struct {
					InnerSeq struct {
						OID asn1.ObjectIdentifier
					}
				} `asn1:"set"`
			}
		} `asn1:"tag:0"`
	}
}

type IPAddrBlocks struct {
	Seq struct {
		IPAddressFamily struct {
			addressFamily   []byte
			IPAddressChoice struct {
				Seq struct {
					AddressesOrRanges struct {
						Zero struct {
							IPAddress asn1.BitString
						} `asn1:"tag:0"`
						Seq struct {
							IPAddressRange struct {
								Min asn1.BitString
								Max asn1.BitString
							} `asn1:"set"`
						} `asn1:"tag:1"`
					}
				} `asn1:"set"`
			}
		}
	} `asn1:"sequence"`
}

func main() {
	// 从sequence开始
	iab := IPAddrBlocks{}
	var bb = []byte{
		0x30, 0x0e, 0x30, 0x0c, 0x04, 0x02, 0x00, 0x01, 0x30, 0x06, 0x03, 0x04, 0x02, 0xb9, 0xa6, 0xfc}
	fmt.Println(bb)
	asn1.Unmarshal(bb, &iab)
	fmt.Printf("%+v\r\n", iab)

	decoded := Tsr{}
	b, _ := ioutil.ReadFile(`E:\Go\go-study\data\reply.tsr`)
	asn1.Unmarshal(b, &decoded)
	fmt.Printf("%+v", decoded)
}
