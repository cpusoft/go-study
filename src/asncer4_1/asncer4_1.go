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
type AnsPointInt struct {
	Min int
	Max int
}

type AsnPoint2 struct {
	AsnPointName AnsPointName2 `asn1:"optional,tag:0"`
}

type AnsPointName2 struct {
	AsnName  int   `asn1:"optional"`
	AsnNames []int `asn1:"optional"`
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
		if asnName.Tag == TagInteger {
			fmt.Println("GetAsns(): is TagInteger")
			b := asnName.Bytes
			asn := big.NewInt(0).SetBytes(b)
			fmt.Println(asn)
		} else if asnName.Tag == TagSequence {
			fmt.Println("GetAsns(): is TagSequence")
			var ansPointInt AnsPointInt
			_, err := Unmarshal(asnName.FullBytes, &ansPointInt)
			if err != nil {
				fmt.Println("GetAsns(): Unmarshal fail:", err)
				return asnPoint, err
			}
			fmt.Println(ansPointInt)
		}
	}
	return asnPoint, nil
}
func GetAsns2(value []byte) {
	fmt.Println("GetAsns(): value:", convert.PrintBytesOneLine(value))

	var asnPoint AnsPointName2
	_, err := Unmarshal(value, &asnPoint)
	if err != nil {
		fmt.Println("GetAsns(): Unmarshal fail:", err)
		return
	}

	fmt.Println("GetAsns(): asnPoint:", jsonutil.MarshalJson(asnPoint))

	return
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

type InfoAccessAsn1 struct {
	InfoAccessOid      ObjectIdentifier
	InfoAccessByteAsn1 []byte `asn1:"implicit,tag:6"`
	//Value string `asn1:"implicit,tag:6"`
}

func ParseInfoAccessAsn1ByAsn1(data []byte) (infoAccessAsn1s []InfoAccessAsn1, err error) {
	fmt.Println("ParseInfoAccessAsn1ByAsn1(): len(data):", len(data))
	infoAccessAsn1s = make([]InfoAccessAsn1, 0)
	_, err = Unmarshal(data, &infoAccessAsn1s)
	if err != nil {
		fmt.Println("ParseInfoAccessAsn1ByAsn1(): Unmarshal data fail, len(data):", len(data), err)
		return nil, err
	}
	fmt.Println("ParseInfoAccessAsn1ByAsn1(): infoAccessAsn1s:", jsonutil.MarshalJson(infoAccessAsn1s))
	return infoAccessAsn1s, nil
}

// RFC 5280 4.2.1.4
type policy struct {
	Policy   ObjectIdentifier
	Policy2s []Policy2
	// policyQualifiers omitted
}

type Policy2 struct {
	Policy ObjectIdentifier
	Url    string
}

func GetPolicies(value []byte) {
	policies := make([]policy, 0)
	_, err := Unmarshal(value, &policies)
	if err != nil {
		return
	}
	fmt.Println(len(policies))

	for i := range policies {
		oid1 := policies[i].Policy.String()
		fmt.Println(" oid1:", oid1)

		policy2 := policies[i].Policy2s
		for j := range policy2 {
			fmt.Println("policy2:", policy2[j])
		}
	}
	return
}

type RsaModel struct {
	Name string `json:"name"`
	// "85:89:43:5d:71:af:...."
	Modulus  string `json:"modulus"`
	Exponent uint64 `json:"exponent"`
}
type Sha256RsaModel struct {
	Name string `json:"name"`
	// may empty
	// "85:89:43:5d:71:af:...."
	Sha256 string `json:"sha256"`
}

/*
	type AlgorithmIdentifier struct {
		Algorithm  asn1.ObjectIdentifier
		Parameters asn1.RawValue `asn1:"optional"`
	}
*/
func ParseSignatureInnerAlgorithmByAsn1(signatureAlgorithm AlgorithmIdentifier) (signatureInnerAlgorithm Sha256RsaModel, err error) {
	fmt.Println(signatureAlgorithm.Algorithm.String())
	if signatureAlgorithm.Algorithm.String() == `1.2.840.113549.1.1.11` {
		signatureInnerAlgorithm.Name = `sha256WithRSAEncryption`
	}
	if signatureAlgorithm.Parameters.Tag == TagNull {
		signatureInnerAlgorithm.Sha256 = ""
	}
	return signatureInnerAlgorithm, nil
}

func main() {
	files := []string{
		`G:\Download\cert\asncer4_1\sx.cer`,
	}
	for _, file := range files {
		b, err := fileutil.ReadFileToBytes(file)
		if err != nil {
			fmt.Println(file, err)
			return
		}
		certificate := Certificate{}
		_, err = Unmarshal(b, &certificate)
		signatureInnerAlgorithm := certificate.TBSCertificate.SignatureAlgorithm
		fmt.Println("signatureInnerAlgorithm:", jsonutil.MarshalJson(signatureInnerAlgorithm))
		//signatureInnerAlgorithm, _ := ParseSignatureInnerAlgorithmByAsn1(certificate.SignatureAlgorithm)
		//fmt.Println("signatureInnerAlgorithm:", jsonutil.MarshalJson(signatureInnerAlgorithm))

		publickKey := certificate.TBSCertificate.PublicKey
		fmt.Println("publickKey:", jsonutil.MarshalJson(publickKey))

		return

		//fmt.Println("certificate:", jsonutil.MarshallJsonIndent(certificate), len(res), err)
		//fmt.Println(len(certificate.TBSCertificate.Extensions))
		for i := range certificate.TBSCertificate.Extensions {
			extension := &certificate.TBSCertificate.Extensions[i]
			fmt.Println(extension.Oid.String())

			if extension.Oid.String() == "1.3.6.1.5.5.7.1.7" {
				// IpBlocks
				//ipAddrBlocks, err := asn1cert.ParseToIpAddressBlocks(extension.Value)
				//fmt.Println("1.3.6.1.5.5.7.1.7:", jsonutil.MarshalJson(ipAddrBlocks), err)
				//fmt.Println("ParseToIpAddressBlocks(): value:", convert.PrintBytesOneLine(extension.Value))
			} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.8" {
				// Asns
				//fmt.Println("\n\n\n---------------------------------------------")
				//fmt.Println("GetAsns(): asn:", convert.PrintBytesOneLine(extension.Value))
				//	GetAsns2(extension.Value)
				//fmt.Println("---------------------------------------------\n\n\n")
			} else if extension.Oid.String() == "2.5.29.31" {
				// Crl
				//fmt.Println("GetCrldp(): value:", convert.PrintBytesOneLine(extension.Value))
				//	seqs, err := GetCrldp(extension.Value)
				//fmt.Println("2.5.29.31:", seqs, err)
			} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.1" {
				// authorityInfoAccess cerModel.AiaModel,
				// cerModel.SiaModel
				//fmt.Println("\n\n\n---------------------------------------------")
				//infoAccessAsn1s, err := ParseInfoAccessAsn1ByAsn1(extension.Value)
				//fmt.Println("AiaModel:", infoAccessAsn1s, err)
				//for i := range infoAccessAsn1s {
				//	fmt.Println("oid:", infoAccessAsn1s[i].InfoAccessOid.String(),
				//		"  value:", string(infoAccessAsn1s[i].InfoAccessByteAsn1))
				//}
				//fmt.Println("---------------------------------------------\n\n\n")
			} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.11" {
				// cerModel.SiaModel
				//fmt.Println("\n\n\n---------------------------------------------")
				//infoAccessAsn1s, err := ParseInfoAccessAsn1ByAsn1(extension.Value)
				//fmt.Println("SiaModel:", infoAccessAsn1s, err)
				//for i := range infoAccessAsn1s {
				//	fmt.Println("oid:", infoAccessAsn1s[i].InfoAccessOid.String(),
				//		"  value:", string(infoAccessAsn1s[i].InfoAccessByteAsn1))
				//}
				//fmt.Println("---------------------------------------------\n\n\n")

			} else if extension.Oid.String() == "2.5.29.32" {
				// Policies
				fmt.Println("\n\n\n\n===========================")
				GetPolicies(extension.Value)
				fmt.Println("===========================\n\n\n\n")
				return
			}
		}
	}

}
