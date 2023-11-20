package main

import (
	asn1 "encoding/asn1"
	"fmt"
	"time"

	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
)

type CrlParse struct {
	// 1.2.840.113549.1.7.2 signedData

	Version            uint64
	CrlOidParse        CrlOidParse
	CrlCommonNameParse CrlCommonNameParse
	ThisUpdateTime     time.Time
	NextUpdateTime     time.Time
	CrlAkiNumberParse  CrlAkiNumberParse
	CrlSha256Parse     CrlOidParse
}
type CrlOidParse struct {
	Oid asn1.ObjectIdentifier
	Nul asn1.RawValue
}

type CrlCommonNameParse struct {
	Oid        asn1.ObjectIdentifier
	CommonName string
}
type CrlAkiNumberParse struct {
	CrlAkiParse    CrlAkiParse
	CrlNumberParse CrlNumberParse
}
type CrlAkiParse struct {
	Oid asn1.ObjectIdentifier
	Aki []byte
}
type CrlNumberParse struct {
	Oid       asn1.ObjectIdentifier
	CrlNumber uint64
}

type Certificate struct {
	TBSCertificate     TBSCertificate
	SignatureAlgorithm AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

type TBSCertificate struct {
	Version            int `asn1:"optional,explicit,default:0,tag:0"`
	SerialNumber       asn1.RawValue
	SignatureAlgorithm AlgorithmIdentifier
	Issuer             RDNSequence
	Validity           Validity
	Subject            RDNSequence
	PublicKey          PublicKeyInfo
}

type AlgorithmIdentifier struct {
	Algorithm asn1.ObjectIdentifier
}

type RDNSequence []RelativeDistinguishedNameSET

type RelativeDistinguishedNameSET []AttributeTypeAndValue

type AttributeTypeAndValue struct {
	Type  asn1.ObjectIdentifier
	Value interface{}
}

type Validity struct {
	NotBefore, NotAfter time.Time
}

type PublicKeyInfo struct {
	Algorithm AlgorithmIdentifier
	PublicKey asn1.BitString
}

func main() {
	file := `E:\Go\go-study\src\file2\1.cer`
	b, err := fileutil.ReadFileToBytes(file)
	if err != nil {
		fmt.Println(file, err)
		return
	}
	fmt.Println(len(b))
	s := convert.PrintBytes(b, 8)
	fmt.Println(s)

	//var roaAllParse asn1.RawValue
	certificate := Certificate{}

	//roaAllParse := make([]RoaAllParse, 0)
	//roaAllParse := RoaAllParse{}
	asn1.Unmarshal(b, &certificate)
	fmt.Println("certificate:", jsonutil.MarshalJson(certificate))

}
