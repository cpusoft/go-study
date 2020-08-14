package main

import (
	_ "crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/cpusoft/goutil/bitutil"
	"github.com/cpusoft/goutil/convert"
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

type IpBlock struct {
	AddressFamily []byte
	//	IPAddressRange []asn1.BitString `asn1:"optional`
	//	IPMaxMin       []IPMaxMin       `asn1:"optional`
	IPAddressRange []asn1.RawValue //[]IPMaxMin
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
				var leastZeroByte uint8 = bitutil.LeftAndFillZero(uint8(unusedBitLen - 1))
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
				var leastOneByte uint8 = bitutil.LeftAndFillOne(uint8(unusedBitLen - 1))
				//		fmt.Printf("max: leastOneByte:%x,%d,%b\n", leastOneByte, leastOneByte, leastOneByte)
				address[addressLen-1] = address[addressLen-1] | leastOneByte
				//		fmt.Printf("max: address: %x,%d,%b\n", address[addressLen-1], address[addressLen-1], address[addressLen-1])
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
				var leastZeroByte uint8 = bitutil.LeftAndFillZero(uint8(unusedBitLen - 1))
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
				var leastOneByte uint8 = bitutil.LeftAndFillOne(uint8(unusedBitLen - 1))
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
func GetIpBlocks(value []byte) (ip IpBlock, err error) {
	/*
		ipBlocks2 := make([]IpBlock2, 0)
		_, err = asn1.Unmarshal(value, &ipBlocks2)
		fmt.Println("ips2:", len(ipBlocks2), err)
		for i := range ipBlocks2 {
			fmt.Println("ips2:", i, "  AddressFamily:", convert.Bytes2Uint64(ipBlocks2[i].AddressFamily))
			fmt.Println(len(ipBlocks2[i].IPMaxMin))

		}

		fmt.Println("``````````````````````````````")
	*/
	ipBlocks := make([]IpBlock, 0)
	_, err = asn1.Unmarshal(value, &ipBlocks)
	fmt.Println("ips:", len(ipBlocks))
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

type Asn struct {
	AsnOrAsnRange asn1.RawValue `asn1:"explicit,optional,tag:0`
	Rdi           asn1.RawValue `asn1:"explicit,optional,tag:1`
}

type AsnOrAsnRange struct {
	Asn      int   `asn1:"optional`
	AsnRange []int `asn1:"optional`
}

func GetAsns(value []byte) {
	asn := Asn{}
	_, err := asn1.Unmarshal(value, &asn)
	fmt.Println("asn:", convert.PrintBytes(asn.AsnOrAsnRange.Bytes, 8), err)

	asnOrAsnRanges := make([]asn1.RawValue, 0)
	_, err = asn1.Unmarshal(asn.AsnOrAsnRange.Bytes, &asnOrAsnRanges)
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
