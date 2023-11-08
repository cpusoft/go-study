package main

import (
	_ "crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/cpusoft/goutil/asn1util/asn1cert"
	"github.com/cpusoft/goutil/bitutil"
	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
)

type Certificate struct {
	//Raw                asn1.RawContent
	TBSCertificate     TbsCertificate
	SignatureAlgorithm AlgorithmIdentifier
	SignatureValue     asn1.BitString
}

type TbsCertificate struct {
	Raw                asn1.RawContent
	Version            int `asn1:"optional,explicit,default:0,tag:0"`
	SerialNumber       *big.Int
	SignatureAlgorithm AlgorithmIdentifier
	Issuer             asn1.RawValue
	Validity           Validity
	Subject            asn1.RawValue
	PublicKey          PublicKeyInfo
	Extensions         []Extension `asn1:"optional,explicit,tag:3"`
}
type AlgorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue `asn1:"optional"`
}
type Validity struct {
	NotBefore, NotAfter time.Time
}
type PublicKeyInfo struct {
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

// CRLDistributionPoints ::= SEQUENCE SIZE (1..MAX) OF DistributionPoint
//
//	DistributionPoint ::= SEQUENCE {
//	    distributionPoint       [0]     DistributionPointName OPTIONAL,
//	    reasons                 [1]     ReasonFlags OPTIONAL,
//	    cRLIssuer               [2]     GeneralNames OPTIONAL }
//
//	DistributionPointName ::= CHOICE {
//	    fullName                [0]     GeneralNames,
//	    nameRelativeToCRLIssuer [1]     RelativeDistinguishedName }
//
// RFC 5280, 4.2.1.14
type distributionPoint struct {
	DistributionPoint distributionPointName `asn1:"optional,tag:0"`
	Reason            asn1.BitString        `asn1:"optional,tag:1"`
	CRLIssuer         asn1.RawValue         `asn1:"optional,tag:2"`
}

type distributionPointName struct {
	FullName     []asn1.RawValue `asn1:"optional,tag:0"`
	RelativeName asn1.RawValue   `asn1:"optional,tag:1"`
}

func GetCrldp(value []byte) ([]string, error) {
	var cdp []distributionPoint
	_, err := asn1.Unmarshal(value, &cdp)
	if err != nil {
		return nil, err
	}

	cls := make([]string, 0)
	for _, dp := range cdp {
		// Per RFC 5280, 4.2.1.13, one of distributionPoint or cRLIssuer may be empty.
		if len(dp.DistributionPoint.FullName) == 0 {
			continue
		}

		for _, fullName := range dp.DistributionPoint.FullName {
			if fullName.Tag == 6 {
				cls = append(cls, string(fullName.Bytes))
			}
		}
	}
	return cls, nil
}

// RFC 5280 4.2.1.4
type policy struct {
	Policy asn1.ObjectIdentifier
	// policyQualifiers omitted
}

func GetPolicies(value []byte) ([]string, error) {
	policies := make([]policy, 0)
	_, err := asn1.Unmarshal(value, &policies)
	if err != nil {
		return nil, err
	}
	fmt.Println(len(policies))
	tmp := make([]string, len(policies))
	for i := range policies {
		tmp[i] = policies[i].Policy.String()
	}
	return tmp, err
}

func decodeAddressPrefixAndMinMax(ipFamily int, addressShouldLenByte []byte, unusedByte []byte, address []byte, ipAddressType string) (ipAddress string, err error) {
	addressShouldLen := int(convert.Bytes2Uint64(addressShouldLenByte))
	unusedBitLen := int(convert.Bytes2Uint64(unusedByte))
	addressLen := len(address)
	fmt.Println("decodeAddressPrefixAndMinMax(): ipFamily:", ipFamily,
		"   addressShouldLen:", addressShouldLen,
		"   unusedByte:", unusedByte, "   unusedBitLen:", unusedBitLen,
		"   address:", address, "   addressLen:", addressLen)
	if ipFamily == 1 {
		// ipv4 ipaddress prefx
		ipv4Address := ""
		if ipAddressType == "min" {
			if unusedBitLen > 0 {
				//leastZeroByte := bitutil.LeftAndFillZero(uint8(unusedBitLen - 1))
				leastZeroByte := bitutil.Shift0xffLeftFillZero(uint8(unusedBitLen - 1))
				fmt.Println("min: before address:", address[addressLen-1], "  leastZeroByte:", leastZeroByte)
				address[addressLen-1] = address[addressLen-1] & leastZeroByte
				fmt.Println("min: after address:", address[addressLen-1], "  address:", address)
			}
			for i := 0; i < 4; i++ {
				if i < len(address) {
					ipv4Address += fmt.Sprintf("%d", address[i])
				} else {
					ipv4Address += fmt.Sprintf("%d", 0)
				}
				if i < 3 {
					ipv4Address += "."
				}
			}
			fmt.Println("ipAddress min ipv4:", ipv4Address)
		} else if ipAddressType == "max" {
			if unusedBitLen > 0 {
				//leastOneByte := bitutil.LeftAndFillOne(uint8(unusedBitLen - 1))
				leastOneByte := bitutil.Shift0x00LeftFillOne(uint8(unusedBitLen - 1))
				fmt.Println("max: before address:", address[addressLen-1], "  leastOneByte:", leastOneByte)
				address[addressLen-1] = address[addressLen-1] | leastOneByte
				fmt.Println("max: after address:", address[addressLen-1], "  address:", address)
			}
			for i := 0; i < 4; i++ {
				if i < len(address) {
					ipv4Address += fmt.Sprintf("%d", address[i])
				} else {
					ipv4Address += fmt.Sprintf("%d", 0xff)
				}
				if i < 3 {
					ipv4Address += "."
				}
			}
			fmt.Println("ipAddress max ipv4:", ipv4Address)
		}

	} else if ipFamily == 2 {

		ipv6Address := ""
		if ipAddressType == "min" {
			if unusedBitLen > 0 {
				//var leastZeroByte uint8 = bitutil.LeftAndFillZero(uint8(unusedBitLen - 1))
				var leastZeroByte uint8 = bitutil.Shift0xffLeftFillZero(uint8(unusedBitLen - 1))
				//		fmt.Printf("min: leastZeroByte:%x,%d,%b\n", leastZeroByte, leastZeroByte, leastZeroByte)
				address[addressLen-1] = address[addressLen-1] & leastZeroByte
				//		fmt.Printf("min: address: %x,%d,%b\n", address[addressLen-1], address[addressLen-1], address[addressLen-1])
			}
			var i int
			for i = 0; i < 16; i++ {
				//	fmt.Printf("min: %d,%d\n", i, addressLen)
				// 2a01:8::
				if i < addressLen {
					ipv6Address += fmt.Sprintf("%02x", address[i])
					//		fmt.Printf("min: ipv6Address:%d, %s\n", i, ipv6Address)
					if i%2 == 1 {
						ipv6Address += ":"
						//			fmt.Printf("min: ipv6Address +: %d, %s\n", i, ipv6Address)
					}
				} else {
					break
				}
			}
			// end with ::
			if i == len(address) {
				ipv6Address += ":"
			}
			fmt.Printf("ipAddress min Ipv6:%s\n", ipAddress)
		} else if ipAddressType == "max" {
			if unusedBitLen > 0 {
				//var leastOneByte uint8 = bitutil.LeftAndFillOne(uint8(unusedBitLen - 1))
				var leastOneByte uint8 = bitutil.Shift0x00LeftFillOne(uint8(unusedBitLen - 1))
				//		fmt.Printf("max: leastOneByte:%x,%d,%b\n", leastOneByte, leastOneByte, leastOneByte)
				address[addressLen-1] = address[addressLen-1] | leastOneByte
				//		fmt.Printf("max: address: %x,%d,%b\n", address[addressLen-1], address[addressLen-1], address[addressLen-1])
			}
			for i := 0; i < 16; i++ {
				//2a01:17:ffff:ffff:ffff:ffff:ffff:ffff
				if i < len(address) {
					ipv6Address += fmt.Sprintf("%02x", address[i])
					if i%2 == 1 {
						ipv6Address += ":"
					}
				} else {
					ipv6Address += "ff"
					if i%2 == 1 && i < 15 {
						ipv6Address += ":"
					}
				}
			}
			fmt.Printf("ipAddress max Ipv6:%s\n", ipAddress)
		} else if ipAddressType == "ipPrefix" {
			// ipv6 ipaddress prefx
			prefix := 8*(addressShouldLen-1) - unusedBitLen
			//	fmt.Println("prefix :=  8*(addressShouldLen-1) - unusedBitLen:  %d := 8 *(%d-1)-  %d \r\n",
			//		prefix, addressShouldLen, unusedBitLen)

			//printBytes("address:", address)

			for i := 0; i < len(address); i++ {
				ipv6Address += fmt.Sprintf("%02x", address[i])
				if i%2 == 1 && i < len(address) {
					ipv6Address += ":"
				}
			}
			//Complement digit
			if len(address)%2 == 1 {
				ipv6Address += "00"
			}
			ipv6Address += "/" + fmt.Sprintf("%d", prefix)
			fmt.Printf("ipAddress ipPrefix Ipv6:%s\n", ipv6Address)
		}
	}
	return
}

// ipFamily: ipv4:1, ipv6:2
// ipAddressType: min, max, ipPrefix
func decodeAddressPrefix(ipFamily int, addressShouldLenByte []byte, unusedByte []byte, address []byte, ipAddressType string) (ipAddress string, err error) {
	addressShouldLen := int(convert.Bytes2Uint64(addressShouldLenByte))
	addressLen := len(address)
	unusedBitLen := int(convert.Bytes2Uint64(unusedByte))
	//fmt.Println("ipFamily:", ipFamily, "   addressShouldLenByte:", addressShouldLenByte, "  unusedByte:", unusedByte,
	//	"   addressShouldLen:", addressShouldLen, "   unusedBitLen:", unusedBitLen, "   addressLen:", addressLen)

	if ipFamily == 1 {
		// ipv4 ipaddress prefx
		ipv4Address := ""
		if ipAddressType == "min" {
			if unusedBitLen > 0 {
				//var leastZeroByte uint8 = bitutil.LeftAndFillZero(uint8(unusedBitLen - 1))
				var leastZeroByte uint8 = bitutil.Shift0xffLeftFillZero(uint8(unusedBitLen - 1))
				//		fmt.Printf("min: leastZeroByte:%x,%d,%b\n", leastZeroByte, leastZeroByte, leastZeroByte)
				address[addressLen-1] = address[addressLen-1] & leastZeroByte
				//		fmt.Printf("min: address: %x,%d,%b\n", address[addressLen-1], address[addressLen-1], address[addressLen-1])
			}
			for i := 0; i < 4; i++ {
				if len(address) > i {
					ipv4Address += fmt.Sprintf("%d", address[i])
				} else {
					ipv4Address += fmt.Sprintf("%d", 0)
				}
				if i < 3 {
					ipv4Address += "."
				}
			}
			fmt.Printf("ipAddress min ipv4:%s\n", ipv4Address)
		} else if ipAddressType == "max" {
			if unusedBitLen > 0 {
				//var leastOneByte uint8 = bitutil.LeftAndFillOne(uint8(unusedBitLen - 1))
				var leastOneByte uint8 = bitutil.Shift0x00LeftFillOne(uint8(unusedBitLen - 1))
				fmt.Printf("max: leastOneByte:%x,%d,%b\n", leastOneByte, leastOneByte, leastOneByte)
				address[addressLen-1] = address[addressLen-1] | leastOneByte
				fmt.Printf("max: address: %x,%d,%b\n", address[addressLen-1], address[addressLen-1], address[addressLen-1])
			}
			for i := 0; i < 4; i++ {
				if len(address) > i {
					ipv4Address += fmt.Sprintf("%d", address[i])
				} else {
					ipv4Address += fmt.Sprintf("%d", 0xff)
				}
				if i < 3 {
					ipv4Address += "."
				}
			}
			fmt.Printf("ipAddress max ipv4:%s\n", ipv4Address)
		} else if ipAddressType == "ipPrefix" {
			for i := 0; i < len(address); i++ {
				ipv4Address += fmt.Sprintf("%d", address[i])
				if i < len(address)-1 {
					ipv4Address += "."
				}
			}
			prefix := 8*(addressShouldLen-1) - unusedBitLen
			//	fmt.Println("prefix := 8*(addressShouldLen-1) - unusedBitLen:  %d := 8 *(%d-1)-  %d \r\n",
			//		prefix, addressShouldLen, unusedBitLen)
			ipv4Address += "/" + fmt.Sprintf("%d", prefix)
			fmt.Printf("ipAddress ipPrefix ipv4:%s\n", ipv4Address)

		}

	} else if ipFamily == 2 {

		ipv6Address := ""
		if ipAddressType == "min" {
			if unusedBitLen > 0 {
				//var leastZeroByte uint8 = bitutil.LeftAndFillZero(uint8(unusedBitLen - 1))
				var leastZeroByte uint8 = bitutil.Shift0xffLeftFillZero(uint8(unusedBitLen - 1))
				//		fmt.Printf("min: leastZeroByte:%x,%d,%b\n", leastZeroByte, leastZeroByte, leastZeroByte)
				address[addressLen-1] = address[addressLen-1] & leastZeroByte
				//		fmt.Printf("min: address: %x,%d,%b\n", address[addressLen-1], address[addressLen-1], address[addressLen-1])
			}
			var i int
			for i = 0; i < 16; i++ {
				//	fmt.Printf("min: %d,%d\n", i, addressLen)
				// 2a01:8::
				if i < addressLen {
					ipv6Address += fmt.Sprintf("%02x", address[i])
					//		fmt.Printf("min: ipv6Address:%d, %s\n", i, ipv6Address)
					if i%2 == 1 {
						ipv6Address += ":"
						//			fmt.Printf("min: ipv6Address +: %d, %s\n", i, ipv6Address)
					}
				} else {
					break
				}
			}
			// end with ::
			if i == len(address) {
				ipv6Address += ":"
			}
			fmt.Printf("ipAddress min Ipv6:%s\n", ipAddress)
		} else if ipAddressType == "max" {
			if unusedBitLen > 0 {
				//var leastOneByte uint8 = bitutil.LeftAndFillOne(uint8(unusedBitLen - 1))
				var leastOneByte uint8 = bitutil.Shift0x00LeftFillOne(uint8(unusedBitLen - 1))
				//		fmt.Printf("max: leastOneByte:%x,%d,%b\n", leastOneByte, leastOneByte, leastOneByte)
				address[addressLen-1] = address[addressLen-1] | leastOneByte
				//		fmt.Printf("max: address: %x,%d,%b\n", address[addressLen-1], address[addressLen-1], address[addressLen-1])
			}
			for i := 0; i < 16; i++ {
				//2a01:17:ffff:ffff:ffff:ffff:ffff:ffff
				if i < len(address) {
					ipv6Address += fmt.Sprintf("%02x", address[i])
					if i%2 == 1 {
						ipv6Address += ":"
					}
				} else {
					ipv6Address += "ff"
					if i%2 == 1 && i < 15 {
						ipv6Address += ":"
					}
				}
			}
			fmt.Printf("ipAddress max Ipv6:%s\n", ipAddress)
		} else if ipAddressType == "ipPrefix" {
			// ipv6 ipaddress prefx
			prefix := 8*(addressShouldLen-1) - unusedBitLen
			//	fmt.Println("prefix :=  8*(addressShouldLen-1) - unusedBitLen:  %d := 8 *(%d-1)-  %d \r\n",
			//		prefix, addressShouldLen, unusedBitLen)

			//printBytes("address:", address)

			for i := 0; i < len(address); i++ {
				ipv6Address += fmt.Sprintf("%02x", address[i])
				if i%2 == 1 && i < len(address) {
					ipv6Address += ":"
				}
			}
			//Complement digit
			if len(address)%2 == 1 {
				ipv6Address += "00"
			}
			ipv6Address += "/" + fmt.Sprintf("%d", prefix)
			fmt.Printf("ipAddress ipPrefix Ipv6:%s\n", ipv6Address)
		}
	}
	return
}

type IPMaxMin struct {
	IpMin asn1.BitString
	IpMax asn1.BitString
}

/*
ok

	type IpBlock struct {
		AddressFamily []byte
		//	IPAddressRange []asn1.BitString `asn1:"optional`
		//	IPMaxMin       []IPMaxMin       `asn1:"optional`
		IPAddressRange []asn1.RawValue //[]IPMaxMin
	}
*/
type IpBlock struct {
	AddressFamily       []byte
	IpAddressBitStrings []asn1.BitString `asn1:"optional,tag:16,class:0`
	//IPMaxMin []IPMaxMin `asn1:"optional`
	//IPAddressRange []asn1.RawValue //[]IPMaxMin
}

type IpBlockRaw struct {
	AddressFamily []byte
	IPAddressRaws []asn1.RawValue `asn1:"optional,tag:16,class:0` //[]IPMaxMin
	//	IPAddressRaws IPMaxMin      `asn1:"optional,tag:16,class:0`
}

/*
ips: 0   AddressFamily: 1
ipAddress ipPrefix ipv4:103.121.40/22
ipFamily: 1  <nil>
ips: 1   AddressFamily: 2
ipAddress ipPrefix Ipv6:2403:63c0:/32
ipFamily: 2  <nil>

	sbgp-ipAddrBlock: critical
	    IPv4:
	      143.137.108.0/22
	      168.181.76.0/22
	      170.150.160.0/22
	      170.244.120.0/22
	      170.245.184.0/22
	      186.194.140.0/22
	      200.53.128.0/18
	      200.57.80.0/20
	      200.77.224.0/20
	      201.159.128.0/20
	      201.175.0.0-201.175.47.255
	    IPv6:
	      2001:1270::/32

	sbgp-autonomousSysNum: critical
	    Autonomous System Numbers:
	      22011
	      22908

	type RawValue struct {
	        Class, Tag int
	        IsCompound bool
	        Bytes      []byte
	        FullBytes  []byte // 包括标签和长度
	}
*/
func GetIpBlocks(value []byte) (ip IpBlock, err error) {

	ipBlockRaws := make([]IpBlockRaw, 0)
	_, err = asn1.Unmarshal(value, &ipBlockRaws)
	fmt.Println("GetIpBlocks():len(ipBlockRaws):", len(ipBlockRaws), "  value:", convert.PrintBytesOneLine(value))
	for _, ipBlockRaws := range ipBlockRaws {
		ipAddressRaws := ipBlockRaws.IPAddressRaws
		fmt.Println("range ipAddressRaws:", ipAddressRaws) //, "  value:", convert.PrintBytesOneLine(value), "  ipBlockRaws", ipBlockRaws)
		for _, ipAddressRaw := range ipAddressRaws {
			fmt.Println("\r\n=============\r\nrange ipAddressRaw: ipAddressRaw:", ipAddressRaw, "\r\n\t class:", ipAddressRaw.Class, " tag:", ipAddressRaw.Tag,
				"  IsCompound:", ipAddressRaw.IsCompound,
				"  Bytes:", convert.PrintBytesOneLine(ipAddressRaw.Bytes),
				"  FullBytes:", convert.PrintBytesOneLine(ipAddressRaw.FullBytes))
			if !ipAddressRaw.IsCompound {
				ipAddressBitString := asn1.BitString{}
				_, err = asn1.Unmarshal(ipAddressRaw.FullBytes, &ipAddressBitString)
				fmt.Println("-----get ipAddressBitString:", ipAddressBitString)
			} else {
				ipAddressMinMax := make([]asn1.BitString, 0)
				_, err = asn1.Unmarshal(ipAddressRaw.FullBytes, &ipAddressMinMax)
				fmt.Println("-----get ipAddressMinMax:", ipAddressMinMax)

				ipMaxMin := IPMaxMin{}
				_, err = asn1.Unmarshal(ipAddressRaw.FullBytes, &ipMaxMin)
				fmt.Println("-----get ipMaxMin:", ipMaxMin)

				ipAddressMin := ipMaxMin.IpMin
				unusedLenMin := (len(ipAddressMin.Bytes))*8 - ipAddressMin.BitLength
				//fillZero := bitutil.LeftAndFillZero(uint8(4*8 - len(ipAddressMin.Bytes)*8 + unusedLenMin))
				fillZero := bitutil.Shift0xffLeftFillZero(uint8(4*8 - len(ipAddressMin.Bytes)*8 + unusedLenMin))
				min := ipAddressMin.Bytes // & fillZero
				fmt.Println("    ipAddressMin:", ipAddressMin, "   unusedLenMin:", unusedLenMin,
					"   Bytes:", ipAddressMin.Bytes, " BitLength:", ipAddressMin.BitLength,
					"   fillZero:", fillZero, "   min:", min)
				//min, _ := decodeAddressPrefixAndMinMax(1, ipAddressMin.Bytes[0:1], ipAddressMin.Bytes[1:2], ipAddressMin.Bytes[2:], "min")

				ipAddressMax := ipMaxMin.IpMax
				unusedLenMax := (len(ipAddressMax.Bytes))*8 - ipAddressMax.BitLength
				//fillOne := bitutil.LeftAndFillOne(uint8(4*8 - len(ipAddressMin.Bytes)*8 + unusedLenMax))
				fillOne := bitutil.Shift0x00LeftFillOne(uint8(4*8 - len(ipAddressMin.Bytes)*8 + unusedLenMax))
				max := ipAddressMax.Bytes // | fillZero
				fmt.Println("    ipAddressMax:", ipAddressMax, "   unusedLenMax:", unusedLenMax,
					"   Bytes:", ipAddressMax.Bytes, " BitLength:", ipAddressMax.BitLength,
					"   fillOne:", fillOne, "   max:", max)
				//max, _ := decodeAddressPrefixAndMinMax(1, ipAddressMax.Bytes[0:1], ipAddressMax.Bytes[1:2], ipAddressMax.Bytes[2:], "max")
				//fmt.Println("    max:", max)

			}
		}

	}

	ipBlocks := make([]IpBlock, 0)
	_, err = asn1.Unmarshal(value, &ipBlocks)
	fmt.Println("GetIpBlocks():len(ipBlocks):", len(ipBlocks), "  value:", convert.PrintBytesOneLine(value), "  ipBlocks", ipBlocks)

	for i := range ipBlocks {
		ipFamily := int(convert.Bytes2Uint64(ipBlocks[i].AddressFamily))
		ipAddressBitStrings := ipBlocks[i].IpAddressBitStrings

		fmt.Println("range ipBlocks, i:", i, "  AddressFamily:", ipFamily, "  ipAddressBitStrings:", ipAddressBitStrings)

		for _, ipAddrBlock := range ipAddressBitStrings {
			fmt.Println("range ipAddrBlocks: ipAddrBlock:", ipAddrBlock,
				"\r\n\tipAddrBlock.Bytes:", convert.PrintBytesOneLine(ipAddrBlock.Bytes),
				"   ipAddrBlock.BitLength :", ipAddrBlock.BitLength,
				"   ipAddrBlock.RightAlign:", convert.PrintBytesOneLine(ipAddrBlock.RightAlign()))
		}

	}

	return
}

/*
	ok

func GetIpBlocks(value []byte) (ip IpBlock, err error) {

		ipBlocks := make([]IpBlock, 0)
		_, err = asn1.Unmarshal(value, &ipBlocks)
		fmt.Println("GetIpBlocks():len(ipBlocks):", len(ipBlocks), "  value:", convert.PrintBytesOneLine(value))
		for i := range ipBlocks {
			ipFamily := int(convert.Bytes2Uint64(ipBlocks[i].AddressFamily))
			fmt.Println("ips:", i, "  AddressFamily:", ipFamily)
			for j := range ipBlocks[i].IPAddressRange {
				//fmt.Println(convert.PrintBytes(ipBlocks[i].IPAddressRange[j].Bytes, 8))

				ipAddress := asn1.BitString{}
				_, err = asn1.Unmarshal(ipBlocks[i].IPAddressRange[j].FullBytes, &ipAddress)
				if len(ipAddress.Bytes) > 0 {
					//fmt.Println("fullBytes:", convert.PrintBytes(ipBlocks[i].IPAddressRange[j].FullBytes, 8))
					//fmt.Println("Bytes:", convert.PrintBytes(ipAddress.Bytes, 8))
					unused := ipBlocks[i].IPAddressRange[j].FullBytes[2:3]
					addressShouldLength := ipBlocks[i].IPAddressRange[j].FullBytes[1:2]
					addressPrefix, err := decodeAddressPrefix(ipFamily, addressShouldLength, unused, ipAddress.Bytes, "ipPrefix")
					fmt.Println("ipFamily:", ipFamily, addressPrefix, err)
				}

				ipAddresses := make([]asn1.RawValue, 0)
				_, err = asn1.Unmarshal(ipBlocks[i].IPAddressRange[j].FullBytes, &ipAddresses)
				if len(ipAddresses) > 0 {
					fmt.Println(len(ipAddresses))
					x := 0
					//fmt.Println("min :", convert.PrintBytes(ipAddresses[x].FullBytes, 8))
					unused := ipAddresses[x].FullBytes[2:3]
					addressShouldLength := ipAddresses[x].FullBytes[1:2]
					ipAddress := asn1.BitString{}
					_, err = asn1.Unmarshal(ipAddresses[x].FullBytes, &ipAddress)
					addressPrefix, err := decodeAddressPrefix(ipFamily, addressShouldLength, unused, ipAddress.Bytes, "min")
					fmt.Println("min : ipFamily:", ipFamily, addressPrefix, err)

					x++
					//fmt.Println("max :", convert.PrintBytes(ipAddresses[x].FullBytes, 8))
					unused = ipAddresses[x].FullBytes[2:3]
					addressShouldLength = ipAddresses[x].FullBytes[1:2]
					ipAddress = asn1.BitString{}
					_, err = asn1.Unmarshal(ipAddresses[x].FullBytes, &ipAddress)

					addressPrefix, err = decodeAddressPrefix(ipFamily, addressShouldLength, unused, ipAddress.Bytes, "max")
					fmt.Println("max : ipFamily:", ipFamily, addressPrefix, err)

				}

			}
		}

		return ip, nil
	}
*/
type Asn struct {
	AsnOrAsnRange asn1.RawValue `asn1:"explicit,optional,tag:0`
	Rdi           asn1.RawValue `asn1:"explicit,optional,tag:1`
}

type AsnOrAsnRange struct {
	Asn      int   `asn1:"optional`
	AsnRange []int `asn1:"optional`
}

type AsnBlock struct {
	Asn []asn1.RawValue
}
type AsnInt struct {
	Asns []uint64
}

func GetAsns(value []byte) {
	fmt.Println("GetAsns(): value:", convert.PrintBytesOneLine(value))

	asnRaws := make([]asn1.RawValue, 0)
	_, err := asn1.Unmarshal(value, &asnRaws)
	fmt.Println("GetAsns(): asn1 Unmarshal asnRaws:", asnRaws, " asnRoaws:", jsonutil.MarshalJson(asnRaws),
		"  len(asnRaws):", len(asnRaws), err)
	for i, asnRawi := range asnRaws {
		fmt.Println("GetAsns(): range asnRaws, i:", i, " asnRawi:", asnRawi, jsonutil.MarshalJson(asnRawi), err)
		if !asnRawi.IsCompound {
			var asn int
			_, err = asn1.Unmarshal(asnRawi.Bytes, &asn)
			fmt.Println("GetAsns(): not IsCompound, range asnRaws, i:", i, " asn:", asn, jsonutil.MarshalJson(asn), err)
		} else {
			if asnRawi.Class != 2 {
				asns := make([]int, 0)
				_, err = asn1.Unmarshal(asnRawi.Bytes, &asns)
				fmt.Println("GetAsns(): IsCompound, Class!=2, range asnRaws, i:", i, " asns:", asns, jsonutil.MarshalJson(asns), err)
			} else {
				asnRawis := make([]asn1.RawValue, 0)
				_, err = asn1.Unmarshal(asnRawi.Bytes, &asnRawis)
				fmt.Println("GetAsns():  Class==2, range asnRaws, i:", i, " asnRawis:", asnRawis, jsonutil.MarshalJson(asnRawis), err)
				for j, asnRawj := range asnRawis {
					fmt.Println("\r\nGetAsns(): get asnRawj:", asnRawj, jsonutil.MarshalJson(asnRawj))
					if asnRawj.IsCompound {
						asns := make([]uint64, 0)
						_, err = asn1.Unmarshal(asnRawj.Bytes, &asns)
						fmt.Println("GetAsns(): IsCompound, i:", i, " j:", j, " asns:", asns, jsonutil.MarshalJson(asns), err)
					} else {
						var asn AsnInt
						_, err = asn1.Unmarshal(asnRawj.Bytes, &asn)
						fmt.Println("GetAsns(): not IsCompound, i:", i, " j:", j, " asn:", asn, jsonutil.MarshalJson(asn), err)

					}
				}
			}
		}
	}

	asnBlocks := make([]AsnBlock, 0)
	_, err = asn1.Unmarshal(value, &asnBlocks)
	fmt.Println("GetAsns(): asn1 Unmarshal asnBlocks:", jsonutil.MarshalJson(asnBlocks), err)
	for _, asnBlock := range asnBlocks {
		fmt.Println("GetAsns(): range asnBlocks, asnBlock:", jsonutil.MarshalJson(asnBlock))
	}
	/*
		asnOrAsnRanges := make([]asn1.RawValue, 0)
		_, err = asn1.Unmarshal(value, &asnOrAsnRanges)
		//fmt.Println("asnOrAsnRanges:", asnOrAsnRanges, err)
		for i := range asnOrAsnRanges {
			//fmt.Println("asn, i:", i, asnOrAsnRanges[i].Class, asnOrAsnRanges[i].Tag, asnOrAsnRanges[i].IsCompound, convert.PrintBytes(asnOrAsnRanges[i].Bytes, 8))
			if !asnOrAsnRanges[i].IsCompound {
				asn := convert.Bytes2Uint64(asnOrAsnRanges[i].Bytes)
				fmt.Println("asn:", asn)

			} else {
				asns := make([]asn1.RawValue, 0)
				_, err = asn1.Unmarshal(asnOrAsnRanges[i].FullBytes, &asns)
				//	fmt.Println("len(asns):", len(asns), err)
				for j := range asns {
					asn := convert.Bytes2Uint64(asns[j].Bytes)
					fmt.Println("asns:", asn)
				}
			}
		}
	*/
	return
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
	files := []string{
		`F:\share\我的坚果云\Go\common\go-study\src\asncer4\00Z.cer`,
		`F:\share\我的坚果云\Go\common\go-study\src\asncer4\c8c59.cer`,
		`F:\share\我的坚果云\Go\common\go-study\src\asncer4\75414d.cer`,
		`F:\share\我的坚果云\Go\common\go-study\src\asncer4\034644.cer`,
	}
	for _, file := range files {
		b, err := fileutil.ReadFileToBytes(file)
		if err != nil {
			fmt.Println(file, err)
			return
		}
		certificate := Certificate{}
		_, err = asn1.Unmarshal(b, &certificate)
		//fmt.Println("certificate:", jsonutil.MarshallJsonIndent(certificate), len(res), err)
		//fmt.Println(len(certificate.TBSCertificate.Extensions))
		for i := range certificate.TBSCertificate.Extensions {
			extension := &certificate.TBSCertificate.Extensions[i]
			fmt.Println(extension.Oid.String())
			if extension.Oid.String() == "2.5.29.14" {
				// subjectKeyIdentifier
				fmt.Println("2.5.29.14:")
				fmt.Println(GetOctectString(extension.Value))
			} else if extension.Oid.String() == "2.5.29.35" {
				// authorityKeyIdentifier
				fmt.Println("2.5.29.35:")
				fmt.Println(GetOctectStringSequenceString(extension.Value))
			} else if extension.Oid.String() == "2.5.29.19" {
				// basicConstraints
				fmt.Println("2.5.29.19:", extension.Critical)
				fmt.Println("2.5.29.19:")
				fmt.Println(GetOctectStringSequenceBool(extension.Value))
			} else if extension.Oid.String() == "2.5.29.15" {
				// keyUsage
				fmt.Println("2.5.29.15", extension.Critical)

				usageValue, err := GetOctectStringBitString(extension.Value)
				fmt.Println("2.5.29.15:", usageValue, err)

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
				fmt.Println("1.3.6.1.5.5.7.1.1:", len(seqs), err)
				for i := range seqs {
					fmt.Println(seqs[i].Oid, string(seqs[i].Value))
				}
			} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.11" {
				// subjectInfoAccess
				seqs, err := GetOctectStringSequenceOidString(extension.Value)
				fmt.Println("1.3.6.1.5.5.7.1.11:", len(seqs), err)
				for i := range seqs {
					fmt.Println(seqs[i].Oid, string(seqs[i].Value))
				}
			} else if extension.Oid.String() == "2.5.29.31" {
				// Crl
				seqs, err := GetCrldp(extension.Value)
				fmt.Println("2.5.29.31:", seqs, err)
			} else if extension.Oid.String() == "2.5.29.32" {
				// Policies
				seqs, err := GetPolicies(extension.Value)
				fmt.Println("2.5.29.32:", seqs, err)
			} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.7" {
				// IpBlocks
				//seqs, err := GetIpBlocks(extension.Value)
				ipAddrBlocks, err := asn1cert.ParseToIpAddressBlocks(extension.Value)
				fmt.Println("1.3.6.1.5.5.7.1.7:", jsonutil.MarshalJson(ipAddrBlocks), err)
			} else if extension.Oid.String() == "1.3.6.1.5.5.7.1.8" {
				// Asns
				fmt.Println("1.3.6.1.5.5.7.1.8:")
				GetAsns(extension.Value)

			}
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
