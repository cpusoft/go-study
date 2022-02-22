package main

import (
	"encoding/asn1"
	"encoding/hex"
	"fmt"

	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/jsonutil"
)

type ASID int64
type ROAIPAddress struct {
	Address   asn1.BitString `json:"address"`
	MaxLength int            `asn1:"optional,default:-1" json:"maxLength"`
}

type IPAddress asn1.BitString

// asID as in rfc6482
type RouteOriginAttestation struct {
	//	AsID         ASID                 `json:"asID"`
	IpAddrBlock1 []IpAddrBlock1        `json:"ipAddrBlockss"`
	FileHashAlg  asn1.ObjectIdentifier `json:"fileHashAlg"`
	FileList     []FileAndHashParse    `json:"fileList"`
}
type IpAddrBlock1 struct {
	IpAddrBlock2 []IpAddrBlock2 `json:"ipAddrBlocks"`
}
type IpAddrBlock2 struct {
	IpAddrBlock3 []IpAddrBlock3 `json:"ipAddrBlocks"`
}
type IpAddrBlock3 struct {
	Tag1 asn1.BitString `json:"Tag1"`
	Tag2 asn1.BitString `json:"Tag2"`
}
type FileAndHashParse struct {
	File string         `asn1:"ia5" json:"file"`
	Hash asn1.BitString `json:"hash"`
}
type ROAIPAddressFamily struct {
	Version byte          `json:"version"`
	Content asn1.RawValue `json:"content"`
	//	AddressFamily []asn1.RawContent `json:"addressFamily"`
	//	Sha256        []asn1.RawContent `json:"sha256"`
	//	Sig           []asn1.RawContent `json:"sig"`
}

type Sig struct {
	//	AsID         ASID                 `json:"asID"`
	IpAddrBlock1 []asn1.RawValue `json:"ipAddrBlockss"`
}

func main() {
	sigStr := `30819C3014A1123010300E04010230090307002001067C208C300B06096086480165030402013077303416106234325F697076365F6C6F612E706E6704209516DD64BE7C1725B9FCA117120E58E8D842A5206873399B3DDFFC91C4B6ACF0303F161B6234325F736572766963655F646566696E6974696F6E2E6A736F6E04200AE1394722005CD92F4C6AA024D5D6B3E2E67D629F11720D9478A633A117A1C7`
	sigBytes, err := hex.DecodeString(sigStr)
	fmt.Println(convert.PrintBytesOneLine(sigBytes), err)

	sig := make([]asn1.RawValue, 0)
	_, err = Unmarshal(sigBytes, &sig)
	fmt.Println("ParseRoaModelByOpensslResults(): sig:", jsonutil.MarshalJson(sig), err)

}
