package main

import (
	_ "crypto/tls"
	"crypto/x509"
	_ "encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
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

type IPAddress string
type IPAddressRange struct {
	min IPAddress
	max IPAddress
}
type IPAddressOrRange struct {
	addressPrefix IPAddress
	addressRange  IPAddressRange
}
type IPAddressChoice struct {
	//inherit           nil
	addressesOrRanges []IPAddressOrRange
}

func parseCer(file string) ([]byte, error) {
	//300E300C040200013006030402B9A6FC     16
	/*  ` 0x0c 0x04 0x02 0x00 0x01 0x30 0x06 0x03 0x04 0x02 0xb9 0xa6 0xfc
	300E300C040200013006030402B9A6FC
	300D300B04020001300503030084FC

	oidValue:
	0x30 0x0e      0x30是SEQUENCE类型固定的， 0e是后面长度
		0x30 0x0c  0x30是SEQUENCE类型固定的， 0c是后面长度
			0x04 0x02 0x00 0x01     0x04, 0x02, 0x00, 0x01, // address family: IPv4    对比：0x04, 0x02, 0x00, 0x02, // address family: IPv6
				0x30 0x06
					0x03 0x04
					 0x02 0xb9 0xa6 0xfc
					      185.166.252/22
	type: 48
	len: 14
	unused: 48
	oidIP:
	0x0c 0x04 0x02 0x00 0x01 0x30 0x06 0x03 0x04 0x02 0xb9 0xa6 0xfc
	`*/

	rootCa := `E:\Go\go-study\src\main\secruity\3.cer`
	caBlock, err := ioutil.ReadFile(rootCa)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(len(caBlock))

	cert, err := x509.ParseCertificate(caBlock)
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}

	fmt.Println(*cert.SerialNumber)
	fmt.Println(cert.NotBefore.Format("2006-01-02 15:04:05"))
	fmt.Println(cert.NotAfter.Format("2006-01-02 15:04:05"))
	fmt.Printf("subject: %+v\r\n", cert.Subject)

	fmt.Printf("issuer: %+v\r\n", cert.Issuer)

	fmt.Printf("Extensions: %+v\r\n", cert.Extensions)
	fmt.Printf("ExtraExtensions: %+v\r\n", cert.ExtraExtensions)
	oidKey := "1.3.6.1.5.5.7.1.7"
	var oidValue []byte
	var critical bool
	for _, extension := range cert.Extensions {
		oid := extension.Id
		if oidKey == oid.String() {
			oidValue = extension.Value
			critical = extension.Critical
			break
		}
	}
	if len(oidValue) == 0 {
		fmt.Println("not found oid:", oidValue)
		return nil, errors.New("not found oid")
	}
	fmt.Println("critical:", critical)
	fmt.Println("oidValue:")
	for _, ip := range oidValue {
		fmt.Print(fmt.Sprintf("0x%02x ", ip))
	}
	fmt.Println("")
	oidType := oidValue[0]
	oidLen := oidValue[1]
	oidRealValue := oidValue[2:]
	fmt.Print(fmt.Sprintf("type:%d (0x%02x)\r\n", oidType, oidType))
	fmt.Print(fmt.Sprintf("len::%d (0x%02x)\r\n", oidLen, oidLen))
	fmt.Println("oidRealValue:")
	for _, ip := range oidRealValue {
		fmt.Print(fmt.Sprintf("0x%02x ", ip))
	}
	return oidRealValue, nil
}

func main() {
	oidByte, err := parseCer(`E:\Go\go-study\src\main\secruity\3.cer`)
	if err != nil {
		return
	}
	for _, ob := range oidByte {
		fmt.Print(fmt.Sprintf("0x%02x ", ob))
	}
}
