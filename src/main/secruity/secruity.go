package main

import (
	_ "crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	_ "encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/PromonLogicalis/asn1"
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

func parseOid(data []byte) {
	ctx := asn1.NewContext()

	// Use BER for encoding and decoding.
	ctx.SetDer(false, false)

	// Add a CHOICE
	/*
		ctx.AddChoice("value", []asn1.Choice{
			{
				Type:    reflect.TypeOf(""),
				Options: "tag:0",
			},
			{
				Type:    reflect.TypeOf(int(0)),
				Options: "tag:1",
			},
		})
	*/

	type Message struct {
		Id    int
		Value interface{} `asn1:"choice:value"`
	}

	// Encode
	/*
		msg := Message{
			Id:    1000,
			Value: "this is a value",
		}
	*/
	msg := Message{
		Id:    1000,
		Value: 999,
	}
	data, err := ctx.Encode(msg)
	if err != nil {
		fmt.Println(err)
	}
	for _, d := range data {
		fmt.Print(fmt.Sprintf("0x%02x ", d))
	}
	// Decode
	decodedMsg := Message{}
	_, err = ctx.Decode(data, &decodedMsg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("%+v\n", decodedMsg)
}

func parseCer(file string) ([]byte, error) {
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

	for _, extension := range cert.Extensions {
		oid := extension.Id
		if oidKey == oid.String() {
			parseExtension(extension)
			fmt.Println("new extension~~~~~~~~~~~~~~~~~~~~~~~~")
		}
	}
	return nil, nil
}
func parseExtension(extension pkix.Extension) ([]byte, error) {
	oidValue := extension.Value
	critical := extension.Critical
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
	/*
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
	*/
	//0x30 0x0c  0x30是SEQUENCE类型固定的， 0c是后面长度, 从这里开始
	ipAddrBlocksType := oidRealValue[0]
	ipAddrBlocksLen := oidRealValue[1]
	ipAddrBlocksValue := oidRealValue[2:]
	fmt.Print(fmt.Sprintf("ipAddrBlocksType:%d (0x%02x)\r\n", ipAddrBlocksType, ipAddrBlocksType))
	fmt.Print(fmt.Sprintf("ipAddrBlocksLen:%d (0x%02x)\r\n", ipAddrBlocksLen, ipAddrBlocksLen))

	//0x04 0x02 0x00 0x01     0x04, 0x02, 0x00, 0x01, // address family: IPv4    对比：0x04, 0x02, 0x00, 0x02, // address family: IPv6
	addressFamilyType := ipAddrBlocksValue[0]
	addressFamilyLen := ipAddrBlocksValue[1]
	addressFamilyValue := ipAddrBlocksValue[2 : 2+addressFamilyLen]
	fmt.Print(fmt.Sprintf("addressFamilyType:%d (0x%02x)\r\n", addressFamilyType, addressFamilyType))
	fmt.Print(fmt.Sprintf("addressFamilyLen:%d (0x%02x)\r\n", addressFamilyLen, addressFamilyLen))
	fmt.Println("addressFamilyValue:")
	for _, af := range addressFamilyValue {
		fmt.Print(fmt.Sprintf("0x%02x ", af))
	}
	const (
		ipv4 = 0x01
		ipv6 = 0x02
	)
	var iptype int
	if addressFamilyValue[addressFamilyLen-1] == ipv4 {
		iptype = ipv4
	} else if addressFamilyValue[addressFamilyLen-1] == ipv6 {
		iptype = ipv6
	} else {
		return nil, errors.New("error iptype")
	}
	fmt.Println("iptype:", iptype)
	/*
		0x30 0x06
			0x03 0x04
				 0x02 0xb9 0xa6 0xfc
			      185.166.252/22
	*/
	IPAddressChoice := ipAddrBlocksValue[2+addressFamilyLen:]
	fmt.Println("IPAddressChoice:")
	for _, ia := range IPAddressChoice {
		fmt.Print(fmt.Sprintf("0x%02x ", ia))
	}
	fmt.Println("")

	const (
		nul      = 0x05
		sequence = 0x30
	)
	var ipaddressChoiceType int
	if IPAddressChoice[0] == nul {
		ipaddressChoiceType = nul
	} else if IPAddressChoice[0] == sequence {
		ipaddressChoiceType = sequence
	} else {
		return nil, errors.New("error iptype")
	}
	fmt.Println(fmt.Sprintf("ipaddressChoiceType: 0x%02x", ipaddressChoiceType))

	if ipaddressChoiceType == nul {
		return nil, nil
	}

	return nil, nil
}

func main() {
	oidByte, err := parseCer(`E:\Go\go-study\src\main\secruity\1.cer`)
	if err != nil {
		return
	}
	for _, ob := range oidByte {
		fmt.Print(fmt.Sprintf("0x%02x ", ob))
	}
}
