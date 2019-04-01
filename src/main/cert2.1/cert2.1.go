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
type IPAddressRange struct {
	Min asn1.BitString `json:"min"`
	Max asn1.BitString `json:"max"`
}
type IPAddressOrRange struct {
	AddressPrefix asn1.BitString `asn1:"optional,tag:0" json:"addressPrefix"`
	AddressRange  IPAddressRange `asn1:"optional,tag:1" json:"addressRange"`
}

//asn1.NullRawValue
type IPAddressChoice struct {
	//Inherit           []byte             `asn1:"optional,tag:0" json:"inherit"`
	AddressesOrRanges []IPAddressOrRange `asn1:"optional,tag:1" json:"addressesOrRanges"`
}

type IPAddressFamily struct {
	AddressFamily   []byte          `asn1:"tag:0" json:"addressFamily"`
	IPAddressChoice IPAddressChoice `asn1:"tag:1" json:"ipAddressChoice"`
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
type Seqs struct {
	seq Seq `asn1:"tag:16" json:"ipPrefx"`
}

type Seq struct {
	IPPrefx []asn1.BitString `asn1:"tag:3" json:"ipPrefx"`
}
type BitStr struct {
	IPPrefx asn1.BitString `asn1:"tag:3" json:"ipPrefx"`
}
type Tags struct {
	Set struct {
		IPPrefx asn1.BitString `asn1:"tag:3" json:"ipPrefx"`
	} `asn1:"set"`
}

type Sets struct {
	FamilType []byte           `json:"familType"`
	IPPrefx   []asn1.BitString `json:"ipPrefx"`
}

func main() {
	// 从sequence开始

	var ipPrefx asn1.BitString
	ipPrefxByte := []byte{
		0x03, 0x04, 0x02, 0x5B, 0xDB, 0x74}
	asn1.Unmarshal(ipPrefxByte, &ipPrefx)
	fmt.Printf("%+v\r\n", ipPrefx)

	seq := []asn1.BitString{}
	seqByte := []byte{
		0x30, 0x06, 0x03, 0x04, 0x02, 0x5B, 0xDB, 0x74}
	asn1.Unmarshal(seqByte, &seq)
	fmt.Printf("%+v\r\n", seq)

	var family []byte
	familyByte := []byte{0x04, 0x02, 0x00, 0x01}
	asn1.Unmarshal(familyByte, &family)
	fmt.Printf("%+v\r\n", family)

	//0x30, 0x0e,0x30, 0x0C,
	seqs := []Sets{}
	seqByte = []byte{
		0x30, 0x0e, 0x30, 0x0C, 0x04, 0x02, 0x00, 0x01, 0x30, 0x06, 0x03, 0x04, 0x02, 0x5B, 0xDB, 0x74}
	asn1.Unmarshal(seqByte, &seqs)
	fmt.Printf("%+v\r\n", seqs)

	decoded := Tsr{}
	b, _ := ioutil.ReadFile(`E:\Go\go-study\data\reply.tsr`)
	asn1.Unmarshal(b, &decoded)
	fmt.Printf("%+v", decoded)
}
