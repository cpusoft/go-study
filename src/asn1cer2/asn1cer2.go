package main

import (
	"encoding/asn1"
	"fmt"
	"time"

	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
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
	SubjectKeyIdentifier   ObjectIdentifierAndRawValue
	AuthorityKeyIdentifier ObjectIdentifierAndRawValue
	KeyUsage               ObjectIdentifierAndBoolAndRawValue
	BasicConstraints       ObjectIdentifierAndBoolAndRawValue
	CRLDistributionPoints  ObjectIdentifierAndRawValue
	AuthorityInfoAccess    ObjectIdentifierAndRawValue
	CertificatePolicies    ObjectIdentifierAndBoolAndRawValue
	SubjectInfoAccess      ObjectIdentifierAndRawValue
	IpAddrBlocks           ObjectIdentifierAndBoolAndRawValue
}

type ObjectIdentifierAndRawValue struct {
	Type     asn1.ObjectIdentifier
	RawValue asn1.RawValue
}
type ObjectIdentifierAndBoolAndRawValue struct {
	Type     asn1.ObjectIdentifier
	Bool     bool
	RawValue asn1.RawValue
}

type RawValue struct {
	RawValue asn1.RawValue
}
type IPAddrBlocks []IPAddressFamily

type IPAddressFamily struct {
	IPAddressChoices []IPAddressChoices
}

type IPAddressChoices struct {
	AddressFamily     asn1.RawValue
	AddressesOrRanges []IPAddressOrRange
}

type IPAddressOrRange struct {
	AddressPrefix IPAddress     `asn:"optional"`
	IPAddresses   []IPAddresses `asn:"optional"`
}

type IPAddresses struct {
	Min IPAddress
	Max IPAddress
}
type IPAddress []byte

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
	file := `E:\Go\go-study\src\asncer3\3.cer`
	b, err := fileutil.ReadFileToBytes(file)
	if err != nil {
		fmt.Println(file, err)
		return
	}
	certificate := Certificate{}
	asn1.Unmarshal(b, &certificate)
	fmt.Println("certificate:", jsonutil.MarshallJsonIndent(certificate))

	cerParse := CerParse{}
	asn1.Unmarshal(certificate.TBSCertificate.CerRawValue.Bytes, &cerParse)

	subjectKeyIdentifier := make([]byte, 0)
	asn1.Unmarshal(cerParse.SubjectKeyIdentifier.RawValue.Bytes, &subjectKeyIdentifier)
	fmt.Println("subjectKeyIdentifier:", convert.Bytes2String(subjectKeyIdentifier))

	fmt.Println("keyUsage.bool:", cerParse.KeyUsage.Bool)
	keyUsage := asn1.BitString{}
	asn1.Unmarshal(cerParse.KeyUsage.RawValue.Bytes, &keyUsage)
	fmt.Println("keyUsage:", convert.Bytes2Uint64(keyUsage.Bytes))

	fmt.Println("basicConstraints.bool:", cerParse.BasicConstraints.Bool)
	basicConstraints := make([]bool, 0)
	asn1.Unmarshal(cerParse.BasicConstraints.RawValue.Bytes, &basicConstraints)
	fmt.Println("basicConstraints:", (basicConstraints))

	cdp := make([]RawValue, 0)
	asn1.Unmarshal(cerParse.CRLDistributionPoints.RawValue.Bytes, &cdp)
	cdp1 := asn1.RawValue{}
	asn1.Unmarshal(cdp[0].RawValue.Bytes, &cdp1)
	cdp2 := asn1.RawValue{}
	asn1.Unmarshal(cdp1.Bytes, &cdp2)
	fmt.Println("crldp:", string(cdp2.Bytes))

	aias := make([]ObjectIdentifierAndRawValue, 0)
	asn1.Unmarshal(cerParse.AuthorityInfoAccess.RawValue.Bytes, &aias)
	fmt.Println("authorityInfoAccess :", string(aias[0].RawValue.Bytes))

	cps := make([]ObjectIdentifierAndRawValue, 0)
	asn1.Unmarshal(cerParse.CertificatePolicies.RawValue.Bytes, &cps)
	cp := ObjectIdentifierAndRawValue{}
	asn1.Unmarshal(cps[0].RawValue.Bytes, &cp)
	fmt.Println("cp :", string(cp.RawValue.Bytes))

	sias := make([]ObjectIdentifierAndRawValue, 0)
	asn1.Unmarshal(cerParse.SubjectInfoAccess.RawValue.Bytes, &sias)
	for i := range sias {
		fmt.Println("sia:", sias[i].Type, "  url:", string(sias[i].RawValue.Bytes))
	}

	fmt.Println("ipAddrBlocks.RawValue:", convert.Bytes2String(cerParse.IpAddrBlocks.RawValue.Bytes))
	ipAddrBlocksRawValue := make([]asn1.RawValue, 0)
	asn1.Unmarshal(cerParse.IpAddrBlocks.RawValue.Bytes, &ipAddrBlocksRawValue)
	fmt.Println("ipAddrBlocksRawValue:", jsonutil.MarshallJsonIndent(ipAddrBlocksRawValue))

	for i := range ipAddrBlocksRawValue {
		fmt.Println("i:", i, convert.Bytes2String(ipAddrBlocksRawValue[i].Bytes))

		ipAddressFamily := IPAddressChoices{}
		asn1.Unmarshal(ipAddrBlocksRawValue[i].Bytes, &ipAddressFamily)
		fmt.Println("ipAddressFamily:", i, jsonutil.MarshallJsonIndent(ipAddressFamily))
	}
	ipAddrBlocks := IPAddrBlocks{}
	asn1.Unmarshal(cerParse.IpAddrBlocks.RawValue.Bytes, &ipAddrBlocks)
	fmt.Println("ipAddrBlocks:", jsonutil.MarshallJsonIndent(ipAddrBlocks))
}
