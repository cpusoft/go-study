package main

import (
	_ "crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	_ "encoding/pem"
	"errors"
	"fmt"
	"github.com/PromonLogicalis/asn1"
	"io/ioutil"
	"strconv"
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
		fmt.Println(err)
		return err
	}
	fmt.Println(len(caBlock))

	cert, err := x509.ParseCertificate(caBlock)
	if err != nil {
		fmt.Println("err:", err)
		return err
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
			err := parseExtension(extension)
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Println("new extension~~~~~~~~~~~~~~~~~~~~~~~~")
		}
	}
	return nil
}
func parseExtension(extension pkix.Extension) error {
	extensionValue := extension.Value
	critical := extension.Critical
	if len(extensionValue) == 0 {
		fmt.Println("not found oid:", extensionValue)
		return errors.New("not found oid")
	}
	fmt.Println("critical:", critical)
	printBytes("extensionValue:", extensionValue)
	// sequences 整个的数组
	ipAddrBlocksType := extensionValue[0]
	ipAddrBlocksLen := extensionValue[1]
	ipAddrBlocksValue := extensionValue[2:]
	printAsn("ipAddrBlocks", ipAddrBlocksType, ipAddrBlocksLen, ipAddrBlocksValue)

	tmpBlock := ipAddrBlocksValue
	//循环数组，
	for {
		fmt.Println("new block==========================")
		ipAddrBlockType := tmpBlock[0]
		ipAddrBlockLen := tmpBlock[1]
		ipAddrBlockValue := tmpBlock[2 : 2+ipAddrBlockLen]
		printAsn("ipAddrBlock", ipAddrBlockType, ipAddrBlockLen, ipAddrBlockValue)
		err := parseIpAddrBlock(tmpBlock)
		if err != nil {
			fmt.Println(err)
			return err
		}
		tmpBlock = tmpBlock[2+ipAddrBlockLen:]
		if len(tmpBlock) == 0 {
			break
		}
	}
	fmt.Println("end ==========================")
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

func parseIpAddrBlock(ipAddrBlock []byte) error {
	//ipaddressFamily: 包括addressFamily+ipAddressChoice
	ipAddressFamilyType := ipAddrBlock[0]
	ipAddressFamilyLen := ipAddrBlock[1]
	ipAddressFamilyValue := ipAddrBlock[2 : 2+ipAddressFamilyLen]
	printAsn("ipAddressFamily", ipAddressFamilyType, ipAddressFamilyLen, ipAddressFamilyValue)

	//读取addressFamily，注意是从ipAddrBlock开始截取的
	addressFamilyType := ipAddrBlock[2]
	addressFamilyLen := ipAddrBlock[3]
	addressFamilyValue := ipAddrBlock[4 : 4+addressFamilyLen]
	printAsn("addressFamily", addressFamilyType, addressFamilyLen, addressFamilyValue)
	var ipType int
	if addressFamilyValue[addressFamilyLen-1] == ipv4 {
		ipType = ipv4
	} else if addressFamilyValue[addressFamilyLen-1] == ipv6 {
		ipType = ipv6
	} else {
		return errors.New("error iptype")
	}
	fmt.Println("get ipType from addressFamilyValue (ipv4 = 0x01,   ipv6 = 0x02): ", ipType)

	//读取ipAddressChoice，注意是从ipAddrBlock开始--即addressFamilyValue节尾--截取的
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
		fmt.Println("new addressesOrRange--------------")
		addressesOrRangeType := tmpAddressOrRange[0]
		addressesOrRangeLen := tmpAddressOrRange[1]
		addressesOrRangeValue := tmpAddressOrRange[2 : 2+addressesOrRangeLen]
		printAsn("addressesOrRange", addressesOrRangeType, addressesOrRangeLen, addressesOrRangeValue)
		err := parseAddressesOrRange(tmpAddressOrRange, ipType)
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

func parseAddressesOrRange(addressesOrRange []byte, ipType int) error {
	//每个addressesOrRange有可能是addressPrefix(0x03开头)，也有可能是addressRange(0x30开头)
	if len(addressesOrRange) == 0 {
		return errors.New("lenght of addressesOrRange is zero")
	}
	addressesOrRangeOneType := addressesOrRange[0]
	addressesOrRangeOneLen := addressesOrRange[1]
	addressesOrRangeOneValue := addressesOrRange[2 : 2+addressesOrRangeOneLen]
	printAsn("addressesOrRangeOne", addressesOrRangeOneType, addressesOrRangeOneLen, addressesOrRangeOneValue)
	// 注意这里传入的是addressesOrRangeOneValue
	if addressesOrRangeOneType == bitstring {
		parseAddressPrefix(addressesOrRangeOneValue, addressesOrRangeOneLen, ipType)
	} else if addressesOrRangeOneType == sequence {
		parseAddressRange(addressesOrRangeOneValue, addressesOrRangeOneLen, ipType)
	} else {
		return errors.New("addressesOrRangeOneType is error")
	}
	return nil
}

func parseAddressPrefix(addressPrefix []byte, addressesOrRangeOneLen byte, ipType int) error {
	// 传入的第0位是unusedbit位
	// 03 03 04 b010              addressPrefix    172.16/12
	//第2位标明长度，-1后(unused bit占用了1位)，为ip地址应该的长度: 标明长度3，应该长度为2
	// 第3位，固定的unused bit位： 为4
	// unusedbit = 32- 应该的长度*8 - prefix  =32-2*8-prefix
	// prefix = 32- 应该的长度*8 - unusedbit = 32 - 2*8 - 4 = 12
	fmt.Println("parseAddressPrefix():  ipType:", ipType)
	printBytes("addressPrefix:", addressPrefix)
	addressShouldLen, _ := strconv.Atoi(fmt.Sprintf("%d", addressesOrRangeOneLen))
	unusedBitLen, _ := strconv.Atoi(fmt.Sprintf("%d", addressPrefix[0]))

	address := addressPrefix[1:]
	if ipType == ipv4 {
		// ipv4 的CIDR 表示法
		prefix := ipv4len - 8*(addressShouldLen-1) - unusedBitLen
		fmt.Println(fmt.Sprintf("prefix := ipv4len - 8*(addressShouldLen-1) - unusedBitLen:  %d := %d - 8 *(%d-1)-  %d \r\n",
			prefix, ipv4len, addressShouldLen, unusedBitLen))

		printBytes("address:", address)

		ipv4Address := ""
		for i := 0; i < len(address); i++ {
			ipv4Address += fmt.Sprintf("%d", address[i])
			if i < len(address)-1 {
				ipv4Address += "."
			}
		}
		ipv4Address += "/" + fmt.Sprintf("%d", prefix)
		fmt.Println(ipv4Address)
	} else if ipType == ipv6 {
		// ipv6的前缀表示法，和ipv4不一样
		prefix := 8*(addressShouldLen-1) - unusedBitLen
		fmt.Println(fmt.Sprintf("prefix :=  8*(addressShouldLen-1) - unusedBitLen:  %d := 8 *(%d-1)-  %d \r\n",
			prefix, addressShouldLen, unusedBitLen))

		printBytes("address:", address)

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
		fmt.Println(ipv6Address)
	}

	return nil
}
func parseAddressRange(addressRange []byte, addressesOrRangeOneLen byte, ipType int) error {
	//传入的是两个sequence，第一个是min，第二个是max
	// Value值，跳过了unused bit位，所以是从3开始，并且长度-1
	fmt.Println("parseAddressRange():  ipType:", ipType)
	minType := addressRange[0]
	minLen := addressRange[1]
	minValue := addressRange[2 : 2+minLen]
	printAsn("min", minType, minLen, minValue)

	tmp := addressRange[2+minLen:]
	maxType := tmp[0]
	maxLen := tmp[1]
	maxValue := tmp[2 : 2+maxLen]
	printAsn("max", maxType, maxLen, maxValue)

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

		fmt.Println("minAddr:", minAddr, "maxAddr", maxAddr)

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
		fmt.Println("minAddr:", minAddr, "maxAddr", maxAddr)

	}
	return nil
}

func main() {
	err := parseCer(`E:\Go\go-study\src\main\secruity\range_ipv6.cer`)
	if err != nil {
		return
	}
	var ar = []byte{0x01, 0x20, 0x01, 0x00, 0x00, 0x02}
	parseAddressPrefix(ar, 0x6, ipv6)
}

func printAsn(name string, typ byte, ln byte, byt []byte) {
	fmt.Println(fmt.Sprintf(name+"Type:0x%02x (%d)", typ, typ))
	fmt.Println(fmt.Sprintf(name+"Len:0x%02x (%d)", ln, ln))
	printBytes(name+"Value:", byt)
}

func printBytes(name string, byt []byte) {
	fmt.Println(name)
	for _, i := range byt {
		fmt.Print(fmt.Sprintf("0x%02x ", i))
	}
	fmt.Println("")
}
