package main

import (
	"fmt"
	"net"
	"net/netip"

	"github.com/IncSW/geoip2"
)

// 域名 → 解析IP → 判断RIR
func getDomainRIR(domain string) (string, error) {
	// 1. 域名解析IP
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", err
	}
	if len(ips) == 0 {
		return "", fmt.Errorf("no IP found")
	}

	// 取第一个IPv4
	var targetIP net.IP
	for _, ip := range ips {
		if ip.To4() != nil {
			targetIP = ip
			break
		}
	}
	if targetIP == nil {
		return "", fmt.Errorf("no IPv4 found")
	}

	// 2. 快速判断IP属于哪个RIR
	ip, _ := netip.AddrFromSlice(targetIP.To4())
	rir := geoip2.DetectRIR(ip)

	return rir.String(), nil
}

func main() {
	domains := []string{
		"baidu.com",
		"github.com",
		"google.com",
		"facebook.com",
		"amazon.co.uk",
	}

	for _, d := range domains {
		rir, err := getDomainRIR(d)
		if err != nil {
			fmt.Printf("%-15s %v\n", d, err)
			continue
		}
		fmt.Printf("%-15s RIR: %s\n", d, rir)
	}
}
