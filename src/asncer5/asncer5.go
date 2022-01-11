package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"

	asn1utilasn1 "github.com/cpusoft/goutil/asn1util/asn1base"
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

	ipBlocks := MakeIPs1(false)
	ipblocksExtension, err := EncodeIPAddressBlock(ipBlocks)
	fmt.Println("ipblocksExtension:", jsonutil.MarshalJson(ipblocksExtension), err)
	fmt.Println("ipblocksExtension hex:", convert.PrintBytesOneLine(ipblocksExtension.Value))
	ipblocksDec, err := DecodeIPAddressBlock(ipblocksExtension.Value)
	fmt.Println("ipblocksDec:", jsonutil.MarshalJson(ipblocksDec), err)
	fmt.Println("-------------------")
	/*
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
	*/
	hexStr := `002001067C208C`
	by, err := hex.DecodeString(hexStr)
	fmt.Println(convert.PrintBytesOneLine(by), err)
	bi, err := asn1utilasn1.ParseBitString(by)
	fmt.Println("bitlength:", bi.BitLength, "  :", convert.PrintBytesOneLine(bi.Bytes))
	fmt.Println(err)
	ipNet, err := DecodeIPNet(2, bi)
	fmt.Println(ipNet, err)

}

func DecodeIPNet(addrFamily int, addr asn1utilasn1.BitString) (*net.IPNet, error) {
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
