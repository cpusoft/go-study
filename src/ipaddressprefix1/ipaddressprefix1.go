package main

import (
	"encoding/hex"
	"fmt"
	"net"

	"github.com/cpusoft/goutil/convert"
)

func main() {
	// asn1cer5
	/*
		_, net1, _ := net.ParseCIDR("0.0.0.0/0")
		_, net2, _ := net.ParseCIDR("::/0")
		ip1 := net.ParseIP("192.168.0.1")
		ip2 := net.ParseIP("192.168.0.3")

		ips := []IPCertificateInformation{
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
		}

		for _, ip := range ips {
			ipBytes, err := ip.ASN1()
			if err != nil {
				fmt.Println(err)
			}

			asn1R := asn1.RawValue{FullBytes: ipBytes}
			fmt.Println(convert.PrintBytesOneLine(ipBytes))
			fmt.Println(asn1R)
		}
	*/
	_, net1, _ := net.ParseCIDR("10.32.0.0/12")
	_, net2, _ := net.ParseCIDR("10.64.0.0/16")
	_, net3, _ := net.ParseCIDR("10.1.0.0/16")
	_, net4, _ := net.ParseCIDR("2001:0:2::/47")
	_, net5, _ := net.ParseCIDR("2001:0:200::/39")
	_, net6, _ := net.ParseCIDR("2a05:6680::/29") // 03 05 03 2A 05 66 80
	_, net7, _ := net.ParseCIDR("2a0f:c1c0::/32") // 03 05 00 2A 0F C1 C0
	ipInfos := []IPCertificateInformation{
		&IPNet{
			IPNet: net1,
		},
		&IPNet{
			IPNet: net2,
		},
		&IPNet{
			IPNet: net3,
		},
		&IPNet{
			IPNet: net4,
		},
		&IPNet{
			IPNet: net5,
		},
		&IPNet{
			IPNet: net6,
		},
		&IPNet{
			IPNet: net7,
		},
	}
	var ipBytes []byte
	var err error
	for _, ip := range ipInfos {
		ipBytes, err = ip.ASN1()
		fmt.Println(hex.Dump(ipBytes))
		fmt.Println(convert.PrintBytesOneLine(ipBytes), err)
	}

	_, net8, _ := net.ParseCIDR("2a0f:c1c0::/32") // 03 05 00 2A 0F C1 C0
	ipNet8 := &IPNet{
		IPNet: net8,
	}

	ipBytes, err = ipNet8.ASN1()
	fmt.Println(convert.PrintBytesOneLine(ipBytes), err)
}
