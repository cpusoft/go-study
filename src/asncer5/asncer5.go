package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"

	asn1base "github.com/cpusoft/goutil/asn1util/asn1base"
	"github.com/cpusoft/goutil/convert"
	"github.com/cpusoft/goutil/jsonutil"
)

// https://github.com/cloudflare/cfrpki
func main() {
	/*
		hexStr := `04010230090307002001067C208C`
		b, err := hex.DecodeString(hexStr)
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		ips, err := DecodeIPAddressBlock(b)
		fmt.Println("ips:", jsonutil.MarshalJson(ips), err)
	*/
	/*  ipaddress
	ipBlocks := MakeIPs1(false)
	ipblocksExtension, err := EncodeIPAddressBlock(ipBlocks)
	fmt.Println("ipblocksExtension:", jsonutil.MarshalJson(ipblocksExtension), err)
	fmt.Println("ipblocksExtension hex:", convert.PrintBytesOneLine(ipblocksExtension.Value))
	ipblocksDec, err := DecodeIPAddressBlock(ipblocksExtension.Value)
	fmt.Println("ipblocksDec:", jsonutil.MarshalJson(ipblocksDec), err)
	fmt.Println("-------------------")

		ipblocksDec, err := DecodeIPAddressBlock(ipblocksExtension.Value)
		fmt.Println("ipblocksDec hex:", convert.PrintBytesOneLine(ipblocksExtension.Value))
		fmt.Println("ipblocksDec:", jsonutil.MarshalJson(ipblocksDec), err)

		n, err := asn1parse.ParseBytes(ipblocksExtension.Value)
		fmt.Println("node:", jsonutil.MarshalJson(n), err)
	*/
	/*
			hexStr := `3010300E04010230090307002001067C208C`
			by, err := hex.DecodeString(hexStr)
			ipblocksDec, err = DecodeIPAddressBlock(by)
			fmt.Println("ipblocksDec2:", jsonutil.MarshalJson(ipblocksDec), err)
			n, err := asn1parse.ParseBytes(by)
			fmt.Println("node2:", jsonutil.MarshalJson(n), err)

		hexStr := `002001067C208C`
		by, err := hex.DecodeString(hexStr)
		fmt.Println(convert.PrintBytesOneLine(by), err)
		bi, err := asn1base.ParseBitString(by)
		fmt.Println("bitlength:", bi.BitLength, "  :", convert.PrintBytesOneLine(bi.Bytes))
		fmt.Println(err)
		ipNet, err := DecodeIPNet(2, bi)
		fmt.Println(ipNet, err)
	*/

	// iprange

	ipBlocks := MakeIPs2(false)
	ipblocksExtension, err := EncodeIPAddressBlock(ipBlocks)
	fmt.Println("ipblocksExtension:", jsonutil.MarshalJson(ipblocksExtension), err)
	fmt.Println("ipblocksExtension hex:", convert.PrintBytesOneLine(ipblocksExtension.Value))
	ipblocksDec, err := DecodeIPAddressBlock(ipblocksExtension.Value)
	fmt.Println("ipblocksDec:", jsonutil.MarshalJson(ipblocksDec), err)
	fmt.Println("-------------------")

	/*
		hexStr := `300C03040067F5A503040067F5A6`
		hexStr:=
		by, err := hex.DecodeString(hexStr)

	*/
	hexStr := `300e030500c0a80001030500c0a80003`
	by, err := hex.DecodeString(hexStr)
	fmt.Println(convert.PrintBytesOneLine(by))
	type AddrRangeBase struct {
		Min asn1base.BitString
		Max asn1base.BitString
	}
	var addrRangeBase AddrRangeBase
	_, err = asn1base.Unmarshal(by, &addrRangeBase)
	fmt.Println("addrRangeBase:", convert.PrintBytesOneLine(by), jsonutil.MarshalJson(addrRangeBase), err)
	ipNetMin, err := DecodeIPNet(1, addrRangeBase.Min)
	ipNetMax, err := DecodeIPNet(1, addrRangeBase.Max)
	fmt.Println("ipNetMin:", ipNetMin.String(), "  ipNetMax:", ipNetMax.String(), err)

	//min, err := DecodeIPNet(2, addrRange.Min)
	//max, err := DecodeIPNet(2, addrRange.Max)
	//fmt.Println("addrRange:", min.String(), max.String(), err)
	/*
		//30 31 30 20 04 02 00 01 30 1A 03 04 02 2D  40 B8 03 04 02 67 1B C8 30 0C 03 04 00 67 F5 A5  03 04 00 67 F5 A6 30 0D 04 02 00 02 30 07 03 05  00 24 07 79 00
			hexStr = `3031302004020001301A0304022D40B8030402671BC8300C03040067F5A503040067F5A6300D04020002300703050024077900`
			by, err = hex.DecodeString(hexStr)
			fmt.Println(convert.PrintBytesOneLine(by), err)
			ipblocksDec, err = DecodeIPAddressBlock(by)
			fmt.Println("ipblocksDec:", jsonutil.MarshalJson(ipblocksDec), err)
	*/
}

func DecodeIPNet(addrFamily int, addr asn1base.BitString) (*net.IPNet, error) {
	var size int
	if addrFamily == 1 {
		size = 4
	} else if addrFamily == 2 {
		size = 16
	} else {
		return nil, errors.New("Not an IP address")
	}
	ipaddr := make([]byte, size)
	fmt.Println(len(ipaddr), "addr.Bytes:", convert.PrintBytesOneLine(addr.Bytes))
	copy(ipaddr, addr.Bytes)
	fmt.Println("ipaddr:", convert.PrintBytesOneLine(ipaddr))
	mask := net.CIDRMask(addr.BitLength, size*8)
	fmt.Println("mask:", mask)
	return &net.IPNet{
		IP:   net.IP(ipaddr),
		Mask: mask,
	}, nil

}

func MakeIPs(null bool) []IPCertificateInformation {
	if null {
		return []IPCertificateInformation{
			&IPAddressNull{
				Family: 1,
			},
		}
	}

	_, net1, _ := net.ParseCIDR("0.0.0.0/0")
	_, net2, _ := net.ParseCIDR("::/0")
	ip1 := net.ParseIP("192.168.0.1")
	ip2 := net.ParseIP("192.168.0.3")

	return []IPCertificateInformation{
		&IPNet{
			IPNet: net1,
		},
		&IPNet{
			IPNet: net2,
		},
		&IPAddressRange{
			Min: ip1,
			Max: ip2,
		},
		//&IPAddressNull{Family: 1,},
	}
}

func MakeIPs1(null bool) []IPCertificateInformation {
	if null {
		return []IPCertificateInformation{
			&IPAddressNull{
				Family: 1,
			},
		}
	}

	addr1, net1, err := net.ParseCIDR("192.168.100.1/24")
	fmt.Println("Make:  addr1:", addr1, net1, err)
	return []IPCertificateInformation{
		&IPNet{
			IPNet: net1,
		},

		//&IPAddressNull{Family: 1,},
	}
}

func MakeIPs2(null bool) []IPCertificateInformation {
	if null {
		return []IPCertificateInformation{
			&IPAddressNull{
				Family: 1,
			},
		}
	}

	ip1 := net.ParseIP("192.168.0.1")
	ip2 := net.ParseIP("192.168.0.3")

	return []IPCertificateInformation{

		&IPAddressRange{
			Min: ip1,
			Max: ip2,
		},
		//&IPAddressNull{Family: 1,},
	}
}
