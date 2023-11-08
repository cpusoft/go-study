package main

import (
	_ "crypto/x509"
	"fmt"
	"math/big"
	"time"

	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
)

type Certificate struct {
	//Raw                asn1.RawContent
	TBSCertificate     TbsCertificate
	SignatureAlgorithm AlgorithmIdentifier
	SignatureValue     BitString
}

type TbsCertificate struct {
	Raw                RawContent
	Version            int `asn1:"optional,explicit,default:0,tag:0"`
	SerialNumber       *big.Int
	SignatureAlgorithm AlgorithmIdentifier
	Issuer             RawValue
	Validity           Validity
	Subject            RawValue
	PublicKey          PublicKeyInfo
	Extensions         []Extension `asn1:"optional,explicit,tag:3"`
}
type AlgorithmIdentifier struct {
	Algorithm  ObjectIdentifier
	Parameters RawValue `asn1:"optional"`
}
type Validity struct {
	NotBefore, NotAfter time.Time
}
type PublicKeyInfo struct {
	Raw       RawContent
	Algorithm AlgorithmIdentifier
	PublicKey BitString
}

type SeqExtension struct {
	Raw   RawContent
	Oid   ObjectIdentifier
	Value []byte `asn1:"implicit,tag:6"`
	//Value string `asn1:"implicit,tag:6"`
}
type Extension struct {
	Raw      RawContent
	Oid      ObjectIdentifier
	Critical bool `asn1:"optional"`
	Value    []byte
}

type AsnPoint struct {
	AsnPointName AnsPointName `asn1:"optional,tag:0"`
}

type AnsPointName struct {
	AsnNames []RawValue //`asn1:"optional,tag:0"`
}

func GetAsns(value []byte) (AsnPoint, error) {
	fmt.Println("GetAsns(): value:", convert.PrintBytesOneLine(value))

	var asnPoint AsnPoint
	_, err := Unmarshal(value, &asnPoint)
	if err != nil {
		fmt.Println("GetAsns(): Unmarshal fail:", err)
		return asnPoint, err
	}

	fmt.Println("GetAsns(): asnPoint:", jsonutil.MarshalJson(asnPoint))
	for _, asnName := range asnPoint.AsnPointName.AsnNames {
		b := asnName.Bytes
		asn := big.NewInt(0).SetBytes(b)
		fmt.Println(asn)
	}
	return asnPoint, nil
}

type distributionPoint struct {
	DistributionPoint distributionPointName `asn1:"optional,tag:0"`
	Reason            BitString             `asn1:"optional,tag:1"`
	CRLIssuer         RawValue              `asn1:"optional,tag:2"`
}

type distributionPointName struct {
	FullName     []RawValue `asn1:"optional,tag:0"`
	RelativeName RawValue   `asn1:"optional,tag:1"`
}

func GetCrldp(value []byte) ([]string, error) {
	var cdp []distributionPoint
	_, err := Unmarshal(value, &cdp)
	if err != nil {
		return nil, err
	}

	cls := make([]string, 0)
	for _, dp := range cdp {
		// Per RFC 5280, 4.2.1.13, one of distributionPoint or cRLIssuer may be empty.
		if len(dp.DistributionPoint.FullName) == 0 {
			fmt.Println("GetCrldp(): fullName is empty")
			continue
		}
		fmt.Println("GetCrldp(): len( dp.DistributionPoint.FullName):", len(dp.DistributionPoint.FullName))
		for _, fullName := range dp.DistributionPoint.FullName {
			if fullName.Tag == 6 {
				cls = append(cls, string(fullName.Bytes))
			}
		}
	}
	return cls, nil
}

func main() {
	files := []string{
		`asncer4\00Z.cer`,
	}
	for _, file := range files {
		b, err := fileutil.ReadFileToBytes(file)
		if err != nil {
			fmt.Println(file, err)
			return
		}
		certificate := Certificate{}
		_, err = Unmarshal(b, &certificate)
		//fmt.Println("certificate:", jsonutil.MarshallJsonIndent(certificate), len(res), err)
		//fmt.Println(len(certificate.TBSCertificate.Extensions))
		for i := range certificate.TBSCertificate.Extensions {
			extension := &certificate.TBSCertificate.Extensions[i]
			fmt.Println(extension.Oid.String())

			if extension.Oid.String() == "1.3.6.1.5.5.7.1.7" {
				// IpBlocks
				//ipAddrBlocks, err := asn1cert.ParseToIpAddressBlocks(extension.Value)
				//fmt.Println("1.3.6.1.5.5.7.1.7:", jsonutil.MarshalJson(ipAddrBlocks), err)
				fmt.Println("ParseToIpAddressBlocks(): value:", convert.PrintBytesOneLine(extension.Value))
			} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.8" {
				// Asns
				fmt.Println("GetAsns(): value:", convert.PrintBytesOneLine(extension.Value))
				GetAsns(extension.Value)

			} else if extension.Oid.String() == "2.5.29.31" {
				// Crl
				fmt.Println("GetCrldp(): value:", convert.PrintBytesOneLine(extension.Value))
				seqs, err := GetCrldp(extension.Value)
				fmt.Println("2.5.29.31:", seqs, err)
			}
		}
	}

}
