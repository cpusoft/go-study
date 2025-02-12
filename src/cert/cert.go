package main

import (
	_ "crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	_ "encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cpusoft/goutil/jsonutil"
)

// cer 文件
// IPAddress
type IPAddressRange struct {
	Min string `json:"min"`
	Max string `json:"max"`
}
type IPAddressOrRange struct {
	AddressPrefix string         `json:"addressPrefix"`
	AddressRange  IPAddressRange `json:"addressRange"`
}

// ASN
type ASRange struct {
	Min uint64 `json:"min"`
	Max uint64 `json:"max"`
}
type ASIdOrRange struct {
	ASId    uint64  `json:"asId"`
	ASRange ASRange `json:"asRange"`
}
type CerInfo struct {
	Version               int                `json:"version"`
	SN                    string             `json:"sn"`
	NotBefore             string             `json:"notBefore"`
	NotAfter              string             `json:"notAfter"`
	BasicConstraintsValid bool               `json:"basicConstraintsValid"`
	IsRoot                bool               `json:"isRoot"`
	DNSNames              []string           `json:"dnsNames"`
	EmailAddresses        []string           `json:"emailAddresses"`
	IPAddresses           []net.IP           `json:"ipAddresses"`
	Subject               string             `json:"subject"`
	SubjectAll            string             `json:"subjectAll"`
	Issuer                string             `json:"issuer"`
	IssuerAll             string             `json:"issuerAll"`
	Ski                   []byte             `json:"ski"`
	Aki                   []byte             `json:"aki"`
	CRLdp                 []string           `json:"crldp"`
	Aia                   []string           `json:"aia"`
	IPAddressOrRange      []IPAddressOrRange `json:"ipAddressOrRange"`
	AsNum                 []ASIdOrRange      `json:"asNum"`
	Rdi                   []ASIdOrRange      `json:"rdi"`
}

type CrlRevokedCert struct {
	SN             string `jsong:"sn"`
	RevocationTime string `json:"revocationTime"`
}

type CrlInfo struct {
	Version         int              `json:"version"`
	Issuer          string           `json:"issuer"`
	ThisUpdate      string           `json:"thisUpdate"`
	NextUpdate      string           `json:"nextUpdate"`
	HasExpired      string           `json:"hasExpired"`
	CrlRevokedCerts []CrlRevokedCert `json:"CrlRevokedCerts"`
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

func parseCer(file string) error {
	//300E300C040200013006030402B9A6FC     16
	/*  ` 0x0c 0x04 0x02 0x00 0x01 0x30 0x06 0x03 0x04 0x02 0xb9 0xa6 0xfc
	300E300C040200013006030402B9A6FC
	300D300B04020001300503030084FC

	oidValue:
	0x30 0x0e      0x30是SEQUENCE类型固定的， 0e是后面长度
		0x30 0x0c  0x30是SEQUENCE类型固定的， 0c是后面长度, 从这里开始
			0x04 0x02 0x00 0x01     0x04, 0x02, 0x00, 0x01, // address family: IPv4    对比：0x04, 0x02, 0x00, 0x02, // address family: IPv6
				0x30 0x06
					0x03 0x04
					 0x02 0xb9 0xa6 0xfc
					      185.166.252/22
	type: 48
	len: 14
	oidIP:
	0x30 0x0c 0x04 0x02 0x00 0x01 0x30 0x06 0x03 0x04 0x02 0xb9 0xa6 0xfc
	`*/

	rootCa := file
	caBlock, err := ioutil.ReadFile(rootCa)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return err
	}

	cert, err := x509.ParseCertificate(caBlock)
	if err != nil {
		fmt.Println("ParseCertificate err:", err)
		return err
	}

	cerInfo := CerInfo{}
	cerInfo.SN = fmt.Sprintf("%x", cert.SerialNumber)
	cerInfo.Version = cert.Version
	cerInfo.DNSNames = cert.DNSNames
	cerInfo.EmailAddresses = cert.EmailAddresses
	cerInfo.IPAddresses = cert.IPAddresses
	cerInfo.BasicConstraintsValid = cert.BasicConstraintsValid
	cerInfo.IsRoot = cert.IsCA
	cerInfo.NotBefore = cert.NotBefore.Format("2006-01-02 15:04:05")
	cerInfo.NotAfter = cert.NotAfter.Format("2006-01-02 15:04:05")
	cerInfo.Subject = cert.Subject.CommonName
	cerInfo.SubjectAll, _ = getDNFromName(cert.Subject, "/")
	cerInfo.Issuer = cert.Issuer.CommonName
	cerInfo.IssuerAll, _ = getDNFromName(cert.Issuer, "/")
	cerInfo.IPAddressOrRange = make([]IPAddressOrRange, 0)
	cerInfo.AsNum = make([]ASIdOrRange, 0)
	cerInfo.Rdi = make([]ASIdOrRange, 0)

	fmt.Println("cerInfo=:", jsonutil.MarshalJson(cerInfo))

	//	fmt.Printf("serialNumber=%s\r\n", cert.Subject.SerialNumber)
	//	fmt.Printf("serialNumber=%s\r\n", cert.Issuer.SerialNumber)
	//	fmt.Printf("SN=%v\r\n", cert.SerialNumber.Uint64())

	//使用者密钥标识符
	cerInfo.Ski = cert.SubjectKeyId
	//颁发机构密钥标识符
	cerInfo.Aki = cert.AuthorityKeyId
	//CRL分发点
	cerInfo.CRLdp = cert.CRLDistributionPoints
	//颁发机构信息访问
	cerInfo.Aia = cert.IssuingCertificateURL

	oidIpAddressKey := "1.3.6.1.5.5.7.1.7"
	oidASKey := "1.3.6.1.5.5.7.1.8"
	for _, extension := range cert.Extensions {
		oid := extension.Id
		if oidIpAddressKey == oid.String() {
			err := parseIpAddressExtension(extension, &cerInfo.IPAddressOrRange)
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else if oidASKey == oid.String() {
			err := parseAsnExtension(extension, &cerInfo.AsNum, &cerInfo.Rdi)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	jsonCer, _ := json.Marshal(cerInfo)

	fmt.Printf("%+v", string(jsonCer))
	/*
		fmt.Printf("subject: /CN=%s/serialNumber=%s\r\n", cert.Subject.CommonName, cert.Subject.SerialNumber)
		fmt.Printf("issuer: /CN=%s/serialNumber=%s\r\n", cert.Issuer.CommonName, cert.Issuer.SerialNumber)


		publicInfo := cert.RawSubjectPublicKeyInfo
		printBytes("publicInfo", publicInfo)
		publicKey := cert.PublicKey
		//printBytes("publicKey", publicKey)
		fmt.Println(publicKey)

		aki := cert.AuthorityKeyId
		printBytes("aki", aki)
		fmt.Println("len(ExtKeyUsage):", len(cert.ExtKeyUsage))
		for _, eku := range cert.ExtKeyUsage {
			fmt.Println("%v\r\n", eku)
		}
		//颁发机构信息访问
		fmt.Println("aia:", cert.IssuingCertificateURL)






	*/
	return nil
}

func parseCrl(file string) error {
	crl := file
	crlBlock, err := ioutil.ReadFile(crl)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return err
	}
	certList, err := x509.ParseCRL(crlBlock)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tbsCertList := certList.TBSCertList
	crlInfo := CrlInfo{}
	crlInfo.Version = tbsCertList.Version
	crlInfo.Issuer, _ = getDNFromRDNSeq(tbsCertList.Issuer, "/")
	crlInfo.ThisUpdate = tbsCertList.ThisUpdate.Local().Format("2006-01-02 15:04:05") //LoadLocation("Local")
	crlInfo.NextUpdate = tbsCertList.NextUpdate.Local().Format("2006-01-02 15:04:05")
	crlInfo.HasExpired = strconv.FormatBool(certList.HasExpired(time.Now()))
	//exts := tbsCertList.Extensions
	crlInfo.CrlRevokedCerts = make([]CrlRevokedCert, 0)
	revokedCerts := tbsCertList.RevokedCertificates
	for _, revokedCert := range revokedCerts {
		crlRevokedCert := CrlRevokedCert{}
		crlRevokedCert.SN = fmt.Sprintf("%x", revokedCert.SerialNumber)
		crlRevokedCert.RevocationTime = revokedCert.RevocationTime.Local().Format("2006-01-02 15:04:05")
		crlInfo.CrlRevokedCerts = append(crlInfo.CrlRevokedCerts, crlRevokedCert)
	}

	jsonCrl, _ := json.Marshal(crlInfo)
	fmt.Printf("%+v", string(jsonCrl))

	return nil
}

/*
manifest.asn
https://datatracker.ietf.org/doc/rfc6486/

-- Declaration for c->asn compatibility
--
--

DEFINITIONS IMPLICIT TAGS ::=
-- imports
IMPORTS AlgorithmIdentifier FROM Algorithms IN Algorithms.asn,

	Extensions IPAddressOrRangeA Attribute FROM extensions IN extensions.asn,
	Certificate Version FROM certificate IN certificate.asn,
	Name FROM name IN name.asn;

-- Manifest Specification

	Manifest ::= SEQUENCE
	 {
	 version         [0] Manifestversion DEFAULT v1,
	 manifestNumber  INTEGER,
	 thisUpdate      GeneralizedTime,
	 nextUpdate      GeneralizedTime,
	 fileHashAlg     OBJECT IDENTIFIER,
	 fileList        SEQUENCE SIZE (0..MAX) OF FileAndHash
	 }

	 Manifestversion ::= INTEGER { v1(0) } (v1)

FileAndHash ::= SEQUENCE

	{
	file        IA5String,
	hash        BIT STRING
	}
*/
func parseMft(file string) error {

	return nil
}
func main() {
	var file string
	if len(os.Args) == 2 {
		file = os.Args[1]
	} else {

		//file = `E:\Go\go-study\src\main\cert\ROUTER-0000FBF0_new.cer`
		//file = `E:\Go\go-study\src\main\cert\ROUTER-00010000_new.cer`
		//file = `E:\Go\go-study\src\main\cert\err1.cer`
		//file = `E:\Go\go-study\src\main\cert\H.cer`
		//file = `E:\Go\go-study\src\main\secruity\1.cer`
		file = `E:\Go\go-study\src\main\secruity\range_ipv6.cer`
		file = `ca.cer`
		//fmt.Println("usage: ./cert 1.cer")
		//return
	}
	//`E:\Go\go-study\src\main\cert\root.cer`
	fmt.Println(file)
	var err error
	certFile := strings.ToLower(file)
	if strings.HasSuffix(certFile, ".cer") {
		err = parseCer(certFile)
	} else if strings.HasSuffix(certFile, ".crl") {
		err = parseCrl(certFile)
	} else if strings.HasSuffix(certFile, ".mft") {
		err = parseMft(certFile)
	}
	if err != nil {
		fmt.Println(err)
	}

}

func parseIpAddressExtension(extension pkix.Extension, ipAddressOrRanges *[]IPAddressOrRange) error {
	extensionValue := extension.Value
	//critical := extension.Critical
	if len(extensionValue) == 0 {
		fmt.Println("not found oid:", extensionValue)
		return errors.New("not found oid")
	}
	return parseIpAddressExtensionValue(extensionValue, ipAddressOrRanges)
}

func parseIpAddressExtensionValue(extensionValue []byte, ipAddressOrRanges *[]IPAddressOrRange) error {
	// sequences 整个的数组
	ipAddrBlocksType := extensionValue[0]
	ipAddrBlocksLen := extensionValue[1]
	ipAddrBlocksValue := extensionValue[2:]
	printAsn("ipAddrBlocks", ipAddrBlocksType, ipAddrBlocksLen, ipAddrBlocksValue)

	tmpBlock := ipAddrBlocksValue
	//循环数组，
	for {
		ipAddrBlockType := tmpBlock[0]
		ipAddrBlockLen := tmpBlock[1]
		ipAddrBlockValue := tmpBlock[2 : 2+ipAddrBlockLen]
		printAsn("ipAddrBlock", ipAddrBlockType, ipAddrBlockLen, ipAddrBlockValue)
		err := parseIpAddrBlock(tmpBlock, ipAddressOrRanges)
		if err != nil {
			fmt.Println(err)
			return err
		}
		tmpBlock = tmpBlock[2+ipAddrBlockLen:]
		if len(tmpBlock) == 0 {
			break
		}
	}
	return nil
}

const (
	ipv4    = 0x01
	ipv6    = 0x02
	ipv4len = 32
	ipv6len = 128
)
const (
	nul       = 0x05
	bitstring = 0x03
	sequence  = 0x30
)

func parseIpAddrBlock(ipAddrBlock []byte, ipAddressOrRanges *[]IPAddressOrRange) error {
	//ipaddressFamily: 包括addressFamily+ipAddressChoice
	ipAddressFamilyType := ipAddrBlock[0]
	ipAddressFamilyLen := ipAddrBlock[1]
	ipAddressFamilyValue := ipAddrBlock[2 : 2+ipAddressFamilyLen]
	printAsn("ipAddressFamily", ipAddressFamilyType, ipAddressFamilyLen, ipAddressFamilyValue)

	//读取addressFamily，类型
	/*
		//
		// 04 03 0001  01 addressFamily: ：  IPv4 Unicast:
		//		前两位是，必有，a two-octet Address Family Identifier (AFI) https://www.iana.org/assignments/address-family-numbers/address-family-numbers.xhtml
					 0001是 ipv4；
		//      后一位是 ，可选， a one-octet Subsequent Address Family Identifier (SAFI) https://www.iana.org/assignments/safi-namespace/safi-namespace.xhtml
					01是 unicast
		 前面是必有，后面是可选
	*/
	addressFamilyType := ipAddressFamilyValue[0]
	addressFamilyLen := ipAddressFamilyValue[1]
	addressFamilyValue := ipAddressFamilyValue[2 : 2+addressFamilyLen]
	printAsn("addressFamily", addressFamilyType, addressFamilyLen, addressFamilyValue)
	//fmt.Println("addressFamilyValue[addressFamilyLen-1]", addressFamilyValue[addressFamilyLen-1])

	ipType := -1
	if addressFamilyValue[1] == ipv4 {
		ipType = ipv4
	} else if addressFamilyValue[1] == ipv6 {
		ipType = ipv6
	}
	if ipType == -1 {
		return errors.New("error iptype")
	}
	//读取ipAddressChoice，注意是从ipAddrBlock开始--
	//即addressFamilyValue节尾--截取的。 2是ipAddressFamily的头，2是addressFamily的头，然后再加上addressFamilyLen
	ipAddressChoice := ipAddrBlock[4+addressFamilyLen:]
	ipAddressChoiceType := ipAddressChoice[0]
	ipAddressChoiceLen := ipAddressChoice[1]
	ipAddressChoiceValue := ipAddressChoice[2 : 2+ipAddressChoiceLen]
	printAsn("ipAddressChoice", ipAddressChoiceType, ipAddressChoiceLen, ipAddressChoiceValue)

	if ipAddressChoiceType == nul {
		return nil
	} else if ipAddressChoiceType == sequence {
		// ok, continue
	} else {
		return errors.New("error IPAddressChoic")
	}

	// ipAddressChoice包括addressesOrRanges的数组，每个addressesOrRange有可能是addressPrefix(0x03开头)，也有可能是addressRange(0x30开头)
	// 循环读取数组
	tmpAddressOrRange := ipAddressChoiceValue
	for {
		//addressesOrRangeType := tmpAddressOrRange[0]
		addressesOrRangeLen := tmpAddressOrRange[1]
		//addressesOrRangeValue := tmpAddressOrRange[2 : 2+addressesOrRangeLen]
		//printAsn("addressesOrRange", addressesOrRangeType, addressesOrRangeLen, addressesOrRangeValue)
		err := parseAddressesOrRange(tmpAddressOrRange, ipType, ipAddressOrRanges)
		if err != nil {
			fmt.Println(err)
			return err
		}
		tmpAddressOrRange = tmpAddressOrRange[2+addressesOrRangeLen:]
		if len(tmpAddressOrRange) == 0 {
			break
		}
	}
	return nil
}

const (
	addressPrefix = iota
	addressRange
)

func parseAddressesOrRange(addressesOrRange []byte, ipType int, ipAddressOrRanges *[]IPAddressOrRange) error {
	//每个addressesOrRange有可能是addressPrefix(0x03开头)，也有可能是addressRange(0x30开头)
	if len(addressesOrRange) == 0 {
		return errors.New("lenght of addressesOrRange is zero")
	}
	addressesOrRangeOneType := addressesOrRange[0]
	addressesOrRangeOneLen := addressesOrRange[1]
	addressesOrRangeOneValue := addressesOrRange[2 : 2+addressesOrRangeOneLen]
	//printAsn("addressesOrRangeOne", addressesOrRangeOneType, addressesOrRangeOneLen, addressesOrRangeOneValue)
	// 注意这里传入的是addressesOrRangeOneValue
	if addressesOrRangeOneType == bitstring {
		parseAddressPrefix(addressesOrRangeOneValue, addressesOrRangeOneLen, ipType, ipAddressOrRanges)
	} else if addressesOrRangeOneType == sequence {
		parseAddressRange(addressesOrRangeOneValue, addressesOrRangeOneLen, ipType, ipAddressOrRanges)
	} else {
		return errors.New("addressesOrRangeOneType is error")
	}
	return nil
}

func parseAddressPrefix(addressPrefix []byte, addressesOrRangeOneLen byte, ipType int, ipAddressOrRanges *[]IPAddressOrRange) error {
	//  03 03 04 b010              addressPrefix    172.16/12
	//传入的第0位是unusedbit位
	//第2位标明长度，-1后(unused bit占用了1位)，为ip地址应该的长度: 标明长度3，应该长度为2
	// 第3位，固定的unused bit位： 为4
	// unusedbit = 32- 应该的长度*8 - prefix  =32-2*8-prefix
	// prefix = 32- 应该的长度*8 - unusedbit = 32 - 2*8 - 4 = 12
	//printBytes("addressPrefix:", addressPrefix)
	addressShouldLen, _ := strconv.Atoi(fmt.Sprintf("%d", addressesOrRangeOneLen))
	unusedBitLen, _ := strconv.Atoi(fmt.Sprintf("%d", addressPrefix[0]))

	address := addressPrefix[1:]
	ipAddress := ""

	if ipType == ipv4 {
		// ipv4 的CIDR 表示法
		prefix := ipv4len - 8*(addressShouldLen-1) - unusedBitLen
		//fmt.Println(fmt.Sprintf("prefix := ipv4len - 8*(addressShouldLen-1) - unusedBitLen:  %d := %d - 8 *(%d-1)-  %d \r\n",
		//	prefix, ipv4len, addressShouldLen, unusedBitLen))

		//printBytes("address:", address)

		ipv4Address := ""
		for i := 0; i < len(address); i++ {
			ipv4Address += fmt.Sprintf("%d", address[i])
			if i < len(address)-1 {
				ipv4Address += "."
			}
		}
		ipv4Address += "/" + fmt.Sprintf("%d", prefix)
		ipAddress = ipv4Address
		//fmt.Println(ipv4Address)
	} else if ipType == ipv6 {
		// ipv6的前缀表示法，和ipv4不一样
		prefix := 8*(addressShouldLen-1) - unusedBitLen
		//		fmt.Println(fmt.Sprintf("prefix :=  8*(addressShouldLen-1) - unusedBitLen:  %d := 8 *(%d-1)-  %d \r\n",
		//			prefix, addressShouldLen, unusedBitLen))

		//printBytes("address:", address)

		ipv6Address := ""
		for i := 0; i < len(address); i++ {
			ipv6Address += fmt.Sprintf("%02x", address[i])
			if i%2 == 1 && i < len(address)-1 {
				ipv6Address += ":"
			}
		}
		//补齐位数
		if len(address)%2 == 1 {
			ipv6Address += "00"
		}
		ipv6Address += "/" + fmt.Sprintf("%d", prefix)
		ipAddress = ipv6Address
		//fmt.Println(ipv6Address)
	}
	ipAddressOrRange := IPAddressOrRange{}
	ipAddressOrRange.AddressPrefix = ipAddress
	*ipAddressOrRanges = append(*ipAddressOrRanges, ipAddressOrRange)
	//jsonCer, _ := json.Marshal(ipAddressOrRanges)
	//fmt.Printf("in parseAddressPrefix(): %+v", string(jsonCer))

	return nil
}
func parseAddressRange(addressRange []byte, addressesOrRangeOneLen byte, ipType int, ipAddressOrRanges *[]IPAddressOrRange) error {
	//传入的是两个sequence，第一个是min，第二个是max
	// Value值，跳过了unused bit位，所以是从3开始，并且长度-1
	//fmt.Println("parseAddressRange():  ipType:", ipType)
	//minType := addressRange[0]
	minLen := addressRange[1]
	minValue := addressRange[2 : 2+minLen]
	//printAsn("min", minType, minLen, minValue)

	tmp := addressRange[2+minLen:]
	//maxType := tmp[0]
	maxLen := tmp[1]
	maxValue := tmp[2 : 2+maxLen]
	//printAsn("max", maxType, maxLen, maxValue)

	ipAddressRange := IPAddressRange{}

	if ipType == ipv4 {
		minAddr := ""
		for i := 0; i < len(minValue); i++ {
			minAddr += fmt.Sprintf("%d.", minValue[i])
		}
		for i := 0; i < 4-len(minValue); i++ {
			minAddr += fmt.Sprintf("%d.", 0)
		}
		minAddr = minAddr[0 : len(minAddr)-1]

		maxAddr := ""
		for i := 0; i < len(maxValue); i++ {
			maxAddr += fmt.Sprintf("%d.", maxValue[i])
		}
		for i := 0; i < 4-len(maxValue); i++ {
			maxAddr += fmt.Sprintf("%d.", 255)
		}
		maxAddr = maxAddr[0 : len(maxAddr)-1]

		ipAddressRange.Max = maxAddr
		ipAddressRange.Min = minAddr

		//fmt.Println("minAddr:", minAddr, "maxAddr", maxAddr)

	} else if ipType == ipv6 {
		// 先拼出整个的ipv6地址，min的用0填充，max用255填充，最后再每4位加:
		addrTmp := ""
		minAddr := ""
		for i := 0; i < len(minValue); i++ {
			addrTmp += fmt.Sprintf("%02x", minValue[i])
		}
		for i := 0; i < 16-len(minValue); i++ {
			addrTmp += fmt.Sprintf("%02x", 0)
		}
		for i := 0; i < len(addrTmp); i += 4 {
			minAddr += (addrTmp[i:i+4] + ":")
		}
		minAddr = minAddr[0 : len(minAddr)-1]

		addrTmp = ""
		maxAddr := ""
		for i := 0; i < len(maxValue); i++ {
			addrTmp += fmt.Sprintf("%02x", maxValue[i])
		}
		for i := 0; i < 16-len(maxValue); i++ {
			addrTmp += fmt.Sprintf("%02x", 255)
		}
		for i := 0; i < len(addrTmp); i += 4 {
			maxAddr += (addrTmp[i:i+4] + ":")
		}
		maxAddr = maxAddr[0 : len(maxAddr)-1]
		//fmt.Println("minAddr:", minAddr, "maxAddr", maxAddr)

		ipAddressRange.Max = maxAddr
		ipAddressRange.Min = minAddr

	}
	ipAddressOrRange := IPAddressOrRange{}
	ipAddressOrRange.AddressRange = ipAddressRange
	*ipAddressOrRanges = append(*ipAddressOrRanges, ipAddressOrRange)

	//jsonCer, _ := json.Marshal(ipAddressOrRanges)
	//fmt.Printf("in parseAddressRange(): %+v", string(jsonCer))

	return nil
}

func parseAsnExtension(extension pkix.Extension, asNum *[]ASIdOrRange, rdi *[]ASIdOrRange) error {
	extensionValue := extension.Value
	//critical := extension.Critical
	if len(extensionValue) == 0 {
		fmt.Println("not found oid:", extensionValue)
		return errors.New("not found oid")
	}
	return parseAsnExtensionValue(extensionValue, asNum, rdi)
}

const (
	ASNUM = byte(0xa0)
	RDI   = byte(0xa1)
)

func parseAsnExtensionValue(extensionValue []byte, asNum *[]ASIdOrRange, rdi *[]ASIdOrRange) error {
	//asIdentifiersType := extensionValue[0]
	//asIdentifiersLen := extensionValue[1]
	asIdentifiersValue := extensionValue[2:]

	//printAsn("AsnExtensionValue", asIdentifiersType, asIdentifiersLen, asIdentifiersValue)
	tmpBlock := asIdentifiersValue
	//循环数组，其实就两组：asnum 和rdi
	for i := 0; i <= 1; i++ {
		asIdentifierType := tmpBlock[0]
		asIdentifierLen := tmpBlock[1]
		asIdentifierValue := tmpBlock[2 : 2+asIdentifierLen]
		//printAsn("asIdentifier", asIdentifierType, asIdentifierLen, asIdentifierValue)
		var err error
		if asIdentifierType == ASNUM {
			err = parseASNum(asIdentifierValue, asNum)
		} else if asIdentifierType == RDI {
			err = parseRdi(asIdentifierValue, rdi)
		}
		if err != nil {
			fmt.Println(err)
			return err
		}
		tmpBlock = tmpBlock[2+asIdentifierLen:]
		if len(tmpBlock) == 0 {
			break
		}
	}

	return nil
}

func parseASNum(asIdentifier []byte, asIdOrRanges *[]ASIdOrRange) error {
	//ASIdOrRange: 包括ASId+ASRange
	//asIdsOrRangesType := asIdentifier[0]
	asIdsOrRangesLen := asIdentifier[1]
	asIdsOrRangesValue := asIdentifier[2 : 2+asIdsOrRangesLen]
	//printAsn("asIdsOrRanges", asIdsOrRangesType, asIdsOrRangesLen, asIdsOrRangesValue)

	asIdOrRange := ASIdOrRange{}

	//注意，这里max和min前面还有个sequence类型
	asIdOrRangesType := asIdsOrRangesValue[0]
	asIdOrRangesLen := asIdsOrRangesValue[1]
	asIdOrRangesValue := asIdsOrRangesValue[2 : 2+asIdOrRangesLen]
	//printAsn("asIdOrRanges", asIdOrRangesType, asIdOrRangesLen, asIdOrRangesValue)
	if asIdOrRangesType == 0x30 {
		asRange := ASRange{}

		asnMinLen := asIdOrRangesValue[1]
		asnMinValue := asIdOrRangesValue[2 : 2+asnMinLen]
		//fmt.Println("asnMinValue", asnMinValue)
		asRange.Min = bytesConvertToUint64(asnMinValue)

		//fmt.Println(asIdOrRangesLen, asnMinLen)
		if asIdOrRangesLen > 2+asnMinLen {
			//asnMaxLen := asIdsOrRangesValue[2+asnMinLen+1]
			asnMaxValue := asIdOrRangesValue[2+asnMinLen+2:]
			asRange.Max = bytesConvertToUint64(asnMaxValue)
		}

		asIdOrRange.ASRange = asRange

	} else if asIdOrRangesType == 0x02 {
		asIdOrRange.ASId = bytesConvertToUint64(asIdOrRangesValue)
		//fmt.Println(asIdOrRange.ASId)
	} else {
		return errors.New("error iptype")
	}
	*asIdOrRanges = append(*asIdOrRanges, asIdOrRange)

	return nil
}

// 当前好还没有实际值  routing domain identifiers (in the rdi element)  RFC1142
func parseRdi(asIdentifier []byte, asIdOrRanges *[]ASIdOrRange) error {
	//ASIdOrRange: 包括ASId+ASRange
	//rdiType := asIdentifier[0]
	//rdiLen := asIdentifier[1]
	//rdiValue := asIdentifier[2 : 2+rdiLen]
	//printAsn("rdi", rdiType, rdiLen, rdiValue)
	/*
		asIdOrRange := ASIdOrRange{}

		//注意，这里max和min前面还有个sequence类型
		asIdOrRangesType := asIdsOrRangesValue[0]
		asIdOrRangesLen := asIdsOrRangesValue[1]
		asIdOrRangesValue := asIdsOrRangesValue[2 : 2+asIdOrRangesLen]
		printAsn("asIdOrRanges", asIdOrRangesType, asIdOrRangesLen, asIdOrRangesValue)

		asIdOrRange.ASId = bytesConvertToUint64(asIdOrRangesValue)
		fmt.Println(asIdOrRange.ASId)
		*asIdOrRanges = append(*asIdOrRanges, asIdOrRange)
	*/
	return nil
}

func printBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func printAsn(name string, typ byte, ln byte, byt []byte) {
	//fmt.Println(fmt.Sprintf(name+"Type:0x%02x (%d)", typ, typ))
	//fmt.Println(fmt.Sprintf(name+"Len:0x%02x (%d)", ln, ln))
	//printBytes(name+"Value:", byt)
}

func printBytes(name string, byt []byte) {
	fmt.Println(name)
	for _, i := range byt {
		fmt.Print(fmt.Sprintf("0x%02x ", i))
	}
	fmt.Println("")
}

func bytesConvertToUint64(bytes []byte) uint64 {
	//fmt.Println("bytesConvertToUint64()")
	//fmt.Println(bytes)
	lens := 8 - len(bytes)
	//fmt.Println(lens)
	bb := make([]byte, lens)
	//fmt.Println(bb)
	bb = append(bb, bytes...)
	//fmt.Println(bb)
	return binary.BigEndian.Uint64(bb)
}
