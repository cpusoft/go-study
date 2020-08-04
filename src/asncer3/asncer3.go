package main

import (
	_ "crypto/x509"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"

	"asn1"
)

type certificate struct {
	//Raw                asn1.RawContent
	TBSCertificate     tbsCertificate
	SignatureAlgorithm AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

type tbsCertificate struct {
	Raw                asn1.RawContent
	Version            int `asn1:"optional,explicit,default:0,tag:0"`
	SerialNumber       *big.Int
	SignatureAlgorithm AlgorithmIdentifier
	Issuer             asn1.RawValue
	Validity           validity
	Subject            asn1.RawValue
	PublicKey          publicKeyInfo
	Extensions         []Extension `asn1:"optional,explicit,tag:3"`
}
type AlgorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue `asn1:"optional"`
}
type validity struct {
	NotBefore, NotAfter time.Time
}
type publicKeyInfo struct {
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

/*

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
	Bool     bool `asn:"optional"`
	RawValue asn1.RawValue
}
type OidAndBoolAndBytes struct {
	Oid   asn1.ObjectIdentifier
	Bool  bool `asn:"optional"`
	Value []byte
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


*/
func GetOctectString(value []byte) (string, error) {
	tmp := make([]byte, 0)
	_, err := asn1.Unmarshal(value, &tmp)
	if err != nil {
		return "", err
	}
	return convert.Bytes2String(tmp), nil
}
func GetOctectStringSequenceString(value []byte) (string, error) {
	raws := make([]asn1.RawValue, 0)
	_, err := asn1.Unmarshal(value, &raws)
	if err != nil {
		return "", err
	}
	if len(raws) > 0 {
		return convert.Bytes2String(raws[0].Bytes), nil
	} else {
		return "", errors.New("it is no sequence of []byte")
	}
}

func GetOctectStringSequenceBool(value []byte) (bool, error) {
	bools := make([]bool, 0)
	_, err := asn1.Unmarshal(value, &bools)
	if err != nil {
		return false, err
	}
	if len(bools) > 0 {
		return bools[0], nil
	} else {
		return false, errors.New("it is no sequence of []bool")
	}
}
func GetOctectStringBitString(value []byte) (asn1.BitString, error) {
	bitString := asn1.BitString{}
	_, err := asn1.Unmarshal(value, &bitString)
	if err != nil {
		return bitString, err
	}
	return bitString, nil

}

type SeqExtension struct {
	Raw   asn1.RawContent
	Oid   asn1.ObjectIdentifier
	Value []byte `asn1:"implicit,tag:6"`
	//Value string `asn1:"implicit,tag:6"`
}

func GetOctectStringSequenceOidString(value []byte) ([]SeqExtension, error) {

	seqExtensions := make([]SeqExtension, 0)
	_, err := asn1.Unmarshal(value, &seqExtensions)
	fmt.Println(len(seqExtensions))
	if err != nil {
		return nil, err
	}
	return seqExtensions, nil

}

type SeqString0 struct {
	Value []SeqString06 // `asn1:"implicit,tag:0"`
}
type SeqString06 struct {
	Value SeqString6 `asn1:"implicit,tag:0"`
}
type SeqString6 struct {
	Value asn1.RawValue `asn1:"implicit,tag:6"`
}

func GetOctectStringSeqSeqString(value []byte) (string, error) {

	raws := make([]asn1.RawValue, 0)
	_, err := asn1.Unmarshal(value, &raws)
	fmt.Println(len(raws), err)
	fmt.Println(convert.Bytes2String(raws[0].Bytes))

	seqString0 := SeqString0{}
	_, err := asn1.Unmarshal(value, &seqString0)

	fmt.Println(seqString0)
	if err != nil {
		return "", err
	}

	cdp := make([]asn1.RawValue, 0)
	asn1.Unmarshal(value, &cdp)
	cdp1 := asn1.RawValue{}
	asn1.Unmarshal(cdp[0].Bytes, &cdp1)
	cdp2 := asn1.RawValue{}
	asn1.Unmarshal(cdp1.Bytes, &cdp2)
	fmt.Println(convert.Bytes2String(cdp2.Bytes), string(cdp2.Bytes))

	cdp3 := make([]byte, 0)
	asn1.Unmarshal(cdp2.Bytes, &cdp3)
	fmt.Println("crldp:", string(cdp3))
	return string(cdp3), nil

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

x509
const (
	KeyUsageDigitalSignature KeyUsage = 1 << iota  //1 << 0 which is  0000 0001
	KeyUsageContentCommitment                      //1 << 1 which is  0000 0010
	KeyUsageKeyEncipherment                        //1 << 2 which is  0000 0100
	KeyUsageDataEncipherment                       //1 << 3 which is  0000 1000
	KeyUsageKeyAgreement                           //1 << 4 which is  0001 0000
	KeyUsageCertSign                               //1 << 5 which is  0010 0000
	KeyUsageCRLSign                                //1 << 6 which is  0100 0000
	KeyUsageEncipherOnly                           //1 << 7 which is  1000 0000
	KeyUsageDecipherOnly                           //1 <<8 which is 1 0000 0000
)

*/
func main() {
	//file := `E:\Go\go-study\src\asncer1\0.cer`
	file := `E:\Go\go-study\src\asncer3\3.cer`
	b, err := fileutil.ReadFileToBytes(file)
	if err != nil {
		fmt.Println(file, err)
		return
	}
	certificate := certificate{}
	asn1.Unmarshal(b, &certificate)
	fmt.Println("certificate:", jsonutil.MarshallJsonIndent(certificate))
	fmt.Println(len(certificate.TBSCertificate.Extensions))
	for i := range certificate.TBSCertificate.Extensions {
		extension := &certificate.TBSCertificate.Extensions[i]
		fmt.Println(extension.Oid.String())
		if extension.Oid.String() == "2.5.29.14" {
			// subjectKeyIdentifier
			fmt.Println(GetOctectString(extension.Value))
		} else if extension.Oid.String() == "2.5.29.35" {
			// authorityKeyIdentifier
			fmt.Println(GetOctectStringSequenceString(extension.Value))
		} else if extension.Oid.String() == "2.5.29.19" {
			// basicConstraints
			fmt.Println(extension.Critical)
			fmt.Println(GetOctectStringSequenceBool(extension.Value))
		} else if extension.Oid.String() == "2.5.29.15" {
			// keyUsage
			fmt.Println(extension.Critical)

			usageValue, err := GetOctectStringBitString(extension.Value)
			fmt.Println(usageValue, err)

			var tmp int
			// usageValue: 0000011
			// 从左边开始数，从0开始计数，即第5,6位为1, 则对应KeyUsageCertSign  KeyUsageCRLSign
			for i := 0; i < 9; i++ {
				//当为1时挪动，即看是第几个进行挪动
				//fmt.Println(i, usageValue.At(i))
				if usageValue.At(i) != 0 {
					tmp |= 1 << uint(i)
				}
			}
			// 先写死吧
			usage := int(tmp)
			usageStr := "Certificate Sign, CRL Sign"
			fmt.Println(usage)
			fmt.Println(usageStr)
			/*
				fmt.Println(x509.KeyUsageDigitalSignature)
				fmt.Println(x509.KeyUsageContentCommitment)
				fmt.Println(x509.KeyUsageKeyEncipherment)
				fmt.Println(x509.KeyUsageDataEncipherment)
				fmt.Println(x509.KeyUsageKeyAgreement)
				fmt.Println(x509.KeyUsageCertSign)
				fmt.Println(x509.KeyUsageCRLSign)
				fmt.Println(x509.KeyUsageEncipherOnly)
				fmt.Println(x509.KeyUsageDecipherOnly)
			*/

		} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.1" {
			// authorityInfoAccess
			seqs, err := GetOctectStringSequenceOidString(extension.Value)
			fmt.Println(len(seqs), err)
			for i := range seqs {
				fmt.Println(seqs[i].Oid, string(seqs[i].Value))
			}
		} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.11" {
			// subjectInfoAccess
			seqs, err := GetOctectStringSequenceOidString(extension.Value)
			fmt.Println(len(seqs), err)
			for i := range seqs {
				fmt.Println(seqs[i].Oid, string(seqs[i].Value))
			}
		} else if extension.Oid.String() == "2.5.29.31" {
			// cRLDistributionPoints
			seqs, err := GetOctectStringSeqSeqString(extension.Value)
			fmt.Println(seqs, err)
		}
	}
	/*
		cerParseExt := CerParseExt{}
		asn1.Unmarshal(certificate.TBSCertificate.CerRawValue.Bytes, &cerParseExt)
		fmt.Println("cerParseExt:", jsonutil.MarshallJsonIndent(cerParseExt))


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
	*/
}
