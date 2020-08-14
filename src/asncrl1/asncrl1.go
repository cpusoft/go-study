package main

import (
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
)

type CertificateList struct {
	TBSCertList        TBSCertificateList
	SignatureAlgorithm AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

// TBSCertificateList represents the ASN.1 structure of the same name. See RFC
// 5280, section 5.1.
type TBSCertificateList struct {
	Raw                 asn1.RawContent
	Version             int `asn1:"optional,default:0"`
	Signature           AlgorithmIdentifier
	Issuer              RDNSequence
	ThisUpdate          time.Time
	NextUpdate          time.Time            `asn1:"optional"`
	RevokedCertificates []RevokedCertificate `asn1:"optional"`
	Extensions          []Extension          `asn1:"tag:0,optional,explicit"`
}
type AlgorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue `asn1:"optional"`
}

type RDNSequence []RelativeDistinguishedNameSET
type RelativeDistinguishedNameSET []AttributeTypeAndValue

type AttributeTypeAndValue struct {
	Type  asn1.ObjectIdentifier
	Value interface{}
}
type Extension struct {
	Oid      asn1.ObjectIdentifier
	Critical bool `asn1:"optional"`
	Value    []byte
}

// RevokedCertificate represents the ASN.1 structure of the same name. See RFC
// 5280, section 5.1.
type RevokedCertificate struct {
	SerialNumber   *big.Int
	RevocationTime time.Time
	Extensions     []Extension `asn1:"optional"`
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
func GetOctectUint64(value []byte) (uint64, error) {
	tmp := asn1.RawValue{}
	_, err := asn1.Unmarshal(value, &tmp)
	if err != nil {
		return 0, err
	}
	return convert.Bytes2Uint64(tmp.Bytes), nil
}

func main() {
	var file string
	file = `E:\Go\go-study\src\asncrl1\1.crl`
	file = `E:\Go\go-study\src\asncrl1\2.crl`
	b, err := fileutil.ReadFileToBytes(file)
	if err != nil {
		fmt.Println(file, err)
		return
	}
	certificate := CertificateList{}
	asn1.Unmarshal(b, &certificate)
	fmt.Println("certificate:", jsonutil.MarshallJsonIndent(certificate))
	fmt.Println(len(certificate.TBSCertList.Extensions))
	for i := range certificate.TBSCertList.Extensions {
		extension := &certificate.TBSCertList.Extensions[i]
		fmt.Println(extension.Oid.String())
		if extension.Oid.String() == "2.5.29.35" {
			// authorityKeyIdentifier
			fmt.Println(GetOctectStringSequenceString(extension.Value))
		} else if extension.Oid.String() == "2.5.29.20" {
			// subjectKeyIdentifier
			fmt.Println(GetOctectUint64(extension.Value))
		}
	}
}
