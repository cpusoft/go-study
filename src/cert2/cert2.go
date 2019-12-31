package main

import (
	_ "crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	_ "encoding/base64"
	_ "encoding/binary"
	"encoding/json"
	_ "encoding/pem"
	_ "errors"
	"fmt"
	"io/ioutil"
	_ "net"
	_ "os"
	_ "strconv"
	"strings"
	_ "time"
)

/*
id-pe-ipAddrBlocks      OBJECT IDENTIFIER ::= { id-pe 7 }

   IPAddrBlocks        ::= SEQUENCE OF IPAddressFamily

   IPAddressFamily     ::= SEQUENCE {    -- AFI & optional SAFI --
      addressFamily        OCTET STRING (SIZE (2..3)),
      ipAddressChoice      IPAddressChoice }

   IPAddressChoice     ::= CHOICE {
      inherit              NULL, -- inherit from issuer --
      addressesOrRanges    SEQUENCE OF IPAddressOrRange }

   IPAddressOrRange    ::= CHOICE {
      addressPrefix        IPAddress,
      addressRange         IPAddressRange }

   IPAddressRange      ::= SEQUENCE {
      min                  IPAddress,
      max                  IPAddress }

   IPAddress           ::= BIT STRING
*/
type IPAddressRange struct {
	Min asn1.BitString `json:"min"`
	Max asn1.BitString `json:"max"`
}
type IPAddressOrRange struct {
	AddressPrefix asn1.BitString `asn1:"optional,tag:0" json:"addressPrefix"`
	AddressRange  IPAddressRange `asn1:"optional,tag:1" json:"addressRange"`
}

//asn1.NullRawValue
type IPAddressChoice struct {
	//Inherit           []byte             `asn1:"optional,tag:0" json:"inherit"`
	AddressesOrRanges []IPAddressOrRange `asn1:"optional,tag:1" json:"addressesOrRanges"`
}

type IPAddressFamily struct {
	AddressFamily   []byte          `asn1:"tag:0" json:"addressFamily"`
	IPAddressChoice IPAddressChoice `asn1:"tag:1" json:"ipAddressChoice"`
}
type IPAddressFamilys struct {
	//	AddressFamily   []byte          `asn1:"optional,tag:0" json:"addressFamily"`
	IPAddressFamily []IPAddressFamily `asn1:"tag:0" json:"ipAddressFamily"`
}

type Tmp struct {
	r asn1.BitString `asn1:"optional,tag:0" json:"r"`
}

type IPAddrBlocks struct {
	IPAddressFamilys IPAddressFamilys `asn1:"tag:4" json:"ipAddressFamilys"`
}

/*
go 进行证书解析
参考x509.go  : 1137到1320行，如何处理的，是
先根据RFC 定义CRLDistributionPoints结构, 支持结构体嵌套，然后通过asn1.Unmarshal来进行解析即可
case 31:
				// RFC 5280, 4.2.1.13

				// CRLDistributionPoints ::= SEQUENCE SIZE (1..MAX) OF DistributionPoint
				//
				// DistributionPoint ::= SEQUENCE {
				//     distributionPoint       [0]     DistributionPointName OPTIONAL,
				//     reasons                 [1]     ReasonFlags OPTIONAL,
				//     cRLIssuer               [2]     GeneralNames OPTIONAL }
				//
				// DistributionPointName ::= CHOICE {
				//     fullName                [0]     GeneralNames,
				//     nameRelativeToCRLIssuer [1]     RelativeDistinguishedName }
// RFC 5280, 4.2.1.14
type distributionPoint struct {
	DistributionPoint distributionPointName `asn1:"optional,tag:0"`
	Reason            asn1.BitString        `asn1:"optional,tag:1"`
	CRLIssuer         asn1.RawValue         `asn1:"optional,tag:2"`
}
type distributionPointName struct {
	FullName     asn1.RawValue    `asn1:"optional,tag:0"`
	RelativeName pkix.RDNSequence `asn1:"optional,tag:1"`
}

				var cdp []distributionPoint
				if rest, err := asn1.Unmarshal(e.Value, &cdp); err != nil {
					return nil, err
				} else if len(rest) != 0 {
					return nil, errors.New("x509: trailing data after X.509 CRL distribution point")
				}
				for _, dp := range cdp {
					// Per RFC 5280, 4.2.1.13, one of distributionPoint or cRLIssuer may be empty.
					if len(dp.DistributionPoint.FullName.Bytes) == 0 {
						continue
					}

					var n asn1.RawValue
					if _, err := asn1.Unmarshal(dp.DistributionPoint.FullName.Bytes, &n); err != nil {
						return nil, err
					}
					// Trailing data after the fullName is
					// allowed because other elements of
					// the SEQUENCE can appear.

					if n.Tag == 6 {
						out.CRLDistributionPoints = append(out.CRLDistributionPoints, string(n.Bytes))
					}
				}
*/
// cer 文件
// IPAddress
/*
type IPAddressRange struct {
	Min string `json:"min"`
	Max string `json:"max"`
}
*/
type IPAddrBlocksnew struct {
	FamilType []byte           `json:"familType"`
	IPPrefx   []asn1.BitString `json:"ipPrefx"`
}

func main() {

	file := `E:\Go\go-study\data\3.cer`
	file = `E:\Go\go-study\data\range_ipv4.cer`
	caBlock, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}

	cert, err := x509.ParseCertificate(caBlock)
	if err != nil {
		fmt.Println("ParseCertificate err:", err)
		return
	}
	alg := cert.SignatureAlgorithm.String()
	fmt.Println("alg:", alg)
	palg := cert.PublicKeyAlgorithm.String()
	fmt.Println("palg:", palg)

	//https://datatracker.ietf.org/doc/rfc3779/?include_text=1
	oidIpAddressKey := "1.3.6.1.5.5.7.1.7"
	oidASKey := "1.3.6.1.5.5.7.1.8"

	for _, extension := range cert.Extensions {
		oid := extension.Id

		if oidIpAddressKey == oid.String() {

			ipAddrBlocksnew := []IPAddrBlocksnew{}

			iak := extension.Value
			printBytes("oidIpAddressKey", iak)
			if rest, err := asn1.Unmarshal(iak, &ipAddrBlocksnew); err != nil {
				fmt.Println(err)
				return
			} else if len(rest) != 0 {
				fmt.Println("x509: rest is not len0")
				return
			}
			fmt.Printf("ipAddrBlocksnew %+v\r\n", ipAddrBlocksnew)
			jsonCer, _ := json.Marshal(ipAddrBlocksnew)
			fmt.Println(string(jsonCer))

		} else if oidASKey == oid.String() {
			iak := extension.Value
			printBytes("oidASKey", iak)
		}
	}

}

func getDNFromName(namespace pkix.Name, sep string) (string, error) {
	return getDNFromRDNSeq(namespace.ToRDNSequence(), sep)
}

func getDNFromRDNSeq(rdns pkix.RDNSequence, sep string) (string, error) {
	subject := []string{}
	for _, s := range rdns {
		for _, i := range s {
			if v, ok := i.Value.(string); ok {
				if name, ok := oid[i.Type.String()]; ok {
					// <oid name>=<value>
					subject = append(subject, fmt.Sprintf("%s=%s", name, v))
				} else {
					// <oid>=<value> if no <oid name> is found
					subject = append(subject, fmt.Sprintf("%s=%s", i.Type.String(), v))
				}
			} else {
				// <oid>=<value in default format> if value is not string
				subject = append(subject, fmt.Sprintf("%s=%v", i.Type.String, v))
			}
		}
	}
	return sep + strings.Join(subject, sep), nil
}

var oid = map[string]string{
	"2.5.4.3":                    "CN",
	"2.5.4.4":                    "SN",
	"2.5.4.5":                    "serialNumber",
	"2.5.4.6":                    "C",
	"2.5.4.7":                    "L",
	"2.5.4.8":                    "ST",
	"2.5.4.9":                    "streetAddress",
	"2.5.4.10":                   "O",
	"2.5.4.11":                   "OU",
	"2.5.4.12":                   "title",
	"2.5.4.17":                   "postalCode",
	"2.5.4.42":                   "GN",
	"2.5.4.43":                   "initials",
	"2.5.4.44":                   "generationQualifier",
	"2.5.4.46":                   "dnQualifier",
	"2.5.4.65":                   "pseudonym",
	"0.9.2342.19200300.100.1.25": "DC",
	"1.2.840.113549.1.9.1":       "emailAddress",
	"0.9.2342.19200300.100.1.1":  "userid",
	"2.5.29.20":                  "CRL Number",
}

func printBytes(name string, byt []byte) {
	fmt.Println(name)
	for _, i := range byt {
		fmt.Print(fmt.Sprintf("0x%02x ", i))
	}
	fmt.Println("")
}
