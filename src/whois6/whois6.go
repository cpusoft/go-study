package main

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/openrdap/rdap"
)

func main() {
	// 1. 域名解析为 IP（RIR 是 IP 层面的概念）
	ips, err := net.LookupIP("zdns.cn")
	if err != nil || len(ips) == 0 {
		panic(err)
	}
	ip := ips[0].String()

	// 2. RDAP 查询（自动路由到对应 RIR）
	client := &rdap.Client{}
	ipNet, err := client.QueryIP(ip)
	if err != nil {
		panic(err)
	}

	// 3. 提取 RIR 与网段信息
	fmt.Printf("IP: %s\n", ip)
	fmt.Printf("RIR Handle: %s\n", ipNet.Handle) // e.g. NET-8-8-8-0-1
	fmt.Printf("Parent ParentHandle: %s\n", ipNet.ParentHandle)
	fmt.Printf("Name: %s\n", ipNet.Name) // e.g. GOGL
	fmt.Printf("StartAddress: %s\n", ipNet.StartAddress)
	fmt.Printf("EndAddress: %s\n", ipNet.EndAddress)
	rir := extractRIRFromLinks(ipNet.Links)
	fmt.Printf("IP: %s -> RIR: %s\n", ip, rir)
	fmt.Printf("Handle: %s\n", ipNet.Handle)
}
func extractRIRFromLinks(links []rdap.Link) string {
	for _, link := range links {
		fmt.Println("Link Rel:", link.Rel, " Href:", link.Href)
		if link.Rel == "self" || link.Rel == "" {
			u, err := url.Parse(link.Href)
			if err != nil {
				continue
			}
			host := strings.ToLower(u.Hostname())
			// RDAP 基础 URL 映射
			switch {
			case strings.Contains(host, "arin.net"):
				return "ARIN"
			case strings.Contains(host, "ripe.net"):
				return "RIPE"
			case strings.Contains(host, "apnic.net"):
				return "APNIC"
			case strings.Contains(host, "lacnic.net"):
				return "LACNIC"
			case strings.Contains(host, "afrinic.net"):
				return "AFRINIC"
			}
		}
	}
	return "UNKNOWN"
}
