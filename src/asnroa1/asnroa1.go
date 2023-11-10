package main

import (
	"encoding/asn1"
	"fmt"
	"math/big"
	"time"

	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
)

type CertificateList struct {
	Oid  asn1.ObjectIdentifier
	Seqs []asn1.RawValue `asn1:"optional,explicit,default:0,tag:0""`
}

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
	MaxLength int            `asn1:"optional,default:-1" json:"maxLength"`
}
type Sha256 struct {
	Oid  asn1.ObjectIdentifier
	Null asn1.RawValue
}

func main() {
	var file string
	file = `F:\share\我的坚果云\Go\common\go-study\src\asnroa1\1.roa`
	file = `F:\share\我的坚果云\Go\common\go-study\src\asnroa1\asn0.roa`
	file = `F:\share\我的坚果云\Go\common\go-study\src\asnroa1\ok.roa`
	//file = `F:\share\我的坚果云\Go\common\go-study\src\asnroa1\fail1.roa`
	b, err := fileutil.ReadFileToBytes(file)
	if err != nil {
		fmt.Println(file, err)
		return
	}

	certificate := CertificateList{}
	asn1.Unmarshal(b, &certificate)
	fmt.Println("file:", file, "   certificate:", jsonutil.MarshallJsonIndent(certificate))

	fmt.Println("len(Seqs):", len(certificate.Seqs))

	// int 3 ???
	fmt.Println("seqs[0]", convert.Bytes2Uint64(certificate.Seqs[0].Bytes), "\n")

	// sha256
	sha256 := Sha256{}
	asn1.Unmarshal(certificate.Seqs[1].Bytes, &sha256)
	fmt.Println("seqs[1]:Oids:", jsonutil.MarshallJsonIndent(sha256.Oid))

	// roa
	oid1 := asn1.ObjectIdentifier{} //make([]asn1.RawValue, 0)
	reset, err := asn1.Unmarshal(certificate.Seqs[2].Bytes, &oid1)
	fmt.Println("seqs[2]:oid:", jsonutil.MarshallJsonIndent(oid1), err)

	raw := asn1.RawValue{}
	_, err = asn1.Unmarshal(reset, &raw)
	//fmt.Println("seqs[2]:", convert.PrintBytes(raw.Bytes, 8))

	raw1 := asn1.RawValue{}
	_, err = asn1.Unmarshal(raw.Bytes, &raw1)
	//fmt.Println("seqs[2]:", convert.PrintBytes(raw1.Bytes, 8))

	routeOriginAttestation := RouteOriginAttestation{}
	_, err = asn1.Unmarshal(raw1.Bytes, &routeOriginAttestation)
	fmt.Println("seqs[2]:roa:", jsonutil.MarshallJsonIndent(routeOriginAttestation))

	// ee cer
	cer := Certificate{}
	asn1.Unmarshal(certificate.Seqs[3].Bytes, &cer)
	fmt.Println("certificate:", jsonutil.MarshallJsonIndent(cer))
	fmt.Println(len(cer.TBSCertificate.Extensions))

}

type Certificate struct {
	//Raw                asn1.RawContent
	TBSCertificate     TbsCertificate
	SignatureAlgorithm AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

type TbsCertificate struct {
	Raw                asn1.RawContent
	Version            int `asn1:"optional,explicit,default:0,tag:0"`
	SerialNumber       *big.Int
	SignatureAlgorithm AlgorithmIdentifier
	Issuer             asn1.RawValue
	Validity           Validity
	Subject            asn1.RawValue
	PublicKey          PublicKeyInfo
	Extensions         []Extension `asn1:"optional,explicit,tag:3"`
}
type AlgorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue `asn1:"optional"`
}
type Validity struct {
	NotBefore, NotAfter time.Time
}
type PublicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm AlgorithmIdentifier
	PublicKey asn1.BitString
}

type Extension struct {
	Raw      asn1.RawContent
	Oid      asn1.ObjectIdentifier
	Critical bool `asn1:"optional"`
	Value    []byte
}
