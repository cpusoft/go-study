package main

import (
	"fmt"
	"time"

	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"

	"asn1"
)

type Certificate struct {
	TBSCertificate TBSCertificate

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

	CerRawValue asn1.RawValue
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
	NotBefore time.Time `asn1:"generalized"`
	NotAfter  time.Time `asn1:"generalized"`
}

type PublicKeyInfo struct {
	Algorithm AlgorithmIdentifier
	PublicKey asn1.BitString
}

type CerParse struct {
	SubjectKeyIdentifier AttributeTypeAndBitString
	/*
		AuthorityKeyIdentifier RDNSequence
		KeyUsage               RDNSequence
		BasicConstraints       RDNSequence
		CRLDistributionPoints  RDNSequence
		CertificatePolicies    RDNSequence
		SubjectInfoAccess      RDNSequence
		IpAddrBlocks           RDNSequence
	*/
}

type AttributeTypeAndBitString struct {
	Type asn1.ObjectIdentifier
	V    asn1.RawValue
}

type OctectString struct {
	V asn1.RawValue
}

/*
Go type                | ASN.1 universal tag
-----------------------|--------------------
bool                   | BOOLEAN
All int and uint types | INTEGER
*big.Int               | INTEGER
string                 | OCTET STRING
[]byte                 | OCTET STRING
asn1.Oid               | OBJECT INDETIFIER
asn1.Null              | NULL
Any array or slice     | SEQUENCE OF
Any struct             | SEQUENCE
*/
func main() {
	file := `E:\Go\go-study\src\asncer1\0.cer`
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
	fmt.Println("certificate:", jsonutil.MarshallJsonIndent(certificate))

	b = certificate.TBSCertificate.CerRawValue.Bytes
	s = convert.PrintBytes(b, 8)
	fmt.Println(s)
	cerParse := CerParse{}
	asn1.Unmarshal(b, &cerParse)
	fmt.Println("cerParse:", jsonutil.MarshallJsonIndent(cerParse))

	b = cerParse.SubjectKeyIdentifier.V.Bytes
	s = convert.PrintBytes(b, 8)
	fmt.Println(s)
	octectString := make([]byte, 0)
	asn1.Unmarshal(b, &octectString)
	fmt.Println("octectString:", jsonutil.MarshallJsonIndent(octectString))
	s = convert.PrintBytes(octectString, 8)
	fmt.Println(s)
	/*
		b = octectString.Bytes
		s = convert.PrintBytes(b, 8)
		fmt.Println(s)
	*/
}
