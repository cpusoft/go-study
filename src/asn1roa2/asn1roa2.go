package main

import (
	"encoding/asn1"
	"fmt"

	_ "github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
)

type CertificateList struct {
	Oid  asn1.ObjectIdentifier
	Seqs asn1.RawValue `asn1:"optional,explicit,default:0,tag:0""`
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
	file = `G:\Download\cert\asnroa2\1.roa`
	file = `G:\Download\cert\asnroa2\asn0.roa`
	file = `G:\Download\cert\asnroa2\ok.roa`
	file = `G:\Download\cert\asnroa2\fail1.roa`
	b, err := fileutil.ReadFileToBytes(file)
	if err != nil {
		fmt.Println(file, err)
		return
	}
	certificate := CertificateList{}
	rest, err := asn1.Unmarshal(b, &certificate)
	fmt.Println("certificate:", jsonutil.MarshallJsonIndent(certificate), len(rest), err)
	/*
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
	*/
}
